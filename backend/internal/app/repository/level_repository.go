package repository

import (
	"backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LevelRepository struct {
	DB *gorm.DB
}

func NewLevelRepository(db *gorm.DB) *LevelRepository {
	return &LevelRepository{DB: db}
}

func (r *LevelRepository) Create(level *models.Level) error {
	return r.DB.Create(level).Error
}

func (r *LevelRepository) GetByID(id uuid.UUID) (*models.Level, error) {
	var level models.Level
	if err := r.DB.First(&level, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &level, nil
}

func (r *LevelRepository) GetAll() ([]models.Level, error) {
	var levels []models.Level
	if err := r.DB.Order("priority ASC").Find(&levels).Error; err != nil {
		return nil, err
	}
	return levels, nil
}

func (r *LevelRepository) Update(level *models.Level) error {
	return r.DB.Save(level).Error
}

func (r *LevelRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.Level{}, "id = ?", id).Error
}
