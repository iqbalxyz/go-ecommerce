package models

import (
	"time"
)

type User struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string    `json:"-" gorm:"size:255;not null"`
	FirstName    string    `json:"first_name" gorm:"size:100;not null"`
	LastName     string    `json:"last_name" gorm:"size:100"`
	Role         string    `json:"role" gorm:"size:20;not null;default:customer"`
	IsActive     bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
