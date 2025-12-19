package services

import (
	"context"
	"errors"
	"fmt"

	"restaurant-backend/internal/dto"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService handles user management operations
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// ListUsers retrieves all users for a restaurant
func (s *UserService) ListUsers(ctx context.Context, restaurantID uint) ([]models.User, error) {
	users, err := s.userRepo.GetByRestaurantIDWithContext(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Clear password hashes
	for i := range users {
		users[i].PasswordHash = ""
	}

	return users, nil
}

// GetUser retrieves a user by ID for a specific restaurant
func (s *UserService) GetUser(ctx context.Context, id uint, restaurantID uint) (*models.User, error) {
	user, err := s.userRepo.GetByIDWithContext(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify user belongs to the restaurant (multi-tenancy check)
	if user.RestaurantID != restaurantID {
		return nil, errors.New("user not found")
	}

	// Clear password hash
	user.PasswordHash = ""

	return user, nil
}

// CreateUser creates a new user with validation and password hashing
func (s *UserService) CreateUser(ctx context.Context, createDTO *dto.CreateUserDTO, restaurantID uint) (*models.User, error) {
	// Validate role (KAM not allowed here)
	if createDTO.Role == "KAM" {
		return nil, errors.New("KAM role cannot be created through this endpoint")
	}

	// Check email uniqueness within restaurant
	existingUser, err := s.userRepo.GetByEmailWithContext(ctx, createDTO.Email, restaurantID)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists in this restaurant")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Set defaults for optional fields
	timezone := createDTO.Timezone
	if timezone == "" {
		timezone = "UTC"
	}

	language := createDTO.Language
	if language == "" {
		language = "en"
	}

	preferences := createDTO.Preferences
	if preferences == "" {
		preferences = "{}"
	}

	// Create user
	user := &models.User{
		RestaurantID: restaurantID,
		Email:        createDTO.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    createDTO.FirstName,
		LastName:     createDTO.LastName,
		Role:         createDTO.Role,
		Phone:        createDTO.Phone,
		Timezone:     timezone,
		Language:     language,
		Preferences:  preferences,
		IsActive:     true,
	}

	if err := s.userRepo.CreateWithContext(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Clear password hash before returning
	user.PasswordHash = ""

	return user, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, id uint, updateDTO *dto.UpdateUserDTO, restaurantID uint) (*models.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByIDWithContext(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify user belongs to the restaurant
	if user.RestaurantID != restaurantID {
		return nil, errors.New("user not found")
	}

	// Validate role if provided (KAM not allowed)
	if updateDTO.Role != "" && updateDTO.Role == "KAM" {
		return nil, errors.New("cannot change role to KAM")
	}

	// Update fields
	if updateDTO.FirstName != "" {
		user.FirstName = updateDTO.FirstName
	}
	if updateDTO.LastName != "" {
		user.LastName = updateDTO.LastName
	}
	if updateDTO.Role != "" {
		user.Role = updateDTO.Role
	}
	if updateDTO.Phone != "" {
		user.Phone = updateDTO.Phone
	}
	if updateDTO.Timezone != "" {
		user.Timezone = updateDTO.Timezone
	}
	if updateDTO.Language != "" {
		user.Language = updateDTO.Language
	}
	if updateDTO.Preferences != "" {
		user.Preferences = updateDTO.Preferences
	}

	// Save updated user
	if err := s.userRepo.UpdateWithContext(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Clear password hash
	user.PasswordHash = ""

	return user, nil
}

// DeleteUser deletes a user (soft delete)
func (s *UserService) DeleteUser(ctx context.Context, id uint, restaurantID uint) error {
	// Get user to verify ownership
	user, err := s.userRepo.GetByIDWithContext(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify user belongs to the restaurant
	if user.RestaurantID != restaurantID {
		return errors.New("user not found")
	}

	// Delete user
	if err := s.userRepo.DeleteWithContext(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ToggleUserStatus toggles the active status of a user
func (s *UserService) ToggleUserStatus(ctx context.Context, id uint, restaurantID uint, isActive bool) error {
	// Get user to verify ownership
	user, err := s.userRepo.GetByIDWithContext(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify user belongs to the restaurant
	if user.RestaurantID != restaurantID {
		return errors.New("user not found")
	}

	// Update status
	if err := s.userRepo.UpdateUserStatus(ctx, id, isActive); err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	return nil
}
