package handlers

import (
	"backend/internal/app/services"
	"backend/internal/app/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApprovalHandler struct {
	Service *services.ApprovalService
}

func NewApprovalHandler(svc *services.ApprovalService) *ApprovalHandler {
	return &ApprovalHandler{Service: svc}
}

type ApprovalInput struct {
	Decision string `json:"decision" binding:"required"`
	Notes    string `json:"notes"`
}

func (h *ApprovalHandler) HandleApproval(c *gin.Context) {
	requestIDStr := c.Param("id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request ID")
		return
	}

	var input ApprovalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid body")
		return
	}

	approvedIDStr, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	approvedID, err := uuid.Parse(approvedIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid approver")
		return
	}

	if err := h.Service.HandleApproval(requestID, approvedID, input.Decision, input.Notes, time.Now()); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "approval recorded", gin.H{
		"request_id": requestID,
		"decision":   input.Decision,
		"notes":      input.Notes,
	})
}
