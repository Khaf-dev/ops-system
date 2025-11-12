package models

import "github.com/google/uuid"

type Activity struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name     string    `json:"name"`
	IsActive bool      `json:"is_active"`
}
