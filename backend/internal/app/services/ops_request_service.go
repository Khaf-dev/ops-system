package services

import (
	constants "backend/internal/app/constants"
	"backend/internal/app/dto"
	"backend/internal/app/models"
	"backend/internal/app/repository"
	"backend/internal/app/utils"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OpsRequestService struct {
	Repo *repository.OpsRequestRepository
}

func NewOpsRequestService(repo *repository.OpsRequestRepository) *OpsRequestService {
	return &OpsRequestService{Repo: repo}
}

func safeName(u *models.User) string {
	if u == nil {
		return ""
	}
	return u.Name
}

// generic relation name helper
func safeStr[T any](v *T, getter func(*T) string) string {
	if v == nil {
		return ""
	}
	return getter(v)
}

func (s *OpsRequestService) ToDTO(m *models.OpsRequest) dto.OpsRequestDTO {
	return dto.OpsRequestDTO{
		ID:              m.ID,
		RequesterID:     m.RequesterID,
		RequesterName:   safeName(m.Requester),
		SiteID:          m.SiteID,
		SiteName:        safeStr(m.Site, func(s *models.Site) string { return s.Name }),
		RequestTypeID:   m.RequestTypeID,
		RequestTypeName: safeStr(m.RequestType, func(r *models.RequestType) string { return r.Name }),
		ActivityID:      m.ActivityID,
		ActivityName:    safeStr(m.Activity, func(a *models.Activity) string { return a.Name }),
		LeaderName:      m.LeaderName,
		RequestDate:     m.RequestDate,
		Location:        m.Location,
		Amount:          m.Amount,
		Description:     m.Description,
		Status:          string(m.Status),
		ApprovedByID:    m.ApprovedByID,
		ApprovedByName:  safeName(m.ApprovedBy),
		Latitude:        m.Latitude,
		Longitude:       m.Longitude,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func (s *OpsRequestService) CreateOpsRequest(req *models.OpsRequest) (*dto.OpsRequestDTO, error) {
	if req.RequesterID == uuid.Nil {
		return nil, errors.New("requester_id is required")
	}
	if req.SiteID == uuid.Nil {
		return nil, errors.New("site_id is required")
	}
	if req.RequestTypeID == uuid.Nil {
		return nil, errors.New("request_type_id is required")
	}
	if req.ActivityID == uuid.Nil {
		return nil, errors.New("activity_id is required")
	}
	if req.RequestDate == nil {
		return nil, errors.New("request_date is required")
	}
	if req.Amount < 0 {
		return nil, errors.New("amount cannot be negative")
	}

	req.Status = constants.RequestPending

	if err := s.Repo.Create(req); err != nil {
		return nil, err
	}

	m, err := s.Repo.GetByID(req.ID,
		"Requester", "ApprovedBy", "Site", "RequestType", "Activity",
		"Approvals", "Attachments",
	)
	if err != nil {
		return nil, fmt.Errorf("created but failed to reload: %w", err)
	}

	dtoRes := s.ToDTO(m)
	return &dtoRes, nil
}

func (s *OpsRequestService) GetByIDDTO(id uuid.UUID) (*dto.OpsRequestDTO, error) {
	m, err := s.Repo.GetByID(id,
		"Requester", "ApprovedBy", "Site", "RequestType", "Activity",
		"Approvals", "Attachments",
	)
	if err != nil {
		return nil, utils.ErrNotFound
	}
	dtoRes := s.ToDTO(m)
	return &dtoRes, nil
}

func (s *OpsRequestService) List(userRole string, userID uuid.UUID, limit, offset int) (*repository.PagedResult[dto.OpsRequestDTO], error) {
	var raw *repository.PagedResult[models.OpsRequest]
	var err error

	switch userRole {
	case "admin":
		raw, err = s.Repo.ListAll(limit, offset)
	default:
		raw, err = s.Repo.ListByRequester(userID, limit, offset)
	}

	if err != nil {
		return nil, err
	}

	items := make([]dto.OpsRequestDTO, 0, len(raw.Items))
	for i := range raw.Items {
		items = append(items, s.ToDTO(&raw.Items[i]))
	}

	return &repository.PagedResult[dto.OpsRequestDTO]{
		Items: items,
		Total: raw.Total,
	}, nil
}

func (s *OpsRequestService) UpdateOpsRequest(
	id uuid.UUID,
	userID uuid.UUID,
	role string,
	input dto.UpdateOpsRequest,
) error {

	req, err := s.Repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound
		}
		return err
	}

	// NON ADMIN cuma bisa edit punya sendiri
	if role != "admin" && req.RequesterID != userID {
		return utils.ErrForbidden
	}

	// APPLY FIELDS
	if input.LeaderName != nil {
		req.LeaderName = *input.LeaderName
	}
	if input.RequestDate != nil {
		req.RequestDate = input.RequestDate
	}
	if input.Location != nil {
		req.Location = *input.Location
	}
	if input.Amount != nil {
		req.Amount = *input.Amount
	}
	if input.Description != nil {
		req.Description = *input.Description
	}

	// STATUS = ADMIN ONLY
	if input.Status != nil {
		if role != "admin" {
			return utils.ErrForbidden
		}

		switch *input.Status {
		case constants.RequestPending,
			constants.RequestCanceled,
			constants.RequestApproved,
			constants.RequestRejected,
			constants.RequestInReview:
			req.Status = *input.Status
		default:
			return errors.New("invalid status value")
		}
	}

	// ADMIN : bisa update  reference
	if role == "admin" {
		if input.SiteID != nil {
			req.SiteID = *input.SiteID
		}
		if input.RequestTypeID != nil {
			req.RequestTypeID = *input.RequestTypeID
		}
		if input.ActivityID != nil {
			req.ActivityID = *input.ActivityID
		}
	}
	return s.Repo.Update(req)
}

func (s *OpsRequestService) DeleteOpsRequest(id uuid.UUID, userID uuid.UUID, role string) error {
	req, err := s.Repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound
		}
		return err
	}

	if role != "admin" && req.RequesterID != userID {
		return utils.ErrForbidden
	}
	return s.Repo.Delete(id)
}
