package handlers

import (
	"net/http"
	"strconv"

	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/repositories"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// OrderHandler handles order-related requests
type OrderHandler struct {
	orderService *services.OrderService
	orderRepo    *repositories.OrderRepository
}

// NewOrderHandler creates a new OrderHandler instance
func NewOrderHandler(
	orderService *services.OrderService,
	orderRepo *repositories.OrderRepository,
) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		orderRepo:    orderRepo,
	}
}

// CreateOrder handles order creation
// @Summary Create Order
// @Description Create a new order with items
// @Tags orders
// @Accept json
// @Produce json
// @Param request body services.CreateOrderRequest true "Order data"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string
// @Router /api/v1/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req services.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), &req, restaurantID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrder handles getting an order by ID
// @Summary Get Order
// @Description Get an order by ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} models.Order
// @Failure 404 {object} map[string]string
// @Router /api/v1/orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := h.orderRepo.GetByIDWithContext(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ListOrders handles listing orders
// @Summary List Orders
// @Description List orders for the restaurant
// @Tags orders
// @Produce json
// @Param user_id query int false "Filter by user ID"
// @Success 200 {array} models.Order
// @Router /api/v1/orders [get]
func (h *OrderHandler) ListOrders(c *gin.Context) {
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	// Check if user_id query parameter is provided
	userIDParam := c.Query("user_id")
	if userIDParam != "" {
		userID, err := strconv.ParseUint(userIDParam, 10, 32)
		if err == nil {
			orders, err := h.orderRepo.GetByUserIDWithContext(c.Request.Context(), restaurantID.(uint), uint(userID))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, orders)
			return
		}
	}

	// Otherwise, get all orders for the restaurant
	orders, err := h.orderRepo.GetByRestaurantIDWithContext(c.Request.Context(), restaurantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// UpdateOrderStatus handles updating order status
// @Summary Update Order Status
// @Description Update the status of an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body services.UpdateOrderStatusRequest true "Status update data"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/orders/{id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var req services.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.UpdateOrderStatusWithCtx(c.Request.Context(), uint(id), &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
