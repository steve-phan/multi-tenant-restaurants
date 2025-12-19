package handlers

import (
	"net/http"
	"strconv"

	"restaurant-backend/internal/ctx"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard statistics requests
type DashboardHandler struct {
	dashboardService *services.DashboardService
}

// NewDashboardHandler creates a new DashboardHandler instance
func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetDashboardStats handles retrieving overall dashboard statistics
// @Summary Get Dashboard Stats
// @Description Get overall dashboard statistics for the restaurant
// @Tags dashboard
// @Produce json
// @Param period query string false "Time period (today, week, month, year)" default(month)
// @Success 200 {object} services.DashboardStats
// @Failure 500 {object} map[string]string
// @Router /api/v1/dashboard/stats [get]
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	// Get restaurant ID from context
	restaurantID, ok := ctx.GetRestaurantID(c.Request.Context())
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	// Get period from query parameter (default to "month")
	period := c.DefaultQuery("period", "month")

	stats, err := h.dashboardService.GetDashboardStats(c.Request.Context(), restaurantID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetRecentOrders handles retrieving recent orders
// @Summary Get Recent Orders
// @Description Get the most recent orders for the restaurant
// @Tags dashboard
// @Produce json
// @Param limit query int false "Number of orders to retrieve (max 100)" default(10)
// @Success 200 {array} models.Order
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/dashboard/recent-orders [get]
func (h *DashboardHandler) GetRecentOrders(c *gin.Context) {
	// Get restaurant ID from context
	restaurantID, ok := ctx.GetRestaurantID(c.Request.Context())
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	// Get limit from query parameter (default to 10)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	orders, err := h.dashboardService.GetRecentOrders(c.Request.Context(), restaurantID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetAnalytics handles retrieving analytics data
// @Summary Get Analytics
// @Description Get analytics data for a specific period
// @Tags dashboard
// @Produce json
// @Param period query string false "Time period (today, week, month, year)" default(month)
// @Success 200 {object} services.AnalyticsData
// @Failure 500 {object} map[string]string
// @Router /api/v1/dashboard/analytics [get]
func (h *DashboardHandler) GetAnalytics(c *gin.Context) {
	// Get restaurant ID from context
	restaurantID, ok := ctx.GetRestaurantID(c.Request.Context())
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	// Get period from query parameter (default to "month")
	period := c.DefaultQuery("period", "month")

	analytics, err := h.dashboardService.GetAnalytics(c.Request.Context(), restaurantID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}
