package models

import (
	"time"
)

// User represents a user (admin, staff, client, or KAM)
// KAM users belong to the Platform Organization (restaurant_id = PlatformOrganizationID)
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RestaurantID uint      `gorm:"index;not null" json:"restaurant_id"` // Required - KAMs belong to Platform Organization
	Email        string    `gorm:"not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Role         string    `gorm:"type:varchar(20);not null" json:"role"` // Admin, Staff, Client, KAM (Key Account Manager)
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	Phone        string    `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Timezone     string    `gorm:"type:varchar(50);default:'UTC'" json:"timezone"`
	Language     string    `gorm:"type:varchar(10);default:'en'" json:"language"`
	Preferences  string    `gorm:"type:jsonb;default:'{}'" json:"preferences,omitempty"` // JSON string for preferences
	AvatarURL    string    `json:"avatar_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	Restaurant *Restaurant `gorm:"foreignKey:RestaurantID" json:"restaurant,omitempty"`
}

// IsKAM checks if user is a KAM
func (u *User) IsKAM() bool {
	return u.Role == "KAM"
}

// IsPlatformUser checks if user belongs to the platform organization
func (u *User) IsPlatformUser() bool {
	return u.RestaurantID == PlatformOrganizationID
}
