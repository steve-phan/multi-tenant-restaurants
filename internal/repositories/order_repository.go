package repositories

import (
	"context"
	"restaurant-backend/internal/models"

	"gorm.io/gorm"
)

// OrderRepository handles order-related database operations
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new OrderRepository instance
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create creates a new order
func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

// CreateWithContext creates a new order using the provided context
func (r *OrderRepository) CreateWithContext(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// GetByID retrieves an order by ID (RLS ensures tenant isolation)
func (r *OrderRepository) GetByID(id uint) (*models.Order, error) {
	var order models.Order
	if err := r.db.Preload("OrderItems").Preload("OrderItems.MenuItem").Preload("User").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByIDWithContext retrieves an order by ID using the provided context
func (r *OrderRepository) GetByIDWithContext(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).Preload("OrderItems").Preload("OrderItems.MenuItem").Preload("User").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByRestaurantID retrieves all orders for a restaurant (RLS ensures tenant isolation)
func (r *OrderRepository) GetByRestaurantID(restaurantID uint) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Where("restaurant_id = ?", restaurantID).
		Preload("OrderItems").
		Preload("OrderItems.MenuItem").
		Preload("User").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// GetByRestaurantIDWithContext retrieves orders for a restaurant using the provided context
func (r *OrderRepository) GetByRestaurantIDWithContext(ctx context.Context, restaurantID uint) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.WithContext(ctx).Where("restaurant_id = ?", restaurantID).
		Preload("OrderItems").
		Preload("OrderItems.MenuItem").
		Preload("User").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// GetByUserID retrieves all orders for a user (RLS ensures tenant isolation)
func (r *OrderRepository) GetByUserID(restaurantID uint, userID uint) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Where("restaurant_id = ? AND user_id = ?", restaurantID, userID).
		Preload("OrderItems").
		Preload("OrderItems.MenuItem").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// GetByUserIDWithContext retrieves orders for a user using the provided context
func (r *OrderRepository) GetByUserIDWithContext(ctx context.Context, restaurantID uint, userID uint) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.WithContext(ctx).Where("restaurant_id = ? AND user_id = ?", restaurantID, userID).
		Preload("OrderItems").
		Preload("OrderItems.MenuItem").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// Update updates an existing order
func (r *OrderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

// UpdateWithContext updates an order using the provided context
func (r *OrderRepository) UpdateWithContext(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// UpdateStatus updates only the status of an order
func (r *OrderRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateStatusWithContext updates the status of an order using the provided context
func (r *OrderRepository) UpdateStatusWithContext(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}
