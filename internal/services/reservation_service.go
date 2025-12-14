package services

import (
	"errors"
	"time"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repositories"
)

// ReservationService handles reservation business logic
type ReservationService struct {
	reservationRepo *repositories.ReservationRepository
}

// NewReservationService creates a new ReservationService instance
func NewReservationService(reservationRepo *repositories.ReservationRepository) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
	}
}

// CreateReservationRequest represents reservation creation request
type CreateReservationRequest struct {
	UserID         uint      `json:"user_id" binding:"required"`
	TableNumber    string    `json:"table_number" binding:"required"`
	StartTime      time.Time `json:"start_time" binding:"required"`
	EndTime        time.Time `json:"end_time" binding:"required"`
	NumberOfGuests int       `json:"number_of_guests" binding:"required,min=1"`
	Notes          string    `json:"notes"`
}

// CreateReservation creates a new reservation with availability checking
func (s *ReservationService) CreateReservation(req *CreateReservationRequest, restaurantID uint) (*models.Reservation, error) {
	// Validate time range
	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("end time must be after start time")
	}

	if req.StartTime.Before(time.Now()) {
		return nil, errors.New("reservation cannot be in the past")
	}

	// Check table availability
	isAvailable, err := s.checkTableAvailability(restaurantID, req.TableNumber, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}

	if !isAvailable {
		return nil, errors.New("table is not available at the requested time")
	}

	// Create reservation
	reservation := &models.Reservation{
		RestaurantID:   restaurantID,
		UserID:         req.UserID,
		TableNumber:    req.TableNumber,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		NumberOfGuests: req.NumberOfGuests,
		Status:         "pending",
		Notes:          req.Notes,
	}

	if err := s.reservationRepo.Create(reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}

// UpdateReservationStatusRequest represents reservation status update request
type UpdateReservationStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed cancelled completed"`
}

// UpdateReservationStatus updates the status of a reservation
func (s *ReservationService) UpdateReservationStatus(reservationID uint, req *UpdateReservationStatusRequest) (*models.Reservation, error) {
	reservation, err := s.reservationRepo.GetByID(reservationID)
	if err != nil {
		return nil, errors.New("reservation not found")
	}

	reservation.Status = req.Status

	if err := s.reservationRepo.Update(reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}

// checkTableAvailability checks if a table is available at the given time range
func (s *ReservationService) checkTableAvailability(restaurantID uint, tableNumber string, startTime, endTime time.Time) (bool, error) {
	// Get existing reservations for this table in the time range
	conflictingReservations, err := s.reservationRepo.GetByTableAndTime(restaurantID, tableNumber, startTime, endTime)
	if err != nil {
		return false, err
	}

	// If there are any conflicting reservations, table is not available
	return len(conflictingReservations) == 0, nil
}
