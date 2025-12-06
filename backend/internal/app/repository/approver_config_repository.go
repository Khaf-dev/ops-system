package repository

import (
	"backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApproverConfigRepository struct {
	DB *gorm.DB
}

func NewApproverConfigRepository(db *gorm.DB) *ApproverConfigRepository {
	return &ApproverConfigRepository{DB: db}
}

// ListByRequestType returns configs ordered by level, priority
func (r *ApproverConfigRepository) ListByRequestType(requestTypeID uuid.UUID) ([]models.ApproverConfig, error) {
	var cfgs []models.ApproverConfig
	if err := r.DB.Where("request_type = ?", requestTypeID).
		Order("level ASC").
		Order("priority ASC").
		Find(&cfgs).Error; err != nil {
		return nil, err
	}
	return cfgs, nil
}

func (r *ApproverConfigRepository) Create(cfg *models.ApproverConfig) error {
	return r.DB.Create(cfg).Error
}

func (r *ApproverConfigRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.ApproverConfig{}, "id = ?", id).Error
}
