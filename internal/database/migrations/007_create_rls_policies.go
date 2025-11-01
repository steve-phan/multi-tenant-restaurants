package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// CreateRLSPolicies migration creates RLS policies for all tenant-isolated tables
type CreateRLSPolicies struct {
	BaseMigration
}

// NewCreateRLSPolicies creates a new migration
func NewCreateRLSPolicies() *CreateRLSPolicies {
	return &CreateRLSPolicies{
		BaseMigration: BaseMigration{
			version: 7,
			name:    "create_rls_policies",
		},
	}
}

// Up creates RLS policies
func (m *CreateRLSPolicies) Up(db *gorm.DB) error {
	// First, ensure the application user role exists
	if err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'restaurant_app_user') THEN
				CREATE ROLE restaurant_app_user;
			END IF;
		END
		$$;
	`).Error; err != nil {
		return fmt.Errorf("failed to create restaurant_app_user role: %w", err)
	}

	// Grant necessary permissions
	if err := db.Exec("GRANT USAGE ON SCHEMA public TO restaurant_app_user").Error; err != nil {
		return fmt.Errorf("failed to grant schema usage: %w", err)
	}

	if err := db.Exec("GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO restaurant_app_user").Error; err != nil {
		return fmt.Errorf("failed to grant table permissions: %w", err)
	}

	if err := db.Exec("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO restaurant_app_user").Error; err != nil {
		return fmt.Errorf("failed to grant default privileges: %w", err)
	}

	// Create policies for each table
	policies := map[string]string{
		"users":            "(restaurant_id = current_setting('app.current_restaurant', true)::INTEGER) OR (restaurant_id = 1 AND current_setting('app.current_user_role', true) IN ('KAM', 'Admin'))",
		"menu_categories":  "restaurant_id = current_setting('app.current_restaurant', true)::INTEGER",
		"menu_items":       "restaurant_id = current_setting('app.current_restaurant', true)::INTEGER",
		"menu_item_images": "restaurant_id = current_setting('app.current_restaurant', true)::INTEGER",
		"reservations":     "restaurant_id = current_setting('app.current_restaurant', true)::INTEGER",
		"orders":           "restaurant_id = current_setting('app.current_restaurant', true)::INTEGER",
		"order_items":      "restaurant_id = current_setting('app.current_restaurant', true)::INTEGER",
	}

	for table, condition := range policies {
		policyName := fmt.Sprintf("isolate_%s", table)

		// Drop policy if it exists
		db.Exec(fmt.Sprintf("DROP POLICY IF EXISTS %s ON %s", policyName, table))

		// Create policy with both USING and WITH CHECK
		sql := fmt.Sprintf(
			"CREATE POLICY %s ON %s FOR ALL TO restaurant_app_user USING (%s) WITH CHECK (%s)",
			policyName,
			table,
			condition,
			condition,
		)

		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("failed to create policy for %s: %w", table, err)
		}
	}

	return nil
}

// Down drops all RLS policies
func (m *CreateRLSPolicies) Down(db *gorm.DB) error {
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
		policyName := fmt.Sprintf("isolate_%s", table)
		if err := db.Exec(fmt.Sprintf("DROP POLICY IF EXISTS %s ON %s", policyName, table)).Error; err != nil {
			return fmt.Errorf("failed to drop policy for %s: %w", table, err)
		}
	}

	// Note: We don't drop the role as it might be used by other parts of the system
	return nil
}

