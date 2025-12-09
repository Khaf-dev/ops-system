package repository

import (
	"backend/internal/app/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApprovalLogRepository struct {
	DB *gorm.DB
}

func NewApprovalLogRepository(db *gorm.DB) *ApprovalLogRepository {
	return &ApprovalLogRepository{DB: db}
}

func (r *ApprovalLogRepository) Add(flowID, stepID *uuid.UUID, action string, byUserID *uuid.UUID, note string) error {
	log := &models.ApprovalLog{
		ID:        uuid.New(),
		FlowID:    *flowID,
		StepID:    stepID,
		Action:    action,
		ByUserID:  byUserID,
		Note:      note,
		CreatedAt: time.Now(),
	}
	return r.DB.Create(log).Error
}
