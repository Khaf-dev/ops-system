package repository

import (
	"backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApprovalStepRepository struct {
	DB *gorm.DB
}

func NewApprovalStepRepository(db *gorm.DB) *ApprovalStepRepository {
	return &ApprovalStepRepository{DB: db}
}

func (r *ApprovalStepRepository) GetByFlowID(flowID uuid.UUID) ([]models.ApprovalStep, error) {
	var steps []models.ApprovalStep
	if err := r.DB.Where("flow_id = ?", flowID).
		Order("step_number ASC, created_at ASC").
		Find(&steps).Error; err != nil {
		return nil, err
	}
	return steps, nil
}

func (r *ApprovalStepRepository) GetPendingStepsByFlowAndNumber(flowID uuid.UUID, stepNumber int) ([]models.ApprovalStep, error) {
	var steps []models.ApprovalStep
	if err := r.DB.Where("flow_id = ? AND step_number = ? AND status = ?", flowID, stepNumber, "pending").
		Find(&steps).Error; err != nil {
		return nil, err
	}
	return steps, nil
}

func (r *ApprovalStepRepository) Update(step *models.ApprovalStep) error {
	return r.DB.Save(step).Error
}

func (r *ApprovalStepRepository) BulkCreate(steps []models.ApprovalStep) error {
	return r.DB.Create(&steps).Error
}
