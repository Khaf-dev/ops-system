package models

import (
	"time"

	"github.com/google/uuid"
)

// ApprovalStep is a single step inside approval flow
type ApprovalStep struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FlowID     uuid.UUID `gorm:"type:uuid;not null;index" json:"flow_id"`
	StepNumber int       `gorm:"not null;index" json:"step_number"`
	// if UserID is present: single approver. If GroupName present: group step
	UserID    *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	GroupName string     `gorm:"size:128" json:"group_name,omitempty"`
	//Mode : "AND" / "OR" (inherited from config)
	Mode       string     `gorm:"size:10;default:'AND'" json:"mode,omitempty"`
	Status     string     `gorm:"size:32;default:'pending'" json:"status"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	Notes      string     `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`

	//relations
	Flow *ApprovalFlow `gorm:"foreignKey:FlowID" json:"flow,omitempty"`
	User *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
