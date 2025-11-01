package database

import (
	"fmt"
	"log"

	"restaurant-backend/internal/config"
	"restaurant-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB, cfg *config.Config) error {
	// Auto-migrate all models in correct order
	// Restaurant must be created first since other models reference it
	// Note: Restaurant has a relationship to User (KAM), so we need to create
	// Restaurant table first without foreign keys, then User, then add FK later

	// Step 1: Create Restaurant table first WITHOUT foreign key constraints
	// (KAM foreign key will be added after User table exists)
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS restaurants (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			address TEXT,
			phone TEXT,
			email TEXT,
			status VARCHAR(20) DEFAULT 'pending',
			is_active BOOLEAN DEFAULT false,
			kam_id BIGINT,
			activated_by BIGINT,
			activated_at TIMESTAMPTZ,
			contact_name TEXT,
			contact_email TEXT,
			contact_phone TEXT,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create Restaurant table: %w", err)
	}

	// Create unique index on email
	db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_restaurants_email ON restaurants(email)`)

	// Step 2: Create User table manually (to avoid GORM trying to add Restaurant FK)
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			restaurant_id BIGINT NOT NULL,
			email TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			first_name TEXT,
			last_name TEXT,
			role VARCHAR(20) NOT NULL,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			CONSTRAINT fk_restaurants_users FOREIGN KEY (restaurant_id) REFERENCES restaurants(id)
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create User table: %w", err)
	}

	// Create index on restaurant_id for RLS
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_restaurant_id ON users(restaurant_id)`)

	// Step 3: Now that both Restaurant and User exist, use AutoMigrate for remaining tables
	// GORM can now properly handle all relationships
	if err := db.AutoMigrate(
		&models.Menu{},
		&models.MenuItem{},
		&models.Reservation{},
		&models.Order{},
		&models.OrderItem{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate remaining models: %w", err)
	}

	// Step 4: Add foreign key constraint from Restaurant.KAMID to User.ID
	// (This was skipped in step 1 to avoid circular dependency)
	db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_restaurants_kam') THEN
				ALTER TABLE restaurants 
				ADD CONSTRAINT fk_restaurants_kam 
				FOREIGN KEY (kam_id) REFERENCES users(id);
			END IF;
		END $$;
	`)

	// Step 5: Ensure the restaurants_id_seq sequence is synchronized
	// This fixes the issue where the sequence might be out of sync after manual inserts
	// Set sequence to max(id) + 1 to ensure next insert doesn't conflict
	db.Exec(`
		DO $$
		DECLARE
			max_id BIGINT;
		BEGIN
			SELECT COALESCE(MAX(id), 0) INTO max_id FROM restaurants;
			-- Set sequence to max_id + 1 to ensure next value doesn't conflict
			PERFORM setval('restaurants_id_seq', GREATEST(max_id, 1));
		END $$;
	`)

	// Initialize platform organization and bootstrap admin
	if err := bootstrapPlatformOrganization(db, cfg); err != nil {
		return fmt.Errorf("failed to bootstrap platform organization: %w", err)
	}

	// Enable Row Level Security on all relevant tables
	if err := enableRLS(db); err != nil {
		return fmt.Errorf("failed to enable RLS: %w", err)
	}

	// Create RLS policies
	if err := createRLSPolicies(db); err != nil {
		return fmt.Errorf("failed to create RLS policies: %w", err)
	}

	fmt.Println("Database migrations completed successfully")
	return nil
}

// enableRLS enables Row Level Security on all tenant-isolated tables
func enableRLS(db *gorm.DB) error {
	tables := []string{
		"users",
		"menus",
		"menu_items",
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

// createRLSPolicies creates RLS policies for all tenant-isolated tables
func createRLSPolicies(db *gorm.DB) error {
	// First, ensure the application user role exists
	db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'restaurant_app_user') THEN
				CREATE ROLE restaurant_app_user;
			END IF;
		END
		$$;
	`)

	// Grant necessary permissions
	db.Exec("GRANT USAGE ON SCHEMA public TO restaurant_app_user;")
	db.Exec("GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO restaurant_app_user;")
	db.Exec("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO restaurant_app_user;")

	// Create policies for each table
	// Note: Users table needs special handling for platform organization (KAM users)
	policies := map[string]string{
		"users":        "(restaurant_id = current_setting('app.current_restaurant', true)::INTEGER) OR (restaurant_id = 1 AND current_setting('app.current_user_role', true) IN ('KAM', 'Admin'))",
		"menus":        "restaurant_id = current_setting('app.current_restaurant', true)::INTEGER",
		"menu_items":   "restaurant_id = current_setting('app.current_restaurant', true)::INTEGER",
		"reservations": "restaurant_id = current_setting('app.current_restaurant', true)::INTEGER",
		"orders":       "restaurant_id = current_setting('app.current_restaurant', true)::INTEGER",
		"order_items":  "restaurant_id = current_setting('app.current_restaurant', true)::INTEGER",
	}

	for table, condition := range policies {
		policyName := fmt.Sprintf("isolate_%s", table)

		// Drop policy if it exists
		db.Exec(fmt.Sprintf("DROP POLICY IF EXISTS %s ON %s", policyName, table))

		// Create policy
		sql := fmt.Sprintf(
			"CREATE POLICY %s ON %s FOR ALL TO restaurant_app_user USING (%s)",
			policyName,
			table,
			condition,
		)

		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("failed to create policy for %s: %w", table, err)
		}
	}

	return nil
}

// bootstrapPlatformOrganization creates the platform organization and initial admin user
func bootstrapPlatformOrganization(db *gorm.DB, cfg *config.Config) error {
	// Step 1: Create platform organization if it doesn't exist
	var platform models.Restaurant
	if err := db.First(&platform, models.PlatformOrganizationID).Error; err != nil {
		// Platform organization doesn't exist, create it using raw SQL to set ID explicitly
		// This avoids sequence conflicts
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
		db.Exec(`
			SELECT setval('restaurants_id_seq', GREATEST(?, 1))
		`, models.PlatformOrganizationID)

		log.Println("✓ Platform organization created")
	} else {
		log.Println("✓ Platform organization already exists")
	}

	// Always sync sequence after checking/creating platform organization
	// This ensures sequence is at least at PlatformOrganizationID
	db.Exec(`
		DO $$
		DECLARE
			max_id BIGINT;
		BEGIN
			SELECT COALESCE(MAX(id), 0) INTO max_id FROM restaurants;
			-- Set sequence to max_id to ensure next value is max_id + 1
			PERFORM setval('restaurants_id_seq', GREATEST(max_id, 1), false);
		END $$;
	`)

	// Step 2: Create initial admin user if it doesn't exist
	var adminUser models.User
	if err := db.Where("restaurant_id = ? AND role = ?", models.PlatformOrganizationID, "KAM").First(&adminUser).Error; err != nil {
		// No admin user exists, create one
		adminEmail := cfg.BootstrapAdminEmail
		adminPassword := cfg.BootstrapAdminPassword

		// If password is not set, generate a random one or use default (for development only)
		if adminPassword == "" {
			if cfg.Environment == "production" {
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
		if cfg.Environment != "production" && adminPassword == "ChangeMe123!" {
			log.Printf("⚠ IMPORTANT: Default password 'ChangeMe123!' was used. Please change it immediately!")
		}
	} else {
		log.Println("✓ Initial admin user already exists")
	}

	return nil
}
