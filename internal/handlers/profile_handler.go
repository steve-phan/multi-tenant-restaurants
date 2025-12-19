package handlers

import (
	"errors"
	"net/http"

	"restaurant-backend/internal/ctx"
	"restaurant-backend/internal/dto"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// ProfileHandler handles profile management requests
type ProfileHandler struct {
	profileService *services.ProfileService
	s3Service      *services.S3Service
}

// NewProfileHandler creates a new ProfileHandler instance
func NewProfileHandler(profileService *services.ProfileService, s3Service *services.S3Service) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
		s3Service:      s3Service,
	}
}

// GetProfile handles retrieving the current user's profile
// @Summary Get Profile
// @Description Get the current authenticated user's profile
// @Tags profile
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Get user ID from context
	userID, ok := ctx.GetUserID(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	user, err := h.profileService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, services.ErrProfileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile handles updating the current user's profile
// @Summary Update Profile
// @Description Update the current authenticated user's profile
// @Tags profile
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileDTO true "Profile update data"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/profile [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context
	userID, ok := ctx.GetUserID(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	var req dto.UpdateProfileDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.profileService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, services.ErrProfileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePassword handles changing the current user's password
// @Summary Change Password
// @Description Change the current authenticated user's password
// @Tags profile
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordDTO true "Password change data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/profile/password [put]
func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	// Get user ID from context
	userID, ok := ctx.GetUserID(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	var req dto.ChangePasswordDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.profileService.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		if errors.Is(err, services.ErrInvalidPassword) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

// UpdatePreferences handles updating the current user's preferences
// @Summary Update Preferences
// @Description Update the current authenticated user's preferences
// @Tags profile
// @Accept json
// @Produce json
// @Param request body dto.UpdatePreferencesDTO true "Preferences update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/profile/preferences [put]
func (h *ProfileHandler) UpdatePreferences(c *gin.Context) {
	// Get user ID from context
	userID, ok := ctx.GetUserID(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	var req dto.UpdatePreferencesDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.profileService.UpdatePreferences(c.Request.Context(), userID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "preferences updated successfully"})
}

// UploadAvatar handles uploading an avatar for the current user
// @Summary Upload Avatar
// @Description Upload an avatar image for the current authenticated user
// @Tags profile
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "Avatar image file"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/profile/avatar [post]
func (h *ProfileHandler) UploadAvatar(c *gin.Context) {
	// Get user ID from context
	userID, ok := ctx.GetUserID(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	// Get restaurant ID from context
	restaurantID, ok := ctx.GetRestaurantID(c.Request.Context())
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	// Get file from form
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "avatar file is required"})
		return
	}

	// Open file
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer fileContent.Close()

	// Upload to S3 using existing S3Service
	fileName := file.Filename
	fileType := file.Header.Get("Content-Type")
	avatarKey, err := h.s3Service.UploadFile(c.Request.Context(), restaurantID, fileName, fileType, fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload avatar"})
		return
	}

	// Update user's avatar URL (storing the S3 key)
	if err := h.profileService.UpdateAvatar(c.Request.Context(), userID, avatarKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "avatar uploaded successfully",
		"avatar_key": avatarKey,
	})
}
