package repositories

import (
	"restaurant-backend/internal/models"

	"gorm.io/gorm"
)

// MenuItemImageRepository handles menu item image-related database operations
type MenuItemImageRepository struct {
	db *gorm.DB
}

// NewMenuItemImageRepository creates a new MenuItemImageRepository instance
func NewMenuItemImageRepository(db *gorm.DB) *MenuItemImageRepository {
	return &MenuItemImageRepository{db: db}
}

// Create creates a new menu item image
func (r *MenuItemImageRepository) Create(image *models.MenuItemImage) error {
	return r.db.Create(image).Error
}

// GetByID retrieves an image by ID (RLS ensures tenant isolation)
func (r *MenuItemImageRepository) GetByID(id uint) (*models.MenuItemImage, error) {
	var image models.MenuItemImage
	if err := r.db.First(&image, id).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

// GetByMenuItemID retrieves all images for a menu item (RLS ensures tenant isolation)
// Ordered by display_order, then is_primary
func (r *MenuItemImageRepository) GetByMenuItemID(menuItemID uint) ([]models.MenuItemImage, error) {
	var images []models.MenuItemImage
	if err := r.db.Where("menu_item_id = ?", menuItemID).
		Order("is_primary DESC, display_order ASC").
		Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

// Update updates an existing menu item image
func (r *MenuItemImageRepository) Update(image *models.MenuItemImage) error {
	return r.db.Save(image).Error
}

// Delete deletes a menu item image
func (r *MenuItemImageRepository) Delete(id uint) error {
	return r.db.Delete(&models.MenuItemImage{}, id).Error
}

// DeleteByMenuItemID deletes all images for a menu item
func (r *MenuItemImageRepository) DeleteByMenuItemID(menuItemID uint) error {
	return r.db.Where("menu_item_id = ?", menuItemID).Delete(&models.MenuItemImage{}).Error
}

// SetPrimary sets an image as primary and un-sets others for the same menu item
func (r *MenuItemImageRepository) SetPrimary(menuItemID uint, imageID uint) error {
	// First, unset all primary flags for this menu item
	if err := r.db.Model(&models.MenuItemImage{}).
		Where("menu_item_id = ?", menuItemID).
		Update("is_primary", false).Error; err != nil {
		return err
	}

	// Then set the specified image as primary
	return r.db.Model(&models.MenuItemImage{}).
		Where("id = ? AND menu_item_id = ?", imageID, menuItemID).
		Update("is_primary", true).Error
}
