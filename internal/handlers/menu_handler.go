package handlers

import (
	"net/http"
	"strconv"

	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

// MenuHandler handles menu-related requests
type MenuHandler struct {
	menuRepo *repositories.MenuRepository
}

// NewMenuHandler creates a new MenuHandler instance
func NewMenuHandler(menuRepo *repositories.MenuRepository) *MenuHandler {
	return &MenuHandler{
		menuRepo: menuRepo,
	}
}

// CreateMenu handles menu creation
// @Summary Create Menu
// @Description Create a new menu for the restaurant
// @Tags menus
// @Accept json
// @Produce json
// @Param menu body models.Menu true "Menu data"
// @Success 201 {object} models.Menu
// @Failure 400 {object} map[string]string
// @Router /api/v1/menus [post]
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var menu models.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get restaurant ID from context (set by middleware)
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	menu.RestaurantID = restaurantID.(uint)

	if err := h.menuRepo.Create(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, menu)
}

// GetMenu handles getting a menu by ID
// @Summary Get Menu
// @Description Get a menu by ID
// @Tags menus
// @Produce json
// @Param id path int true "Menu ID"
// @Success 200 {object} models.Menu
// @Failure 404 {object} map[string]string
// @Router /api/v1/menus/{id} [get]
func (h *MenuHandler) GetMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu ID"})
		return
	}

	menu, err := h.menuRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// ListMenus handles listing all menus for the restaurant
// @Summary List Menus
// @Description List all menus for the restaurant
// @Tags menus
// @Produce json
// @Success 200 {array} models.Menu
// @Router /api/v1/menus [get]
func (h *MenuHandler) ListMenus(c *gin.Context) {
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	menus, err := h.menuRepo.GetByRestaurantID(restaurantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menus)
}

// UpdateMenu handles updating a menu
// @Summary Update Menu
// @Description Update an existing menu
// @Tags menus
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Param menu body models.Menu true "Menu data"
// @Success 200 {object} models.Menu
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/menus/{id} [put]
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu ID"})
		return
	}

	var menu models.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	menu.ID = uint(id)
	restaurantID, _ := c.Get(middleware.RestaurantIDKey)
	menu.RestaurantID = restaurantID.(uint)

	if err := h.menuRepo.Update(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// DeleteMenu handles deleting a menu
// @Summary Delete Menu
// @Description Delete a menu
// @Tags menus
// @Param id path int true "Menu ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /api/v1/menus/{id} [delete]
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu ID"})
		return
	}

	if err := h.menuRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

