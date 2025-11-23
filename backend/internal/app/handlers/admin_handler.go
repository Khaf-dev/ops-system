package handlers

import (
	"backend/internal/app/models"
	"backend/internal/app/services"
	"backend/internal/app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	AdminSvc *services.AdminService
	LevelSvc *services.LevelService
	UserSvc  *services.UserService
}

func NewAdminHandler(ad *services.AdminService, ls *services.LevelService, us *services.UserService) *AdminHandler {
	return &AdminHandler{
		AdminSvc: ad,
		LevelSvc: ls,
		UserSvc:  us,
	}
}

// GET /admin/users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	users, err := h.AdminSvc.ListUsers()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "ok", users)
}

// POST /admin/levels
func (h *AdminHandler) CreateLevel(c *gin.Context) {
	var l models.Level
	if err := c.ShouldBindJSON(&l); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}
	created, err := h.LevelSvc.Create(&l)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "created", created)
}

// POST /admin/users/:id/levels
func (h *AdminHandler) SetUserLevels(c *gin.Context) {
	idStr := c.Param("id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}
	var payload struct {
		LevelIDs []uuid.UUID `json:"level_ids"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.AdminSvc.SetUserLevels(uid, payload.LevelIDs); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "ok", nil)
}
