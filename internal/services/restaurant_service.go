package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

// RestaurantService handles restaurant business logic
type RestaurantService struct {
	restaurantRepo *repositories.RestaurantRepository
	userRepo       *repositories.UserRepository
	emailService   *EmailService
}

// NewRestaurantService creates a new RestaurantService instance
func NewRestaurantService(
	restaurantRepo *repositories.RestaurantRepository,
	userRepo *repositories.UserRepository,
	emailService *EmailService,
) *RestaurantService {
	return &RestaurantService{
		restaurantRepo: restaurantRepo,
		userRepo:       userRepo,
		emailService:   emailService,
	}
}

// RegisterRestaurantRequest represents restaurant registration request
type RegisterRestaurantRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	Address      string `json:"address" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	ContactName  string `json:"contact_name" binding:"required"`
	ContactEmail string `json:"contact_email" binding:"required,email"`
	ContactPhone string `json:"contact_phone" binding:"required"`
}

// RegisterRestaurant creates a new restaurant in pending status
func (s *RestaurantService) RegisterRestaurant(ctx context.Context, req *RegisterRestaurantRequest) (*models.Restaurant, error) {
	// Check if restaurant with same email already exists
	existing, _ := s.restaurantRepo.GetByEmailWithContext(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("restaurant with this email already exists")
	}

	// Create restaurant with pending status
	// Ensure ID is zero so GORM uses auto-increment
	restaurant := &models.Restaurant{
		ID:           0, // Explicitly set to 0 to ensure auto-increment
		Name:         req.Name,
		Description:  req.Description,
		Address:      req.Address,
		Phone:        req.Phone,
		Email:        req.Email,
		Status:       models.RestaurantStatusPending,
		ContactName:  req.ContactName,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
	}

	if err := s.restaurantRepo.CreateWithContext(ctx, restaurant); err != nil {
		return nil, err
	}

	return restaurant, nil
}

// ActivateRestaurant activates a pending restaurant
// The activating user (KAM) is passed as activatedBy
// If no KAM is assigned, the activating KAM becomes the assigned KAM
// This also creates an admin user for the restaurant and sends a welcome email
func (s *RestaurantService) ActivateRestaurant(ctx context.Context, restaurantID uint, activatedBy uint) (*models.Restaurant, error) {
	// Get restaurant
	restaurant, err := s.restaurantRepo.GetByIDWithContext(ctx, restaurantID)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	// Check if restaurant is already active
	if restaurant.Status == models.RestaurantStatusActive {
		return nil, errors.New("restaurant is already active")
	}

	// Verify the activating user is a KAM
	activatingUser, err := s.userRepo.GetByIDWithContext(ctx, activatedBy)
	if err != nil || activatingUser.Role != "KAM" {
		return nil, errors.New("only KAM users can activate restaurants")
	}

	// Generate secure temporary password for admin user
	tempPassword, err := GenerateSecurePassword()
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create admin user for the restaurant
	adminUser := &models.User{
		RestaurantID: restaurant.ID,
		Email:        restaurant.ContactEmail,
		PasswordHash: string(hashedPassword),
		FirstName:    ExtractFirstName(restaurant.ContactName),
		LastName:     ExtractLastName(restaurant.ContactName),
		Role:         "Admin",
		IsActive:     true,
	}

	// Check if user with this email already exists for this restaurant
	existingUser, _ := s.userRepo.GetByEmailWithContext(ctx, restaurant.ContactEmail, restaurant.ID)
	if existingUser != nil {
		return nil, errors.New("admin user already exists for this restaurant")
	}

	// Create the admin user
	if err := s.userRepo.CreateWithContext(ctx, adminUser); err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	// Activate restaurant
	now := time.Now()
	restaurant.Status = models.RestaurantStatusActive
	restaurant.ActivatedBy = &activatedBy
	restaurant.ActivatedAt = &now

	// If no KAM is assigned yet, assign the activating KAM
	if restaurant.KAMID == nil {
		restaurant.KAMID = &activatedBy
	}

	if err := s.restaurantRepo.UpdateWithContext(ctx, restaurant); err != nil {
		return nil, err
	}

	// Send welcome email with credentials
	// Note: Email failure should not rollback the activation
	if s.emailService != nil {
		if err := s.emailService.SendRestaurantWelcomeEmail(ctx, restaurant, restaurant.ContactEmail, tempPassword); err != nil {
			// Log error but don't fail the activation
			// In production, you might want to queue this for retry
			fmt.Printf("Warning: Failed to send welcome email to %s: %v\n", restaurant.ContactEmail, err)
		}
	}

	return restaurant, nil
}

// UpdateRestaurantStatus updates the status of a restaurant
func (s *RestaurantService) UpdateRestaurantStatus(ctx context.Context, restaurantID uint, status models.RestaurantStatus) (*models.Restaurant, error) {
	restaurant, err := s.restaurantRepo.GetByIDWithContext(ctx, restaurantID)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	restaurant.Status = status

	if err := s.restaurantRepo.UpdateWithContext(ctx, restaurant); err != nil {
		return nil, err
	}

	return restaurant, nil
}

// AssignKAM assigns a Key Account Manager to a restaurant
func (s *RestaurantService) AssignKAM(ctx context.Context, restaurantID uint, kamID uint) (*models.Restaurant, error) {
	// Verify KAM exists and is a KAM
	kam, err := s.userRepo.GetByIDWithContext(ctx, kamID)
	if err != nil || kam.Role != "KAM" {
		return nil, errors.New("invalid KAM")
	}

	// Get restaurant
	restaurant, err := s.restaurantRepo.GetByIDWithContext(ctx, restaurantID)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	// Assign KAM
	restaurant.KAMID = &kamID

	if err := s.restaurantRepo.UpdateWithContext(ctx, restaurant); err != nil {
		return nil, err
	}

	return restaurant, nil
}
