package services

import (
	"backend/internal/app/logic"

	"gorm.io/gorm"
)

type ApprovalOrchestrator struct {
	Flowsvc   *ApprovalFlowService
	Actionsvc *ApprovalActionService
	Configsvc *ApproverConfigService
	Logic     *logic.ApprovalLogic
	DB        *gorm.DB
}

func NewApprovalOrchestrator(
	db *gorm.DB,
	flow *ApprovalFlowService,
	action *ApprovalActionService,
	cfg *ApproverConfigService,
	logic *logic.ApprovalLogic,
) *ApprovalOrchestrator {
	return &ApprovalOrchestrator{
		DB:        db,
		Flowsvc:   flow,
		Actionsvc: action,
		Configsvc: cfg,
		Logic:     logic,
	}
}
