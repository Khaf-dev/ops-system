package models

import (
	"time"

	"github.com/google/uuid"
)

type ApproverConfig struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	RequestTypeID uuid.UUID    `gorm:"type:uuid;not null"`
	RequestType   *RequestType `gorm:"foreignKey:RequestTypeID"`

	Level uint `gorm:"not null"` // e.g : 1,2,3,4 ...

	// Either role-based or spesific user-based approver
	RoleName *string    `gorm:"type:varchar(100)"`
	UserID   uuid.UUID `gorm:"type:uuid"`
	User     *User

	CreatedAt time.Time
	UpdatedAt time.Time
}
