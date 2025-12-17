package guards

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrFlowNotFound   = errors.New("approval flow not found")
	ErrFlowNotActive  = errors.New("approval flow not active")
	ErrNotCurrentStep = errors.New("not current approval step")
	ErrNotApprover    = errors.New("user is not approver")
)

type ApprovalGuard struct {
	FlowRepo ApprovalFlowReader
	StepRepo ApprovalStepReader
	UserRepo UserGroupChecker
}

func NewApprovalGuard(
	flowRepo ApprovalFlowReader,
	stepRepo ApprovalStepReader,
	userRepo UserGroupChecker,
) *ApprovalGuard {
	return &ApprovalGuard{
		FlowRepo: flowRepo,
		StepRepo: stepRepo,
		UserRepo: userRepo,
	}
}

func (g *ApprovalGuard) CanApprove(userID, flowID uuid.UUID) error {

	flow, err := g.FlowRepo.GetByID(flowID)
	if err != nil {
		return ErrFlowNotFound
	}

	if flow.Status != "in_review" {
		return ErrFlowNotActive
	}

	step, err := g.StepRepo.GetCurrentStep(flow.ID, flow.CurrentStep)
	if err != nil {
		return ErrNotCurrentStep
	}

	// ==== validate approver ==== //
	if step.UserID != nil {
		if *step.UserID != userID {
			return ErrNotApprover
		}
		return nil
	}

	if step.GroupName != "" {
		ok, err := g.UserRepo.IsUserInGroup(userID, step.GroupName)
		if err != nil {
			return err
		}
		if !ok {
			return ErrNotApprover
		}
		return nil
	}

	return ErrNotApprover
}
