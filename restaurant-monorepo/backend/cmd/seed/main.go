package main

import (
	"fmt"
	"log"
	"math/rand"
	"restaurant-saas/config"
	"restaurant-saas/database"
	"restaurant-saas/models"
	"restaurant-saas/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	// Initialize Config
	config.LoadConfig()

	// Initialize Logger
	utils.InitLogger()
	defer utils.Logger.Sync()

	// Initialize Database
	database.ConnectDB()
	// Ensure migrations are run (including partitions)
	database.Migrate()

	seedData(database.DB)
}

func seedData(db *gorm.DB) {
	log.Println("Starting seeding...")

	// Create 20 Organizations
	for i := 1; i <= 20; i++ {
		org := models.Organization{
			Name: fmt.Sprintf("Organization %d", i),
		}
		if err := db.Create(&org).Error; err != nil {
			log.Printf("Failed to create org %d: %v", i, err)
			continue
		}

		// Create Admin User for Org
		admin := models.User{
			OrganizationID: &org.ID,
			Email:          fmt.Sprintf("admin%d@org%d.com", i, i),
			PasswordHash:   "$2a$10$ExampleHashForTestingOnly", // In real app, hash properly
			Role:           models.RoleOrgAdmin,
		}
		db.Create(&admin)

		// Create 15 Restaurants for each Org
		for j := 1; j <= 15; j++ {
			restaurant := models.Restaurant{
				OrganizationID: org.ID,
				Name:           fmt.Sprintf("Restaurant %d - Org %d", j, i),
				Address:        fmt.Sprintf("123 Street, City %d", i),
				ContactEmail:   fmt.Sprintf("contact@rest%d-org%d.com", j, i),
			}
			if err := db.Create(&restaurant).Error; err != nil {
				log.Printf("Failed to create restaurant %d for org %d: %v", j, i, err)
				continue
			}

			seedRestaurantData(db, org.ID, restaurant.ID)
		}
		log.Printf("Seeded Organization %d with 15 restaurants", i)
	}

	log.Println("Seeding completed successfully!")
}

func seedRestaurantData(db *gorm.DB, orgID, restaurantID uuid.UUID) {
	// Create Menu Categories
	categories := []string{"Starters", "Mains", "Desserts", "Drinks"}
	var catIDs []uuid.UUID
	var allItems []models.MenuItem

	for k, catName := range categories {
		cat := models.MenuCategory{
			OrganizationID: orgID,
			RestaurantID:   restaurantID,
			Name:           catName,
			DisplayOrder:   k,
		}
		db.Create(&cat)
		catIDs = append(catIDs, cat.ID)

		// Create Menu Items for each Category
		for m := 1; m <= 5; m++ {
			item := models.MenuItem{
				OrganizationID: orgID,
				RestaurantID:   restaurantID,
				CategoryID:     cat.ID,
				Name:           fmt.Sprintf("%s Item %d", catName, m),
				Description:    "Delicious food item description",
				Price:          float64(10 + m*2),
			}
			db.Create(&item)
			allItems = append(allItems, item)
		}
	}

	// Create Tables
	for t := 1; t <= 10; t++ {
		table := models.Table{
			OrganizationID: orgID,
			RestaurantID:   restaurantID,
			Name:           fmt.Sprintf("Table %d", t),
			Capacity:       rand.Intn(6) + 2,
		}
		db.Create(&table)
	}

	// Create some dummy bookings
	for b := 1; b <= 5; b++ {
		booking := models.Booking{
			OrganizationID: orgID,
			RestaurantID:   restaurantID,
			CustomerName:   fmt.Sprintf("Guest %d", b),
			CustomerEmail:  fmt.Sprintf("guest%d@example.com", b),
			StartTime:      time.Now().Add(time.Duration(b) * time.Hour),
			NumberOfGuests: rand.Intn(4) + 1,
			Status:         models.BookingConfirmed,
		}
		db.Create(&booking)
	}

	// Create some dummy orders
	// Need to have items to create orders
	if len(allItems) > 0 {
		for o := 1; o <= 5; o++ {
			order := models.Order{
				OrganizationID: orgID,
				RestaurantID:   restaurantID,
				Status:         models.OrderCompleted,
				TotalAmount:    0,
			}
			db.Create(&order)

			var totalAmount float64
			// Add 1-3 items to order
			numItems := rand.Intn(3) + 1
			for k := 0; k < numItems; k++ {
				// Pick random item
				randomItem := allItems[rand.Intn(len(allItems))]
				qty := rand.Intn(2) + 1

				orderItem := models.OrderItem{
					OrganizationID: orgID,
					OrderID:        order.ID,
					MenuItemID:     randomItem.ID,
					Quantity:       qty,
					UnitPrice:      randomItem.Price,
				}
				db.Create(&orderItem)
				totalAmount += float64(qty) * randomItem.Price
			}

			// Update order total
			db.Model(&order).Update("total_amount", totalAmount)
		}
	}
}
