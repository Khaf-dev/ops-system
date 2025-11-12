package models

import (
	"time"

	"github.com/google/uuid"
)

type Approval struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RequestID  uuid.UUID `gorm:"type:uuid;not null" json:"request_id"`
	ApproverID uuid.UUID `gorm:"type:uuid;not null" json:"approver_id"`
	Decision   string    `gorm:"not null" json:"decision"` // "approved" atau "rejected"
	Notes      string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
