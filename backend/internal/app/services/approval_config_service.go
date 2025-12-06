package services

import (
	"backend/internal/app/models"
	"backend/internal/app/repository"

	"github.com/google/uuid"
)

type ApprovalConfigService struct {
	Repo *repository.ApproverConfigRepository
}

func NewApprovalConfigService(repo *repository.ApproverConfigRepository) *ApprovalConfigService {
	return &ApprovalConfigService{Repo: repo}
}

func (s *ApprovalConfigService) GetByType(typeID uuid.UUID) ([]models.ApproverConfig, error) {
	return s.Repo.GetConfigs(typeID)
}

func (s *ApprovalConfigService) Create(cfg *models.ApproverConfig) error {
	return s.Repo.Create(cfg)
}

func (s *ApprovalConfigService) Delete(typeID uuid.UUID) error {
	return s.Repo.Delete(typeID)
}
