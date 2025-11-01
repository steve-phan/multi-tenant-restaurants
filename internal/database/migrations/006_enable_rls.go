package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// EnableRLS migration enables Row Level Security on tenant-isolated tables
type EnableRLS struct {
	BaseMigration
}

// NewEnableRLS creates a new migration
func NewEnableRLS() *EnableRLS {
	return &EnableRLS{
		BaseMigration: BaseMigration{
			version: 6,
			name:    "enable_rls",
		},
	}
}

// Up enables RLS on all tenant-isolated tables
func (m *EnableRLS) Up(db *gorm.DB) error {
	tables := []string{
		"users",
		"menu_categories",
		"menu_items",
		"menu_item_images",
		"reservations",
		"orders",
		"order_items",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("ALTER TABLE %s ENABLE ROW LEVEL SECURITY", table)).Error; err != nil {
			return fmt.Errorf("failed to enable RLS on %s: %w", table, err)
		}
	}

	return nil
}

// Down disables RLS on all tables
func (m *EnableRLS) Down(db *gorm.DB) error {
	tables := []string{
		"order_items",
		"orders",
		"reservations",
		"menu_item_images",
		"menu_items",
		"menu_categories",
		"users",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("ALTER TABLE %s DISABLE ROW LEVEL SECURITY", table)).Error; err != nil {
			return fmt.Errorf("failed to disable RLS on %s: %w", table, err)
		}
	}

	return nil
}

