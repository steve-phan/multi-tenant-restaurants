package router

import (
	"restaurant-backend/internal/handlers"
	"restaurant-backend/internal/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupPublicMenuRoutes configures public menu routes (no authentication required)
// Clients can view menu items and categories for ordering
func setupPublicMenuRoutes(api *gin.RouterGroup, db *gorm.DB) {
	// Initialize repositories
	categoryRepo := repositories.NewCategoryRepository(db)
	menuItemRepo := repositories.NewMenuItemRepository(db)

	// Initialize handler
	publicMenuHandler := handlers.NewPublicMenuHandler(categoryRepo, menuItemRepo)

	// Public menu routes (no authentication required)
	public := api.Group("/public/restaurants")
	{
		// Get menu item details for ordering
		public.GET("/:restaurant_id/menu-items/:item_id", publicMenuHandler.GetMenuItemPublic)

		// List categories for a restaurant
		public.GET("/:restaurant_id/categories", publicMenuHandler.ListCategoriesPublic)

		// List menu items for a restaurant (optionally filtered by category)
		public.GET("/:restaurant_id/menu-items", publicMenuHandler.ListMenuItemsPublic)
	}
}
