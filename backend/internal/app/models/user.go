package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Email        string    `gorm:"size:120;uniqueIndex;not null" json:"email"`
	Phone        string    `gorm:"size:20" json:"phone"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"size:20;default:user" json:"role"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
