package handlers

import (
	"net/http"
	"strconv"

	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

// CategoryHandler handles menu category-related requests
type CategoryHandler struct {
	categoryRepo *repositories.CategoryRepository
}

// NewCategoryHandler creates a new CategoryHandler instance
func NewCategoryHandler(categoryRepo *repositories.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{
		categoryRepo: categoryRepo,
	}
}

// CreateCategory handles category creation
// @Summary Create Menu Category
// @Description Create a new menu category (e.g., "Hot Food", "Drinks", "Vegans")
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.MenuCategory true "Category data"
// @Success 201 {object} models.MenuCategory
// @Failure 400 {object} map[string]string
// @Router /api/v1/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category models.MenuCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get restaurant ID from context (set by middleware)
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	category.RestaurantID = restaurantID.(uint)

	if err := h.categoryRepo.Create(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategory handles getting a category by ID
// @Summary Get Menu Category
// @Description Get a menu category by ID with its items
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.MenuCategory
// @Failure 404 {object} map[string]string
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	category, err := h.categoryRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// ListCategories handles listing all categories for the restaurant
// @Summary List Menu Categories
// @Description List all menu categories for the restaurant
// @Tags categories
// @Produce json
// @Success 200 {array} models.MenuCategory
// @Router /api/v1/categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	categories, err := h.categoryRepo.GetByRestaurantID(restaurantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// UpdateCategory handles updating a category
// @Summary Update Menu Category
// @Description Update an existing menu category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body models.MenuCategory true "Category data"
// @Success 200 {object} models.MenuCategory
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	var category models.MenuCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.ID = uint(id)

	// Get restaurant ID from context (set by middleware)
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	category.RestaurantID = restaurantID.(uint)

	if err := h.categoryRepo.Update(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory handles deleting a category
// @Summary Delete Menu Category
// @Description Delete a menu category
// @Tags categories
// @Param id path int true "Category ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	if err := h.categoryRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

