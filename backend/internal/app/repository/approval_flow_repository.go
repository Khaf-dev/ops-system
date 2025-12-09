package repository

import (
	"backend/internal/app/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApprovalFlowRepository struct {
	DB *gorm.DB
}

func NewApprovalFlowRepository(db *gorm.DB) *ApprovalFlowRepository {
	return &ApprovalFlowRepository{DB: db}
}

// Create inserts a new flow and sets timestamps if missing.
func (r *ApprovalFlowRepository) Create(flow *models.ApprovalFlow) error {
	if flow == nil {
		return gorm.ErrInvalidData
	}
	if flow.ID == uuid.Nil {
		flow.ID = uuid.New()
	}
	if flow.CreatedAt.IsZero() {
		flow.CreatedAt = time.Now()
	}
	flow.UpdatedAt = time.Now()
	return r.DB.Create(flow).Error
}

// GetByRequestID fetches a flow by request_id and preloads steps (ordered).
// Returns gorm.ErrRecordNotFound if not present.
func (r *ApprovalFlowRepository) GetByRequestID(requestID uuid.UUID) (*models.ApprovalFlow, error) {
	var f models.ApprovalFlow
	err := r.DB.
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_number ASC, created_at ASC")
		}).
		Preload("Request").
		First(&f, "request_id = ?", requestID).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// GetByID - useful helper
func (r *ApprovalFlowRepository) GetByID(id uuid.UUID) (*models.ApprovalFlow, error) {
	var f models.ApprovalFlow
	err := r.DB.
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_number ASC, created_at ASC")
		}).
		Preload("Request").
		First(&f, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *ApprovalFlowRepository) Update(flow *models.ApprovalFlow) error {
	if flow == nil {
		return gorm.ErrInvalidData
	}
	flow.UpdatedAt = time.Now()
	return r.DB.Save(flow).Error
}

func (r *ApprovalFlowRepository) UpdateCurrentStep(flowID uuid.UUID, step int) error {
	res := r.DB.Model(&models.ApprovalFlow{}).
		Where("id = ?", flowID).
		Updates(map[string]interface{}{
			"current_step": step, "updated_at": time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ApprovalFlowRepository) UpdateStatus(flowID uuid.UUID, status string) error {
	res := r.DB.Model(&models.ApprovalFlow{}).
		Where("id = ?", flowID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// MoveSteps (NEW CODE)
func (r *ApprovalFlowRepository) MoveToStep(flowID uuid.UUID, nextStep int) error {
	return r.DB.Model(&models.ApprovalFlow{}).
		Where("id = ?", flowID).
		Updates(map[string]interface{}{
			"current_step": nextStep,
			"status":       "in_review",
			"updated_at":   time.Now(),
		}).Error
}

func (r *ApprovalFlowRepository) MarkApproved(flowID uuid.UUID) error {
	return r.DB.Model(&models.ApprovalFlow{}).
		Where("id = ?", flowID).
		Updates(map[string]interface{}{
			"current_step": gorm.Expr("GREATEST(current_step, (SELECT COUNT(*) FROM approval_steps WHERE flow_id = ?))", flowID),
			"status":       "approved",
			"updated_at":   time.Now(),
		}).Error
}

func (r *ApprovalFlowRepository) MarkRejected(flowID uuid.UUID) error {
	return r.DB.Model(&models.ApprovalFlow{}).
		Where("id = ?", flowID).
		Updates(map[string]interface{}{
			"status":     "rejected",
			"updated_at": time.Now(),
		}).Error
}
