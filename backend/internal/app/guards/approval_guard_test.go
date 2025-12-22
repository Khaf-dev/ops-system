package guards

import (
	"backend/internal/app/constants"
	"backend/internal/app/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	userID := uuid.New()

	guard := NewApprovalGuard(
		&mockFlowRepo{
			flow: &models.ApprovalFlow{
				Status:      constants.RequestInReview,
				CurrentStep: 1,
			},
		},
		&mockStepRepo{
			steps: []models.ApprovalStep{
				{StepNumber: 1, UserID: &userID},
			},
		},
		&mockUserRepo{},
	)

	err := guard.CanApprove(userID, uuid.New())
	assert.NoError(t, err)
}

func TestCanApprove_NotApprove(t *testing.T) {
	userID := uuid.New()
	other := uuid.New()

	guard := NewApprovalGuard(
		&mockFlowRepo{
			flow: &models.ApprovalFlow{
				Status:      constants.RequestInReview,
				CurrentStep: 1,
			},
		},
		&mockStepRepo{
			steps: []models.ApprovalStep{
				{StepNumber: 1, UserID: &other},
			},
		},
		&mockUserRepo{},
	)

	err := guard.CanApprove(userID, uuid.New())
	assert.ErrorIs(t, err, ErrNotApprover)
}

func TestCanApprove_GroupApprover(t *testing.T) { // MASIH ERROR YE GENGS
	userID := uuid.New()

	guard := NewApprovalGuard(
		&mockFlowRepo{
			flow: &models.ApprovalFlow{
				Status:      constants.RequestInReview,
				CurrentStep: 1,
			},
		},
		&mockStepRepo{
			steps: []models.ApprovalStep{
				{StepNumber: 1, GroupName: "finance"},
			},
		},
		&mockUserRepo{
			inGroup: true,
		},
	)

	err := guard.CanApprove(userID, uuid.New())
	assert.NoError(t, err)
}

func TestCanApprove_GroupNotMember(t *testing.T) { // MASIH ERROR BJIR ANJAY
	userID := uuid.New()

	guard := NewApprovalGuard(
		&mockFlowRepo{
			flow: &models.ApprovalFlow{
				Status:      constants.RequestInReview,
				CurrentStep: 1,
			},
		},
		&mockStepRepo{
			steps: []models.ApprovalStep{
				{StepNumber: 1, GroupName: "finance"},
			},
		},
		&mockUserRepo{
			inGroup: false,
		},
	)

	err := guard.CanApprove(userID, uuid.New())
	assert.ErrorIs(t, err, ErrNotApprover)
}

func TestCanApprove_FlowNotActive(t *testing.T) {
	guard := NewApprovalGuard(
		&mockFlowRepo{
			flow: &models.ApprovalFlow{
				Status: constants.RequestApproved,
			},
		},
		&mockStepRepo{},
		&mockUserRepo{},
	)

	err := guard.CanApprove(uuid.New(), uuid.New())
	assert.ErrorIs(t, err, ErrFlowNotActive)
}

func TestCanApprove_NoStep(t *testing.T) {
	guard := NewApprovalGuard(
		&mockFlowRepo{
			flow: &models.ApprovalFlow{
				Status:      constants.RequestInReview,
				CurrentStep: 1,
			},
		},
		&mockStepRepo{
			steps: []models.ApprovalStep{},
		},
		&mockUserRepo{},
	)

	err := guard.CanApprove(uuid.New(), uuid.New())
	assert.ErrorIs(t, err, ErrNotCurrentStep)
}
