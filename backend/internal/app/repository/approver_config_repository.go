package repository

import (
	"backend/internal/app/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApproverConfigRepository struct {
	DB *gorm.DB
}

func NewApproverConfigRepository(db *gorm.DB) *ApproverConfigRepository {
	return &ApproverConfigRepository{DB: db}
}

// ListByRequestType returns configs ordered by level, priority and preloads User & RequestType
func (r *ApproverConfigRepository) ListByRequestType(requestTypeID uuid.UUID) ([]models.ApproverConfig, error) {
	var cfgs []models.ApproverConfig
	if err := r.DB.
		Where("request_type_id = ?", requestTypeID).
		Order("level ASC, priority ASC").
		Preload("User").
		Preload("RequestTypeObj").
		Find(&cfgs).Error; err != nil {
		return nil, err
	}
	return cfgs, nil
}

// GetByID returns single config
func (r *ApproverConfigRepository) GetByID(id uuid.UUID) (*models.ApproverConfig, error) {
	var c models.ApproverConfig
	if err := r.DB.Preload("User").Preload("RequestTypeObj").First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ApproverConfigRepository) Create(cfg *models.ApproverConfig) error {
	if cfg == nil {
		return gorm.ErrInvalidData
	}
	if cfg.ID == uuid.Nil {
		cfg.ID = uuid.New()
	}
	if cfg.CreatedAt.IsZero() {
		cfg.CreatedAt = time.Now()
	}
	cfg.UpdatedAt = time.Now()
	return r.DB.Create(cfg).Error
}

func (r *ApproverConfigRepository) Update(cfg *models.ApproverConfig) error {
	if cfg == nil {
		return gorm.ErrInvalidData
	}
	cfg.UpdatedAt = time.Now()
	return r.DB.Save(cfg).Error
}

func (r *ApproverConfigRepository) Delete(id uuid.UUID) error {
	res := r.DB.Delete(&models.ApproverConfig{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
