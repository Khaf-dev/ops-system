package logic

import (
	"backend/internal/app/constants"
	"backend/internal/app/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildStepsFromConfigs_SingleLevelSingleUser(t *testing.T) {
	logic := NewApprovalLogic()

	uid := uuid.New()
	cfgs := []models.ApproverConfig{
		{
			Level:  1,
			UserID: &uid,
			Mode:   constants.ModeAND,
		},
	}

	steps, err := logic.BuildStepsFromConfigs(cfgs)

	assert.NoError(t, err)
	assert.Len(t, steps, 1)

	step := steps[0]
	assert.Equal(t, 1, step.StepNumber)
	assert.Equal(t, uid, *step.UserID)
	assert.Equal(t, constants.ModeAND, step.Mode)
}

func TestBuildStepsFromConfigs_PriorityOrdering(t *testing.T) {
	logic := NewApprovalLogic()

	u1 := uuid.New()
	u2 := uuid.New()

	cfgs := []models.ApproverConfig{
		{Level: 1, Priority: 2, UserID: &u2},
		{Level: 1, Priority: 1, UserID: &u1},
	}

	steps, err := logic.BuildStepsFromConfigs(cfgs)

	assert.NoError(t, err)
	assert.Len(t, steps, 2)

	assert.Equal(t, u1, *steps[0].UserID)
	assert.Equal(t, u2, *steps[1].UserID)
}

func TestDetermineNextStepNumber_NextPending(t *testing.T) {
	logic := NewApprovalLogic()

	flow := &models.ApprovalFlow{
		CurrentStep: 1,
		Steps: []models.ApprovalStep{
			{StepNumber: 1, Status: constants.RequestApproved},
			{StepNumber: 2, Status: constants.RequestPending},
		},
	}

	next := logic.DetermineNextStepNumber(flow)

	assert.Equal(t, 2, next)
}

func TestDetermineNextStepNumber_FinalApproved(t *testing.T) {
	logic := NewApprovalLogic()

	flow := &models.ApprovalFlow{
		CurrentStep: 2,
		Steps: []models.ApprovalStep{
			{StepNumber: 1, Status: constants.RequestApproved},
			{StepNumber: 2, Status: constants.RequestApproved},
		},
	}

	next := logic.DetermineNextStepNumber(flow)

	assert.Equal(t, 0, next) // finalize
}
