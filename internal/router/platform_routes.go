package router

import (
	"restaurant-backend/internal/handlers"
	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/repositories"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupPlatformRoutes configures platform-level routes (KAM management)
func setupPlatformRoutes(protected *gin.RouterGroup, db *gorm.DB, authService *services.AuthService) {
	// Initialize platform service and handler
	platformRepo := repositories.NewRestaurantRepository(db)
	platformUserRepo := repositories.NewUserRepository(db)
	platformService := services.NewPlatformService(platformRepo, platformUserRepo)
	platformHandler := handlers.NewPlatformHandler(platformService, authService)

	// Platform management routes (KAM/Admin only)
	platform := protected.Group("/platform")
	platform.Use(middleware.RequireKAMOrAdmin())
	{
		platform.POST("/kams", platformHandler.CreateKAM)
		platform.GET("/kams", platformHandler.ListKAMs)
	}
}

