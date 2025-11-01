package repositories

import (
	"restaurant-backend/internal/models"

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
func (r *MenuItemRepository) GetByID(id uint) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	if err := r.db.First(&menuItem, id).Error; err != nil {
		return nil, err
	}
	return &menuItem, nil
}

// GetByMenuID retrieves all menu items for a menu (RLS ensures tenant isolation)
func (r *MenuItemRepository) GetByMenuID(menuID uint) ([]models.MenuItem, error) {
	var menuItems []models.MenuItem
	if err := r.db.Where("menu_id = ?", menuID).Find(&menuItems).Error; err != nil {
		return nil, err
	}
	return menuItems, nil
}

// GetByRestaurantID retrieves all menu items for a restaurant (RLS ensures tenant isolation)
func (r *MenuItemRepository) GetByRestaurantID(restaurantID uint) ([]models.MenuItem, error) {
	var menuItems []models.MenuItem
	if err := r.db.Where("restaurant_id = ?", restaurantID).Find(&menuItems).Error; err != nil {
		return nil, err
	}
	return menuItems, nil
}

// Update updates an existing menu item
func (r *MenuItemRepository) Update(menuItem *models.MenuItem) error {
	return r.db.Save(menuItem).Error
}

// Delete deletes a menu item
func (r *MenuItemRepository) Delete(id uint) error {
	return r.db.Delete(&models.MenuItem{}, id).Error
}

