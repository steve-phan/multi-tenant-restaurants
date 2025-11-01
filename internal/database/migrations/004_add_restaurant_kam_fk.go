package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// AddRestaurantKamFK migration adds foreign key from Restaurant.KAMID to User.ID
type AddRestaurantKamFK struct {
	BaseMigration
}

// NewAddRestaurantKamFK creates a new migration
func NewAddRestaurantKamFK() *AddRestaurantKamFK {
	return &AddRestaurantKamFK{
		BaseMigration: BaseMigration{
			version: 4,
			name:    "add_restaurant_kam_fk",
		},
	}
}

// Up adds the foreign key constraint
func (m *AddRestaurantKamFK) Up(db *gorm.DB) error {
	if err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_restaurants_kam') THEN
				ALTER TABLE restaurants 
				ADD CONSTRAINT fk_restaurants_kam 
				FOREIGN KEY (kam_id) REFERENCES users(id);
			END IF;
		END $$;
	`).Error; err != nil {
		return fmt.Errorf("failed to add KAM foreign key: %w", err)
	}

	return nil
}

// Down removes the foreign key constraint
func (m *AddRestaurantKamFK) Down(db *gorm.DB) error {
	if err := db.Exec(`ALTER TABLE restaurants DROP CONSTRAINT IF EXISTS fk_restaurants_kam`).Error; err != nil {
		return fmt.Errorf("failed to drop KAM foreign key: %w", err)
	}

	return nil
}

