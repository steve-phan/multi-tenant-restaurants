package handlers

import (
	"net/http"
	"restaurant-saas/database"
	"restaurant-saas/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetPublicMenu(c *gin.Context) {
	restaurantID := c.Param("id")

	var categories []models.MenuCategory
	if err := database.DB.Where("restaurant_id = ?", restaurantID).Order("display_order").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	var items []models.MenuItem
	if err := database.DB.Where("restaurant_id = ?", restaurantID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"items":      items,
	})
}

func GetAvailableTables(c *gin.Context) {
	restaurantID := c.Param("id")
	// Simplified availability check logic for MVP
	// In real app, check against existing bookings
	var tables []models.Table
	if err := database.DB.Where("restaurant_id = ?", restaurantID).Find(&tables).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tables"})
		return
	}
	c.JSON(http.StatusOK, tables)
}

func CreateBooking(c *gin.Context) {
	restaurantID := c.Param("id")
	var input models.Booking
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rID, _ := uuid.Parse(restaurantID)

	var restaurant models.Restaurant
	if err := database.DB.Select("organization_id").First(&restaurant, "id = ?", rID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	input.RestaurantID = rID
	input.OrganizationID = restaurant.OrganizationID
	input.Status = models.BookingPending
	if input.StartTime.IsZero() {
		input.StartTime = time.Now() // Default to now if not set (should be set by frontend)
	}

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func CreateOrder(c *gin.Context) {
	restaurantID := c.Param("id")
	var input models.Order
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rID, _ := uuid.Parse(restaurantID)

	var restaurant models.Restaurant
	if err := database.DB.Select("organization_id").First(&restaurant, "id = ?", rID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	input.RestaurantID = rID
	input.OrganizationID = restaurant.OrganizationID

	// Set OrganizationID for all order items
	for i := range input.Items {
		input.Items[i].OrganizationID = restaurant.OrganizationID
	}
	input.Status = models.OrderPending

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}
