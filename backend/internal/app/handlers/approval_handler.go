package handlers

import (
	"backend/internal/app/services"
	"backend/internal/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApprovalHandler struct {
	FlowSvc   *services.ApprovalFlowService
	ActionSvc *services.ApprovalActionService
}

func NewApprovalHandler(
	flowSvc *services.ApprovalFlowService,
	actionSvc *services.ApprovalActionService,
) *ApprovalHandler {
	return &ApprovalHandler{
		FlowSvc:   flowSvc,
		ActionSvc: actionSvc,
	}
}

func (h *ApprovalHandler) StartFlow(c *gin.Context) {
	var req struct {
		RequestID string `json:"request_id" binding:"required,uuid"`
		UserID    string `json:"user_id" binding:"required,uuid"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	requestID, _ := uuid.Parse(req.RequestID)
	userID, _ := uuid.Parse(req.UserID)

	flow, err := h.FlowSvc.StartFlow(requestID, userID)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "flow started", flow)
}

func (h *ApprovalHandler) Approve(c *gin.Context) {
	var req struct {
		FlowID string `json:"flow_id" binding:"required,uuid"`
		UserID string `json:"user_id" binding:"required,uuid"`
		Note   string `json:"note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	flowID, _ := uuid.Parse(req.FlowID)
	userID, _ := uuid.Parse(req.UserID)

	if err := h.ActionSvc.Approve(flowID, userID, req.Note); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Approved", nil)
}

func (h *ApprovalHandler) Reject(c *gin.Context) {
	var req struct {
		FlowID string `json:"flow_id" binding:"required,uuid"`
		UserID string `json:"user_id" binding:"required,uuid"`
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	flowID, _ := uuid.Parse(req.FlowID)
	userID, _ := uuid.Parse(req.UserID)

	if err := h.ActionSvc.Reject(flowID, userID, req.Reason); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Rejected", nil)
}
