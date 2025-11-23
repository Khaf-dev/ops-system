package services

import (
	"backend/internal/app/models"
	"backend/internal/app/repository"
	"backend/internal/app/utils"

	"github.com/google/uuid"
)

type UserService struct {
	Repo *repository.UserRepository
	// opsional untuk repo lainnya kek levels repo atau apapun itu tar coba bisa ditambahin
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) FindByID(id uuid.UUID) (*models.User, error) {
	u, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, utils.ErrNotFound
	}
	return u, nil
}

func (s *UserService) FindByEmail(email string) (*models.User, error) {
	u, err := s.Repo.FindByEmail(email)
	if err != nil {
		return nil, utils.ErrNotFound
	}
	return u, nil
}

func (s *UserService) AssignLevel(userID, levelID uuid.UUID) error {
	// ensure user exists
	if _, err := s.Repo.FindByID(userID); err != nil {
		return utils.ErrNotFound
	}
	// ensure level exists via repo method (if repo has one) - repo.AssignLevel will fail if FK invalid
	return s.Repo.AssignLevel(userID, levelID)
}

func (s *UserService) RemoveLevel(userID, levelID uuid.UUID) error {
	return s.Repo.RemoveLevel(userID, levelID)
}

func (s *UserService) GetUserLevels(userID uuid.UUID) ([]models.UserLevel, error) {
	return s.Repo.GetUserLevels(userID)
}
