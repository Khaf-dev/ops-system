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
	Svc *services.ApprovalService
}

func NewApprovalHandler(svc *services.ApprovalService) *ApprovalHandler {
	return &ApprovalHandler{Svc: svc}
}

type ApprovalInput struct {
	Decision string `json:"decision" binding:"required"` // "approved" or "rejected"
	Notes    string `json:"notes"`
}

// POST /approve/:id
func (h *ApprovalHandler) HandleApproval(c *gin.Context) {
	reqIDStr := c.Param("id")
	reqID, err := uuid.Parse(reqIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request id")
		return
	}

	var input ApprovalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}
	if input.Decision != "approved" && input.Decision != "rejected" {
		utils.ErrorResponse(c, http.StatusBadRequest, "decision must be 'approved' or 'rejected'")
		return
	}

	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	approverID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.Svc.HandleApproval(reqID, approverID, input.Decision, input.Notes, time.Now()); err != nil {
		// graceful mapping
		switch err {
		case utils.ErrNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "not found")
		default:
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "ok", nil)
}
