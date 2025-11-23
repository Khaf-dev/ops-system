package repository

import (
	"backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttachmentRepository struct {
	DB *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{DB: db}
}

func (r *AttachmentRepository) Create(a *models.Attachment) error {
	return r.DB.Create(a).Error
}

func (r *AttachmentRepository) GetByRequestID(id uuid.UUID) ([]models.Attachment, error) {
	var list []models.Attachment
	err := r.DB.Where("request_id = ?", id).Find(&list).Error
	return list, err
}
