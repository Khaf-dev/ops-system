package models

import (
	constants "backend/internal/app/constants"
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type OpsRequest struct {
	BaseModel

	RequesterID   uuid.UUID `gorm:"type:uuid;not null;index" json:"requester_id"`
	SiteID        uuid.UUID `gorm:"type:uuid;not null;index" json:"site_id"`
	RequestTypeID uuid.UUID `gorm:"type:uuid;not null;index" json:"request_type_id"` // Jenis Pengajuan
	ActivityID    uuid.UUID `gorm:"type:uuid;not null;index" json:"activity_id"`     // Jenis Kegiatan

	LeaderName  string                  `gorm:"size:150" json:"leader_name,omitempty"`
	RequestDate *time.Time              `gorm:"column:request_date" json:"request_date,omitempty"` // Tanggal Field
	Location    string                  `gorm:"size:255" json:"location,omitempty"`
	Amount      float64                 `gorm:"type:numeric(12,2)" json:"amount"`
	Description string                  `gorm:"type:text" json:"description,omitempty"`
	Status      constants.RequestStatus `gorm:"type:varchar(20);default:pending" json:"status"`

	ApprovedByID         *uuid.UUID `gorm:"type:uuid" json:"approved_by_id,omitempty"`
	CurrentApproverID    *uuid.UUID `gorm:"type:uuid" json:"current_approver_id,omitempty"`
	CurrentApprovalLevel int        `gorm:"default:0" json:"current_approval_level"`
	FinalApprovedAt      *time.Time `json:"final_approved_at,omitempty"`

	Latitude  *float64 `gorm:"type:numeric(10,6);check:latitude >= -90 AND latitude <= 90" json:"latitude,omitempty"`
	Longitude *float64 `gorm:"type:numeric(10,6);check:longitude >= -180 AND longitude <= 180" json:"longitude,omitempty"`

	// attachments -> handled in attachments table
	Requester   *User        `gorm:"foreignKey:RequesterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"requester,omitempty"`
	ApprovedBy  *User        `gorm:"foreignKey:ApprovedByID" json:"approved_by,omitempty"`
	Site        *Site        `gorm:"foreignKey:SiteID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"site,omitempty"`
	RequestType *RequestType `gorm:"foreignKey:RequestTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"request_type,omitempty"`
	Activity    *Activity    `gorm:"foreignKey:ActivityID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"activity,omitempty"`
	// Type      RequestType `gorm:"foreignKey:TypeID" json:"type"`
	// Attachments []Attachment `gorm:"foreignKey:RequestID" json:"attachments"`
	Approvals   []Approval   `gorm:"foreignKey:RequestID" json:"approvals,omitempty"`
	Attachments []Attachment `gorm:"foreignKey:RequestID" json:"attachments,omitempty"`

	ReferenceCode string `gorm:"size:50;uniqueIndex" json:"reference_code"`
}
