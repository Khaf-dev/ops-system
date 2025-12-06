package repository

import (
	"backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PageResult[T any] struct {
	Items   []T   `json:"items"`
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasNext bool  `json:"has_next"`
}

type OpsRequestRepository struct {
	DB *gorm.DB
}

func NewOpsRequestRepository(db *gorm.DB) *OpsRequestRepository {
	return &OpsRequestRepository{DB: db}
}

// -------- PRELOAD HELPERS -------- //

var opsRequestPreloads = []string{
	"Requester",
	"Site",
	"RequestType",
	"Activity",
	"Approvals",
	"Attachments",
}

func applyPreloads(q *gorm.DB, preloads []string) *gorm.DB {
	for _, p := range preloads {
		q = q.Preload(p)
	}
	return q
}

// -------- CRUD -------- //

func (r *OpsRequestRepository) Create(req *models.OpsRequest) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(req).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *OpsRequestRepository) Update(req *models.OpsRequest) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(req).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *OpsRequestRepository) GetByID(id uuid.UUID, preload ...string) (*models.OpsRequest, error) {
	var req models.OpsRequest

	q := r.DB.Model(&models.OpsRequest{})
	if len(preload) > 0 {
		q = applyPreloads(q, preload)
	} else {
		q = applyPreloads(q, opsRequestPreloads)
	}

	if err := q.First(&req, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

// -------- LIST FUNCTIONS -------- //
func (r *OpsRequestRepository) ListAll(limit, offset int) (*PageResult[models.OpsRequest], error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var (
		list  []models.OpsRequest
		total int64
	)

	// Consistent query
	query := r.DB.Model(&models.OpsRequest{})
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	q := applyPreloads(query, opsRequestPreloads).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset)

	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}

	return &PageResult[models.OpsRequest]{
		Items:   list,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasNext: int64(offset+limit) < total,
	}, nil
}

func (r *OpsRequestRepository) ListByRequester(userID uuid.UUID, limit, offset int) (*PageResult[models.OpsRequest], error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var (
		list  []models.OpsRequest
		total int64
	)

	query := r.DB.Model(&models.OpsRequest{}).
		Where("requester_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	q := applyPreloads(query, opsRequestPreloads).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset)

	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}

	return &PageResult[models.OpsRequest]{
		Items:   list,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasNext: int64(offset+limit) < total,
	}, nil
}

// -------- DELETE -------- //
func (r *OpsRequestRepository) Delete(id uuid.UUID) error {
	res := r.DB.Delete(&models.OpsRequest{}, "id = ?", id)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
