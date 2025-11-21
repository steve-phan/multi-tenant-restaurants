package main

import (
	"restaurant-saas/config"
	"restaurant-saas/database"
	"restaurant-saas/middleware"
	"restaurant-saas/routes"
	"restaurant-saas/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize Config
	config.LoadConfig()

	// Initialize Logger
	utils.InitLogger()
	defer utils.Logger.Sync()

	// Initialize Database
	database.ConnectDB()
	database.Migrate()

	// Initialize Gin
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.ZapLogger())

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.AppConfig.ClientOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Setup Routes
	routes.SetupRoutes(r)

	// Start Server
	port := config.AppConfig.Port
	utils.Logger.Info("Server starting", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		utils.Logger.Fatal("Server failed to start", zap.Error(err))
	}
}
