package repositories

import (
	"restaurant-backend/internal/models"

	"gorm.io/gorm"
)

// MenuRepository handles menu-related database operations
type MenuRepository struct {
	db *gorm.DB
}

// NewMenuRepository creates a new MenuRepository instance
func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

// Create creates a new menu
func (r *MenuRepository) Create(menu *models.Menu) error {
	return r.db.Create(menu).Error
}

// GetByID retrieves a menu by ID (RLS ensures tenant isolation)
func (r *MenuRepository) GetByID(id uint) (*models.Menu, error) {
	var menu models.Menu
	if err := r.db.Preload("MenuItems").First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

// GetByRestaurantID retrieves all menus for a restaurant (RLS ensures tenant isolation)
func (r *MenuRepository) GetByRestaurantID(restaurantID uint) ([]models.Menu, error) {
	var menus []models.Menu
	if err := r.db.Where("restaurant_id = ?", restaurantID).
		Preload("MenuItems").
		Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// Update updates an existing menu
func (r *MenuRepository) Update(menu *models.Menu) error {
	return r.db.Save(menu).Error
}

// Delete deletes a menu
func (r *MenuRepository) Delete(id uint) error {
	return r.db.Delete(&models.Menu{}, id).Error
}

