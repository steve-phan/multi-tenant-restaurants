package repositories

import (
	"restaurant-backend/internal/models"

	"gorm.io/gorm"
)

// CategoryRepository handles menu category-related database operations
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new CategoryRepository instance
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create creates a new menu category
func (r *CategoryRepository) Create(category *models.MenuCategory) error {
	return r.db.Create(category).Error
}

// GetByID retrieves a category by ID (RLS ensures tenant isolation)
func (r *CategoryRepository) GetByID(id uint) (*models.MenuCategory, error) {
	var category models.MenuCategory
	if err := r.db.Preload("MenuItems").Order("display_order ASC").First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// GetByRestaurantID retrieves all categories for a restaurant (RLS ensures tenant isolation)
// Ordered by display_order
func (r *CategoryRepository) GetByRestaurantID(restaurantID uint) ([]models.MenuCategory, error) {
	var categories []models.MenuCategory
	if err := r.db.Where("restaurant_id = ?", restaurantID).
		Preload("MenuItems", "is_available = ?", true).Order("display_order ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// Update updates an existing category
func (r *CategoryRepository) Update(category *models.MenuCategory) error {
	return r.db.Save(category).Error
}

// Delete deletes a category
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.MenuCategory{}, id).Error
}

