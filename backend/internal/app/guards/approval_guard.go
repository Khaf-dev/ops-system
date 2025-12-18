package guards

import (
	"backend/internal/app/constants"
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
	if err != nil || flow == nil {
		return ErrFlowNotFound
	}

	if flow.Status != constants.RequestInReview {
		return ErrFlowNotActive
	}

	step, err := g.StepRepo.GetStepsByStepNumber(flow.ID, flow.CurrentStep)
	if err != nil || step == nil {
		return ErrNotCurrentStep
	}

	if step.StepNumber != flow.CurrentStep {
		return ErrNotCurrentStep
	}

	// ============= invalid config
	if step.UserID != nil && step.GroupName != "" {
		return errors.New("invalid approval step config")
	}

	// ============= direct user approver
	if step.UserID != nil {
		if *step.UserID != userID {
			return ErrNotApprover
		}
		return nil
	}

	// ============= group approver
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
