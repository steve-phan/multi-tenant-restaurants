package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	ServerPort  string
	Environment string
	LogLevel    string

	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// AWS configuration
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string

	// S3 configuration
	S3BucketName string

	// JWT configuration
	JWTSecret     string
	JWTExpiration int // in hours

	// CORS configuration
	CORSAllowedOrigins []string

	// Brevo Email configuration
	BrevoAPIKey      string
	BrevoSenderEmail string
	BrevoSenderName  string
	FrontendURL      string

	// Bootstrap configuration (for initial admin user)
	BootstrapAdminEmail    string
	BootstrapAdminPassword string
}

// Load reads configuration from environment variables
// It first tries to load from .env file if it exists, then reads from environment
func Load() (*Config, error) {
	// Try to load .env file (ignore error if it doesn't exist)
	// In development mode or if ENVIRONMENT is not set, always try to load .env
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" {
		if err := godotenv.Load(); err != nil {
			// Only log if .env file doesn't exist, don't fail
			log.Printf("Note: .env file not found or couldn't be loaded: %v", err)
		} else {
			log.Println("✓ Loaded configuration from .env file")
		}
	} else {
		// In production, still try to load .env but don't log if it fails
		// Environment variables should be set explicitly in production
		_ = godotenv.Load()
	}

	cfg := &Config{
		ServerPort:             getEnv("SERVER_PORT", "8080"),
		Environment:            getEnv("ENVIRONMENT", "development"),
		LogLevel:               getEnv("LOG_LEVEL", "info"),
		DBHost:                 getEnv("DB_HOST", "localhost"),
		DBPort:                 getEnv("DB_PORT", "5432"),
		DBUser:                 getEnv("DB_USER", "postgres"),
		DBPassword:             getEnv("DB_PASSWORD", ""),
		DBName:                 getEnv("DB_NAME", "restaurant_db"),
		DBSSLMode:              getEnv("DB_SSL_MODE", "disable"),
		AWSRegion:              getEnv("AWS_REGION", "us-east-1"),
		AWSAccessKeyID:         getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey:     getEnv("AWS_SECRET_ACCESS_KEY", ""),
		S3BucketName:           getEnv("S3_BUCKET_NAME", ""),
		JWTSecret:              getEnv("JWT_SECRET", ""),
		JWTExpiration:          getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		BrevoAPIKey:            getEnv("BREVO_API_KEY", ""),
		BrevoSenderEmail:       getEnv("BREVO_SENDER_EMAIL", "noreply@restaurant-platform.local"),
		BrevoSenderName:        getEnv("BREVO_SENDER_NAME", "Restaurant Platform"),
		FrontendURL:            getEnv("FRONTEND_URL", "http://localhost:3000"),
		BootstrapAdminEmail:    getEnv("BOOTSTRAP_ADMIN_EMAIL", "admin@platform.local"),
		BootstrapAdminPassword: getEnv("BOOTSTRAP_ADMIN_PASSWORD", ""),
	}

	// Validate required fields
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}
	if cfg.JWTSecret == "" && cfg.Environment == "production" {
		return nil, fmt.Errorf("JWT_SECRET is required in production")
	}

	// Parse CORS origins
	corsOrigins := getEnv("CORS_ALLOWED_ORIGINS", "*")
	if corsOrigins != "*" {
		cfg.CORSAllowedOrigins = []string{corsOrigins}
	} else {
		cfg.CORSAllowedOrigins = []string{"*"}
	}

	return cfg, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	value := getEnv(key, fmt.Sprintf("%d", defaultValue))
	var intValue int
	if _, err := fmt.Sscanf(value, "%d", &intValue); err != nil {
		return defaultValue
	}
	return intValue
}
