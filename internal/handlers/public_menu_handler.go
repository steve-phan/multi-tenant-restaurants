package handlers

import (
	"net/http"
	"strconv"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

// PublicMenuHandler handles public menu-related requests (no authentication required)
type PublicMenuHandler struct {
	categoryRepo *repositories.CategoryRepository
	menuItemRepo *repositories.MenuItemRepository
}

// NewPublicMenuHandler creates a new PublicMenuHandler instance
func NewPublicMenuHandler(
	categoryRepo *repositories.CategoryRepository,
	menuItemRepo *repositories.MenuItemRepository,
) *PublicMenuHandler {
	return &PublicMenuHandler{
		categoryRepo: categoryRepo,
		menuItemRepo: menuItemRepo,
	}
}

// GetMenuItemPublic handles getting a menu item by ID for public access
// @Summary Get Menu Item (Public)
// @Description Get menu item details for ordering (no authentication required)
// @Tags public-menu
// @Produce json
// @Param restaurant_id path int true "Restaurant ID"
// @Param item_id path int true "Menu Item ID"
// @Success 200 {object} models.MenuItem
// @Failure 404 {object} map[string]string
// @Router /api/v1/public/restaurants/{restaurant_id}/menu-items/{item_id} [get]
func (h *PublicMenuHandler) GetMenuItemPublic(c *gin.Context) {
	restaurantID, err := strconv.ParseUint(c.Param("restaurant_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item ID"})
		return
	}

	menuItem, err := h.menuItemRepo.GetByIDPublic(uint(itemID), uint(restaurantID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "menu item not found"})
		return
	}

	c.JSON(http.StatusOK, menuItem)
}

// ListCategoriesPublic handles listing categories for a restaurant (public access)
// @Summary List Menu Categories (Public)
// @Description List all menu categories for a restaurant (no authentication required)
// @Tags public-menu
// @Produce json
// @Param restaurant_id path int true "Restaurant ID"
// @Success 200 {array} models.MenuCategory
// @Router /api/v1/public/restaurants/{restaurant_id}/categories [get]
func (h *PublicMenuHandler) ListCategoriesPublic(c *gin.Context) {
	restaurantID, err := strconv.ParseUint(c.Param("restaurant_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	categories, err := h.categoryRepo.GetByRestaurantID(uint(restaurantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// ListMenuItemsPublic handles listing menu items for a restaurant/category (public access)
// @Summary List Menu Items (Public)
// @Description List menu items for a restaurant, optionally filtered by category (no authentication required)
// @Tags public-menu
// @Produce json
// @Param restaurant_id path int true "Restaurant ID"
// @Param category_id query int false "Category ID filter"
// @Success 200 {array} models.MenuItem
// @Router /api/v1/public/restaurants/{restaurant_id}/menu-items [get]
func (h *PublicMenuHandler) ListMenuItemsPublic(c *gin.Context) {
	restaurantID, err := strconv.ParseUint(c.Param("restaurant_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	// Check if category_id query parameter is provided
	categoryIDParam := c.Query("category_id")
	if categoryIDParam != "" {
		categoryID, err := strconv.ParseUint(categoryIDParam, 10, 32)
		if err == nil {
			// Get items for specific category (need to verify category belongs to restaurant)
			menuItems, err := h.menuItemRepo.GetByCategoryID(uint(categoryID))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			// Filter by restaurant_id to ensure proper access
			var filteredItems []models.MenuItem
			for _, item := range menuItems {
				if item.RestaurantID == uint(restaurantID) {
					filteredItems = append(filteredItems, item)
				}
			}
			c.JSON(http.StatusOK, filteredItems)
			return
		}
	}

	// Otherwise, get all menu items for the restaurant
	menuItems, err := h.menuItemRepo.GetByRestaurantID(uint(restaurantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menuItems)
}
