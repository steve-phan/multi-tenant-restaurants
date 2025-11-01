package models

import (
	"time"
)

// MenuCategory represents a menu category (e.g., "Hot Food", "Drinks", "Vegans")
type MenuCategory struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RestaurantID uint      `gorm:"index;not null" json:"restaurant_id"` // Crucial for RLS
	Name         string    `gorm:"not null" json:"name"`
	Description  string    `json:"description"`
	DisplayOrder int       `gorm:"default:0;not null" json:"display_order"` // Order for sorting categories
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	Restaurant Restaurant `gorm:"foreignKey:RestaurantID"`
	MenuItems  []MenuItem `gorm:"foreignKey:CategoryID"`
}

// TableName specifies the table name for MenuCategory
func (MenuCategory) TableName() string {
	return "menu_categories"
}
