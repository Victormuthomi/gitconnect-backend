package routes

import (
	"gitconnect-backend/controllers"
	"gitconnect-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func ProfileRoutes(router *gin.Engine) {
	// Public route: Get all profiles
	router.GET("/api/profiles", controllers.GetProfiles)

	// Protected routes
	protected := router.Group("/api/profiles").Use(middlewares.AuthMiddleware()) // Apply AuthMiddleware to this group
	{
		// Create a new profile
		protected.POST("/", controllers.CreateProfile)

		// Get a single profile by ID (protected)
		protected.GET("/:id", controllers.GetProfile)

		// Update a profile (protected)
		protected.PUT("/:id", controllers.UpdateProfile)
	}
}

