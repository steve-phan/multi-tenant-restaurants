package models

import (
	"time"
)

// MenuItem represents a menu item
type MenuItem struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RestaurantID uint      `gorm:"index;not null" json:"restaurant_id"` // Crucial for RLS
	MenuID       uint      `gorm:"index;not null" json:"menu_id"`
	Name         string    `gorm:"not null" json:"name"`
	Description  string    `json:"description"`
	Price        float64   `gorm:"not null" json:"price"`
	ImageURL     string    `json:"image_url"`
	Category     string    `json:"category"`
	IsAvailable  bool      `gorm:"default:true" json:"is_available"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	Restaurant Restaurant  `gorm:"foreignKey:RestaurantID"`
	Menu       Menu        `gorm:"foreignKey:MenuID"`
	OrderItems []OrderItem `gorm:"foreignKey:MenuItemID"`
}

