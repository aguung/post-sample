package entity

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID  uint   `gorm:"not null" json:"user_id"`
	User    User   `json:"user"`
	Title   string `gorm:"not null" json:"title"`
	Content string `gorm:"not null" json:"content"`
}
