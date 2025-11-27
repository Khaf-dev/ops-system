package services

import (
	"backend/internal/app/models"
	"backend/internal/app/repository"
	"backend/internal/app/utils"

	"github.com/google/uuid"
)

type LevelService struct {
	Repo *repository.LevelRepository
}

func NewLevelService(repo *repository.LevelRepository) *LevelService {
	return &LevelService{Repo: repo}
}

func (s *LevelService) Create(level *models.Level) (*models.Level, error) {
	if level.Name == "" {
		return nil, utils.NewBadRequest
	}
	if level.Rank <= 0 {
		level.Rank = 1
	}
	if err := s.Repo.Create(level); err != nil {
		return nil, err
	}
	return level, nil
}

func (s *LevelService) GetAll() ([]models.Level, error) {
	return s.Repo.GetAll()
}

func (s *LevelService) GetByID(id uuid.UUID) (*models.Level, error) {
	l, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, utils.ErrNotFound
	}
	return l, nil
}

func (s *LevelService) Update(level *models.Level) error {
	_, err := s.Repo.GetByID(level.ID)
	if err != nil {
		return utils.ErrNotFound
	}
	return s.Repo.Update(level)
}

func (s *LevelService) Delete(id uuid.UUID) error {
	_, err := s.Repo.GetByID(id)
	if err != nil {
		return utils.ErrNotFound
	}
	return s.Repo.Delete(id)
}
