package repositories

import (
	"context"
	"restaurant-backend/internal/models"

	"gorm.io/gorm"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// CreateWithContext creates a new user using the provided context
func (r *UserRepository) CreateWithContext(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by ID (RLS ensures tenant isolation)
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByIDWithContext retrieves a user by ID using the provided context
func (r *UserRepository) GetByIDWithContext(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email and restaurant ID
func (r *UserRepository) GetByEmail(email string, restaurantID uint) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ? AND restaurant_id = ?", email, restaurantID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmailWithContext retrieves a user by email and restaurant ID using the provided context
func (r *UserRepository) GetByEmailWithContext(ctx context.Context, email string, restaurantID uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ? AND restaurant_id = ?", email, restaurantID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmailGlobalWithContext retrieves a user by email across all restaurants (useful for login)
func (r *UserRepository) GetByEmailGlobalWithContext(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Restaurant").Where("email = ? AND is_active = ?", email, true).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByRestaurantID retrieves all users for a restaurant (RLS ensures tenant isolation)
func (r *UserRepository) GetByRestaurantID(restaurantID uint) ([]models.User, error) {
	var users []models.User
	if err := r.db.Where("restaurant_id = ?", restaurantID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetByRestaurantIDWithContext retrieves all users for a restaurant using the provided context
func (r *UserRepository) GetByRestaurantIDWithContext(ctx context.Context, restaurantID uint) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(ctx).Where("restaurant_id = ?", restaurantID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetKAMs retrieves all KAM (Key Account Manager) users
func (r *UserRepository) GetKAMs() ([]models.User, error) {
	var users []models.User
	// KAMs belong to platform organization, so we query by restaurant_id = PlatformOrganizationID
	if err := r.db.Where("role = ? AND restaurant_id = ? AND is_active = ?", "KAM", models.PlatformOrganizationID, true).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetKAMsWithContext retrieves all KAM users using the provided context
func (r *UserRepository) GetKAMsWithContext(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(ctx).Where("role = ? AND restaurant_id = ? AND is_active = ?", "KAM", models.PlatformOrganizationID, true).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// UpdateWithContext updates a user using the provided context
func (r *UserRepository) UpdateWithContext(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// DeleteWithContext deletes a user using the provided context
func (r *UserRepository) DeleteWithContext(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// UpdateUserStatus updates the is_active status of a user
func (r *UserRepository) UpdateUserStatus(ctx context.Context, id uint, isActive bool) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("is_active", isActive).Error
}

// UpdateUserPassword updates the password hash of a user
func (r *UserRepository) UpdateUserPassword(ctx context.Context, userID uint, hashedPassword string) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("password_hash", hashedPassword).Error
}

// GetByEmailAnyRestaurant checks if email exists in any restaurant (for uniqueness check)
func (r *UserRepository) GetByEmailAnyRestaurant(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
