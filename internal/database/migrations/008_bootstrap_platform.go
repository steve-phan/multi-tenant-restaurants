package migrations

import (
	"fmt"
	"log"

	"restaurant-backend/internal/config"
	"restaurant-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// BootstrapPlatform migration creates platform organization and initial admin user
type BootstrapPlatform struct {
	BaseMigration
	cfg *config.Config
}

// NewBootstrapPlatform creates a new migration
func NewBootstrapPlatform(cfg *config.Config) *BootstrapPlatform {
	return &BootstrapPlatform{
		BaseMigration: BaseMigration{
			version: 8,
			name:    "bootstrap_platform",
		},
		cfg: cfg,
	}
}

// Up creates platform organization and initial admin
func (m *BootstrapPlatform) Up(db *gorm.DB) error {
	// Step 1: Create platform organization if it doesn't exist
	var platform models.Restaurant
	if err := db.First(&platform, models.PlatformOrganizationID).Error; err != nil {
		// Platform organization doesn't exist, create it
		err := db.Exec(`
			INSERT INTO restaurants (id, name, description, status, is_active, email, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING
		`, models.PlatformOrganizationID,
			"Platform Organization",
			"Platform-level organization for KAM and system administrators",
			models.RestaurantStatusActive,
			true,
			"platform@system.local").Error

		if err != nil {
			return fmt.Errorf("failed to create platform organization: %w", err)
		}

		// Sync sequence to ensure next restaurant gets ID >= 2
		// Set sequence to max(platform_id, current_max) + 1 to avoid conflicts
		db.Exec(`
			DO $$
			DECLARE
				max_id BIGINT;
			BEGIN
				SELECT COALESCE(MAX(id), 0) INTO max_id FROM restaurants;
				PERFORM setval('restaurants_id_seq', GREATEST(max_id, 1), true);
			END $$;
		`)

		log.Println("✓ Platform organization created")
	} else {
		log.Println("✓ Platform organization already exists")
	}

	// Always sync sequence after checking/creating platform organization
	// Set sequence to max_id + 1 (or at least 2) to ensure next restaurant gets a unique ID
	db.Exec(`
			DO $$
			DECLARE
				max_id BIGINT;
			BEGIN
				SELECT COALESCE(MAX(id), 0) INTO max_id FROM restaurants;
				-- Set sequence to max_id + 1, but at least 2 (since 1 is reserved for platform org)
				-- The third parameter 'true' means use this value (nextval will return max_id + 1)
				PERFORM setval('restaurants_id_seq', GREATEST(max_id, 1), true);
			END $$;
		`)

	// Step 2: Create initial admin user if it doesn't exist
	var adminUser models.User
	if err := db.Where("restaurant_id = ? AND role = ?", models.PlatformOrganizationID, "KAM").First(&adminUser).Error; err != nil {
		// No admin user exists, create one
		adminEmail := m.cfg.BootstrapAdminEmail
		adminPassword := m.cfg.BootstrapAdminPassword

		// If password is not set, generate a random one or use default (for development only)
		if adminPassword == "" {
			if m.cfg.Environment == "production" {
				return fmt.Errorf("BOOTSTRAP_ADMIN_PASSWORD is required in production")
			}
			// Development default - should be changed immediately
			adminPassword = "ChangeMe123!"
			log.Println("⚠ WARNING: Using default admin password. Set BOOTSTRAP_ADMIN_PASSWORD in production!")
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash admin password: %w", err)
		}

		adminUser = models.User{
			RestaurantID: models.PlatformOrganizationID,
			Email:        adminEmail,
			PasswordHash: string(hashedPassword),
			FirstName:    "Platform",
			LastName:     "Administrator",
			Role:         "KAM",
			IsActive:     true,
		}

		if err := db.Create(&adminUser).Error; err != nil {
			return fmt.Errorf("failed to create initial admin user: %w", err)
		}

		log.Printf("✓ Initial admin user created: %s", adminEmail)
		if m.cfg.Environment != "production" && adminPassword == "ChangeMe123!" {
			log.Printf("⚠ IMPORTANT: Default password 'ChangeMe123!' was used. Please change it immediately!")
		}
	} else {
		log.Println("✓ Initial admin user already exists")
	}

	return nil
}

// Down removes platform organization and admin (use with caution)
func (m *BootstrapPlatform) Down(db *gorm.DB) error {
	// Remove admin users first
	if err := db.Where("restaurant_id = ?", models.PlatformOrganizationID).Delete(&models.User{}).Error; err != nil {
		return fmt.Errorf("failed to delete admin users: %w", err)
	}

	// Remove platform organization
	if err := db.Where("id = ?", models.PlatformOrganizationID).Delete(&models.Restaurant{}).Error; err != nil {
		return fmt.Errorf("failed to delete platform organization: %w", err)
	}

	log.Println("⚠ Platform organization and admin users removed")
	return nil
}
