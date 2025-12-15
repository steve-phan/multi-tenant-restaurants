package handlers

import (
	"net/http"
	"strconv"

	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repositories"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// RestaurantHandler handles restaurant-related requests
type RestaurantHandler struct {
	restaurantService *services.RestaurantService
	restaurantRepo    *repositories.RestaurantRepository
}

// NewRestaurantHandler creates a new RestaurantHandler instance
func NewRestaurantHandler(
	restaurantService *services.RestaurantService,
	restaurantRepo *repositories.RestaurantRepository,
) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantService: restaurantService,
		restaurantRepo:    restaurantRepo,
	}
}

// RegisterRestaurant handles restaurant registration (public endpoint)
// @Summary Register Restaurant
// @Description Register a new restaurant (will be in pending status)
// @Tags restaurants
// @Accept json
// @Produce json
// @Param request body services.RegisterRestaurantRequest true "Restaurant registration data"
// @Success 201 {object} models.Restaurant
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/restaurants/register [post]
func (h *RestaurantHandler) RegisterRestaurant(c *gin.Context) {
	var req services.RegisterRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurant, err := h.restaurantService.RegisterRestaurant(c.Request.Context(), &req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "restaurant with this email already exists" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Restaurant registered successfully. Awaiting activation by Key Account Manager.",
		"restaurant": restaurant,
	})
}

// ListRestaurants handles listing restaurants (KAM/Admin only)
// @Summary List Restaurants
// @Description List all restaurants (filtered by status and KAM if provided)
// @Tags restaurants
// @Produce json
// @Param status query string false "Filter by status (pending, active, inactive, suspended)"
// @Param kam_id query int false "Filter by KAM ID"
// @Success 200 {array} models.Restaurant
// @Failure 403 {object} map[string]string
// @Router /api/v1/restaurants [get]
func (h *RestaurantHandler) ListRestaurants(c *gin.Context) {
	var status *models.RestaurantStatus
	var kamID *uint

	// Get status filter
	statusParam := c.Query("status")
	if statusParam != "" {
		s := models.RestaurantStatus(statusParam)
		status = &s
	}

	// Get KAM ID filter
	kamIDParam := c.Query("kam_id")
	if kamIDParam != "" {
		if id, err := strconv.ParseUint(kamIDParam, 10, 32); err == nil {
			uid := uint(id)
			kamID = &uid
		}
	}

	restaurants, err := h.restaurantRepo.ListWithContext(c.Request.Context(), status, kamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restaurants)
}

// GetRestaurant handles getting a restaurant by ID
// @Summary Get Restaurant
// @Description Get a restaurant by ID
// @Tags restaurants
// @Produce json
// @Param id path int true "Restaurant ID"
// @Success 200 {object} models.Restaurant
// @Failure 404 {object} map[string]string
// @Router /api/v1/restaurants/{id} [get]
func (h *RestaurantHandler) GetRestaurant(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	restaurant, err := h.restaurantRepo.GetByIDWithContext(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "restaurant not found"})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

// ActivateRestaurant handles restaurant activation (KAM/Admin only)
// @Summary Activate Restaurant
// @Description Activate a pending restaurant. The KAM from the token will be set as activated_by and kam_id (if not already assigned)
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path int true "Restaurant ID"
// @Success 200 {object} models.Restaurant
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/v1/restaurants/{id}/activate [post]
func (h *RestaurantHandler) ActivateRestaurant(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	// Get user who is activating (from context set by auth middleware)
	// This user must be a KAM (enforced by middleware)
	activatedBy, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user context not found"})
		return
	}

	// Activate restaurant - no request body needed, KAM ID comes from token
	restaurant, err := h.restaurantService.ActivateRestaurant(c.Request.Context(), uint(id), activatedBy.(uint))
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "restaurant not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "only KAM users can activate restaurants" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "restaurant is already active" {
			statusCode = http.StatusConflict // 409 Conflict - already active
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Restaurant activated successfully",
		"restaurant": restaurant,
	})
}

type UpdateRestaurantStatusRequest struct {
	Status models.RestaurantStatus `form:"status" binding:"required,oneof=pending active inactive suspended"`
}

// UpdateRestaurantStatus handles updating restaurant status (KAM/Admin only)
// @Summary Update Restaurant Status
// @Description Update the status of a restaurant
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path int true "Restaurant ID"
// @Param status body map[string]string true "Status update" SchemaExample({"status": "active"})
// @Success 200 {object} models.Restaurant
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/restaurants/{id}/status [put]
func (h *RestaurantHandler) UpdateRestaurantStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	var req UpdateRestaurantStatusRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid the status value. It must be one of: pending, active, inactive, suspended."})
		return
	}

	restaurant, err := h.restaurantService.UpdateRestaurantStatus(c.Request.Context(), uint(id), req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

// ListPendingRestaurants handles listing pending restaurants (KAM/Admin only)
// @Summary List Pending Restaurants
// @Description List all restaurants awaiting activation
// @Tags restaurants
// @Produce json
// @Success 200 {array} models.Restaurant
// @Router /api/v1/restaurants/pending [get]
func (h *RestaurantHandler) ListPendingRestaurants(c *gin.Context) {
	restaurants, err := h.restaurantRepo.ListPendingWithContext(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restaurants)
}

// AssignKAM handles assigning a KAM to a restaurant (KAM/Admin only)
// @Summary Assign KAM
// @Description Assign a Key Account Manager to a restaurant
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path int true "Restaurant ID"
// @Param request body map[string]uint true "KAM assignment" SchemaExample({"kam_id": 1})
// @Success 200 {object} models.Restaurant
// @Failure 400 {object} map[string]string
// @Router /api/v1/restaurants/{id}/assign-kam [put]
func (h *RestaurantHandler) AssignKAM(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	var req map[string]uint
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kamID, exists := req["kam_id"]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kam_id is required"})
		return
	}

	restaurant, err := h.restaurantService.AssignKAM(c.Request.Context(), uint(id), kamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}
