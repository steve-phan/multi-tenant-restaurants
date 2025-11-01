package handlers

import (
	"net/http"
	"restaurant-backend/internal/dto"
	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/repositories"
	"restaurant-backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CategoryHandler handles menu category-related requests
type CategoryHandler struct {
	categoryRepo    *repositories.CategoryRepository
	categoryService *services.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler instance
func NewCategoryHandler(categoryRepo *repositories.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{
		categoryRepo:    categoryRepo,
		categoryService: services.NewCategoryService(categoryRepo),
	}
}

// CreateCategory handles category creation
// @Summary Create Menu Category
// @Description Create a new menu category (e.g., "Hot Food", "Drinks", "Vegans")
// @Tags categories
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoryRequest true "Category data"
// @Success 201 {object} models.MenuCategory
// @Failure 400 {object} map[string]string
// @Router /api/v1/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	// Bind request
	var req dto.CreateCategoryRequest
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

	// Create category using service
	category, err := h.categoryService.CreateCategory(&req, restaurantID.(uint))
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "category name already taken" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
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
// @Description Update an existing menu category (only provided fields will be updated)
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Category update data (only provided fields will be updated)"
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

	// Bind update request
	var req dto.UpdateCategoryRequest
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

	// Update category using service (with ownership validation)
	category, err := h.categoryService.UpdateCategory(uint(id), &req, restaurantID.(uint))
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "category not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
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
