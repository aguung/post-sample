package entity

import (
	"time"

	"gorm.io/gorm"
)

type Role int

const (
	RoleAdmin Role = 1
	RoleUser  Role = 2
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Email    string `gorm:"uniqueIndex:idx_email_unique,where:deleted_at IS NULL;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Role     Role   `gorm:"not null;default:2" json:"role"`

	Profile Profile `json:"profile,omitempty"`
	Posts   []Post  `json:"posts,omitempty"`
}
