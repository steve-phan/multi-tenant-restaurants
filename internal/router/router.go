package router

import (
	"restaurant-backend/internal/config"
	"restaurant-backend/internal/handlers"
	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/repositories"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter configures and returns the Gin router
func SetupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	// Use gin.New() instead of Default() to skip default logger
	r := gin.New()

	// Add middlewares
	r.Use(middleware.RequestLogger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware(cfg))

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	emailService := services.NewEmailService(cfg)
	authService := services.NewAuthService(db, cfg, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "restaurant-backend",
		})
	})

	// Public API routes
	api := r.Group("/api/v1")
	{
		// Setup authentication routes
		setupAuthRoutes(api, authHandler)

		// Setup public menu routes (no authentication required for viewing menu)
		setupPublicMenuRoutes(api, db)
	}

	// Protected API routes
	protected := api.Group("")
	protected.Use(middleware.RequireAuth(authService))
	protected.Use(middleware.SetTenantContext(db))
	{
		// Setup business routes (menus, orders, reservations)
		setupBusinessRoutes(protected, db)

		// Setup restaurant routes (includes public registration)
		setupRestaurantRoutes(api, protected, db, emailService)

		// Setup platform routes (KAM management)
		setupPlatformRoutes(protected, db, authService)

		// Setup image routes (S3)
		setupImageRoutes(protected, cfg)
	}

	return r
}

// corsMiddleware handles CORS
func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if len(cfg.CORSAllowedOrigins) == 1 && cfg.CORSAllowedOrigins[0] == "*" {
			// When wildcard is configured but credentials are needed,
			// echo back the requesting origin instead of using "*"
			if origin != "" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			}
		} else {
			// Check if origin is in allowed list
			for _, allowedOrigin := range cfg.CORSAllowedOrigins {
				if origin == allowedOrigin {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
