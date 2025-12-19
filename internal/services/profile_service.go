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

// ProfileService handles profile management operations
type ProfileService struct {
	userRepo *repositories.UserRepository
}

// NewProfileService creates a new ProfileService instance
func NewProfileService(userRepo *repositories.UserRepository) *ProfileService {
	return &ProfileService{
		userRepo: userRepo,
	}
}

// GetProfile retrieves the current user's profile
func (s *ProfileService) GetProfile(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.userRepo.GetByIDWithContext(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// Clear password hash
	user.PasswordHash = ""

	return user, nil
}

// UpdateProfile updates the current user's profile
func (s *ProfileService) UpdateProfile(ctx context.Context, userID uint, updateDTO *dto.UpdateProfileDTO) (*models.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByIDWithContext(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields
	if updateDTO.FirstName != "" {
		user.FirstName = updateDTO.FirstName
	}
	if updateDTO.LastName != "" {
		user.LastName = updateDTO.LastName
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

	// Save updated user
	if err := s.userRepo.UpdateWithContext(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Clear password hash
	user.PasswordHash = ""

	return user, nil
}

// ChangePassword changes the current user's password
func (s *ProfileService) ChangePassword(ctx context.Context, userID uint, changeDTO *dto.ChangePasswordDTO) error {
	// Get user with password hash
	user, err := s.userRepo.GetByIDWithContext(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(changeDTO.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changeDTO.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	if err := s.userRepo.UpdateUserPassword(ctx, userID, string(hashedPassword)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// UpdatePreferences updates the current user's preferences
func (s *ProfileService) UpdatePreferences(ctx context.Context, userID uint, prefsDTO *dto.UpdatePreferencesDTO) error {
	// Get existing user
	user, err := s.userRepo.GetByIDWithContext(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Update preferences
	user.Preferences = prefsDTO.Preferences

	// Save updated user
	if err := s.userRepo.UpdateWithContext(ctx, user); err != nil {
		return fmt.Errorf("failed to update preferences: %w", err)
	}

	return nil
}

// UpdateAvatar updates the current user's avatar URL
func (s *ProfileService) UpdateAvatar(ctx context.Context, userID uint, avatarURL string) error {
	// Get existing user
	user, err := s.userRepo.GetByIDWithContext(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Update avatar URL
	user.AvatarURL = avatarURL

	// Save updated user
	if err := s.userRepo.UpdateWithContext(ctx, user); err != nil {
		return fmt.Errorf("failed to update avatar: %w", err)
	}

	return nil
}
