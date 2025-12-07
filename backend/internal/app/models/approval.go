package models

import (
	"backend/internal/app/constants"
	"time"

	"github.com/google/uuid"
)

type Approval struct {
	ID         uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RequestID  uuid.UUID          `gorm:"type:uuid;not null;index" json:"request_id"`
	ApproverID uuid.UUID          `gorm:"type:uuid;not null;index" json:"approver_id"`
	Decision   constants.Decision `gorm:"type:varchar(20);not null;default:'pending'" json:"decision"` // "approved" atau "rejected"
	Notes      string             `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt  time.Time          `gorm:"autoCreateTime" json:"created_at"`

	// relations (optional)
	Request  *OpsRequest `gorm:"foreignKey:RequestID;constraint:OnDelete:CASCADE" json:"request,omitempty"`
	Approver *User       `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
}
