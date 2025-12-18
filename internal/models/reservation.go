package models

import (
	"time"
)

// Reservation represents a table reservation
type Reservation struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	RestaurantID   uint      `gorm:"index;not null" json:"restaurant_id"` // Crucial for RLS
	UserID         uint      `gorm:"index;not null" json:"user_id"`
	TableNumber    string    `gorm:"not null" json:"table_number"`
	StartTime      time.Time `gorm:"not null" json:"start_time"`
	EndTime        time.Time `gorm:"not null" json:"end_time"`
	NumberOfGuests int       `gorm:"not null" json:"number_of_guests"`
	Status         string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, confirmed, cancelled, completed
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relationships
	Restaurant Restaurant `gorm:"foreignKey:RestaurantID"`
	User       User       `gorm:"foreignKey:UserID"`
}
