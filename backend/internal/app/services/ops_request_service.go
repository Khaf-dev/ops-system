package services

import (
	"backend/internal/app/dto"
	"backend/internal/app/models"
	"backend/internal/app/repository"
	"errors"

	"github.com/google/uuid"
)

type OpsRequestService struct {
	Repo *repository.OpsRequestRepository
}

func NewOpsRequestService(repo *repository.OpsRequestRepository) *OpsRequestService {
	return &OpsRequestService{Repo: repo}
}

func mapToDTO(m *models.OpsRequest) dto.OpsRequestDTO {
	var requesterName, siteName, approvedByName string
	if m.Requester != nil {
		requesterName = m.Requester.Name
	}
	if m.Site != nil {
		siteName = m.Site.Name
	}
	if m.ApproveBy != nil {
		approvedByName = m.ApproveBy.Name
	}

	return dto.OpsRequestDTO{
		ID:             m.ID,
		RequesterID:    m.RequesterID,
		RequesterName:  requesterName,
		SiteID:         m.SiteID,
		SiteName:       siteName,
		RequestType:    m.RequestType,
		ActivityName:   m.ActivityName,
		LeaderName:     m.LeaderName,
		RequestDate:    m.RequestDate,
		Location:       m.Location,
		Amount:         m.Amount,
		Description:    m.Description,
		Status:         m.Status,
		ApprovedByID:   m.ApprovedByID,
		ApprovedByName: approvedByName,
		Latitude:       m.Latitude,
		Longitude:      m.Longitude,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func (s *OpsRequestService) CreateOpsRequest(req *models.OpsRequest) error {
	if req.RequesterID == uuid.Nil {
		return errors.New("invalid requester ID")
	}
	if req.RequestType == "" {
		return errors.New("request_type must be specified")
	}
	if req.ActivityName == "" {
		return errors.New("activity_name is required")
	}
	if req.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	req.Status = "pending"
	// TODO : bisa nanti ditambahin notifikasi/log/approval flow disini yow
	return s.Repo.Create(req)
}

func (s *OpsRequestService) GetByIDDTO(id uuid.UUID) (*dto.OpsRequestDTO, error) {
	m, err := s.Repo.GetByID(id, "Requester", "ApproveBy", "Site", "Approvals")
	if err != nil {
		return nil, errors.New("not found")
	}
	dto := mapToDTO(m)
	return &dto, nil
}

// list returns DTO paged result based on rule
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

	// map items to DTo
	items := make([]dto.OpsRequestDTO, 0, len(raw.Items))
	for i := range raw.Items {
		items = append(items, mapToDTO(&raw.Items[i]))
	}
	return &repository.PagedResult[dto.OpsRequestDTO]{Items: items, Total: raw.Total}, nil
}

func (s *OpsRequestService) UpdateOpsRequest(id uuid.UUID, userID uuid.UUID, role string, updated *models.OpsRequest) error {
	// load existing
	existing, err := s.Repo.GetByID(id, "")
	if err != nil {
		return errors.New("not found")
	}
	// ownership check: requester or admin
	if role != "admin" && existing.RequesterID != userID {
		return errors.New("forbidden")
	}

	//apply allowed updates (extend)
	if updated.RequestType != "" {
		existing.RequestType = updated.RequestType
	}
	if updated.ActivityName != "" {
		existing.ActivityName = updated.ActivityName
	}
	if updated.LeaderName != "" {
		existing.LeaderName = updated.LeaderName
	}
	if updated.RequestDate != nil {
		existing.RequestDate = updated.RequestDate
	}
	if updated.Location != "" {
		existing.Location = updated.Location
	}
	if updated.Amount != 0 {
		existing.Amount = updated.Amount
	}
	if updated.Description != "" {
		existing.Description = updated.Description
	}
	if updated.Status != "" && role == "admin" {
		//only admin can change status freely
		existing.Status = updated.Status
	}
	// save
	return s.Repo.Update(existing)

}

func (s *OpsRequestService) DeleteOpsRequest(id uuid.UUID, userID uuid.UUID, role string) error {
	existing, err := s.Repo.GetByID(id, "")
	if err != nil {
		return errors.New("not found")
	}
	// allow delete only requester or admin
	if role != "admin" && existing.RequesterID != userID {
		return errors.New("forbidden")
	}
	return s.Repo.Delete(id)
}
