package router

import (
	"backend/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterApprovalRoutes(
	r *gin.RouterGroup,
	h *handlers.ApprovalHandler,
) {
	approval := r.Group("/approvals")
	{
		approval.POST("/start", h.StartFlow) // start approval flow
		approval.POST("/approve", h.Approve) // approve current step
		approval.POST("/reject", h.Reject)   // reject flow
	}
}
