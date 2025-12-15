package services

import (
	"context"
	"errors"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repositories"
)

// OrderService handles order business logic
type OrderService struct {
	orderRepo     *repositories.OrderRepository
	orderItemRepo *repositories.OrderItemRepository
	menuItemRepo  *repositories.MenuItemRepository
}

// NewOrderService creates a new OrderService instance
func NewOrderService(
	orderRepo *repositories.OrderRepository,
	orderItemRepo *repositories.OrderItemRepository,
	menuItemRepo *repositories.MenuItemRepository,
) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		menuItemRepo:  menuItemRepo,
	}
}

// OrderItemRequest represents an item in an order request
type OrderItemRequest struct {
	MenuItemID uint   `json:"menu_item_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	Notes      string `json:"notes"`
}

// CreateOrderRequest represents order creation request
type CreateOrderRequest struct {
	UserID uint               `json:"user_id" binding:"required"`
	Items  []OrderItemRequest `json:"items" binding:"required,min=1"`
	Notes  string             `json:"notes"`
}

// CreateOrder creates a new order with items
func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest, restaurantID uint) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	// Validate menu items and calculate total
	var totalAmount float64
	orderItems := make([]models.OrderItem, 0, len(req.Items))

	for _, itemReq := range req.Items {
		// Get menu item to validate and get price
		menuItem, err := s.menuItemRepo.GetByIDWithContext(ctx, itemReq.MenuItemID)
		if err != nil {
			return nil, errors.New("menu item not found")
		}

		// Validate menu item belongs to restaurant (RLS ensures this)
		if menuItem.RestaurantID != restaurantID {
			return nil, errors.New("menu item does not belong to restaurant")
		}

		// Check availability
		if !menuItem.IsAvailable {
			return nil, errors.New("menu item is not available")
		}

		// Calculate item total
		itemTotal := menuItem.Price * float64(itemReq.Quantity)
		totalAmount += itemTotal

		// Create order item
		orderItem := models.OrderItem{
			MenuItemID: itemReq.MenuItemID,
			Quantity:   itemReq.Quantity,
			Price:      menuItem.Price,
			Notes:      itemReq.Notes,
		}
		orderItems = append(orderItems, orderItem)
	}

	// Create order
	order := &models.Order{
		RestaurantID: restaurantID,
		UserID:       req.UserID,
		Status:       "pending",
		TotalAmount:  totalAmount,
		Notes:        req.Notes,
		OrderItems:   orderItems,
	}

	// Set restaurant ID for all order items
	for i := range order.OrderItems {
		order.OrderItems[i].RestaurantID = restaurantID
	}

	if err := s.orderRepo.CreateWithContext(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

// UpdateOrderStatusRequest represents order status update request
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed preparing ready completed cancelled"`
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(orderID uint, req *UpdateOrderStatusRequest) (*models.Order, error) {
	order, err := s.orderRepo.GetByIDWithContext(context.Background(), orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	order.Status = req.Status

	if err := s.orderRepo.UpdateWithContext(context.Background(), order); err != nil {
		return nil, err
	}

	return order, nil
}

// UpdateOrderStatusWithCtx updates order status using provided context
func (s *OrderService) UpdateOrderStatusWithCtx(ctx context.Context, orderID uint, req *UpdateOrderStatusRequest) (*models.Order, error) {
	order, err := s.orderRepo.GetByIDWithContext(ctx, orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	order.Status = req.Status

	if err := s.orderRepo.UpdateWithContext(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}
