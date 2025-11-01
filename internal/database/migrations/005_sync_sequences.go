package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// SyncSequences migration synchronizes PostgreSQL sequences
type SyncSequences struct {
	BaseMigration
}

// NewSyncSequences creates a new migration
func NewSyncSequences() *SyncSequences {
	return &SyncSequences{
		BaseMigration: BaseMigration{
			version: 5,
			name:    "sync_sequences",
		},
	}
}

// Up syncs the restaurants_id_seq sequence
func (m *SyncSequences) Up(db *gorm.DB) error {
	if err := db.Exec(`
		DO $$
		DECLARE
			max_id BIGINT;
		BEGIN
			SELECT COALESCE(MAX(id), 0) INTO max_id FROM restaurants;
			-- Set sequence to max_id + 1 (or at least 1) to ensure next value doesn't conflict
			-- The third parameter 'true' means use the value (don't call nextval first)
			PERFORM setval('restaurants_id_seq', GREATEST(max_id, 1), true);
		END $$;
	`).Error; err != nil {
		return fmt.Errorf("failed to sync sequences: %w", err)
	}

	return nil
}

// Down is a no-op for sequence syncing (sequences will auto-adjust)
func (m *SyncSequences) Down(db *gorm.DB) error {
	// No-op: sequences will auto-adjust on inserts
	return nil
}
