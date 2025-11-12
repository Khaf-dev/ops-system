package models

import (
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ActorID    uuid.UUID  `gorm:"type:uuid;index" json:"actor_id,omitempty"`
	Action     string     `json:"action"`
	TargetType string     `json:"target_type,omitempty"`
	TargetID   *uuid.UUID `gorm:"type:uuid" json:"target_id,omitempty"`
	Metadata   []byte     `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
}
