package models

import (
	"time"
)

// MenuItemImage represents an image for a menu item
type MenuItemImage struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RestaurantID uint      `gorm:"index;not null" json:"restaurant_id"` // Crucial for RLS
	MenuItemID   uint      `gorm:"index;not null" json:"menu_item_id"`
	ImageURL     string    `gorm:"not null" json:"image_url"`
	DisplayOrder int       `gorm:"default:0;not null" json:"display_order"` // Order for sorting images
	IsPrimary    bool      `gorm:"default:false" json:"is_primary"`          // Primary/first image
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	Restaurant Restaurant `gorm:"foreignKey:RestaurantID"`
	MenuItem   MenuItem   `gorm:"foreignKey:MenuItemID"`
}

