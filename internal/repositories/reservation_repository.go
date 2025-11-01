package repositories

import (
	"restaurant-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

// ReservationRepository handles reservation-related database operations
type ReservationRepository struct {
	db *gorm.DB
}

// NewReservationRepository creates a new ReservationRepository instance
func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

// Create creates a new reservation
func (r *ReservationRepository) Create(reservation *models.Reservation) error {
	return r.db.Create(reservation).Error
}

// GetByID retrieves a reservation by ID (RLS ensures tenant isolation)
func (r *ReservationRepository) GetByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := r.db.Preload("User").First(&reservation, id).Error; err != nil {
		return nil, err
	}
	return &reservation, nil
}

// GetByRestaurantID retrieves all reservations for a restaurant (RLS ensures tenant isolation)
func (r *ReservationRepository) GetByRestaurantID(restaurantID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := r.db.Where("restaurant_id = ?", restaurantID).
		Preload("User").
		Order("start_time ASC").
		Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

// GetByDate retrieves reservations for a specific date
func (r *ReservationRepository) GetByDate(restaurantID uint, date time.Time) ([]models.Reservation, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var reservations []models.Reservation
	if err := r.db.Where("restaurant_id = ? AND start_time >= ? AND start_time < ?", restaurantID, startOfDay, endOfDay).
		Preload("User").
		Order("start_time ASC").
		Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

// GetByTableAndTime retrieves reservations for a specific table and time range
func (r *ReservationRepository) GetByTableAndTime(restaurantID uint, tableNumber string, startTime, endTime time.Time) ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := r.db.Where(
		"restaurant_id = ? AND table_number = ? AND status != 'cancelled' AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?) OR (start_time >= ? AND start_time < ?))",
		restaurantID, tableNumber, startTime, startTime, endTime, endTime, startTime, endTime,
	).Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

// Update updates an existing reservation
func (r *ReservationRepository) Update(reservation *models.Reservation) error {
	return r.db.Save(reservation).Error
}

// Delete deletes a reservation (soft delete by setting status to cancelled)
func (r *ReservationRepository) Delete(id uint) error {
	return r.db.Model(&models.Reservation{}).Where("id = ?", id).Update("status", "cancelled").Error
}

