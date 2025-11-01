package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// CreateRestaurantsTable migration
type CreateRestaurantsTable struct {
	BaseMigration
}

// NewCreateRestaurantsTable creates a new migration
func NewCreateRestaurantsTable() *CreateRestaurantsTable {
	return &CreateRestaurantsTable{
		BaseMigration: BaseMigration{
			version: 1,
			name:    "create_restaurants_table",
		},
	}
}

// Up creates the restaurants table
func (m *CreateRestaurantsTable) Up(db *gorm.DB) error {
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS restaurants (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			address TEXT,
			phone TEXT,
			email TEXT,
			status VARCHAR(20) DEFAULT 'pending',
			is_active BOOLEAN DEFAULT false,
			kam_id BIGINT,
			activated_by BIGINT,
			activated_at TIMESTAMPTZ,
			contact_name TEXT,
			contact_email TEXT,
			contact_phone TEXT,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create restaurants table: %w", err)
	}

	// Create unique index on email
	if err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_restaurants_email ON restaurants(email)`).Error; err != nil {
		return fmt.Errorf("failed to create email index: %w", err)
	}

	return nil
}

// Down drops the restaurants table
func (m *CreateRestaurantsTable) Down(db *gorm.DB) error {
	// Drop foreign key constraint first if it exists
	db.Exec(`ALTER TABLE restaurants DROP CONSTRAINT IF EXISTS fk_restaurants_kam`)

	// Drop index
	db.Exec(`DROP INDEX IF EXISTS idx_restaurants_email`)

	// Drop table
	if err := db.Exec(`DROP TABLE IF EXISTS restaurants CASCADE`).Error; err != nil {
		return fmt.Errorf("failed to drop restaurants table: %w", err)
	}

	return nil
}
