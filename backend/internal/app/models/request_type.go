package models

import "github.com/google/uuid"

type RequestType struct {
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name              string    `json:"name"`
	IsActive          bool      `json:"is_active"`
	RequiredLevelRank int       `gorm:"default:1" json:"required_level_rank"`
	// TODO : buat MinApprovalLevel
}
