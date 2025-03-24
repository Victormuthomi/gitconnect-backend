package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitconnect-backend/config"
	"gitconnect-backend/models"
	"gitconnect-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Register a new user
// @Description Creates a new user account with a hashed password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "User Data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/register [post]
func Register(c *gin.Context) {
	var input models.User

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	// Save user to DB
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Automatically create a profile linked to the new user
	profile := models.Profile{
		UserID: user.ID,
	}

	if err := config.DB.Create(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
		return
	}

	// Return user and profile info
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
		"profile": gin.H{
			"id":     profile.ID,
			"user_id": profile.UserID,
		},
	})
}


// @Summary Login user
// @Description Authenticates a user and returns a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body models.User true "User Credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	var input models.User
	var user models.User

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	result := config.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// ✅ Return token + user ID + username
	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"token":    token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

