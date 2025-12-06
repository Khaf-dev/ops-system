package repository

import (
	"backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApprovalFlowRepository struct {
	DB *gorm.DB
}

func NewApprovalFlowRepository(db *gorm.DB) *ApprovalFlowRepository {
	return &ApprovalFlowRepository{DB: db}
}

func (r *ApprovalFlowRepository) Create(flow *models.ApprovalFlow) error {
	return r.DB.Create(flow).Error
}

func (r *ApprovalFlowRepository) GetByRequestID(requestID uuid.UUID) (*models.ApprovalFlow, error) {
	var f models.ApprovalFlow
	if err := r.DB.Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_number ASC")
	}).First(&f, "request_id = ?", requestID).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *ApprovalFlowRepository) GetByID(id uuid.UUID) (*models.ApprovalFlow, error) {
	var f models.ApprovalFlow
	if err := r.DB.Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_number ASC")
	}).First(&f, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *ApprovalFlowRepository) Update(flow *models.ApprovalFlow) error {
	return r.DB.Save(flow).Error
}
