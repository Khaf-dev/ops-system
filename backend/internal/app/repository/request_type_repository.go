package repository

import (
	"backend/internal/app/models"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type RequestTypeRepository interface {
	GetAll(ctx context.Context, onlyActive bool) ([]models.RequestType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.RequestType, error)
	Create(ctx context.Context, rt *models.RequestType) error
	Update(ctx context.Context, rt *models.RequestType) error
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
}

type requestTypeRepository struct {
	db *sql.DB
}

func NewRequestTypeRepository(db *sql.DB) RequestTypeRepository {
	return &requestTypeRepository{db: db}
}

func (r *requestTypeRepository) GetAll(ctx context.Context, onlyActive bool) ([]models.RequestType, error) {
	query := `SELECT id, name, is_active FROM request_types`
	if onlyActive {
		query += ` WHERE is_active = TRUE`
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.RequestType
	for rows.Next() {
		var rt models.RequestType
		if err := rows.Scan(&rt.ID, &rt.Name, &rt.IsActive); err != nil {
			return nil, err
		}
		list = append(list, rt)
	}
	return list, nil
}

func (r *requestTypeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.RequestType, error) {
	query := `SELECT id, name, is_active FROM request_types WHERE id=$1`

	var rt models.RequestType
	err := r.db.QueryRowContext(ctx, query, id).Scan(&rt.ID, &rt.Name, &rt.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &rt, nil
}

func (r *requestTypeRepository) Create(ctx context.Context, rt *models.RequestType) error {
	query := `
        INSERT INTO request_types (id, name, is_active)
        VALUES ($1, $2, $3)
    `
	if rt.ID == uuid.Nil {
		rt.ID = uuid.New()
	}

	_, err := r.db.ExecContext(ctx, query, rt.ID, rt.Name, rt.IsActive)
	return err
}

func (r *requestTypeRepository) Update(ctx context.Context, rt *models.RequestType) error {
	query := `
        UPDATE request_types
        SET name=$1, is_active=$2
        WHERE id=$3
    `
	_, err := r.db.ExecContext(ctx, query, rt.Name, rt.IsActive, rt.ID)
	return err
}

func (r *requestTypeRepository) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	query := `
        UPDATE request_types
        SET is_active=$1
        WHERE id=$2
    `
	_, err := r.db.ExecContext(ctx, query, active, id)
	return err
}
