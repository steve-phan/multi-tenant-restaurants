package handlers

import (
	"net/http"

	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// PlatformHandler handles platform-level operations (KAM management)
type PlatformHandler struct {
	platformService *services.PlatformService
	authService     *services.AuthService
}

// NewPlatformHandler creates a new PlatformHandler instance
func NewPlatformHandler(
	platformService *services.PlatformService,
	authService *services.AuthService,
) *PlatformHandler {
	return &PlatformHandler{
		platformService: platformService,
		authService:     authService,
	}
}

// CreateKAM handles KAM user creation (KAM/Admin only)
// @Summary Create KAM
// @Description Create a new Key Account Manager user (only by existing KAMs/Admins)
// @Tags platform
// @Accept json
// @Produce json
// @Param request body services.CreateKAMRequest true "KAM creation data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/v1/platform/kams [post]
func (h *PlatformHandler) CreateKAM(c *gin.Context) {
	var req services.CreateKAMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get creator user ID from context
	createdBy, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user context not found"})
		return
	}

	// Create KAM user structure
	user, err := h.platformService.CreateKAM(&req, createdBy.(uint))
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "only platform KAMs or Admins can create new KAM users" {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	user.PasswordHash = string(hashedPassword)

	// Create user in database
	if err := h.platformService.CreateKAMUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user: " + err.Error()})
		return
	}

	// Clear password hash from response
	user.PasswordHash = ""

	c.JSON(http.StatusCreated, user)
}

// ListKAMs handles listing all KAM users
// @Summary List KAMs
// @Description List all Key Account Manager users
// @Tags platform
// @Produce json
// @Success 200 {array} models.User
// @Failure 403 {object} map[string]string
// @Router /api/v1/platform/kams [get]
func (h *PlatformHandler) ListKAMs(c *gin.Context) {
	kams, err := h.platformService.ListKAMs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kams)
}

