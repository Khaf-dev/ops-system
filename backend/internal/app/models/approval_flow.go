package models

import (
	"time"

	"github.com/google/uuid"
)

// ApprovalFlow is the orchestration entity per ops_request
type ApprovalFlow struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RequestID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"request_id"`
	//Current step number (1..N)
	CurrentStep int `gorm:"default:0" json:"current_step"`
	//Status : Pending / in_review / approved / rejected / cancelled
	Status      string     `gorm:"size:32;default:'pending'" json:"status"`
	CreatedByID *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	//relations
	Request *OpsRequest    `gorm:"foreignKey:RequestID" json:"request,omitempty"`
	Steps   []ApprovalStep `gorm:"foreignKey:FlowID;constraint:OnDelete:CASCADE" json:"step,omitempty"`
}
