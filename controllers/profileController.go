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

// UploadProfilePicture handles profile picture uploads
func UploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	// Maximum file size for the profile picture (e.g., 5MB)
	const MaxFileSize = 5 << 20 // 5MB

	// Parse the multipart form (this will extract the file from the request)
	err := r.ParseMultipartForm(MaxFileSize)
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("profile_picture")
	if err != nil {
		http.Error(w, "Error getting file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Define the allowed file types (e.g., JPEG, PNG)
	allowedFileTypes := []string{"image/jpeg", "image/png"}
	// Get the file's MIME type
	contentType := r.Header.Get("Content-Type")

	// Validate the file type
	if !contains(allowedFileTypes, contentType) {
		http.Error(w, "Invalid file type. Only JPEG and PNG are allowed.", http.StatusBadRequest)
		return
	}

	// Generate a unique filename to avoid overwriting
	// You can use a UUID, timestamp, or user ID here
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), "profile.jpg") // Simple naming, change as needed

	// Define the path where you want to store the file
	filePath := filepath.Join("uploads", fileName)

	// Create the file on the server
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Copy the content of the uploaded file to the new file
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Error saving file content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message and the file path
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Profile picture uploaded successfully",
		"file_path": filePath, // You can also store this URL in the database to associate it with the user's profile
	}
	json.NewEncoder(w).Encode(response)
}

// Utility function to check if a string exists in a list
func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

