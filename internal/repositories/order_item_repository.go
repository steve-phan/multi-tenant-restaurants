package repositories

import (
	"context"
	"restaurant-backend/internal/models"

	"gorm.io/gorm"
)

// OrderItemRepository handles order item-related database operations
type OrderItemRepository struct {
	db *gorm.DB
}

// NewOrderItemRepository creates a new OrderItemRepository instance
func NewOrderItemRepository(db *gorm.DB) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

// Create creates a new order item
func (r *OrderItemRepository) Create(orderItem *models.OrderItem) error {
	return r.db.Create(orderItem).Error
}

// CreateWithContext creates a new order item using the provided context
func (r *OrderItemRepository) CreateWithContext(ctx context.Context, orderItem *models.OrderItem) error {
	return r.db.WithContext(ctx).Create(orderItem).Error
}

// CreateBatch creates multiple order items in a transaction
func (r *OrderItemRepository) CreateBatch(orderItems []models.OrderItem) error {
	return r.db.Create(&orderItems).Error
}

// CreateBatchWithContext creates multiple order items using the provided context
func (r *OrderItemRepository) CreateBatchWithContext(ctx context.Context, orderItems []models.OrderItem) error {
	return r.db.WithContext(ctx).Create(&orderItems).Error
}

// GetByOrderID retrieves all order items for an order (RLS ensures tenant isolation)
func (r *OrderItemRepository) GetByOrderID(orderID uint) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	if err := r.db.Where("order_id = ?", orderID).
		Preload("MenuItem").
		Find(&orderItems).Error; err != nil {
		return nil, err
	}
	return orderItems, nil
}

// GetByOrderIDWithContext retrieves order items for an order using the provided context
func (r *OrderItemRepository) GetByOrderIDWithContext(ctx context.Context, orderID uint) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).
		Preload("MenuItem").
		Find(&orderItems).Error; err != nil {
		return nil, err
	}
	return orderItems, nil
}
