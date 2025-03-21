package main

import (
	"log"
	"os"

	_ "gitconnect-backend/docs" // Import Swagger docs
	"gitconnect-backend/config"
	"gitconnect-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title GitConnect API
// @version 1.0
// @description API documentation for GitConnect
// @termsOfService http://swagger.io/terms/
// @contact.name Victor Muthomi
// @contact.email victor@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api

func main() {
	// Set Gin mode
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

	// Enable CORS for frontend communication
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Change to frontend URL in production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

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

	// Bind to 0.0.0.0 to allow connections from Docker
	serverAddr := "0.0.0.0:" + port
	log.Printf("✅ Server running on %s", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

