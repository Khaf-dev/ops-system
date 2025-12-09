package models

import (
	"time"

	"github.com/google/uuid"
)

type UserLevel struct {
	UserID     uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"user_id"`
	LevelID    uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"level_id"`
	AssignedAt time.Time `gorm:"autoCreateTime" json:"assigned_at"`

	Level Level `gorm:"foreignKey:LevelID;references:ID" json:"level,omitempty"`
}
