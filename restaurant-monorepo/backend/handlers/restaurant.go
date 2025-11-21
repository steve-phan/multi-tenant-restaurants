package handlers

import (
	"net/http"
	"restaurant-saas/database"
	"restaurant-saas/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateRestaurant(c *gin.Context) {
	// Simplified for MVP: Directly create restaurant
	var input models.Restaurant
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure Org ID is set (in real app, from token or context)
	if input.OrganizationID == uuid.Nil {
		// For MVP testing, allow passing it or generate one if needed
		// But ideally should come from authenticated user's org
	}
	/// restaurant Name : acb , abd, .... xyz, zypz

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func GetRestaurant(c *gin.Context) {
	id := c.Param("id")
	var restaurant models.Restaurant
	if err := database.DB.First(&restaurant, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}
	c.JSON(http.StatusOK, restaurant)
}

func CreateMenuCategory(c *gin.Context) {
	restaurantID := c.Param("id")
	var input models.MenuCategory
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

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func CreateMenuItem(c *gin.Context) {
	restaurantID := c.Param("id")
	var input models.MenuItem
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

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func CreateTable(c *gin.Context) {
	restaurantID := c.Param("id")
	var input models.Table
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

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func GetBookings(c *gin.Context) {
	restaurantID := c.Param("id")
	var bookings []models.Booking
	if err := database.DB.Where("restaurant_id = ?", restaurantID).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func GetOrders(c *gin.Context) {
	restaurantID := c.Param("id")
	var orders []models.Order
	if err := database.DB.Where("restaurant_id = ?", restaurantID).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}
