package logic

import (
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

	steps := make([]models.ApprovalStep, 0)
	currentLevel := -1
	stepNumber := 0

	for _, c := range cfg {
		if c.Level != currentLevel {
			currentLevel = c.Level
			stepNumber++
		}
		s := models.ApprovalStep{
			StepNumber: stepNumber,
			Mode:       c.Mode,
			UserID:     c.UserID,
			GroupName:  c.GroupName,
			CreatedAt:  time.Now(),
		}
		steps = append(steps, s)
	}

	// normalize step numbers 1..N
	for i := range steps {
		steps[i].StepNumber = i + 1
	}
	return steps, nil
}

func (l *ApprovalLogic) IsLastStep(flow *models.ApprovalFlow) bool {
	return flow.CurrentStep >= len(flow.Steps)
}

func (l *ApprovalLogic) DetermineNextStepNumber(flow *models.ApprovalFlow) int {
	total := len(flow.Steps)
	if total == 0 {
		return 0
	}
	if flow.CurrentStep < total {
		return flow.CurrentStep + 1
	}
	return 0
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
