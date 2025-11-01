package handlers

import (
	"net/http"
	"strconv"

	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

// MenuItemImageHandler handles menu item image-related requests
type MenuItemImageHandler struct {
	imageRepo *repositories.MenuItemImageRepository
}

// NewMenuItemImageHandler creates a new MenuItemImageHandler instance
func NewMenuItemImageHandler(imageRepo *repositories.MenuItemImageRepository) *MenuItemImageHandler {
	return &MenuItemImageHandler{
		imageRepo: imageRepo,
	}
}

// CreateMenuItemImage handles menu item image creation
// @Summary Add Image to Menu Item
// @Description Add an image to a menu item
// @Tags menu-item-images
// @Accept json
// @Produce json
// @Param item_id path int true "Menu Item ID"
// @Param image body models.MenuItemImage true "Image data"
// @Success 201 {object} models.MenuItemImage
// @Failure 400 {object} map[string]string
// @Router /api/v1/menu-item-images/:item_id [post]
func (h *MenuItemImageHandler) CreateMenuItemImage(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item ID"})
		return
	}

	var image models.MenuItemImage
	if err := c.ShouldBindJSON(&image); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get restaurant ID from context (set by middleware)
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	image.RestaurantID = restaurantID.(uint)
	image.MenuItemID = uint(itemID)

	// Create the image
	if err := h.imageRepo.Create(&image); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If this image should be primary, set it as primary
	if image.IsPrimary {
		if err := h.imageRepo.SetPrimary(image.MenuItemID, image.ID); err != nil {
			// Log error but don't fail the request - image is already created
		}
	}

	c.JSON(http.StatusCreated, image)
}

// ListMenuItemImages handles listing images for a menu item
// @Summary List Menu Item Images
// @Description List all images for a menu item
// @Tags menu-item-images
// @Produce json
// @Param item_id path int true "Menu Item ID"
// @Success 200 {array} models.MenuItemImage
// @Router /api/v1/menu-item-images/:item_id [get]
func (h *MenuItemImageHandler) ListMenuItemImages(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item ID"})
		return
	}

	images, err := h.imageRepo.GetByMenuItemID(uint(itemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, images)
}

// DeleteMenuItemImage handles deleting a menu item image
// @Summary Delete Menu Item Image
// @Description Delete an image from a menu item
// @Tags menu-item-images
// @Param item_id path int true "Menu Item ID"
// @Param image_id path int true "Image ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Router /api/v1/menu-item-images/:item_id/:image_id [delete]
func (h *MenuItemImageHandler) DeleteMenuItemImage(c *gin.Context) {
	imageID, err := strconv.ParseUint(c.Param("image_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image ID"})
		return
	}

	if err := h.imageRepo.Delete(uint(imageID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SetPrimaryImage handles setting an image as primary
// @Summary Set Primary Image
// @Description Set an image as the primary image for a menu item
// @Tags menu-item-images
// @Param item_id path int true "Menu Item ID"
// @Param image_id path int true "Image ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/menu-item-images/:item_id/:image_id/primary [put]
func (h *MenuItemImageHandler) SetPrimaryImage(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item ID"})
		return
	}

	imageID, err := strconv.ParseUint(c.Param("image_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image ID"})
		return
	}

	if err := h.imageRepo.SetPrimary(uint(itemID), uint(imageID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Primary image updated successfully"})
}

