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

func (r *ApproverConfigRepository) Create(cfg *models.ApproverConfig) error {
	return r.DB.Create(cfg).Error
}

func (r *ApproverConfigRepository) Update(cfg *models.ApproverConfig) error {
	return r.DB.Save(cfg).Error
}

func (r *ApproverConfigRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.ApproverConfig{}, "id = ?", id).Error
}

func (r *ApproverConfigRepository) GetByRequestType(requestTypeID uuid.UUID) ([]models.ApproverConfig, error) {
	var list []models.ApproverConfig
	err := r.DB.
		Where("request_type_id = ?", requestTypeID).
		Order("level ASC").
		Preload("User").
		Find(&list).Error
	return list, err
}

func (r *ApproverConfigRepository) GetByID(id uuid.UUID) (*models.ApproverConfig, error) {
	var cfg models.ApproverConfig
	if err := r.DB.First(&cfg, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}
