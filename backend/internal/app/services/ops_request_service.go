package services

import (
	"backend/internal/app/dto"
	"backend/internal/app/models"
	"backend/internal/app/repository"
	"errors"
	"fmt"

	"github.com/google/uuid"
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

func (s *OpsRequestService) ToDTO(m *models.OpsRequest) dto.OpsRequestDTO {
	return dto.OpsRequestDTO{
		ID:            m.ID,
		RequesterID:   m.RequesterID,
		RequesterName: safeName(m.Requester),
		SiteID:        m.SiteID,
		SiteName: func() string {
			if m.Site != nil {
				return m.Site.Name
			}
			return ""
		}(),
		RequestTypeID: m.RequesterID,
		RequestTypeName: func() string {
			if m.RequestType != nil {
				return m.RequestType.Name
			}
			return ""
		}(),
		ActivityID: m.ActivityID,
		ActivityName: func() string {
			if m.Activity != nil {
				return m.Activity.Name
			}
			return ""
		}(),
		LeaderName:     m.LeaderName,
		RequestDate:    m.RequestDate,
		Location:       m.Location,
		Amount:         m.Amount,
		Description:    m.Description,
		Status:         string(m.Status),
		ApprovedByID:   m.ApprovedByID,
		ApprovedByName: safeName(m.ApprovedBy),
		Latitude:       m.Latitude,
		Longitude:      m.Longitude,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
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

	req.Status = "pending"

	if err := s.Repo.Create(req); err != nil {
		return nil, err
	}

	m, err := s.Repo.GetByID(req.ID,
		"Requester", "ApprovedBy", "Site", "RequestType", "Activity", "Approvals", "Attachments",
	)
	if err != nil {
		return nil, fmt.Errorf("created but failed to reload: %w", err)
	}

	dtoRes := s.ToDTO(m)
	return &dtoRes, nil
}

func (s *OpsRequestService) List(userRole string, userID uuid.UUID, limit, offset int) (*repository.PageResult[dto.OpsRequestDTO], error) {
	var raw *repository.PageResult[models.OpsRequest]
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

	return &repository.PageResult[dto.OpsRequestDTO]{
		Items: items,
		Total: raw.Total,
	}, nil
}


