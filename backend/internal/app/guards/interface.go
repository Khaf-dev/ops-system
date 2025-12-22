package guards

import (
	"backend/internal/app/models"

	"github.com/google/uuid"
)

type ApprovalFlowReader interface {
	GetByID(id uuid.UUID) (*models.ApprovalFlow, error)
}

type ApprovalStepReader interface {
	GetStepsByStepNumber(flowID uuid.UUID, step int) ([]models.ApprovalStep, error)
	GetMaxStepNumber(flowID uuid.UUID) (int, error)
}

type UserGroupChecker interface {
	IsUserInGroup(userID uuid.UUID, group string) (bool, error)
}
