package handlers

import (
	"net/http"
	"strconv"

	"restaurant-backend/internal/dto"
	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/repositories"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// MenuItemHandler handles menu item-related requests
type MenuItemHandler struct {
	menuItemRepo    *repositories.MenuItemRepository
	menuItemService *services.MenuItemService
}

// NewMenuItemHandler creates a new MenuItemHandler instance
func NewMenuItemHandler(menuItemRepo *repositories.MenuItemRepository) *MenuItemHandler {
	return &MenuItemHandler{
		menuItemRepo:    menuItemRepo,
		menuItemService: services.NewMenuItemService(menuItemRepo),
	}
}

// CreateMenuItem handles menu item creation
// @Summary Create Menu Item
// @Description Create a new menu item
// @Tags menu-items
// @Accept json
// @Produce json
// @Param request body dto.CreateMenuItemRequest true "Menu Item data"
// @Success 201 {object} models.MenuItem
// @Failure 400 {object} map[string]string
// @Router /api/v1/menu-items [post]
func (h *MenuItemHandler) CreateMenuItem(c *gin.Context) {
	// Bind request
	var req dto.CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get restaurant ID from context (set by middleware)
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	// Create menu item using service
	menuItem, err := h.menuItemService.CreateMenuItem(&req, restaurantID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, menuItem)
}

// GetMenuItem handles getting a menu item by ID (protected)
// @Summary Get Menu Item (Protected)
// @Description Get a menu item by ID with all details including images
// @Tags menu-items
// @Produce json
// @Param id path int true "Menu Item ID"
// @Success 200 {object} models.MenuItem
// @Failure 404 {object} map[string]string
// @Router /api/v1/menu-items/{id} [get]
func (h *MenuItemHandler) GetMenuItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item ID"})
		return
	}

	menuItem, err := h.menuItemRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "menu item not found"})
		return
	}

	c.JSON(http.StatusOK, menuItem)
}

// ListMenuItems handles listing menu items
// @Summary List Menu Items
// @Description List menu items, optionally filtered by category ID
// @Tags menu-items
// @Produce json
// @Param category_id query int false "Category ID filter"
// @Success 200 {array} models.MenuItem
// @Router /api/v1/menu-items [get]
func (h *MenuItemHandler) ListMenuItems(c *gin.Context) {
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	// Check if category_id query parameter is provided
	categoryIDParam := c.Query("category_id")
	if categoryIDParam != "" {
		categoryID, err := strconv.ParseUint(categoryIDParam, 10, 32)
		if err == nil {
			menuItems, err := h.menuItemRepo.GetByCategoryID(uint(categoryID))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, menuItems)
			return
		}
	}

	// Otherwise, get all menu items for the restaurant
	menuItems, err := h.menuItemRepo.GetByRestaurantID(restaurantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menuItems)
}

// UpdateMenuItem handles updating a menu item
// @Summary Update Menu Item
// @Description Update an existing menu item (only provided fields will be updated)
// @Tags menu-items
// @Accept json
// @Produce json
// @Param id path int true "Menu Item ID"
// @Param request body dto.UpdateMenuItemRequest true "Menu Item update data (only provided fields will be updated)"
// @Success 200 {object} models.MenuItem
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/menu-items/{id} [put]
func (h *MenuItemHandler) UpdateMenuItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item ID"})
		return
	}

	// Bind update request
	var req dto.UpdateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get restaurant ID from context (set by middleware)
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	// Update menu item using service (with ownership validation)
	menuItem, err := h.menuItemService.UpdateMenuItem(uint(id), &req, restaurantID.(uint))
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "menu item not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menuItem)
}

// DeleteMenuItem handles deleting a menu item
// @Summary Delete Menu Item
// @Description Delete a menu item
// @Tags menu-items
// @Param id path int true "Menu Item ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /api/v1/menu-items/{id} [delete]
func (h *MenuItemHandler) DeleteMenuItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item ID"})
		return
	}

	if err := h.menuItemRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
