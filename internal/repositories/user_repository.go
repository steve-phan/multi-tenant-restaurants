package repositories

import (
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

// GetByID retrieves a user by ID (RLS ensures tenant isolation)
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
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

// GetByRestaurantID retrieves all users for a restaurant (RLS ensures tenant isolation)
func (r *UserRepository) GetByRestaurantID(restaurantID uint) ([]models.User, error) {
	var users []models.User
	if err := r.db.Where("restaurant_id = ?", restaurantID).Find(&users).Error; err != nil {
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

// Update updates an existing user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
