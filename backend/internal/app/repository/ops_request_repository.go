package repository

import (
	"backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PagedResult[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
}

type OpsRequestRepository struct {
	DB *gorm.DB
}

func NewOpsRequestRepository(db *gorm.DB) *OpsRequestRepository {
	return &OpsRequestRepository{DB: db}
}

func (r *OpsRequestRepository) Create(req *models.OpsRequest) error {
	return r.DB.Create(req).Error
}

func (r *OpsRequestRepository) Update(req *models.OpsRequest) error {
	return r.DB.Save(req).Error
}

func (r *OpsRequestRepository) GetByID(id uuid.UUID, preload ...string) (*models.OpsRequest, error) {
	var req models.OpsRequest
	q := r.DB
	for _, p := range preload {
		if p == "" {
			continue
		}
		q = q.Preload(p)
	}
	if err := q.First(&req, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *OpsRequestRepository) ListAll(limit, offset int) (*PagedResult[models.OpsRequest], error) {
	if limit <= 0 {
		limit = 20
	}
	if offset <= 0 {
		offset = 0
	}

	var list []models.OpsRequest
	var total int64

	query := r.DB.Model(&models.OpsRequest{})
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("Requester").
		Preload("Site").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&list).Error; err != nil {
		return nil, err
	}

	return &PagedResult[models.OpsRequest]{Items: list, Total: total}, nil
}

func (r *OpsRequestRepository) ListByRequester(userID uuid.UUID, limit, offset int) (*PagedResult[models.OpsRequest], error) {
	if limit <= 0 {
		limit = 20
	}
	if offset <= 0 {
		offset = 0
	}

	var list []models.OpsRequest
	var total int64

	query := r.DB.Model(&models.OpsRequest{}).Where("requester_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Preload("Requester").
		Preload("Site").
		Where("requester_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&list).Error; err != nil {
		return nil, err
	}

	return &PagedResult[models.OpsRequest]{
		Items: list,
		Total: total,
	}, nil
}

func (r *OpsRequestRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.OpsRequest{}, "id = ?", id).Error
}
