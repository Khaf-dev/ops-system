package dto

import (
	constants "backend/internal/app/constant"
	"time"

	"github.com/google/uuid"
)

type UpdateOpsRequest struct {
	LeaderName  *string                  `json:"leader_name,omitempty"`
	RequestDate *time.Time               `json:"request_date,omitempty"`
	Location    *string                  `json:"location,omitempty"`
	Amount      *float64                 `json:"amount,omitempty"`
	Description *string                  `json:"description,omitempty"`
	Status      *constants.RequestStatus `json:"status,omitempty"`

	// admin-only fields
	SiteID        *uuid.UUID `json:"site_id,omitempty"`
	RequestTypeID *uuid.UUID `json:"request_type_id,omitempty"`
	ActivityID    *uuid.UUID `json:"activity_id,omitempty"`
}
