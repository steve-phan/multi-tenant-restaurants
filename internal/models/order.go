package models

import (
	"time"
)

// Order represents an order
type Order struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RestaurantID uint      `gorm:"index;not null" json:"restaurant_id"` // Crucial for RLS
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	Status       string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, confirmed, preparing, ready, completed, cancelled
	TotalAmount  float64   `gorm:"not null" json:"total_amount"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	Restaurant Restaurant  `gorm:"foreignKey:RestaurantID"`
	User       User        `gorm:"foreignKey:UserID"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
}
