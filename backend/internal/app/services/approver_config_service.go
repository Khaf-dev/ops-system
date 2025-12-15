package services

import (
	"backend/internal/app/models"
	"backend/internal/app/repository"

	"github.com/google/uuid"
)

type ApproverConfigService struct {
	Repo *repository.ApproverConfigRepository
}

func NewApproverConfigService(repo *repository.ApproverConfigRepository) *ApproverConfigService {
	return &ApproverConfigService{Repo: repo}
}

func (s *ApproverConfigService) ListByRequestType(requestTypeID uuid.UUID) ([]models.ApproverConfig, error) {
	return s.Repo.ListByRequestType(requestTypeID)
}

func (s *ApproverConfigService) Create(cfg *models.ApproverConfig) error {
	return s.Repo.Create(cfg)
}

func (s *ApproverConfigService) Update(cfg *models.ApproverConfig) error {
	return s.Repo.Update(cfg)
}

func (s *ApproverConfigService) Delete(id uuid.UUID) error {
	return s.Repo.Delete(id)
}
