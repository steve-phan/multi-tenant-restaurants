package router

import (
	"restaurant-backend/internal/handlers"
	"restaurant-backend/internal/repositories"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupUserRoutes configures user management routes
func setupUserRoutes(protected *gin.RouterGroup, db *gorm.DB) {
	// Initialize repository
	userRepo := repositories.NewUserRepository(db)

	// Initialize service
	userService := services.NewUserService(userRepo)

	// Initialize handler
	userHandler := handlers.NewUserHandler(userService)

	// User routes (Admin/Staff access)
	users := protected.Group("/users")
	{
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
		users.POST("", userHandler.CreateUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.PATCH("/:id/status", userHandler.ToggleUserStatus)
	}
}
