package models

import (
	"time"

	"github.com/google/uuid"
)

type Attachment struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RequestID  uuid.UUID `gorm:"type:uuid;not null;index" json:"request_id"`
	FileURL    string    `json:"file_url"`
	FileType   string    `json:"file_type"`
	FileName   string    `json:"file_name,omitempty"`
	FileExt    string    `json:"file_ext,omitempty"`
	MimeType   string    `json:"mime_type,omitempty"`
	FileSize   int64     `json:"file_size,omitempty"`
	Checksum   string    `json:"checksum,omitempty"`
	UploadedAt time.Time `gorm:"autoCreateTime" json:"uploaded_at"`
}
