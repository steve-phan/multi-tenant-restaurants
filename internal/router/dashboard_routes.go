package router

import (
	"restaurant-backend/internal/handlers"
	"restaurant-backend/internal/repositories"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupDashboardRoutes configures dashboard routes
func setupDashboardRoutes(protected *gin.RouterGroup, db *gorm.DB) {
	// Initialize repositories
	orderRepo := repositories.NewOrderRepository(db)
	reservationRepo := repositories.NewReservationRepository(db)

	// Initialize service
	dashboardService := services.NewDashboardService(orderRepo, reservationRepo)

	// Initialize handler
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// Dashboard routes
	dashboard := protected.Group("/dashboard")
	{
		dashboard.GET("/stats", dashboardHandler.GetDashboardStats)
		dashboard.GET("/recent-orders", dashboardHandler.GetRecentOrders)
		dashboard.GET("/analytics", dashboardHandler.GetAnalytics)
	}
}
