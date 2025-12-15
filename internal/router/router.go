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
	r := gin.Default()

	// Add CORS middleware
	r.Use(corsMiddleware(cfg))

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
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
		setupRestaurantRoutes(api, protected, db)

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
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
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
