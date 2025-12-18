package router

import (
	"restaurant-backend/internal/handlers"
	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/repositories"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupRestaurantRoutes configures restaurant-related routes
func setupRestaurantRoutes(api *gin.RouterGroup, protected *gin.RouterGroup, db *gorm.DB, emailService *services.EmailService) {
	// Initialize repositories and services for restaurant routes
	restaurantRepo := repositories.NewRestaurantRepository(db)
	userRepo := repositories.NewUserRepository(db)
	restaurantService := services.NewRestaurantService(restaurantRepo, userRepo, emailService)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService, restaurantRepo)

	// Public restaurant registration route
	restaurantPublic := api.Group("/restaurants")
	{
		restaurantPublic.POST("/register", restaurantHandler.RegisterRestaurant)
	}

	// Protected restaurant management routes (KAM/Admin only)
	restaurants := protected.Group("/restaurants")
	restaurants.Use(middleware.RequireKAMOrAdmin())
	{
		restaurants.GET("", restaurantHandler.ListRestaurants)
		restaurants.GET("/pending", restaurantHandler.ListPendingRestaurants)
		restaurants.GET("/:id", restaurantHandler.GetRestaurant)
		restaurants.POST("/:id/activate", restaurantHandler.ActivateRestaurant)
		restaurants.PATCH("/:id/status", restaurantHandler.UpdateRestaurantStatus)
		restaurants.PUT("/:id/assign-kam", restaurantHandler.AssignKAM)
	}
}
