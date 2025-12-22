package handlers

import (
	"backend/internal/app/guards"
	"backend/internal/app/services"
	"backend/internal/app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApprovalHandler struct {
	FlowSvc   *services.ApprovalFlowService
	ActionSvc *services.ApprovalActionService
	Guard     *guards.ApprovalGuard
}

func NewApprovalHandler(
	flowSvc *services.ApprovalFlowService,
	actionSvc *services.ApprovalActionService,
	guard *guards.ApprovalGuard,
) *ApprovalHandler {
	return &ApprovalHandler{
		FlowSvc:   flowSvc,
		ActionSvc: actionSvc,
		Guard:     guard,
	}
}

// ================== STARTFLOW ================== //
func (h *ApprovalHandler) StartFlow(c *gin.Context) {
	var req struct {
		RequestID string `json:"request_id" binding:"required,uuid"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	requestID, _ := uuid.Parse(req.RequestID)
	userID := c.MustGet("user_id").(uuid.UUID)

	flow, err := h.FlowSvc.StartFlow(requestID, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "flow started", flow)
}

// ================== APPROVE ================== //
func (h *ApprovalHandler) Approve(c *gin.Context) {
	var req struct {
		FlowID string `json:"flow_id" binding:"required,uuid"`
		Note   string `json:"note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	flowID, _ := uuid.Parse(req.FlowID)
	userID := c.MustGet("user_id").(uuid.UUID)

	// ========== GUARD CHECK ========== //
	if err := h.Guard.CanApprove(userID, flowID); err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	// ========== ACTION ========== //
	if err := h.ActionSvc.Approve(
		c.Request.Context(),
		flowID,
		userID,
		req.Note,
	); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "approved", nil)
}

// ================== REJECT ================== //
func (h *ApprovalHandler) Reject(c *gin.Context) {
	var req struct {
		FlowID string `json:"flow_id" binding:"required,uuid"`
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	flowID, _ := uuid.Parse(req.FlowID)
	userID := c.MustGet("user_id").(uuid.UUID)

	// ========== GUARD CHECK ========== //
	if err := h.Guard.CanApprove(userID, flowID); err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	if err := h.ActionSvc.Reject(
		c.Request.Context(),
		flowID,
		userID,
		req.Reason,
	); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "rejected", nil)
}
