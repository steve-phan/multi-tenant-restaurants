package models

import (
	"time"
)

// RestaurantStatus represents the status of a restaurant
type RestaurantStatus string

const (
	RestaurantStatusPending   RestaurantStatus = "pending"
	RestaurantStatusActive    RestaurantStatus = "active"
	RestaurantStatusInactive  RestaurantStatus = "inactive"
	RestaurantStatusSuspended RestaurantStatus = "suspended"
)

// PlatformOrganizationID is the special organization ID for platform-level users (KAMs)
// This is a reserved organization that represents the platform itself
const PlatformOrganizationID uint = 1

// IsPlatformOrganization checks if a restaurant ID is the platform organization
func IsPlatformOrganization(id uint) bool {
	return id == PlatformOrganizationID
}

// Restaurant represents a tenant (restaurant)
type Restaurant struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"not null" json:"name"`
	Description string          `json:"description"`
	Address     string          `json:"address"`
	Phone       string          `json:"phone"`
	Email       string          `gorm:"uniqueIndex" json:"email"`
	Status      RestaurantStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	IsActive    bool            `gorm:"default:false" json:"is_active"` // Deprecated: use Status instead
	
	// KAM (Key Account Manager) fields
	KAMID       *uint      `gorm:"index" json:"kam_id,omitempty"` // Assigned KAM
	ActivatedBy *uint      `json:"activated_by,omitempty"`        // User who activated
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
	
	// Registration details
	ContactName  string    `json:"contact_name"`
	ContactEmail string    `json:"contact_email"`
	ContactPhone string    `json:"contact_phone"`
	
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	Users        []User         `gorm:"foreignKey:RestaurantID"`
	Categories   []MenuCategory `gorm:"foreignKey:RestaurantID"`
	Reservations []Reservation  `gorm:"foreignKey:RestaurantID"`
	Orders       []Order        `gorm:"foreignKey:RestaurantID"`
	KAM          *User          `gorm:"foreignKey:KAMID" json:"kam,omitempty"`
}

