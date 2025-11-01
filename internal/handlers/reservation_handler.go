package handlers

import (
	"net/http"
	"strconv"
	"time"

	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/repositories"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// ReservationHandler handles reservation-related requests
type ReservationHandler struct {
	reservationService *services.ReservationService
	reservationRepo    *repositories.ReservationRepository
}

// NewReservationHandler creates a new ReservationHandler instance
func NewReservationHandler(
	reservationService *services.ReservationService,
	reservationRepo *repositories.ReservationRepository,
) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
		reservationRepo:     reservationRepo,
	}
}

// CreateReservation handles reservation creation
// @Summary Create Reservation
// @Description Create a new table reservation with availability checking
// @Tags reservations
// @Accept json
// @Produce json
// @Param request body services.CreateReservationRequest true "Reservation data"
// @Success 201 {object} models.Reservation
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/reservations [post]
func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	var req services.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	reservation, err := h.reservationService.CreateReservation(&req, restaurantID.(uint))
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "table is not available at the requested time" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reservation)
}

// GetReservation handles getting a reservation by ID
// @Summary Get Reservation
// @Description Get a reservation by ID
// @Tags reservations
// @Produce json
// @Param id path int true "Reservation ID"
// @Success 200 {object} models.Reservation
// @Failure 404 {object} map[string]string
// @Router /api/v1/reservations/{id} [get]
func (h *ReservationHandler) GetReservation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation ID"})
		return
	}

	reservation, err := h.reservationRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reservation not found"})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// ListReservations handles listing reservations
// @Summary List Reservations
// @Description List reservations, optionally filtered by date
// @Tags reservations
// @Produce json
// @Param date query string false "Date filter (YYYY-MM-DD)"
// @Success 200 {array} models.Reservation
// @Router /api/v1/reservations [get]
func (h *ReservationHandler) ListReservations(c *gin.Context) {
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "restaurant_id not found in context"})
		return
	}

	// Check if date query parameter is provided
	dateParam := c.Query("date")
	if dateParam != "" {
		date, err := time.Parse("2006-01-02", dateParam)
		if err == nil {
			reservations, err := h.reservationRepo.GetByDate(restaurantID.(uint), date)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, reservations)
			return
		}
	}

	// Otherwise, get all reservations for the restaurant
	reservations, err := h.reservationRepo.GetByRestaurantID(restaurantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

// UpdateReservation handles updating a reservation
// @Summary Update Reservation
// @Description Update an existing reservation (currently supports status updates)
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path int true "Reservation ID"
// @Param reservation body services.UpdateReservationStatusRequest true "Reservation update data"
// @Success 200 {object} models.Reservation
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/reservations/{id} [put]
func (h *ReservationHandler) UpdateReservation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation ID"})
		return
	}

	var req services.UpdateReservationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation, err := h.reservationService.UpdateReservationStatus(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// DeleteReservation handles deleting a reservation
// @Summary Delete Reservation
// @Description Cancel a reservation (soft delete)
// @Tags reservations
// @Param id path int true "Reservation ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /api/v1/reservations/{id} [delete]
func (h *ReservationHandler) DeleteReservation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation ID"})
		return
	}

	if err := h.reservationRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

