package migrations

import (
	"fmt"

	"restaurant-backend/internal/models"

	"gorm.io/gorm"
)

// CreateTables migration creates all remaining tables using AutoMigrate
type CreateTables struct {
	BaseMigration
}

// NewCreateTables creates a new migration
func NewCreateTables() *CreateTables {
	return &CreateTables{
		BaseMigration: BaseMigration{
			version: 3,
			name:    "create_tables",
		},
	}
}

// Up creates all remaining tables
func (m *CreateTables) Up(db *gorm.DB) error {
	// Migrate MenuCategory first
	if err := db.AutoMigrate(&models.MenuCategory{}); err != nil {
		return fmt.Errorf("failed to migrate MenuCategory: %w", err)
	}

	// Migrate MenuItem
	if err := db.AutoMigrate(&models.MenuItem{}); err != nil {
		return fmt.Errorf("failed to migrate MenuItem: %w", err)
	}

	// Migrate remaining tables
	if err := db.AutoMigrate(
		&models.MenuItemImage{},
		&models.Reservation{},
		&models.Order{},
		&models.OrderItem{},
	); err != nil {
		return fmt.Errorf("failed to migrate remaining models: %w", err)
	}

	return nil
}

// Down drops all tables created by this migration
func (m *CreateTables) Down(db *gorm.DB) error {
	tables := []string{
		"order_items",
		"orders",
		"reservations",
		"menu_item_images",
		"menu_items",
		"menu_categories",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	return nil
}
