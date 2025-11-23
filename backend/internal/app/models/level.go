package models

import (
	"time"

	"github.com/google/uuid"
)

type Level struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Rank      int       `gorm:"not null;uniqueIndex" json:"rank"`
	CreatedAt time.Time `gorm:"autoCreatedTime" json:"created_at"`
}
