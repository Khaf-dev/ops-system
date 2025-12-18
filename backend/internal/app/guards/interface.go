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
	GetStepsByStepNumber(flowID uuid.UUID, stepNumber int) ([]models.ApprovalStep, error)
	GetMaxStepNumber(flowID uuid.UUID) (int, error)
}

type UserGroupChecker interface {
	IsUserInGroup(userID uuid.UUID, group string) (bool, error)
}
