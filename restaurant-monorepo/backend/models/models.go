package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Organization struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Restaurant struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;primaryKey"` // Composite PK for partitioning
	Name           string    `gorm:"not null"`
	Address        string
	ContactEmail   string
	OpeningHours   datatypes.JSON
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserRole string

const (
	RoleSuperAdmin UserRole = "SUPER_ADMIN"
	RoleOrgAdmin   UserRole = "ORG_ADMIN"
	RoleStaff      UserRole = "STAFF"
	RoleCustomer   UserRole = "CUSTOMER"
)

type User struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID *uuid.UUID `gorm:"type:uuid"` // Nullable for Super Admin or Customers not bound to org
	Email          string     `gorm:"uniqueIndex;not null"`
	PasswordHash   string     `gorm:"not null"`
	Role           UserRole   `gorm:"type:varchar(20);default:'CUSTOMER'"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type MenuCategory struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;primaryKey"` // Added for partitioning
	RestaurantID   uuid.UUID `gorm:"type:uuid;not null"`
	Name           string    `gorm:"not null"`
	DisplayOrder   int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type MenuItem struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;primaryKey"` // Added for partitioning
	RestaurantID   uuid.UUID `gorm:"type:uuid;not null"`
	CategoryID     uuid.UUID `gorm:"type:uuid;not null"`
	Name           string    `gorm:"not null"`
	Description    string
	Price          float64 `gorm:"not null"`
	ImageURL       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Table struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;primaryKey"` // Added for partitioning
	RestaurantID   uuid.UUID `gorm:"type:uuid;not null"`
	Name           string    `gorm:"not null"`
	Capacity       int       `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type BookingStatus string

const (
	BookingConfirmed BookingStatus = "CONFIRMED"
	BookingPending   BookingStatus = "PENDING"
	BookingCancelled BookingStatus = "CANCELLED"
)

type Booking struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID uuid.UUID  `gorm:"type:uuid;not null;primaryKey"` // Added for partitioning
	RestaurantID   uuid.UUID  `gorm:"type:uuid;not null"`
	TableID        *uuid.UUID `gorm:"type:uuid"`
	CustomerID     *uuid.UUID `gorm:"type:uuid"` // Nullable if guest booking
	CustomerName   string
	CustomerEmail  string
	StartTime      time.Time `gorm:"not null"`
	EndTime        time.Time
	NumberOfGuests int           `gorm:"not null"`
	Status         BookingStatus `gorm:"type:varchar(20);default:'PENDING'"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type OrderStatus string

const (
	OrderPending   OrderStatus = "PENDING"
	OrderPreparing OrderStatus = "PREPARING"
	OrderReady     OrderStatus = "READY"
	OrderCompleted OrderStatus = "COMPLETED"
	OrderCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID             uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID uuid.UUID   `gorm:"type:uuid;not null;primaryKey"` // Added for partitioning
	RestaurantID   uuid.UUID   `gorm:"type:uuid;not null"`
	CustomerID     *uuid.UUID  `gorm:"type:uuid"`
	TableID        *uuid.UUID  `gorm:"type:uuid"`
	Status         OrderStatus `gorm:"type:varchar(20);default:'PENDING'"`
	TotalAmount    float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Items          []OrderItem `gorm:"foreignKey:OrderID,OrganizationID;references:ID,OrganizationID"`
}

type OrderItem struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;primaryKey"` // Added for partitioning
	OrderID        uuid.UUID `gorm:"type:uuid;not null"`
	MenuItemID     uuid.UUID `gorm:"type:uuid;not null"`
	Quantity       int       `gorm:"not null"`
	UnitPrice      float64   `gorm:"not null"`
}
