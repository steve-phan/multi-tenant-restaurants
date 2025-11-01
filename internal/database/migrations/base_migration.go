package migrations

import (
	"gorm.io/gorm"
)

// BaseMigration provides common functionality for migrations
type BaseMigration struct {
	version int
	name    string
}

// GetVersion returns the migration version
func (m *BaseMigration) GetVersion() int {
	return m.version
}

// GetName returns the migration name
func (m *BaseMigration) GetName() string {
	return m.name
}

// ensureMigrationTable ensures the schema_migrations table exists
func ensureMigrationTable(db *gorm.DB) error {
	if err := db.AutoMigrate(&MigrationVersion{}); err != nil {
		return err
	}
	return nil
}

// isMigrationApplied checks if a migration has already been applied
func isMigrationApplied(db *gorm.DB, version int) (bool, error) {
	var count int64
	if err := db.Model(&MigrationVersion{}).Where("version = ?", version).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// recordMigration marks a migration as applied
func recordMigration(db *gorm.DB, version int, name string) error {
	if err := ensureMigrationTable(db); err != nil {
		return err
	}

	migration := MigrationVersion{
		Version:   version,
		Name:      name,
		AppliedAt: db.NowFunc().Unix(),
	}

	return db.Create(&migration).Error
}

// removeMigration removes a migration record (for rollback)
func removeMigration(db *gorm.DB, version int) error {
	return db.Where("version = ?", version).Delete(&MigrationVersion{}).Error
}

// getAppliedMigrations returns all applied migrations
func getAppliedMigrations(db *gorm.DB) ([]MigrationVersion, error) {
	if err := ensureMigrationTable(db); err != nil {
		return nil, err
	}

	var migrations []MigrationVersion
	if err := db.Order("version ASC").Find(&migrations).Error; err != nil {
		return nil, err
	}
	return migrations, nil
}
