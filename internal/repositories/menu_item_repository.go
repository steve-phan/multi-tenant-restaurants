package repositories

import (
	"restaurant-backend/internal/models"
	"strings"

	"gorm.io/gorm"
)

// MenuItemRepository handles menu item-related database operations
type MenuItemRepository struct {
	db *gorm.DB
}

// NewMenuItemRepository creates a new MenuItemRepository instance
func NewMenuItemRepository(db *gorm.DB) *MenuItemRepository {
	return &MenuItemRepository{db: db}
}

// Create creates a new menu item
func (r *MenuItemRepository) Create(menuItem *models.MenuItem) error {
	return r.db.Create(menuItem).Error
}

// GetByID retrieves a menu item by ID (RLS ensures tenant isolation)
// Includes images ordered by display_order
func (r *MenuItemRepository) GetByID(id uint) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	if err := r.db.Preload("Images").
		Preload("Category").
		First(&menuItem, id).Error; err != nil {
		return nil, err
	}
	// Sort images manually (primary first, then by display_order)
	return &menuItem, nil
}

func (r *MenuItemRepository) GetByName(name string) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	if err := r.db.Where("lower(name) = lower(?)", strings.TrimSpace(name)).First(&menuItem).Error; err != nil {
		return nil, err
	}
	return &menuItem, nil
}

// GetByIDPublic retrieves a menu item by ID for public access (no auth required)
// Requires restaurant_id to ensure proper access
func (r *MenuItemRepository) GetByIDPublic(id uint, restaurantID uint) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	if err := r.db.Where("id = ? AND restaurant_id = ?", id, restaurantID).
		Preload("Images").
		Preload("Category").
		First(&menuItem).Error; err != nil {
		return nil, err
	}
	return &menuItem, nil
}

// GetByCategoryID retrieves all menu items for a category (RLS ensures tenant isolation)
// Ordered by display_order, includes images
func (r *MenuItemRepository) GetByCategoryID(categoryID uint) ([]models.MenuItem, error) {
	var menuItems []models.MenuItem
	if err := r.db.Where("category_id = ?", categoryID).
		Preload("Images").
		Order("display_order ASC").Find(&menuItems).Error; err != nil {
		return nil, err
	}
	return menuItems, nil
}

// GetByRestaurantID retrieves all menu items for a restaurant (RLS ensures tenant isolation)
// Includes images for each item
func (r *MenuItemRepository) GetByRestaurantID(restaurantID uint) ([]models.MenuItem, error) {
	var menuItems []models.MenuItem
	if err := r.db.Where("restaurant_id = ?", restaurantID).
		Preload("Images").
		Preload("Category").
		Order("category_id, display_order ASC").
		Find(&menuItems).Error; err != nil {
		return nil, err
	}
	return menuItems, nil
}

// Update updates an existing menu item using provided updates map (only updates fields in the map)
func (r *MenuItemRepository) Update(id uint, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil // Nothing to update
	}
	return r.db.Model(&models.MenuItem{}).Where("id = ?", id).Updates(updates).Error
}

// Delete deletes a menu item
func (r *MenuItemRepository) Delete(id uint) error {
	return r.db.Delete(&models.MenuItem{}, id).Error
}
