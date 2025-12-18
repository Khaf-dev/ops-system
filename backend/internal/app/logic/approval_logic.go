package logic

import (
	"backend/internal/app/constants"
	"backend/internal/app/models"
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
)

type ApprovalLogic struct{}

func NewApprovalLogic() *ApprovalLogic { return &ApprovalLogic{} }

func (l *ApprovalLogic) BuildStepsFromConfigs(cfg []models.ApproverConfig) ([]models.ApprovalStep, error) {

	if len(cfg) == 0 {
		return nil, errors.New("no approver configs")
	}

	sort.SliceStable(cfg, func(i, j int) bool {
		if cfg[i].Level == cfg[j].Level {
			return cfg[i].Priority < cfg[j].Priority
		}
		return cfg[i].Level < cfg[j].Level
	})

	steps := make([]models.ApprovalStep, 0, len(cfg))

	for _, c := range cfg {
		steps = append(steps, models.ApprovalStep{
			StepNumber: c.Level,
			Mode:       c.Mode,
			UserID:     c.UserID,
			GroupName:  c.GroupName,
			Status:     constants.RequestPending,
			CreatedAt:  time.Now(),
		})
	}

	return steps, nil
}

func (l *ApprovalLogic) MaxStepNumber(steps []models.ApprovalStep) int {
	max := 0
	for _, s := range steps {
		if s.StepNumber > max {
			max = s.StepNumber
		}
	}
	return max
}

func (l *ApprovalLogic) IsLastStep(flow *models.ApprovalFlow) bool {
	maxStep := l.MaxStepNumber(flow.Steps)
	return flow.CurrentStep >= maxStep
}

func (l *ApprovalLogic) DetermineNextStepNumber(flow *models.ApprovalFlow) int {
	maxStep := l.MaxStepNumber(flow.Steps)

	if flow.CurrentStep < maxStep {
		return flow.CurrentStep + 1
	}
	return 0 // no next step
}

func (l *ApprovalLogic) StepsByLevel(steps []models.ApprovalStep, level int) []models.ApprovalStep {

	result := make([]models.ApprovalStep, 0)
	for _, s := range steps {
		if s.StepNumber == level {
			result = append(result, s)
		}
	}
	return result
}

func (l *ApprovalLogic) ValidateApproverForStep(step *models.ApprovalStep, userID uuid.UUID) bool {
	if step.UserID != nil && *step.UserID == userID {
		return true
	}
	if step.GroupName != "" {
		return true // caller must verify group membership
	}
	return false
}

func (l *ApprovalLogic) IsStepCompleted(
	steps []models.ApprovalStep,
	mode constants.StepMode,
) bool {
	
	switch mode {
	case constants.ModeAND:
		for _, s := range steps {
			if s.Status != constants.RequestApproved {
				return false
			}
		}
		return true

	case constants.ModeOR:
		for _, s := range steps {
			if s.Status == constants.RequestApproved {
				return true
			}
		}
		return false
	}
	return false
}
