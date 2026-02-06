package entity

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	Name   string `json:"name"`
	Bio    string `json:"bio"`
}
