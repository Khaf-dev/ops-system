package guards

import (
	"backend/internal/app/models"
	"testing"

	"github.com/google/uuid"
)

type mockFlowRepo struct {
	flow *models.ApprovalFlow
	err  error
}

func (m *mockFlowRepo) GetByID(id uuid.UUID) (*models.ApprovalFlow, error) {
	return m.flow, m.err
}

type mockStepRepo struct {
	steps []models.ApprovalStep
	err   error
}

func (m *mockStepRepo) GetStepsByStepNumber(flowID uuid.UUID, step int) ([]models.ApprovalStep, error) {
	return m.steps, m.err
}

func (m *mockStepRepo) GetMaxStepNumber(flowID uuid.UUID) (int, error) {
	return len(m.steps), nil
}

type mockUserRepo struct {
	inGroup bool
	err     error
}

func (m *mockUserRepo) IsUserInGroup(userID uuid.UUID, group string) (bool, error) {
	return m.inGroup, m.err
}

func TestCanApprove_DirectUser(t *testing.T) {
	
}
