package main

import (
	"log"
	"os"

	_ "gitconnect-backend/docs" // Import Swagger docs
	"gitconnect-backend/config"
	"gitconnect-backend/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title GitConnect API
// @version 1.0
// @description This is the API documentation for GitConnect.
// @termsOfService http://swagger.io/terms/
// @contact.name Victor Muthomi
// @contact.email victor@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api

func main() {
	// Set Gin mode, default to "debug" if not set
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = "debug"
	}
	gin.SetMode(mode)

	// Initialize database
	if err := config.ConnectDatabase(); err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	log.Println("✅ Database connected successfully.")

	// Initialize Gin router
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Fix proxy warning
	router.SetTrustedProxies(nil)

	// Register routes
	routes.AuthRoutes(router)
	routes.PostRoutes(router)
	routes.ProfileRoutes(router)

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Get port from env, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("✅ Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

