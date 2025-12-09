package repository

import (
	"backend/internal/app/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApprovalRepository struct {
	DB *gorm.DB
}

func NewApprovalRepository(db *gorm.DB) *ApprovalRepository {
	return &ApprovalRepository{DB: db}
}

// ========== CREATE / BULK ========== //
func (r *ApprovalRepository) Create(a *models.Approval) error {
	return r.DB.Create(a).Error
}

func (r *ApprovalRepository) BulkInsert(items []models.Approval) error {
	return r.DB.Create(&items).Error
}

// ========== GETTERS (hehe) ========== //
func (r *ApprovalRepository) ListByRequest(requestID uuid.UUID) ([]models.Approval, error) {
	var list []models.Approval
	if err := r.DB.Where("request_id = ?", requestID).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// Get single approval record for a request + spesific step
func (r *ApprovalRepository) GetByStep(requestID uuid.UUID, step int) (*models.Approval, error) {
	var a models.Approval
	if err := r.DB.Where("request_id = ? AND step = ?", requestID, step).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

// ========== UPDATE ========== //
func (r *ApprovalRepository) Update(a *models.Approval) error {
	return r.DB.Save(a).Error
}

// Optimized update only status + notes
func (r *ApprovalRepository) UpdateDecision(id uuid.UUID, decision, notes string) error {
	return r.DB.Model(&models.Approval{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"decision": decision,
			"notes":    notes,
		}).Error
}

// Get the latest pending approval for a request
func (r *ApprovalRepository) FindCurrentPending(requestID uuid.UUID) (*models.Approval, error) {
	var a models.Approval
	err := r.DB.Where("request_id = ? AND decision = ?", requestID, "pending").
		Order("step ASC").
		First(&a).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &a, nil
}
