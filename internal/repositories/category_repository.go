package repositories

import (
	"context"
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

// CreateWithContext creates a new category using the provided context
func (r *CategoryRepository) CreateWithContext(ctx context.Context, category *models.MenuCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// GetByID retrieves a category by ID (RLS ensures tenant isolation)
func (r *CategoryRepository) GetByID(id uint) (*models.MenuCategory, error) {
	var category models.MenuCategory
	if err := r.db.Preload("MenuItems").Order("display_order ASC").First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// GetByIDWithContext retrieves a category by ID using the provided context
func (r *CategoryRepository) GetByIDWithContext(ctx context.Context, id uint) (*models.MenuCategory, error) {
	var category models.MenuCategory
	if err := r.db.WithContext(ctx).Preload("MenuItems").Order("display_order ASC").First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// GetByName retrieves a category by name
func (r *CategoryRepository) GetByName(name string) (*models.MenuCategory, error) {
	var category models.MenuCategory
	if err := r.db.Where("lower(name) = lower(?)", name).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// GetByNameWithContext retrieves a category by name using the provided context
func (r *CategoryRepository) GetByNameWithContext(ctx context.Context, name string) (*models.MenuCategory, error) {
	var category models.MenuCategory
	if err := r.db.WithContext(ctx).Where("lower(name) = lower(?)", name).First(&category).Error; err != nil {
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

// GetByRestaurantIDWithContext retrieves categories for a restaurant using context
func (r *CategoryRepository) GetByRestaurantIDWithContext(ctx context.Context, restaurantID uint) ([]models.MenuCategory, error) {
	var categories []models.MenuCategory
	if err := r.db.WithContext(ctx).Where("restaurant_id = ?", restaurantID).
		Preload("MenuItems", "is_available = ?", true).Order("display_order ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// Update updates an existing category using provided updates map (only updates fields in the map)
func (r *CategoryRepository) Update(id uint, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil // Nothing to update
	}
	return r.db.Model(&models.MenuCategory{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateWithContext updates a category using the provided context
func (r *CategoryRepository) UpdateWithContext(ctx context.Context, id uint, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&models.MenuCategory{}).Where("id = ?", id).Updates(updates).Error
}

// Delete deletes a category
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.MenuCategory{}, id).Error
}

// DeleteWithContext deletes a category using the provided context
func (r *CategoryRepository) DeleteWithContext(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.MenuCategory{}, id).Error
}
