package models

import (
	"time"

	"github.com/google/uuid"
)

type UserLevel struct {
	UserID     uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"user_id"`
	LevelID    uuid.UUID `gorm:"type:uuid;not null" json:"level_id"`
	AssignedAt time.Time `gorm:"autoCreatedTime" json:"assigned_at"`

	Level Level `gorm:"foreignKey:LevelID" json:"level,omitempty"`
}
