package models

import "time"

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PostID    uint      `json:"post_id" gorm:"not null;index"` // Foreign key with index
	UserID    uint      `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID"` // Relation with User
	Content   string    `json:"content" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

