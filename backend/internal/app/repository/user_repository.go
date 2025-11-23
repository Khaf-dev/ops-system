package repository

import (
	"backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.DB.Preload("Levels.Level").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Preload("Levels.Level").First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) AssignLevel(userID uuid.UUID, levelID uuid.UUID) error {
	ul := models.UserLevel{
		ID:      uuid.New(),
		UserID:  userID,
		LevelID: levelID,
	}
	return r.DB.Create(&ul).Error
}

func (r *UserRepository) RemoveLevel(userID uuid.UUID, levelID uuid.UUID) error {
	return r.DB.Delete(&models.UserLevel{}, "user_id = ? AND level_id = ?", userID, levelID).Error
}

func (r *UserRepository) GetUserLevels(userID uuid.UUID) ([]models.UserLevel, error) {
	var levels []models.UserLevel
	if err := r.DB.Preload("Level").
		Where("user_id = ?", userID).
		Order("levels.priority ASC").
		Joins("JOIN levels ON levels.id = user_levels.level_id").
		Find(&levels).Error; err != nil {
		return nil, err
	}
	return levels, nil
}
