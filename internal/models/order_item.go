package models

import (
	"time"
)

// OrderItem represents an item in an order
type OrderItem struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RestaurantID uint      `gorm:"index;not null" json:"restaurant_id"` // Crucial for RLS
	OrderID      uint      `gorm:"index;not null" json:"order_id"`
	MenuItemID   uint      `gorm:"index;not null" json:"menu_item_id"`
	Quantity     int       `gorm:"not null" json:"quantity"`
	Price        float64   `gorm:"not null" json:"price"` // Price at time of order
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	Restaurant Restaurant `gorm:"foreignKey:RestaurantID"`
	Order      Order      `gorm:"foreignKey:OrderID"`
	MenuItem   MenuItem   `gorm:"foreignKey:MenuItemID"`
}

