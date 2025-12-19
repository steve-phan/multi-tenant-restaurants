package database

import (
	"fmt"

	"restaurant-backend/internal/config"
	"restaurant-backend/internal/database/migrations"

	"gorm.io/gorm"
)

// RunMigrations runs all database migrations using the new migration system
// Note: This does NOT bootstrap the platform - use BootstrapPlatform() separately
func RunMigrations(db *gorm.DB, cfg *config.Config) error {
	// Register all migrations in order (excluding bootstrap)
	migrationList := []migrations.Migration{
		migrations.NewCreateRestaurantsTable(),
		migrations.NewCreateUsersTable(),
		migrations.NewCreateTables(),
		migrations.NewAddRestaurantKamFK(),
		migrations.NewSyncSequences(),
		migrations.NewEnableRLS(),
		migrations.NewCreateRLSPolicies(),
		migrations.NewAddUserFields(),
		// Bootstrap is separate - use BootstrapPlatform() instead
	}

	// Create runner and execute migrations
	runner := migrations.NewRunner(db, migrationList)

	if err := runner.Up(); err != nil {
		return fmt.Errorf("migrations failed: %w", err)
	}

	fmt.Println("Database migrations completed successfully")
	return nil
}

// RunMigrationsDown rolls back the last migration
func RunMigrationsDown(db *gorm.DB, cfg *config.Config) error {
	// Register all migrations (needed to find the one to rollback)
	// Note: Bootstrap is not in the migration list anymore
	migrationList := []migrations.Migration{
		migrations.NewCreateRestaurantsTable(),
		migrations.NewCreateUsersTable(),
		migrations.NewCreateTables(),
		migrations.NewAddRestaurantKamFK(),
		migrations.NewSyncSequences(),
		migrations.NewEnableRLS(),
		migrations.NewCreateRLSPolicies(),
	}

	runner := migrations.NewRunner(db, migrationList)

	if err := runner.Down(); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	fmt.Println("Migration rolled back successfully")
	return nil
}

// ShowMigrationStatus shows the status of all migrations
func ShowMigrationStatus(db *gorm.DB, cfg *config.Config) error {
	// Register all migrations (excluding bootstrap)
	migrationList := []migrations.Migration{
		migrations.NewCreateRestaurantsTable(),
		migrations.NewCreateUsersTable(),
		migrations.NewCreateTables(),
		migrations.NewAddRestaurantKamFK(),
		migrations.NewSyncSequences(),
		migrations.NewEnableRLS(),
		migrations.NewCreateRLSPolicies(),
		migrations.NewAddUserFields(),
	}

	runner := migrations.NewRunner(db, migrationList)
	return runner.Status()
}
