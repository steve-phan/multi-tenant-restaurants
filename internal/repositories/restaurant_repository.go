package repositories

import (
	"restaurant-backend/internal/models"

	"gorm.io/gorm"
)

// RestaurantRepository handles restaurant-related database operations
type RestaurantRepository struct {
	db *gorm.DB
}

// NewRestaurantRepository creates a new RestaurantRepository instance
func NewRestaurantRepository(db *gorm.DB) *RestaurantRepository {
	return &RestaurantRepository{db: db}
}

// Create creates a new restaurant
func (r *RestaurantRepository) Create(restaurant *models.Restaurant) error {
	// Ensure ID is zero so GORM uses auto-increment
	// This prevents duplicate key errors if ID was accidentally set
	if restaurant.ID != 0 && restaurant.ID != models.PlatformOrganizationID {
		// Only allow ID=1 for platform organization
		restaurant.ID = 0
	}
	return r.db.Create(restaurant).Error
}

// GetByID retrieves a restaurant by ID
func (r *RestaurantRepository) GetByID(id uint) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	if err := r.db.Preload("KAM").First(&restaurant, id).Error; err != nil {
		return nil, err
	}
	return &restaurant, nil
}

// GetByEmail retrieves a restaurant by email
func (r *RestaurantRepository) GetByEmail(email string) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	if err := r.db.Where("email = ?", email).First(&restaurant).Error; err != nil {
		return nil, err
	}
	return &restaurant, nil
}

// List retrieves all restaurants (for KAM/Admin use)
func (r *RestaurantRepository) List(status *models.RestaurantStatus, kamID *uint) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	query := r.db

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if kamID != nil {
		query = query.Where("kam_id = ?", *kamID)
	}

	if err := query.Preload("KAM").Order("created_at DESC").Find(&restaurants).Error; err != nil {
		return nil, err
	}

	return restaurants, nil
}

// ListPending retrieves all pending restaurants
func (r *RestaurantRepository) ListPending() ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	status := models.RestaurantStatusPending
	if err := r.db.Where("status = ?", status).
		Preload("KAM").
		Order("created_at ASC").
		Find(&restaurants).Error; err != nil {
		return nil, err
	}
	return restaurants, nil
}

// Update updates an existing restaurant
func (r *RestaurantRepository) Update(restaurant *models.Restaurant) error {
	return r.db.Save(restaurant).Error
}

// Delete deletes a restaurant (soft delete by setting status)
func (r *RestaurantRepository) Delete(id uint) error {
	return r.db.Model(&models.Restaurant{}).Where("id = ?", id).
		Update("status", models.RestaurantStatusSuspended).Error
}
