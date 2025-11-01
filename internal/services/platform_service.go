package services

import (
	"errors"
	"fmt"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repositories"
)

// PlatformService handles platform organization and KAM management
type PlatformService struct {
	restaurantRepo *repositories.RestaurantRepository
	userRepo       *repositories.UserRepository
}

// NewPlatformService creates a new PlatformService instance
func NewPlatformService(
	restaurantRepo *repositories.RestaurantRepository,
	userRepo *repositories.UserRepository,
) *PlatformService {
	return &PlatformService{
		restaurantRepo: restaurantRepo,
		userRepo:       userRepo,
	}
}

// InitializePlatformOrganization creates the platform organization if it doesn't exist
func (s *PlatformService) InitializePlatformOrganization() error {
	// Check if platform organization already exists
	platform, err := s.restaurantRepo.GetByID(models.PlatformOrganizationID)
	if err == nil && platform != nil {
		return nil // Already exists
	}

	// Create platform organization
	platform = &models.Restaurant{
		ID:          models.PlatformOrganizationID,
		Name:        "Platform Organization",
		Description: "Platform-level organization for KAM and system administrators",
		Status:      models.RestaurantStatusActive,
		IsActive:    true,
		Email:       "platform@system.local",
	}

	if err := s.restaurantRepo.Create(platform); err != nil {
		return fmt.Errorf("failed to create platform organization: %w", err)
	}

	return nil
}

// CreateKAMRequest represents KAM creation request
type CreateKAMRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// CreateKAM creates a new KAM user (only by existing KAMs/Admins)
func (s *PlatformService) CreateKAM(req *CreateKAMRequest, createdBy uint) (*models.User, error) {
	// Verify creator is KAM or Admin from platform organization
	creator, err := s.userRepo.GetByID(createdBy)
	if err != nil {
		return nil, errors.New("creator user not found")
	}

	if !creator.IsPlatformUser() || (creator.Role != "KAM" && creator.Role != "Admin") {
		return nil, errors.New("only platform KAMs or Admins can create new KAM users")
	}

	// Check if user already exists
	existing, _ := s.userRepo.GetByEmail(req.Email, models.PlatformOrganizationID)
	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create KAM user in platform organization
	user := &models.User{
		RestaurantID: models.PlatformOrganizationID,
		Email:        req.Email,
		PasswordHash: "", // Will be set by calling service
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         "KAM",
		IsActive:     true,
	}

	// Note: Password hashing should be done by calling handler/service
	return user, nil
}

// CreateKAMUser creates a KAM user in the database (password should already be hashed)
func (s *PlatformService) CreateKAMUser(user *models.User) error {
	// Verify this is a KAM user for platform organization
	if user.RestaurantID != models.PlatformOrganizationID {
		return errors.New("KAM users must belong to platform organization")
	}
	if user.Role != "KAM" {
		return errors.New("only KAM role allowed for platform organization")
	}

	// Check if user already exists
	existing, _ := s.userRepo.GetByEmail(user.Email, models.PlatformOrganizationID)
	if existing != nil {
		return errors.New("user with this email already exists")
	}

	// Create user via repository
	return s.userRepo.Create(user)
}

// ListKAMs lists all KAM users
func (s *PlatformService) ListKAMs() ([]models.User, error) {
	users, err := s.userRepo.GetByRestaurantID(models.PlatformOrganizationID)
	if err != nil {
		return nil, err
	}

	// Filter for KAM role
	kams := make([]models.User, 0)
	for _, user := range users {
		if user.Role == "KAM" {
			kams = append(kams, user)
		}
	}

	return kams, nil
}

