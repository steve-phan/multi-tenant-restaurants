package models

import (
	"time"
)

// MenuItem represents a menu item within a category
type MenuItem struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RestaurantID uint      `gorm:"index;not null" json:"restaurant_id"` // Crucial for RLS
	CategoryID   uint      `gorm:"index;not null" json:"category_id"`   // References MenuCategory
	Name         string    `gorm:"not null" json:"name"`
	Description  string    `json:"description"`
	Price        float64   `gorm:"not null" json:"price"`
	ImageURL     string    `json:"image_url"`                               // Deprecated: use Images relationship instead
	DisplayOrder int       `gorm:"default:0;not null" json:"display_order"` // Order for sorting items within category
	IsAvailable  bool      `gorm:"default:true" json:"is_available"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	Restaurant Restaurant      `gorm:"foreignKey:RestaurantID"`
	Category   MenuCategory    `gorm:"foreignKey:CategoryID"`
	Images     []MenuItemImage `gorm:"foreignKey:MenuItemID;order:display_order asc" json:"images,omitempty"`
	OrderItems []OrderItem     `gorm:"foreignKey:MenuItemID"`
}
