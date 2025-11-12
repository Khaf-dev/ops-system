package models

import (
	"time"

	"github.com/google/uuid"
)

type Attachment struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RequestID  uuid.UUID `gorm:"type:uuid;not null"`
	FileURL    string    `json:"file_url"`
	FileType   string    `json:"file_type"`
	UploadedAt time.Time `json:"uploaded_at"`
}
