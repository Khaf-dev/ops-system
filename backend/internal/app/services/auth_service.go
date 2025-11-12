package services

import (
	"backend/config"
	"backend/internal/app/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB  *gorm.DB
	cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{DB: db, cfg: cfg}
}

func (s *AuthService) Register(name, email, phone, password string) (*models.User, error) {
	var u models.User
	if err := s.DB.Where("email = ?", email).First(&u).Error; err == nil {
		return nil, errors.New("email sudah terdaftar")
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hashed),
		Role:         "user",
		CreatedAt:    time.Now(),
	}
	if err := s.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Authenticate(identifier, password string) (*models.User, error) {
	var u models.User
	if err := s.DB.Where("email = ?", identifier).Or("phone = ?", identifier).First(&u).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return &u, nil
}
