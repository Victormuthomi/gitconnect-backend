package controllers

import (
	"fmt"
	"net/http"

	"gitconnect-backend/config"
	"gitconnect-backend/models"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new profile
// @Description Allows an authenticated user to create a new profile
// @Tags Profiles
// @Accept json
// @Produce json
// @Param profile body models.Profile true "Profile Data"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/profiles [post]
func CreateProfile(c *gin.Context) {
	var profile models.Profile
	fmt.Println("📥 Received request to create profile") // Debug log

	// Bind request JSON to profile struct
	if err := c.ShouldBindJSON(&profile); err != nil {
		fmt.Println("❌ Error binding JSON:", err) // Debug log
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the UserID exists in the Users table
	var user models.User
	if err := config.DB.First(&user, profile.UserID).Error; err != nil {
		fmt.Println("❌ Invalid UserID, user does not exist") // Debug log
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID: user does not exist"})
		return
	}

	// Save to database
	if err := config.DB.Create(&profile).Error; err != nil {
		fmt.Println("❌ Failed to create profile:", err) // Debug log
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
		return
	}

	fmt.Println("✅ Profile created successfully") // Debug log
	c.JSON(http.StatusCreated, gin.H{"message": "Profile created successfully", "profile": profile})
}

// @Summary Get all profiles
// @Description Fetch all profiles
// @Tags Profiles
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/profiles [get]
func GetProfiles(c *gin.Context) {
	var profiles []models.Profile
	config.DB.Find(&profiles)
	c.JSON(http.StatusOK, gin.H{"profiles": profiles})
}

// @Summary Get a specific profile
// @Description Fetch a profile by ID
// @Tags Profiles
// @Accept json
// @Produce json
// @Param id path int true "Profile ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/profiles/{id} [get]
func GetProfile(c *gin.Context) {
	var profile models.Profile
	if err := config.DB.First(&profile, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

// @Summary Update a profile
// @Description Update a profile by ID
// @Tags Profiles
// @Accept json
// @Produce json
// @Param id path int true "Profile ID"
// @Param profile body models.Profile true "Updated Profile Data"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/profiles/{id} [put]
func UpdateProfile(c *gin.Context) {
	var profile models.Profile
	if err := config.DB.First(&profile, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&profile)
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated", "profile": profile})
}

// @Summary Delete a profile
// @Description Delete a profile by ID
// @Tags Profiles
// @Accept json
// @Produce json
// @Param id path int true "Profile ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} 
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/profiles/{id} [delete]
func DeleteProfile(c *gin.Context) {
	var profile models.Profile
	// Find profile by ID
	if err := config.DB.First(&profile, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Delete the profile
	if err := config.DB.Delete(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}

