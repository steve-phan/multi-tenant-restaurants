package migrations

import (
	"gorm.io/gorm"
)

// Migration represents a database migration with up and down functions
type Migration interface {
	Up(db *gorm.DB) error
	Down(db *gorm.DB) error
	GetVersion() int
	GetName() string
}

// MigrationVersion tracks which migrations have been applied
type MigrationVersion struct {
	ID        uint   `gorm:"primaryKey"`
	Version   int    `gorm:"uniqueIndex;not null"`
	Name      string `gorm:"not null"`
	AppliedAt int64  `gorm:"not null"`
}

// TableName specifies the table name for MigrationVersion
func (MigrationVersion) TableName() string {
	return "schema_migrations"
}
