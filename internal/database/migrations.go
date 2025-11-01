package database

import (
	"fmt"

	"restaurant-backend/internal/config"
	"restaurant-backend/internal/database/migrations"

	"gorm.io/gorm"
)

// RunMigrations runs all database migrations using the new migration system
func RunMigrations(db *gorm.DB, cfg *config.Config) error {
	// Register all migrations in order
	migrationList := []migrations.Migration{
		migrations.NewCreateRestaurantsTable(),
		migrations.NewCreateUsersTable(),
		migrations.NewCreateTables(),
		migrations.NewAddRestaurantKamFK(),
		migrations.NewSyncSequences(),
		migrations.NewEnableRLS(),
		migrations.NewCreateRLSPolicies(),
		migrations.NewBootstrapPlatform(cfg),
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
	migrationList := []migrations.Migration{
		migrations.NewCreateRestaurantsTable(),
		migrations.NewCreateUsersTable(),
		migrations.NewCreateTables(),
		migrations.NewAddRestaurantKamFK(),
		migrations.NewSyncSequences(),
		migrations.NewEnableRLS(),
		migrations.NewCreateRLSPolicies(),
		migrations.NewBootstrapPlatform(cfg),
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
	// Register all migrations
	migrationList := []migrations.Migration{
		migrations.NewCreateRestaurantsTable(),
		migrations.NewCreateUsersTable(),
		migrations.NewCreateTables(),
		migrations.NewAddRestaurantKamFK(),
		migrations.NewSyncSequences(),
		migrations.NewEnableRLS(),
		migrations.NewCreateRLSPolicies(),
		migrations.NewBootstrapPlatform(cfg),
	}

	runner := migrations.NewRunner(db, migrationList)
	return runner.Status()
}
