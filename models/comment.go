package models

import "time"

type Comment struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	PostID    uint   `json:"post_id"`
	UserID    uint   `json:"user_id"`
	Content   string `json:"content" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

