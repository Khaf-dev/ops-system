package handlers

import (
	"backend/internal/app/services"
	"backend/internal/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ConfigHandler struct {
	Svc *services.ApproverConfigService
}

func NewConfigHandler(svc *services.ApproverConfigService) *ConfigHandler {
	return &ConfigHandler{Svc: svc}
}

func (h *ConfigHandler) ListByRequestType(c *gin.Context) {
	rid, err := uuid.Parse(c.Param("request_type"))
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid request type UUID")
		return
	}

	cfgs, err := h.Svc.ListByRequestType(c, rid)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "ok", cfgs)
}

func (h *ConfigHandler) Create(c *gin.Context) {
	var body services.CreateConfigDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	cfg, err := h.Svc.Create(c, body)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "created_successfully", cfg)
}

func (h *ConfigHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid config id")
		return
	}
	var body services.UpdateConfigDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	cfg, err := h.Svc.Update(c, id, body)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "update_successfully", cfg)
}

func (h *ConfigHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid config id")
		return
	}

	err = h.Svc.Delete(c, id)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, 200, "Deleted", nil)
}
