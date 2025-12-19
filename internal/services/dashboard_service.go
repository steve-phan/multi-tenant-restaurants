package services

import (
	"context"
	"fmt"
	"time"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repositories"
)

// DashboardService handles dashboard statistics operations
type DashboardService struct {
	orderRepo       *repositories.OrderRepository
	reservationRepo *repositories.ReservationRepository
}

// NewDashboardService creates a new DashboardService instance
func NewDashboardService(orderRepo *repositories.OrderRepository, reservationRepo *repositories.ReservationRepository) *DashboardService {
	return &DashboardService{
		orderRepo:       orderRepo,
		reservationRepo: reservationRepo,
	}
}

// DashboardStats represents the overall dashboard statistics
type DashboardStats struct {
	OrderStats       *repositories.OrderStats       `json:"order_stats"`
	ReservationStats *repositories.ReservationStats `json:"reservation_stats"`
	OrdersByStatus   []repositories.OrderStatusCount `json:"orders_by_status"`
}

// GetDashboardStats retrieves overall dashboard statistics for a restaurant
func (s *DashboardService) GetDashboardStats(ctx context.Context, restaurantID uint, period string) (*DashboardStats, error) {
	// Calculate date range based on period
	startDate, endDate := s.calculateDateRange(period)

	// Get order stats
	orderStats, err := s.orderRepo.GetOrderStats(ctx, restaurantID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get order stats: %w", err)
	}

	// Get reservation stats
	reservationStats, err := s.reservationRepo.GetReservationStats(ctx, restaurantID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation stats: %w", err)
	}

	// Get orders by status
	ordersByStatus, err := s.orderRepo.GetOrdersByStatus(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by status: %w", err)
	}

	return &DashboardStats{
		OrderStats:       orderStats,
		ReservationStats: reservationStats,
		OrdersByStatus:   ordersByStatus,
	}, nil
}

// GetRecentOrders retrieves the most recent orders for a restaurant
func (s *DashboardService) GetRecentOrders(ctx context.Context, restaurantID uint, limit int) ([]models.Order, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	orders, err := s.orderRepo.GetRecentOrders(ctx, restaurantID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent orders: %w", err)
	}

	return orders, nil
}

// AnalyticsData represents analytics data for a specific period
type AnalyticsData struct {
	Period           string                          `json:"period"`
	StartDate        string                          `json:"start_date"`
	EndDate          string                          `json:"end_date"`
	OrderStats       *repositories.OrderStats        `json:"order_stats"`
	ReservationStats *repositories.ReservationStats  `json:"reservation_stats"`
}

// GetAnalytics retrieves analytics data for a specific period
func (s *DashboardService) GetAnalytics(ctx context.Context, restaurantID uint, period string) (*AnalyticsData, error) {
	// Calculate date range
	startDate, endDate := s.calculateDateRange(period)

	// Get order stats
	orderStats, err := s.orderRepo.GetOrderStats(ctx, restaurantID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get order stats: %w", err)
	}

	// Get reservation stats
	reservationStats, err := s.reservationRepo.GetReservationStats(ctx, restaurantID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation stats: %w", err)
	}

	return &AnalyticsData{
		Period:           period,
		StartDate:        startDate,
		EndDate:          endDate,
		OrderStats:       orderStats,
		ReservationStats: reservationStats,
	}, nil
}

// calculateDateRange calculates the start and end date based on the period
func (s *DashboardService) calculateDateRange(period string) (string, string) {
	now := time.Now()
	var startDate time.Time

	switch period {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "week":
		// Start from beginning of current week (Sunday)
		weekday := int(now.Weekday())
		startDate = now.AddDate(0, 0, -weekday)
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	case "month":
		// Start from beginning of current month
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case "year":
		// Start from beginning of current year
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	default:
		// Default to current month
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	return startDate.Format(time.RFC3339), endDate.Format(time.RFC3339)
}
