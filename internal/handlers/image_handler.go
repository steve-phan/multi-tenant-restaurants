package handlers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// ImageHandler handles image upload and download
type ImageHandler struct {
	s3Service *services.S3Service
}

// NewImageHandler creates a new ImageHandler instance
func NewImageHandler(s3Service *services.S3Service) *ImageHandler {
	return &ImageHandler{
		s3Service: s3Service,
	}
}

// UploadImage handles image upload
// @Summary Upload Image
// @Description Upload an image file to S3 with tenant isolation
// @Tags images
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/images/upload [post]
func (h *ImageHandler) UploadImage(c *gin.Context) {
	// Get restaurant ID from context
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 10MB limit"})
		return
	}

	// Validate file type
	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type. Allowed: jpg, jpeg, png, gif, webp"})
		return
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	// Determine content type
	contentType := http.DetectContentType([]byte{})
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	}

	// Upload to S3
	ctx := context.Background()
	key, err := h.s3Service.UploadFile(ctx, restaurantID.(uint), file.Filename, contentType, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to upload file: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"key":  key,
		"url":  fmt.Sprintf("/api/v1/images/%s", key), // Relative URL for getting presigned URL
		"size": file.Size,
	})
}

// GetImageURL generates a presigned URL for image access
// @Summary Get Image URL
// @Description Generate a presigned URL for accessing an image
// @Tags images
// @Produce json
// @Param key path string true "S3 Object Key"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/images/{key} [get]
func (h *ImageHandler) GetImageURL(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}

	// Get restaurant ID from context for validation (ensure tenant can only access their own images)
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	// Validate that the key belongs to the restaurant
	expectedPrefix := fmt.Sprintf("restaurant-%d/", restaurantID.(uint))
	if len(key) < len(expectedPrefix) || key[:len(expectedPrefix)] != expectedPrefix {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Generate presigned URL (valid for 1 hour)
	ctx := context.Background()
	url, err := h.s3Service.GeneratePresignedURL(ctx, key, time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":     url,
		"expires": time.Now().Add(time.Hour).Unix(),
	})
}

// DeleteImage handles image deletion
// @Summary Delete Image
// @Description Delete an image from S3
// @Tags images
// @Param key path string true "S3 Object Key"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/v1/images/{key} [delete]
func (h *ImageHandler) DeleteImage(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}

	// Get restaurant ID from context for validation
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	// Validate that the key belongs to the restaurant
	expectedPrefix := fmt.Sprintf("restaurant-%d/", restaurantID.(uint))
	if len(key) < len(expectedPrefix) || key[:len(expectedPrefix)] != expectedPrefix {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Delete from S3
	ctx := context.Background()
	if err := h.s3Service.DeleteFile(ctx, key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file"})
		return
	}

	c.Status(http.StatusNoContent)
}

