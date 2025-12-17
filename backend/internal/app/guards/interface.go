package guards

import (
	"backend/internal/app/models"

	"github.com/google/uuid"
)

type ApprovalFlowReader interface {
	GetByID(id uuid.UUID) (*models.ApprovalFlow, error)
}

type ApprovalStepReader interface {
	GetCurrentStep(flowID uuid.UUID, step int) (*models.ApprovalStep, error)
}

type UserGroupChecker interface {
	IsUserInGroup(userID uuid.UUID, group string) (bool, error)
}
