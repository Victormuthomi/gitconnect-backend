package models

import "time"

// User represents a registered user
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-"` // Exclude password from JSON response
	Profile   *Profile  `json:"profile,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // Use pointer to avoid recursion
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

