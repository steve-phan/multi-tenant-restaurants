package main

// @title Restaurant Management API
// @version 1.0
// @description Multi-tenant Restaurant Management System
// @host localhost:8080
// @BasePath /api/v1

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"restaurant-backend/internal/config"
	"restaurant-backend/internal/database"
	"restaurant-backend/internal/logger"
	"restaurant-backend/internal/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	var migrate = flag.Bool("migrate", false, "Run database migrations (up)")
	var migrateDown = flag.Bool("migrate-down", false, "Rollback last migration (down)")
	var migrateStatus = flag.Bool("migrate-status", false, "Show migration status")
	var bootstrap = flag.Bool("bootstrap", false, "Bootstrap platform organization and admin user")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Logger
	if err := logger.Initialize(cfg.Environment); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Restaurant Backend starting...", zap.String("environment", cfg.Environment))

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	db, err := database.NewConnection(cfg)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		os.Exit(1)
	}

	// Handle migration commands
	if *migrate {
		if err := database.RunMigrations(db, cfg); err != nil {
			logger.Error("Failed to run migrations", zap.Error(err))
			os.Exit(1)
		}
		logger.Info("Migrations completed successfully")
		os.Exit(0)
	}

	if *migrateDown {
		if err := database.RunMigrationsDown(db, cfg); err != nil {
			logger.Error("Failed to rollback migration", zap.Error(err))
			os.Exit(1)
		}
		logger.Info("Migration rollback completed successfully")
		os.Exit(0)
	}

	if *migrateStatus {
		if err := database.ShowMigrationStatus(db, cfg); err != nil {
			logger.Error("Failed to show migration status", zap.Error(err))
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *bootstrap {
		if err := database.BootstrapPlatform(db, cfg); err != nil {
			logger.Error("Failed to bootstrap platform", zap.Error(err))
			os.Exit(1)
		}
		logger.Info("Bootstrap completed successfully")
		os.Exit(0)
	}

	// Setup router
	r := router.SetupRouter(cfg, db)

	// Configure server with graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server listening", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to listen", zap.Error(err))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
