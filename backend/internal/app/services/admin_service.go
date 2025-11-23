package services

import (
	"backend/internal/app/models"
	"backend/internal/app/repository"
	"backend/internal/app/utils"

	"github.com/google/uuid"
)

type AdminService struct {
	UserRepo  *repository.UserRepository
	LevelRepo *repository.LevelRepository
	ReqRepo   *repository.RequestTypeRepository
}

func NewAdminService(u *repository.UserRepository, l *repository.LevelRepository, r *repository.RequestTypeRepository) *AdminService {
	return &AdminService{
		UserRepo:  u,
		LevelRepo: l,
		ReqRepo:   r,
	}
}

// List all users (simple)
func (s *AdminService) ListUsers() ([]models.User, error) {
	// naive : use user repo via GORM
	return s.UserRepo.ListAll()
}

// set user's levels (replace)
func (s *AdminService) SetUserLevels(userID uuid.UUID, levelIDs []uuid.UUID) error {
	// validate user
	if _, err := s.UserRepo.FindByID(userID); err != nil {
		return utils.ErrNotFound
	}
	// naively remove all and add new (repo should implement helpers)
	if err := s.UserRepo.RemoveAllLevels(userID); err != nil {
		return err
	}
	for _, lid := range levelIDs {
		if err := s.UserRepo.AssignLevel(userID, lid); err != nil {
			return err
		}
	}
	return nil
}

// List request types (admin)
func (s *AdminService) ListRequestTypes() ([]models.RequestType, error) {
	return s.ReqRepo.ListAll()
}

// For small admin endpoints like create request type
func (s *AdminService) CreateRequestType(r *models.RequestType) (*models.RequestType, error) {
	if r.Name == "" {
		return nil, utils.NewBadRequest("name required")
	}
	if r.RequiredLevelRank <= 0 {
		r.RequiredLevelRank = 1
	}
	if err := s.ReqRepo.Create(r); err != nil {
		return nil, err
	}
	return r, nil
}
