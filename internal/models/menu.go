package models

import (
	"time"
)

// Menu represents a menu category
type Menu struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RestaurantID uint      `gorm:"index;not null" json:"restaurant_id"` // Crucial for RLS
	Name         string    `gorm:"not null" json:"name"`
	Description  string    `json:"description"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	Restaurant Restaurant `gorm:"foreignKey:RestaurantID"`
	MenuItems  []MenuItem `gorm:"foreignKey:MenuID"`
}

