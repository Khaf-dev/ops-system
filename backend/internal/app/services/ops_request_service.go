package services

import (
	"backend/internal/app/dto"
	"backend/internal/app/models"
	"backend/internal/app/repository"
	"backend/internal/app/utils"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OpsRequestService struct {
	Repo *repository.OpsRequestRepository
}

func NewOpsRequestService(repo *repository.OpsRequestRepository) *OpsRequestService {
	return &OpsRequestService{Repo: repo}
}

/*
	ToDTO exported mapping. So handlers can reuse consistent mapping.
	Users preloaded relations if available
*/

func (s *OpsRequestService) ToDTO(m *models.OpsRequest) dto.OpsRequestDTO {
	var requesterName, siteName, approvedByName, reqTypeName, activityName string
	if m.Requester != nil {
		requesterName = m.Requester.Name
	}
	if m.Site != nil {
		siteName = m.Site.Name
	}
	if m.ApprovedBy != nil {
		approvedByName = m.ApprovedBy.Name
	}
	if m.RequestType != nil {
		reqTypeName = m.RequestType.Name
	}
	if m.Activity != nil {
		activityName = m.Activity.Name
	}

	return dto.OpsRequestDTO{
		ID:              m.ID,
		RequesterID:     m.RequesterID,
		RequesterName:   requesterName,
		SiteID:          m.SiteID,
		SiteName:        siteName,
		RequestTypeID:   m.RequestTypeID,
		RequestTypeName: reqTypeName,
		ActivityID:      m.ActivityID,
		ActivityName:    activityName,
		LeaderName:      m.LeaderName,
		RequestDate:     m.RequestDate,
		Location:        m.Location,
		Amount:          m.Amount,
		Description:     m.Description,
		Status:          m.Status,
		ApprovedByID:    m.ApprovedByID,
		ApprovedByName:  approvedByName,
		Latitude:        m.Latitude,
		Longitude:       m.Longitude,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

// CreateOpsRequest: validate -> create -> reload with preloads -> return DTO
func (s *OpsRequestService) CreateOpsRequest(req *models.OpsRequest) (*dto.OpsRequestDTO, error) {
	if req.RequesterID == uuid.Nil {
		return nil, errors.New("invalid requester ID")
	}
	if req.SiteID == uuid.Nil {
		return nil, errors.New("site_id is required")
	}
	if req.RequestTypeID == uuid.Nil {
		return nil, errors.New("request_type_id must be specified")
	}
	if req.ActivityID == uuid.Nil {
		return nil, errors.New("activity_id is required")
	}
	if req.RequestDate == nil {
		return nil, errors.New("request_date is required")
	}
	if req.Amount < 0 {
		return nil, errors.New("amount must be zero")
	}

	req.Status = "pending"
	//set timestamps if not set
	if req.CreatedAt.IsZero() {
		req.CreatedAt = time.Now()
	}
	req.UpdatedAt = time.Now()

	if err := s.Repo.Create(req); err != nil {
		return nil, err
	}

	// reload with relations for proper DTO
	m, err := s.Repo.GetByID(req.ID, "Requester", "ApprovedBy", "Site", "RequestType", "Activity", "Approvals", "Attachments")
	if err != nil {
		return nil, fmt.Errorf("created but failed to reload: %w", err)
	}

	d := s.ToDTO(m)
	return &d, nil
}

// GetByIDDTO returns DTO for a given id (with preloads)
func (s *OpsRequestService) GetByIDDTO(id uuid.UUID) (*dto.OpsRequestDTO, error) {
	m, err := s.Repo.GetByID(id, "Requester", "ApprovedBy", "Site", "RequestType", "Activity", "Approvals", "Attachments")
	if err != nil {
		return nil, utils.ErrNotFound
	}
	d := s.ToDTO(m)
	return &d, nil
}

// List returns paged DTO depending on role
func (s *OpsRequestService) List(userRole string, userID uuid.UUID, limit, offset int) (*repository.PagedResult[dto.OpsRequestDTO], error) {
	var raw *repository.PagedResult[models.OpsRequest]
	var err error

	if userRole == "admin" {
		raw, err = s.Repo.ListAll(limit, offset)
	} else {
		raw, err = s.Repo.ListByRequester(userID, limit, offset)
	}
	if err != nil {
		return nil, err
	}

	items := make([]dto.OpsRequestDTO, 0, len(raw.Items))
	for i := range raw.Items {
		// ensure each item has preloads - repo already preloads Requester & Site
		items = append(items, s.ToDTO(&raw.Items[i]))
	}
	return &repository.PagedResult[dto.OpsRequestDTO]{Items: items, Total: raw.Total}, nil
}

// UpdateOpsRequest: apply map of updates (only provided keys), ownership check
func (s *OpsRequestService) UpdateOpsRequest(id uuid.UUID, userID uuid.UUID, role string, updates map[string]interface{}) error {
	existing, err := s.Repo.GetByID(id, "")
	if err != nil {
		return utils.ErrNotFound
	}
	// ownership check: requester or admin
	if role != "admin" && existing.RequesterID != userID {
		return utils.ErrForbidden
	}

	//apply allowed keys -> update existing model fields
	//only keys in "updates" will be applied
	if v, ok := updates["leader_name"].(string); ok {
		existing.LeaderName = v
	}
	if v, ok := updates["request_date"].(*time.Time); ok && v != nil {
		existing.RequestDate = v
	} else if v2, ok := updates["request_date"].(time.Time); ok {
		existing.RequestDate = &v2
	}
	if v, ok := updates["location"].(string); ok {
		existing.Location = v
	}
	if v, ok := updates["amount"].(float64); ok {
		// amount must be >= 0 (SQL constraint already enforces)
		existing.Amount = v
	}
	if v, ok := updates["description"].(string); ok {
		existing.Description = v
	}
	// only admin can change status
	if v, ok := updates["status"].(string); ok && role == "admin" {
		existing.Status = v
	}
	if v, ok := updates["site_id"].(uuid.UUID); ok && role == "admin" {
		existing.SiteID = v
	}
	if v, ok := updates["request_type_id"].(uuid.UUID); ok && role == "admin" {
		existing.RequestTypeID = v
	}
	if v, ok := updates["activity_id"].(uuid.UUID); ok && role == "admin" {
		existing.ActivityID = v
	}

	existing.UpdatedAt = time.Now()

	return s.Repo.Update(existing)
}

// DEleteOpsRequest with ownership check
func (s *OpsRequestService) DeleteOpsRequest(id uuid.UUID, userID uuid.UUID, role string) error {
	existing, err := s.Repo.GetByID(id, "")
	if err != nil {
		return utils.ErrNotFound
	}
	if role != "admin" && existing.RequesterID != userID {
		return utils.ErrForbidden
	}
	return s.Repo.Delete(id)
}
