package router

import (
	"restaurant-backend/internal/handlers"
	"restaurant-backend/internal/repositories"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupBusinessRoutes configures business-related routes (categories, menu items, orders, reservations)
func setupBusinessRoutes(protected *gin.RouterGroup, db *gorm.DB) {
	// Initialize repositories
	categoryRepo := repositories.NewCategoryRepository(db)
	menuItemRepo := repositories.NewMenuItemRepository(db)
	reservationRepo := repositories.NewReservationRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	orderItemRepo := repositories.NewOrderItemRepository(db)

	// Initialize services
	reservationService := services.NewReservationService(reservationRepo)
	orderService := services.NewOrderService(orderRepo, orderItemRepo, menuItemRepo)

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	menuItemHandler := handlers.NewMenuItemHandler(menuItemRepo)
	reservationHandler := handlers.NewReservationHandler(reservationService, reservationRepo)
	orderHandler := handlers.NewOrderHandler(orderService, orderRepo)

	// Menu Category routes (Admin/Staff only - for managing categories)
	categories := protected.Group("/categories")
	{
		categories.POST("", categoryHandler.CreateCategory)
		categories.GET("", categoryHandler.ListCategories)
		categories.GET("/:id", categoryHandler.GetCategory)
		categories.PUT("/:id", categoryHandler.UpdateCategory)
		categories.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	// Menu Item routes (Admin/Staff only - for managing items)
	menuItems := protected.Group("/menu-items")
	{
		menuItems.POST("", menuItemHandler.CreateMenuItem)
		menuItems.GET("", menuItemHandler.ListMenuItems)
		menuItems.GET("/:id", menuItemHandler.GetMenuItem)
		menuItems.PUT("/:id", menuItemHandler.UpdateMenuItem)
		menuItems.DELETE("/:id", menuItemHandler.DeleteMenuItem)
	}

	// Menu Item Image routes (Admin/Staff only - for managing item images)
	// Using separate prefix to avoid routing conflicts with /menu-items/:id
	imageRepo := repositories.NewMenuItemImageRepository(db)
	imageHandler := handlers.NewMenuItemImageHandler(imageRepo)
	menuItemImages := protected.Group("/menu-item-images")
	{
		menuItemImages.POST("/:item_id", imageHandler.CreateMenuItemImage)
		menuItemImages.GET("/:item_id", imageHandler.ListMenuItemImages)
		menuItemImages.DELETE("/:item_id/:image_id", imageHandler.DeleteMenuItemImage)
		menuItemImages.PUT("/:item_id/:image_id/primary", imageHandler.SetPrimaryImage)
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
