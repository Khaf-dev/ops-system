package handlers

import (
	"backend/internal/app/services"
	"backend/internal/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApprovalHandler struct {
	Service *services.ApprovalService
}

func NewApprovalHandler(s *services.ApprovalService) *ApprovalHandler {
	return &ApprovalHandler{Service: s}
}

// -----------------------------
// Start Flow
// -----------------------------
func (h *ApprovalHandler) StartFlow(c *gin.Context) {
	reqIDStr := c.Param("request_id")
	startedByStr := c.GetString("user_id") // from auth middleware

	requestID, err := uuid.Parse(reqIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid request_id")
		return
	}

	startedBy, err := uuid.Parse(startedByStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid user_id")
		return
	}

	flow, err := h.Service.StartFlow(requestID, startedBy)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, flow)
}

// -----------------------------
// Approve Step
// -----------------------------
type approveBody struct {
	Note string `json:"note"`
}

func (h *ApprovalHandler) ApproveStep(c *gin.Context) {
	flowIDStr := c.Param("flow_id")
	userStr := c.GetString("user_id")

	flowID, err := uuid.Parse(flowIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid flow_id")
		return
	}

	userID, err := uuid.Parse(userStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid user_id")
		return
	}

	var body approveBody
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, 400, "invalid body")
		return
	}

	err = h.Service.ApproveStep(flowID, userID, body.Note)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"flow_id": flowID,
		"status":  "approved",
	})
}

// -----------------------------
// Reject Step
// -----------------------------
type rejectBody struct {
	Reason string `json:"reason"`
}

func (h *ApprovalHandler) RejectStep(c *gin.Context) {
	flowIDStr := c.Param("flow_id")
	userStr := c.GetString("user_id")

	flowID, err := uuid.Parse(flowIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid flow_id")
		return
	}

	userID, err := uuid.Parse(userStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid user_id")
		return
	}

	var body rejectBody
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, 400, "invalid body")
		return
	}

	err = h.Service.RejectStep(flowID, userID, body.Reason)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"flow_id": flowID,
		"status":  "rejected",
	})
}
