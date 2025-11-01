package repositories

import (
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

// CreateBatch creates multiple order items in a transaction
func (r *OrderItemRepository) CreateBatch(orderItems []models.OrderItem) error {
	return r.db.Create(&orderItems).Error
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

