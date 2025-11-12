package handlers

import (
	"backend/config"
	"backend/internal/app/models"
	"backend/internal/app/services"
	"backend/internal/app/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB  *gorm.DB
	Svc *services.AuthService
	Cfg *config.Config
}

func NewAuthHandler(db *gorm.DB, s *services.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{DB: db, Svc: s, Cfg: cfg}
}

// ================= REGISTER =================
func (h *AuthHandler) Register(c *gin.Context) {
	var body struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}

	user, err := h.Svc.Register(body.Name, body.Email, body.Phone, body.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user.PasswordHash = ""
	utils.SuccessResponse(c, http.StatusCreated, "registered", user)
}

// ================= LOGIN =================
func (h *AuthHandler) Login(c *gin.Context) {
	var body struct {
		Identifier string `json:"identifier" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}

	user, err := h.Svc.Authenticate(body.Identifier, body.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	access, _ := utils.GenerateAccess(h.Cfg, user.ID, user.Role)
	refresh, err := utils.GenerateRefresh()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}

	refreshToken := models.RefreshToken{
		UserID:    user.ID,
		Token:     refresh.Token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := h.DB.Create(&refreshToken).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to store refresh token")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "login success", gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

// ================= REFRESH =================
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing refresh token")
		return
	}

	var stored models.RefreshToken
	if err := h.DB.First(&stored, "token = ?", req.Token).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	var user models.User
	if err := h.DB.First(&user, "id = ?", stored.UserID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "user not found")
		return
	}

	if time.Now().After(stored.ExpiresAt) {
		h.DB.Delete(&stored)
		utils.ErrorResponse(c, http.StatusUnauthorized, "refresh token expired")
		return
	}

	access, err := utils.GenerateAccess(h.Cfg, stored.UserID, user.Role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to generate new access token")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "token refreshed", gin.H{
		"access_token": access,
	})
}

func (h *AuthHandler) getRefreshToken(token string) (*models.RefreshToken, error) {
	var t models.RefreshToken
	if err := h.DB.First(&t, "token = ?", token).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

// ================= LOGOUT =================
func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing refresh token")
		return
	}

	if err := h.DB.Delete(&models.RefreshToken{}, "token = ?", req.Token).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to logout")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "logged out successfully", nil)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	user.PasswordHash = ""
	utils.SuccessResponse(c, http.StatusOK, "ok", user)
}
