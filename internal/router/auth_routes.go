package router

import (
	"restaurant-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

// setupAuthRoutes configures authentication routes
func setupAuthRoutes(api *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		// User registration (for restaurant admins to create staff/users)
		// Note: KAM role is NOT allowed via this endpoint
		auth.POST("/register", authHandler.Register)
	}
}
