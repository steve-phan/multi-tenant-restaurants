package database

import (
	"fmt"
	"log"
	"restaurant-saas/config"
	"restaurant-saas/models"
	"restaurant-saas/utils"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := config.AppConfig.DBUrl
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.Logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	utils.Logger.Info("Connected to database")
}

func Migrate() {
	// Migrate non-partitioned tables
	err := DB.AutoMigrate(
		&models.Organization{},
		&models.User{},
	)
	if err != nil {
		log.Fatal("Migration failed for non-partitioned tables:", err)
	}

	// Setup partitions for other tables
	SetupPartitions()

	// Run AutoMigrate for partitioned tables to ensure schema consistency (e.g. indexes, constraints)
	// Note: GORM might try to create FKs. If they fail due to partitioning constraints, we might need to handle FKs manually.
	err = DB.AutoMigrate(
		&models.Restaurant{},
		&models.MenuCategory{},
		&models.MenuItem{},
		&models.Table{},
		&models.Booking{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatal("Migration failed for partitioned tables:", err)
	}
	utils.Logger.Info("Database migration completed")
}

func SetupPartitions() {
	tables := []struct {
		Name   string
		Schema string
	}{
		{
			Name: "restaurants",
			Schema: `CREATE TABLE IF NOT EXISTS restaurants (
				id uuid DEFAULT gen_random_uuid(),
				organization_id uuid NOT NULL,
				name text NOT NULL,
				address text,
				contact_email text,
				opening_hours jsonb,
				created_at timestamptz,
				updated_at timestamptz,
				PRIMARY KEY (organization_id, id)
			) PARTITION BY HASH (organization_id)`,
		},
		{
			Name: "menu_categories",
			Schema: `CREATE TABLE IF NOT EXISTS menu_categories (
				id uuid DEFAULT gen_random_uuid(),
				organization_id uuid NOT NULL,
				restaurant_id uuid NOT NULL,
				name text NOT NULL,
				display_order int,
				created_at timestamptz,
				updated_at timestamptz,
				PRIMARY KEY (organization_id, id)
			) PARTITION BY HASH (organization_id)`,
		},
		{
			Name: "menu_items",
			Schema: `CREATE TABLE IF NOT EXISTS menu_items (
				id uuid DEFAULT gen_random_uuid(),
				organization_id uuid NOT NULL,
				restaurant_id uuid NOT NULL,
				category_id uuid NOT NULL,
				name text NOT NULL,
				description text,
				price numeric NOT NULL,
				image_url text,
				created_at timestamptz,
				updated_at timestamptz,
				PRIMARY KEY (organization_id, id)
			) PARTITION BY HASH (organization_id)`,
		},
		{
			Name: "tables",
			Schema: `CREATE TABLE IF NOT EXISTS tables (
				id uuid DEFAULT gen_random_uuid(),
				organization_id uuid NOT NULL,
				restaurant_id uuid NOT NULL,
				name text NOT NULL,
				capacity int NOT NULL,
				created_at timestamptz,
				updated_at timestamptz,
				PRIMARY KEY (organization_id, id)
			) PARTITION BY HASH (organization_id)`,
		},
		{
			Name: "bookings",
			Schema: `CREATE TABLE IF NOT EXISTS bookings (
				id uuid DEFAULT gen_random_uuid(),
				organization_id uuid NOT NULL,
				restaurant_id uuid NOT NULL,
				table_id uuid,
				customer_id uuid,
				customer_name text,
				customer_email text,
				start_time timestamptz NOT NULL,
				end_time timestamptz,
				number_of_guests int NOT NULL,
				status varchar(20) DEFAULT 'PENDING',
				created_at timestamptz,
				updated_at timestamptz,
				PRIMARY KEY (organization_id, id)
			) PARTITION BY HASH (organization_id)`,
		},
		{
			Name: "orders",
			Schema: `CREATE TABLE IF NOT EXISTS orders (
				id uuid DEFAULT gen_random_uuid(),
				organization_id uuid NOT NULL,
				restaurant_id uuid NOT NULL,
				customer_id uuid,
				table_id uuid,
				status varchar(20) DEFAULT 'PENDING',
				total_amount numeric,
				created_at timestamptz,
				updated_at timestamptz,
				PRIMARY KEY (organization_id, id)
			) PARTITION BY HASH (organization_id)`,
		},
		{
			Name: "order_items",
			Schema: `CREATE TABLE IF NOT EXISTS order_items (
				id uuid DEFAULT gen_random_uuid(),
				organization_id uuid NOT NULL,
				order_id uuid NOT NULL,
				menu_item_id uuid NOT NULL,
				quantity int NOT NULL,
				unit_price numeric NOT NULL,
				PRIMARY KEY (organization_id, id)
			) PARTITION BY HASH (organization_id)`,
		},
	}

	for _, t := range tables {
		if err := DB.Exec(t.Schema).Error; err != nil {
			utils.Logger.Fatal("Failed to create partitioned table", zap.String("table", t.Name), zap.Error(err))
		}

		// Create partitions
		for i := 0; i < 16; i++ {
			partitionName := fmt.Sprintf("%s_p%d", t.Name, i)
			query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s PARTITION OF %s FOR VALUES WITH (MODULUS 16, REMAINDER %d)", partitionName, t.Name, i)
			if err := DB.Exec(query).Error; err != nil {
				utils.Logger.Fatal("Failed to create partition", zap.String("partition", partitionName), zap.Error(err))
			}
		}
	}
}
