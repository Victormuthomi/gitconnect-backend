package models

import "time"

// Post represents a post in the system
type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Content   string    `json:"content" binding:"required"`
	UserID    uint      `json:"user_id" gorm:"not null;index"` // Foreign key for users
	User      User      `json:"user" gorm:"foreignKey:UserID"` // Establish relation
	Likes     int       `json:"likes" gorm:"default:0"`
	Dislikes  int       `json:"dislikes" gorm:"default:0"`
	Comments  []Comment `json:"comments" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;"` // Comments linked to post
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LikePost increments the like count for the post
func (p *Post) LikePost() {
	p.Likes++
}

// DislikePost increments the dislike count for the post
func (p *Post) DislikePost() {
	p.Dislikes++
}

