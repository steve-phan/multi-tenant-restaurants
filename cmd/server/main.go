package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"restaurant-backend/internal/config"
	"restaurant-backend/internal/database"
	"restaurant-backend/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	var migrate = flag.Bool("migrate", false, "Run database migrations (up)")
	var migrateDown = flag.Bool("migrate-down", false, "Rollback last migration (down)")
	var migrateStatus = flag.Bool("migrate-status", false, "Show migration status")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Handle migration commands
	if *migrate {
		if err := database.RunMigrations(db, cfg); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		os.Exit(0)
	}

	if *migrateDown {
		if err := database.RunMigrationsDown(db, cfg); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		os.Exit(0)
	}

	if *migrateStatus {
		if err := database.ShowMigrationStatus(db, cfg); err != nil {
			log.Fatalf("Failed to show migration status: %v", err)
		}
		os.Exit(0)
	}

	// Setup router
	r := router.SetupRouter(cfg, db)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	fmt.Printf("Server starting on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
