package router

import (
	"restaurant-backend/internal/config"
	"restaurant-backend/internal/handlers"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// setupImageRoutes configures image-related routes (S3)
func setupImageRoutes(protected *gin.RouterGroup, cfg *config.Config) *handlers.ImageHandler {
	// Initialize S3 service (optional, only if configured)
	var s3Service *services.S3Service
	var imageHandler *handlers.ImageHandler

	if cfg.S3BucketName != "" {
		if s3Svc, err := services.NewS3Service(cfg); err == nil {
			s3Service = s3Svc
			imageHandler = handlers.NewImageHandler(s3Service)

			// Image routes (if S3 is configured)
			images := protected.Group("/images")
			{
				images.POST("/upload", imageHandler.UploadImage)
				images.GET("/*key", imageHandler.GetImageURL)
				images.DELETE("/*key", imageHandler.DeleteImage)
			}
		}
		// Log error but don't fail startup if S3 is not configured
		// In production, this should be handled more gracefully
	}

	return imageHandler
}
