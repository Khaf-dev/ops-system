package repository

import (
	"backend/internal/app/models"
	"errors"

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

func (r *UserRepository) RemoveAllLevels(userID uuid.UUID) error {
	return r.DB.Delete(&models.UserLevel{}, "id = ?", userID).Error
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

// IsUserInGroup ini itu akan mengembalikan nilai true klo misalkan user sesuai dengan nama grup nye.
// Implementasi try user_groups table first; fallback to users.role == groupName
func (r *UserRepository) IsUserInGroup(userID uuid.UUID, groupName string) (bool, error) {
	// 1. Try user_groups table (kalo ada nih)
	type row struct {
		Count int64
	}

	var res row
	err := r.DB.
		Table("user_groups").
		Where("user_id = ? AND group_name = ?", userID, groupName).
		Count(&res.Count).Error

	if err == nil {
		return res.Count > 0, nil
	}

	// klo table nya gaada atau DB nya error. fallback to checking users.role
	// if its a true DB schema issu, returning error might be ddesirable, but fallback gives resilience
	if errors.Is(err, gorm.ErrInvalidDB) || err != nil {
		// fallback : check users.role
		var user models.User
		if err2 := r.DB.Select("role").First(&user, "id = ?", userID).Error; err2 != nil {
			if errors.Is(err2, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err2
		}
		if user.Role == groupName {
			return true, nil
		}
		// klo gaada di grup
		return false, nil
	}
	// if user_groups returned some other error, bubble up
	return false, err
}
