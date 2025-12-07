package models

import (
	"backend/internal/app/constants"
	"time"

	"github.com/google/uuid"
)

// ApproverConfig defines which users at what level approve which request
type ApproverConfig struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RequestTypeID uuid.UUID `gorm:"type:uuid;not null;index" json:"request_type_id"`
	//Level is ordinal rank inside approval chain (1 = first)
	Level int `gorm:"not null;index" json:"level"`
	//User assigned to this level (nullable if group-based)
	UserID *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	//Optional : role/group/department to resolve dynamic approvers
	GroupName string `gorm:"size:128" json:"group_name,omitempty"`
	//Mode : "AND" or "OR" for multi-approver step behaviour
	Mode constants.StepMode `gorm:"size:10;default:'AND'" json:"mode,omitempty"`
	//Priority within level (smaller first)
	Priority  int       `gorm:"default:0" json:"priority,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	//relations (optional)
	RequestTypeObj *RequestType `gorm:"foreignKey:RequestTypeID;reference:ID" json:"request_type,omitempty"`
	User           *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
