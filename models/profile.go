package models

import (
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// Profile represents a user's profile.
type Profile struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         uint      `json:"user_id" gorm:"not null;unique;index"`
	FullName       string    `json:"full_name" binding:"required"`
	Bio            string    `json:"bio"`
	Github         string    `json:"github"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// DB is a global pointer to your database. Initialize it in your main or init function.
var DB *gorm.DB

// UpdateProfilePicture updates the profile picture filename for the given user.
func UpdateProfilePicture(userId string, filename string) error {
	// Convert userId (string) to uint.
	id64, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return err
	}
	id := uint(id64)

	var profile Profile
	// Find the profile using the UserID field.
	if err := DB.Where("user_id = ?", id).First(&profile).Error; err != nil {
		return err
	}

	// Update the profile picture field.
	profile.ProfilePicture = filename

	// Save the updated profile.
	if err := DB.Save(&profile).Error; err != nil {
		return err
	}

	return nil
}

// GetProfileByID retrieves the profile by userId.
func GetProfileByID(userId string) (*Profile, error) {
	// Convert userId (string) to uint.
	id64, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return nil, err
	}
	id := uint(id64)

	var profile Profile
	if err := DB.Where("user_id = ?", id).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &profile, nil
}

