package routes

import (
	"restaurant-saas/handlers"
	"restaurant-saas/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Auth
	api.POST("/register", handlers.Register)
	api.POST("/login", handlers.Login)

	// Restaurant Admin (Protected)
	admin := api.Group("/restaurants")
	admin.Use(middleware.AuthMiddleware())
	{
		admin.POST("/", handlers.CreateRestaurant) // Create new restaurant
		admin.GET("/:id", handlers.GetRestaurant)
		admin.PUT("/:id", handlers.GetRestaurant) // Placeholder for update
		admin.POST("/:id/menu/categories", handlers.CreateMenuCategory)
		admin.POST("/:id/menu/items", handlers.CreateMenuItem)
		admin.POST("/:id/tables", handlers.CreateTable)
		admin.GET("/:id/bookings", handlers.GetBookings)
		admin.GET("/:id/orders", handlers.GetOrders)
	}

	// Customer (Public)
	public := api.Group("/public/restaurants")
	{
		public.GET("/:id/menu", handlers.GetPublicMenu)
		public.GET("/:id/tables/available", handlers.GetAvailableTables)
		public.POST("/:id/bookings", handlers.CreateBooking)
		public.POST("/:id/orders", handlers.CreateOrder)
	}
}
