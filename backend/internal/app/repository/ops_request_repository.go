package repository

import (
	"backend/internal/app/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PagedResult[T any] struct {
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

// ---------- PRELOADS ---------- //

var defaultOpsRequestPreloads = []string{
	"Requester",
	"Site",
	"RequestType",
	"Activity",
	"Approvals",
	"Attachments",
}

func preloadAll(q *gorm.DB, preloads []string) *gorm.DB {
	for _, p := range preloads {
		q = q.Preload(p)
	}
	return q
}

// ---------- CRUD ---------- //

func (r *OpsRequestRepository) Create(req *models.OpsRequest) error {
	return r.DB.Create(req).Error
}

func (r *OpsRequestRepository) Update(req *models.OpsRequest) error {
	return r.DB.Save(req).Error
}

func (r *OpsRequestRepository) GetByID(id uuid.UUID, preload ...string) (*models.OpsRequest, error) {
	var req models.OpsRequest

	finalPreloads := defaultOpsRequestPreloads
	if len(preload) > 0 {
		finalPreloads = preload
	}

	q := preloadAll(r.DB.Model(&models.OpsRequest{}), finalPreloads)

	if err := q.First(&req, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

// ---------- PAGINATION HELP ---------- //

func ensurePaging(limit, offset *int) {
	if *limit <= 0 {
		*limit = 20
	}
	if *offset < 0 {
		*offset = 0
	}
}

// ---------- LIST ---------- //

func (r *OpsRequestRepository) ListAll(limit, offset int) (*PagedResult[models.OpsRequest], error) {
	ensurePaging(&limit, &offset)

	var (
		list  []models.OpsRequest
		total int64
	)

	q := r.DB.Model(&models.OpsRequest{})

	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := preloadAll(q, defaultOpsRequestPreloads).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&list).Error; err != nil {
		return nil, err
	}

	return &PagedResult[models.OpsRequest]{
		Items:   list,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasNext: int64(offset+limit) < total,
	}, nil
}

func (r *OpsRequestRepository) ListByRequester(userID uuid.UUID, limit, offset int) (*PagedResult[models.OpsRequest], error) {
	ensurePaging(&limit, &offset)

	var (
		list  []models.OpsRequest
		total int64
	)

	q := r.DB.Model(&models.OpsRequest{}).
		Where("requester_id = ?", userID)

	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := preloadAll(q, defaultOpsRequestPreloads).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&list).Error; err != nil {
		return nil, err
	}

	return &PagedResult[models.OpsRequest]{
		Items:   list,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasNext: int64(offset+limit) < total,
	}, nil
}

// ---------- DELETE ---------- //

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

// ---------- MARK ---------- //
func (r *OpsRequestRepository) MarkApproved(reqID uuid.UUID, approverID uuid.UUID, finalAt time.Time) error {
	return r.DB.Model(&models.OpsRequest{}).
		Where("id = ?", reqID).
		Updates(map[string]interface{}{
			"status":            "approved",
			"approved_by_id":    approverID,
			"final_approved_at": finalAt,
			"updated_at":        time.Now(),
		}).Error
}

func (r *OpsRequestRepository) MarkRejected(reqID uuid.UUID, by uuid.UUID) error {
	return r.DB.Model(&models.OpsRequest{}).
		Where("id = ?", reqID).
		Updates(map[string]interface{}{
			"status":         "rejected",
			"approved_by_id": by,
			"updated_at":     time.Now(),
		}).Error
}
