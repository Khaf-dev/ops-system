package models

import (
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

	RequesterID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"requester_id"`
	SiteID       *uuid.UUID `gorm:"type:uuid;not null;index" json:"site_id"`
	LeaderName   string     `gorm:"size:150" json:"leader_name"`
	RequestDate  *time.Time `json:"request_date,omitempty"` // Tanggal Field
	Location     string     `gorm:"size:255" json:"location"`
	RequestType  string     `gorm:"size:100" json:"request_type"` // Jenis Pengajuan
	ActivityName string     `gorm:"size:150" json:"activity"`     // Jenis Kegiatan
	Amount       float64    `gorm:"type:numeric(12,2)" json:"amount"`
	Description  string     `gorm:"type:text" json:"description"`
	Status       string     `gorm:"size:20;default:pending" json:"status"`

	ApprovedByID *uuid.UUID `gorm:"type:uuid;size:150" json:"approved_by_id,omitempty"`

	Latitude  *float64 `gorm:"type:numeric(10,6)" json:"latitude,omitempty"`
	Longitude *float64 `gorm:"type:numeric(10,6)" json:"longitude,omitempty"`

	// attachments -> handled in attachments table
	Requester *User `gorm:"foreignKey:RequesterID" json:"requester,omitempty"`
	ApproveBy *User `gorm:"foreignKey:ApprovedByID" json:"approved_by,omitempty"`
	Site      *Site `gorm:"foreignKey:SiteID" json:"site,omitempty"`
	// Type      RequestType `gorm:"foreignKey:TypeID" json:"type"`
	// Activity Activity `gorm:"foreignKey:ActivityID" json:"activity"`
	// Attachments []Attachment `gorm:"foreignKey:RequestID" json:"attachments"`
	Approvals []Approval `gorm:"foreignKey:RequestID" json:"approvals,omitempty"`
}
