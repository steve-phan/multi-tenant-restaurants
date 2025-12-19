package router

import (
	"restaurant-backend/internal/config"
	"restaurant-backend/internal/handlers"
	"restaurant-backend/internal/repositories"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupProfileRoutes configures profile management routes
func setupProfileRoutes(protected *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	// Initialize repository
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	profileService := services.NewProfileService(userRepo)

	// Initialize S3 service (optional)
	var s3Service *services.S3Service
	if cfg.S3BucketName != "" {
		if s3Svc, err := services.NewS3Service(cfg); err == nil {
			s3Service = s3Svc
		}
	}

	// Initialize handler
	profileHandler := handlers.NewProfileHandler(profileService, s3Service)

	// Profile routes (authenticated user access)
	profile := protected.Group("/profile")
	{
		profile.GET("", profileHandler.GetProfile)
		profile.PUT("", profileHandler.UpdateProfile)
		profile.PUT("/password", profileHandler.ChangePassword)
		profile.PUT("/preferences", profileHandler.UpdatePreferences)
		if s3Service != nil {
			profile.POST("/avatar", profileHandler.UploadAvatar)
		}
	}
}
