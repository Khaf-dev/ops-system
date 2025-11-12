package dto

import (
	"time"

	"github.com/google/uuid"
)

// OpsRequestDTO is the shape returned to clients (safe + small)
type OpsRequestDTO struct {
	ID             uuid.UUID  `json:"id"`
	RequesterID    uuid.UUID  `json:"requester_id"`
	RequesterName  string     `json:"requester_name,omitempty"`
	SiteID         *uuid.UUID `json:"site_id,omitempty"`
	SiteName       string     `json:"site_name,omitempty"`
	RequestType    string     `json:"request_type"`
	ActivityName   string     `json:"activity_name"`
	LeaderName     string     `json:"leader_name,omitempty"`
	RequestDate    *time.Time `json:"request_date,omitempty"`
	Location       string     `json:"location,omitempty"`
	Amount         float64    `json:"amount"`
	Description    string     `json:"description"`
	Status         string     `json:"status"`
	ApprovedByID   *uuid.UUID `json:"approved_by_id,omitempty"`
	ApprovedByName string     `json:"approved_by_name,omitempty"`
	Latitude       *float64   `json:"latitude,omitempty"`
	Longitude      *float64   `json:"longitude,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
