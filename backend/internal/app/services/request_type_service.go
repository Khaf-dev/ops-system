package services

import (
	"backend/internal/app/models"
	"backend/internal/app/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrRequestTypeNotFound = errors.New("request type not found")
)

type RequestTypeService interface {
	GetAll(ctx context.Context, onlyActive bool) ([]models.RequestType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.RequestType, error)
	Create(ctx context.Context, name string) (*models.RequestType, error)
	Update(ctx context.Context, id uuid.UUID, name string, active bool) error
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
}

type requestTypeService struct {
	repo repository.RequestTypeRepository
}

func NewRequestTypeService(repo repository.RequestTypeRepository) RequestTypeService {
	return &requestTypeService{repo: repo}
}

func (s *requestTypeService) GetAll(ctx context.Context, onlyActive bool) ([]models.RequestType, error) {
	return s.repo.GetAll(ctx, onlyActive)
}

func (s *requestTypeService) GetByID(ctx context.Context, id uuid.UUID) (*models.RequestType, error) {
	rt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rt == nil {
		return nil, ErrRequestTypeNotFound
	}
	return rt, nil
}

func (s *requestTypeService) Create(ctx context.Context, name string) (*models.RequestType, error) {
	rt := &models.RequestType{
		ID:       uuid.New(),
		Name:     name,
		IsActive: true,
	}
	if err := s.repo.Create(ctx, rt); err != nil {
		return nil, err
	}
	return rt, nil
}

func (s *requestTypeService) Update(ctx context.Context, id uuid.UUID, name string, active bool) error {
	rt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if rt == nil {
		return ErrRequestTypeNotFound
	}

	rt.Name = name
	rt.IsActive = active
	return s.repo.Update(ctx, rt)
}

func (s *requestTypeService) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	rt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if rt == nil {
		return ErrRequestTypeNotFound
	}
	return s.repo.SetActive(ctx, id, active)
}
