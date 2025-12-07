package models

import (
	"time"

	"github.com/google/uuid"
)

// ApprovalLog stores audit trail for approval events
type ApprovalLog struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FlowID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"flow_id"`
	StepID    *uuid.UUID `gorm:"type:uuid;index" json:"step_id,omitempty"`
	Action    string     `gorm:"size:64;not null" json:"action"` // e.g step_approved, step_rejected, flow_started
	ByUserID  *uuid.UUID `gorm:"type:uuid;index" json:"by_user_id,omitempty"`
	Note      string     `gorm:"type:text" json:"note,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`

	//relations
	Flow *ApprovalFlow `gorm:"foreignKey:FlowID;constraint:OnDelete:CASCADE" json:"flow,omitempty"`
	Step *ApprovalStep `gorm:"foreignKey:StepID" json:"step,omitempty"`
}
