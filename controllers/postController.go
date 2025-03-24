package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitconnect-backend/config"
	"gitconnect-backend/models"
)

// @Summary Create a new post
// @Description Allows an authenticated user to create a new post
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body models.Post true "Post Data"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/posts [post]
func CreatePost(c *gin.Context) {
	// Get user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assign the authenticated user to the post
	post.UserID = userID.(uint)

	// Save post
	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post created", "post": post})
}

// @Summary Get all posts
// @Description Fetch all posts with user details
// @Tags Posts
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/posts [get]
func GetPosts(c *gin.Context) {
	var posts []models.Post

	// Include user details in the response
	if err := config.DB.Preload("User").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// @Summary Get a single post
// @Description Fetch a post by ID
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/posts/{id} [get]
func GetPost(c *gin.Context) {
	var post models.Post

	// Convert ID param to uint
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Find post
	if err := config.DB.Preload("User").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// @Summary Delete a post
// @Description Deletes a post (Only the author can delete)
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/posts/{id} [delete]
func DeletePost(c *gin.Context) {
	// Get user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var post models.Post

	// Convert ID param to uint
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Find post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if the logged-in user is the post owner
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own post"})
		return
	}

	// Delete post
	if err := config.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

// @Summary Update a post
// @Description Updates an existing post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body models.Post true "Updated Post Data"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/posts/{id} [put]
func UpdatePost(c *gin.Context) {
	var post models.Post

	// Convert ID param to uint
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Find post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Bind request data to the post object
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save updated post
	config.DB.Save(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Post updated", "post": post})
}

// @Summary Like a post
// @Description Allows a user to like a post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/posts/{id}/like [post]
func LikePost(c *gin.Context) {
	var post models.Post
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Find the post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Increment the like counter
	post.LikePost()

	// Save post with updated likes
	if err := config.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post liked", "likes": post.Likes})
}

// @Summary Dislike a post
// @Description Allows a user to dislike a post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/posts/{id}/dislike [post]
func DislikePost(c *gin.Context) {
	var post models.Post
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Find the post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Increment the dislike counter
	post.DislikePost()

	// Save post with updated dislikes
	if err := config.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to dislike post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post disliked", "dislikes": post.Dislikes})
}

// @Summary Comment on a post
// @Description Allows a user to comment on a post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param comment body models.Comment true "Comment Data"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/posts/{id}/comments [post]
func CommentOnPost(c *gin.Context) {
	// Get user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get post ID from params
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var comment models.Comment
	// Bind comment data to the model
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assign the post ID and user ID
	comment.PostID = uint(postID)
	comment.UserID = userID.(uint)

	// Save the comment
	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Comment added", "comment": comment})
}

// @Summary Get all comments for a post
// @Description Fetch all comments for a specific post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/posts/{id}/comments [get]
func GetCommentsForPost(c *gin.Context) {
    postID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    var comments []models.Comment

    // Fetch comments with Post and User details only if necessary
    query := config.DB.Where("post_id = ?", postID)

    // Remove Preload if it's causing issues
    query = query.Preload("User") // If User exists, keep this

    if err := query.Find(&comments).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"comments": comments})
}

