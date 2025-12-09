package repository

import (
	"backend/internal/app/models"
	"errors"
	"strings"

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

func (r *UserRepository) ListAll() ([]models.User, error) {
	var users []models.User
	err := r.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) AssignLevel(userID uuid.UUID, levelID uuid.UUID) error {
	ul := models.UserLevel{
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
	// join levels to order by rank
	if err := r.DB.Preload("Level").
		Joins("JOIN levels ON levels_id = user_levels.level_id").
		Where("user_levels.user_id = ?", userID).
		Order("levels.rank ASC").
		Find(&levels).Error; err != nil {
		return nil, err
	}
	return levels, nil
}

func (r *UserRepository) FindUsersByLevel(levelID uuid.UUID) ([]models.User, error) {
	var users []models.User

	err := r.DB.
		Table("users").
		Select("users.*").
		Joins("JOIN user_levels ul ON ul.user_id = users.id").
		Where("ul.level_id = ?", levelID).
		Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) IsUserInGroup(userID uuid.UUID, groupName string) (bool, error) {
	type cntRow struct {
		Count int64
	}
	var cr cntRow
	err := r.DB.Table("user_groups").
		Where("user_id = ? AND group_name = ?", userID, groupName).
		Count(&cr.Count).Error
	if err == nil {
		return cr.Count > 0, nil
	}

	// if error indicates missing relation or other schema issue, fallback to role check
	// best efforts : check error string for "does not exist"
	if strings.Contains(strings.ToLower(err.Error()), "does not exist") || strings.Contains(strings.ToLower(err.Error()), "undefined_table") {
		// fallback : check user's role
		var u models.User
		if err2 := r.DB.Select("role").First(&u, "id = ?", userID).Error; err2 != nil {
			if errors.Is(err2, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err2
		}
		return u.Role == groupName, nil
	}
	return false, err
}
