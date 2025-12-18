package guards

import (
	"backend/internal/app/constants"
	"backend/internal/app/models"
	"errors"
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
	step *models.ApprovalStep
	err  error
}

func (m *mockStepRepo) GetCurrentStep(flowID uuid.UUID, step int) (*models.ApprovalStep, error) {
	return m.step, m.err
}

type mockUserRepo struct {
	inGroup bool
	err     error
}

func (m *mockUserRepo) IsUserInGroup(userID uuid.UUID, group string) (bool, error) {
	return m.inGroup, m.err
}

func TestCanApprove_UserDirect(t *testing.T) {
	userID := uuid.New()
	flowID := uuid.New()

	guard := NewApprovalGuard(
		&mockFlowRepo{
			flow: &models.ApprovalFlow{
				ID:          flowID,
				Status:      "in_review",
				CurrentStep: 1,
			},
		},
		&mockStepRepo{
			step: &models.ApprovalStep{
				StepNumber: 1,
				UserID:     &userID,
			},
		},
		&mockUserRepo{},
	)

	err := guard.CanApprove(userID, flowID)
	assert.NoError(t, err)
}

func TestCanApprove_NotApprover(t *testing.T) {
	userID := uuid.New()
	otherUser := uuid.New()

	guard := NewApprovalGuard(
		&mockFlowRepo{
			flow: &models.ApprovalFlow{
				Status:      constants.RequestInReview,
				CurrentStep: 1,
			},
		},
		&mockStepRepo{
			step: &models.ApprovalStep{
				StepNumber: 1,
				UserID:     &otherUser,
			},
		},
		&mockUserRepo{},
	)

	err := guard.CanApprove(userID, uuid.New())
	assert.ErrorIs(t, err, ErrNotApprover)
}

func TestCanApprove_GroupApprover(t *testing.T) {
	userID := uuid.New()

	guard := NewApprovalGuard(
		&mockFlowRepo{
			flow: &models.ApprovalFlow{
				Status:      constants.RequestInReview,
				CurrentStep: 1,
			},
		},
		&mockStepRepo{
			step: &models.ApprovalStep{
				StepNumber: 1,
				GroupName:  "finance",
			},
		},
		&mockUserRepo{
			inGroup: true,
		},
	)

	err := guard.CanApprove(userID, uuid.New())
	assert.NoError(t, err)
}

func TestCanApprove_FlowNotActive(t *testing.T) {
	guard := NewApprovalGuard(
		&mockFlowRepo{
			flow: &models.ApprovalFlow{
				Status: "Approved",
			},
		},
		&mockStepRepo{},
		&mockUserRepo{},
	)

	err := guard.CanApprove(uuid.New(), uuid.New())
	assert.ErrorIs(t, err, ErrFlowNotActive)
}

func TestCanApprove_FlowNotFound(t *testing.T) {
	guard := NewApprovalGuard(
		&mockFlowRepo{err: errors.New("db error")},
		&mockStepRepo{},
		&mockUserRepo{},
	)

	err := guard.CanApprove(uuid.New(), uuid.New())
	assert.ErrorIs(t, err, ErrFlowNotFound)
}

func TestCanApprove_StepNil(t *testing.T) {
	guard := NewApprovalGuard(
		&mockFlowRepo{
			flow: &models.ApprovalFlow{
				Status:      "in_review",
				CurrentStep: 1,
			},
		},
		&mockStepRepo{
			step: nil,
		},
		&mockUserRepo{},
	)

	err := guard.CanApprove(uuid.New(), uuid.New())
	assert.ErrorIs(t, err, ErrNotCurrentStep)
}

func TestCanApprove_StepMismatch(t *testing.T) {
	guard := NewApprovalGuard(
		&mockFlowRepo{
			flow: &models.ApprovalFlow{
				Status:      constants.RequestInReview,
				CurrentStep: 2,
			},
		},
		&mockStepRepo{
			step: &models.ApprovalStep{
				StepNumber: 1,
			},
		},
		&mockUserRepo{},
	)

	err := guard.CanApprove(uuid.New(), uuid.New())
	assert.ErrorIs(t, err, ErrNotCurrentStep)
}

func TestCanApprove_GroupRepoError(t *testing.T) {
	guard := NewApprovalGuard(
		&mockFlowRepo{
			flow: &models.ApprovalFlow{
				Status:      "in_review",
				CurrentStep: 1,
			},
		},
		&mockStepRepo{
			step: &models.ApprovalStep{
				GroupName: "finance",
			},
		},
		&mockUserRepo{
			err: errors.New("ldap down"),
		},
	)

	err := guard.CanApprove(uuid.New(), uuid.New())
	assert.Error(t, err)
}
