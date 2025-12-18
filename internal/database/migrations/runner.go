package migrations

import (
	"fmt"
	"sort"

	"gorm.io/gorm"
)

// Runner executes migrations
type Runner struct {
	db         *gorm.DB
	migrations []Migration
}

// NewRunner creates a new migration runner
func NewRunner(db *gorm.DB, migrations []Migration) *Runner {
	// Sort migrations by version
	sorted := make([]Migration, len(migrations))
	copy(sorted, migrations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].GetVersion() < sorted[j].GetVersion()
	})

	return &Runner{
		db:         db,
		migrations: sorted,
	}
}

// Up runs all pending migrations
func (r *Runner) Up() error {
	if err := ensureMigrationTable(r.db); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	applied, err := getAppliedMigrations(r.db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Create map of applied versions for quick lookup
	appliedMap := make(map[int]bool)
	for _, m := range applied {
		appliedMap[m.Version] = true
	}

	// Run pending migrations
	for _, migration := range r.migrations {
		if appliedMap[migration.GetVersion()] {
			fmt.Printf("Skipping migration %d: %s (already applied)\n", migration.GetVersion(), migration.GetName())
			continue
		}

		fmt.Printf("Running migration %d: %s...\n", migration.GetVersion(), migration.GetName())
		if err := migration.Up(r.db); err != nil {
			return fmt.Errorf("failed to run migration %d (%s): %w", migration.GetVersion(), migration.GetName(), err)
		}

		if err := recordMigration(r.db, migration.GetVersion(), migration.GetName()); err != nil {
			return fmt.Errorf("failed to record migration %d: %w", migration.GetVersion(), err)
		}

		fmt.Printf("✓ Migration %d: %s completed\n", migration.GetVersion(), migration.GetName())
	}

	return nil
}

// Down rolls back the last migration
func (r *Runner) Down() error {
	if err := ensureMigrationTable(r.db); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	applied, err := getAppliedMigrations(r.db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(applied) == 0 {
		return fmt.Errorf("no migrations to rollback")
	}

	// Find the last applied migration
	lastApplied := applied[len(applied)-1]

	// Find the migration in our list
	var migrationToRollback Migration
	for _, m := range r.migrations {
		if m.GetVersion() == lastApplied.Version {
			migrationToRollback = m
			break
		}
	}

	if migrationToRollback == nil {
		return fmt.Errorf("migration version %d not found in migration list", lastApplied.Version)
	}

	fmt.Printf("Rolling back migration %d: %s...\n", lastApplied.Version, lastApplied.Name)
	if err := migrationToRollback.Down(r.db); err != nil {
		return fmt.Errorf("failed to rollback migration %d (%s): %w", lastApplied.Version, lastApplied.Name, err)
	}

	if err := removeMigration(r.db, lastApplied.Version); err != nil {
		return fmt.Errorf("failed to remove migration record %d: %w", lastApplied.Version, err)
	}

	fmt.Printf("✓ Migration %d: %s rolled back\n", lastApplied.Version, lastApplied.Name)
	return nil
}

// Status shows the status of all migrations
func (r *Runner) Status() error {
	if err := ensureMigrationTable(r.db); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	applied, err := getAppliedMigrations(r.db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedMap := make(map[int]bool)
	for _, m := range applied {
		appliedMap[m.Version] = true
	}

	fmt.Println("\nMigration Status:")
	fmt.Println("==================")
	for _, migration := range r.migrations {
		status := "pending"
		if appliedMap[migration.GetVersion()] {
			status = "applied"
		}
		fmt.Printf("[%s] %d: %s\n", status, migration.GetVersion(), migration.GetName())
	}

	return nil
}
