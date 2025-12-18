package repository

import (
	"backend/internal/app/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApprovalStepRepository struct {
	DB *gorm.DB
}

func NewApprovalStepRepository(db *gorm.DB) *ApprovalStepRepository {
	return &ApprovalStepRepository{DB: db}
}

func (r *ApprovalStepRepository) Create(step *models.ApprovalStep) error {
	if step.CreatedAt.IsZero() {
		step.CreatedAt = time.Now()
	}
	return r.DB.Create(step).Error
}

func (r *ApprovalStepRepository) Update(step *models.ApprovalStep) error {
	return r.DB.Save(step).Error
}

func (r *ApprovalStepRepository) GetByID(id uuid.UUID) (*models.ApprovalStep, error) {
	var s models.ApprovalStep
	if err := r.DB.First(&s, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ApprovalStepRepository) ListByFlow(flowID uuid.UUID) ([]models.ApprovalStep, error) {
	var steps []models.ApprovalStep
	if err := r.DB.Where("flow_id = ?", flowID).Order("step_number ASC, created_at ASC").Find(&steps).Error; err != nil {
		return nil, err
	}
	return steps, nil
}

func (r *ApprovalStepRepository) GetStepsByStepNumber(flowID uuid.UUID, stepNumber int) ([]models.ApprovalStep, error) {

	var steps []models.ApprovalStep
	err := r.DB.Where("flow_id = ? AND step_number = ?", flowID, stepNumber).
		Order("created_at ASC").
		Find(&steps).Error

	if err != nil {
		return nil, err
	}

	return steps, nil
}

func (r *ApprovalStepRepository) GetMaxStepNumber(flowID uuid.UUID) (int, error) {

	var max int
	err := r.DB.
		Model(&models.ApprovalStep{}).
		Where("flow_id = ?", flowID).
		Select("COALESCE(MAX(step_number), 0)").
		Scan(&max).Error

	return max, err
}
