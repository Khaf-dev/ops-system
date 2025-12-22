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

/*
Can Approve hanya bisa :
1. flow ada & aktif
2. ada step di current step
3. user termasuk salah satu approver (direct/group)
*/

func (g *ApprovalGuard) CanApprove(userID, flowID uuid.UUID) error {
	// ======= flow ======= //
	flow, err := g.FlowRepo.GetByID(flowID)
	if err != nil || flow == nil {
		return ErrFlowNotFound
	}

	if flow.Status != constants.RequestInReview {
		return ErrFlowNotActive
	}

	// ======= steps di current level ======= //
	steps, err := g.StepRepo.GetStepsByStepNumber(flow.ID, flow.CurrentStep)
	if err != nil || len(steps) == 0 {
		return ErrNotCurrentStep
	}

	// ======= cek apakah user termasuk approver ======= //
	for _, step := range steps {

		// invalid config guard
		if step.UserID != nil && step.GroupName != "" {
			return errors.New("invalid approval step config")
		}

		// direct user
		if step.UserID != nil && *step.UserID == userID {
			return nil
		}

		// group approve
		if step.GroupName != "" {
			ok, err := g.UserRepo.IsUserInGroup(userID, step.GroupName)
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
		}
	}

	return ErrNotApprover
}
