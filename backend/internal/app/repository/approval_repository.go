package repository

import (
	"backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApprovalRepository struct {
	DB *gorm.DB
}

func NewApprovalRepository(db *gorm.DB) *ApprovalRepository {
	return &ApprovalRepository{DB: db}
}

func (r *ApprovalRepository) Create(a *models.Approval) error {
	return r.DB.Create(a).Error
}

func (r *ApprovalRepository) ListByRequest(requestID uuid.UUID) ([]models.Approval, error) {
	var list []models.Approval
	if err := r.DB.Where("request_id = ?", requestID).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ApprovalRepository) BulkInsert(items []models.Approval) error {
	return r.DB.Create(&items).Error
}
