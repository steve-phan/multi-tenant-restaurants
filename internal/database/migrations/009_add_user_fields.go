package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// AddUserFields migration
type AddUserFields struct {
	BaseMigration
}

// NewAddUserFields creates a new migration
func NewAddUserFields() *AddUserFields {
	return &AddUserFields{
		BaseMigration: BaseMigration{
			version: 9,
			name:    "add_user_fields",
		},
	}
}

// Up adds new fields to users table
func (m *AddUserFields) Up(db *gorm.DB) error {
	// Add phone field
	if err := db.Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20)
	`).Error; err != nil {
		return fmt.Errorf("failed to add phone column: %w", err)
	}

	// Add timezone field
	if err := db.Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS timezone VARCHAR(50) DEFAULT 'UTC'
	`).Error; err != nil {
		return fmt.Errorf("failed to add timezone column: %w", err)
	}

	// Add language field
	if err := db.Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS language VARCHAR(10) DEFAULT 'en'
	`).Error; err != nil {
		return fmt.Errorf("failed to add language column: %w", err)
	}

	// Add preferences field (JSONB for better querying)
	if err := db.Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS preferences JSONB DEFAULT '{}'::jsonb
	`).Error; err != nil {
		return fmt.Errorf("failed to add preferences column: %w", err)
	}

	// Add avatar_url field
	if err := db.Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url TEXT
	`).Error; err != nil {
		return fmt.Errorf("failed to add avatar_url column: %w", err)
	}

	return nil
}

// Down removes the added fields
func (m *AddUserFields) Down(db *gorm.DB) error {
	// Drop columns
	if err := db.Exec(`
		ALTER TABLE users 
		DROP COLUMN IF EXISTS phone,
		DROP COLUMN IF EXISTS timezone,
		DROP COLUMN IF EXISTS language,
		DROP COLUMN IF EXISTS preferences,
		DROP COLUMN IF EXISTS avatar_url
	`).Error; err != nil {
		return fmt.Errorf("failed to drop user fields: %w", err)
	}

	return nil
}
