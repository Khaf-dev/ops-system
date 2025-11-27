package repository

import (
	"backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestTypeRepository struct {
	DB *gorm.DB
}

func NewRequestTypeRepository(db *gorm.DB) *RequestTypeRepository {
	return &RequestTypeRepository{DB: db}
}

func (r *RequestTypeRepository) Create(rt *models.RequestType) error {
	return r.DB.Create(rt).Error
}

func (r *RequestTypeRepository) GetByID(id uuid.UUID) (*models.RequestType, error) {
	var rt models.RequestType
	err := r.DB.First(&rt, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *RequestTypeRepository) ListAll() ([]models.RequestType, error) {
	var list []models.RequestType
	if err := r.DB.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// optional for update
func (r *RequestTypeRepository) Update(rt *models.RequestType) error {
	return r.DB.Save(rt).Error
}

func (r *RequestTypeRepository) SetActive() error {
	return r.DB.Error // <== TODO : PLEASE FIX THIS BEFORE YOU USE THIS CODES!!!
}
