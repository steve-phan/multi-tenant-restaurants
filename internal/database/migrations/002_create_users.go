package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// CreateUsersTable migration
type CreateUsersTable struct {
	BaseMigration
}

// NewCreateUsersTable creates a new migration
func NewCreateUsersTable() *CreateUsersTable {
	return &CreateUsersTable{
		BaseMigration: BaseMigration{
			version: 2,
			name:    "create_users_table",
		},
	}
}

// Up creates the users table
func (m *CreateUsersTable) Up(db *gorm.DB) error {
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			restaurant_id BIGINT NOT NULL,
			email TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			first_name TEXT,
			last_name TEXT,
			role VARCHAR(20) NOT NULL,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			CONSTRAINT fk_restaurants_users FOREIGN KEY (restaurant_id) REFERENCES restaurants(id)
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create index on restaurant_id for RLS
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_restaurant_id ON users(restaurant_id)`).Error; err != nil {
		return fmt.Errorf("failed to create restaurant_id index: %w", err)
	}

	return nil
}

// Down drops the users table
func (m *CreateUsersTable) Down(db *gorm.DB) error {
	// Drop index
	db.Exec(`DROP INDEX IF EXISTS idx_users_restaurant_id`)

	// Drop table (CASCADE will drop dependent objects)
	if err := db.Exec(`DROP TABLE IF EXISTS users CASCADE`).Error; err != nil {
		return fmt.Errorf("failed to drop users table: %w", err)
	}

	return nil
}
