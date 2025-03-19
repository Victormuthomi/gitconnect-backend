package routes

import (
	"github.com/gin-gonic/gin"
	"gitconnect-backend/controllers"
	"gitconnect-backend/middlewares"
)

func PostRoutes(router *gin.Engine) {
	// Public route: Get all posts
	router.GET("/api/posts", controllers.GetPosts)

	// Protected routes
	protected := router.Group("/api/posts").Use(middlewares.AuthMiddleware()) // Updated to use the correct middleware
	{
		// Create a new post
		protected.POST("", controllers.CreatePost)

		// Update a post
		protected.PUT("/:id", controllers.UpdatePost)

		// Delete a post
		protected.DELETE("/:id", controllers.DeletePost)

		// Like a post
		protected.POST("/:id/like", controllers.LikePost)

		// Dislike a post
		protected.POST("/:id/dislike", controllers.DislikePost)

		// Comment on a post
		protected.POST("/:id/comments", controllers.CommentOnPost)
	}

	// Get a single post
	router.GET("/api/posts/:id", controllers.GetPost)

	// Get comments for a post
	router.GET("/api/posts/:id/comments", controllers.GetCommentsForPost)
}

