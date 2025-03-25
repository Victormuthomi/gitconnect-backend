package models

import "time"

// Profile represents a user's profile
type Profile struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint      `json:"user_id" gorm:"not null;unique;index"`
	User      *User     `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // Prevent circular JSON recursion
	FullName  string    `json:"full_name" binding:"required"`
	Bio       string    `json:"bio"`
	Github    string    `json:"github"`
  ProfilePicture string   `json:"profile_picture"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

