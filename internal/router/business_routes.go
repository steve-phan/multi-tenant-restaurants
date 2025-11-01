package router

import (
	"restaurant-backend/internal/handlers"
	"restaurant-backend/internal/repositories"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupBusinessRoutes configures business-related routes (menus, orders, reservations)
func setupBusinessRoutes(protected *gin.RouterGroup, db *gorm.DB) {
	// Initialize repositories
	menuRepo := repositories.NewMenuRepository(db)
	menuItemRepo := repositories.NewMenuItemRepository(db)
	reservationRepo := repositories.NewReservationRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	orderItemRepo := repositories.NewOrderItemRepository(db)

	// Initialize services
	reservationService := services.NewReservationService(reservationRepo)
	orderService := services.NewOrderService(orderRepo, orderItemRepo, menuItemRepo)

	// Initialize handlers
	menuHandler := handlers.NewMenuHandler(menuRepo)
	menuItemHandler := handlers.NewMenuItemHandler(menuItemRepo)
	reservationHandler := handlers.NewReservationHandler(reservationService, reservationRepo)
	orderHandler := handlers.NewOrderHandler(orderService, orderRepo)

	// Menu routes
	menus := protected.Group("/menus")
	{
		menus.POST("", menuHandler.CreateMenu)
		menus.GET("", menuHandler.ListMenus)
		menus.GET("/:id", menuHandler.GetMenu)
		menus.PUT("/:id", menuHandler.UpdateMenu)
		menus.DELETE("/:id", menuHandler.DeleteMenu)
	}

	// Menu Item routes
	menuItems := protected.Group("/menu-items")
	{
		menuItems.POST("", menuItemHandler.CreateMenuItem)
		menuItems.GET("", menuItemHandler.ListMenuItems)
		menuItems.GET("/:id", menuItemHandler.GetMenuItem)
		menuItems.PUT("/:id", menuItemHandler.UpdateMenuItem)
		menuItems.DELETE("/:id", menuItemHandler.DeleteMenuItem)
	}

	// Reservation routes
	reservations := protected.Group("/reservations")
	{
		reservations.POST("", reservationHandler.CreateReservation)
		reservations.GET("", reservationHandler.ListReservations)
		reservations.GET("/:id", reservationHandler.GetReservation)
		reservations.PUT("/:id", reservationHandler.UpdateReservation)
		reservations.DELETE("/:id", reservationHandler.DeleteReservation)
	}

	// Order routes
	orders := protected.Group("/orders")
	{
		orders.POST("", orderHandler.CreateOrder)
		orders.GET("", orderHandler.ListOrders)
		orders.GET("/:id", orderHandler.GetOrder)
		orders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
	}
}

