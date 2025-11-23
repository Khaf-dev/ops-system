package dto

import (
	"time"

	"github.com/google/uuid"
)

// OpsRequestDTO is the shape returned to clients (safe + small)
type OpsRequestDTO struct {
	ID              uuid.UUID  `json:"id"`
	RequesterID     uuid.UUID  `json:"requester_id"`
	RequesterName   string     `json:"requester_name,omitempty"`
	SiteID          uuid.UUID  `json:"site_id"`
	SiteName        string     `json:"site_name,omitempty"`
	RequestTypeID   uuid.UUID  `json:"request_type_id"`
	RequestTypeName string     `json:"request_type_name,omitempty"`
	ActivityID      uuid.UUID  `json:"activity_id"`
	ActivityName    string     `json:"activity_name,omitempty"`
	LeaderName      string     `json:"leader_name,omitempty"`
	RequestDate     *time.Time `json:"request_date,omitempty"`
	Location        string     `json:"location,omitempty"`
	Amount          float64    `json:"amount"`
	Description     string     `json:"description,omitempty"`
	Status          string     `json:"status"`
	ApprovedByID    *uuid.UUID `json:"approved_by_id,omitempty"`
	ApprovedByName  string     `json:"approved_by_name,omitempty"`
	Latitude        *float64   `json:"latitude,omitempty"`
	Longitude       *float64   `json:"longitude,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
