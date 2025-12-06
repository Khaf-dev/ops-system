package logic

import (
	"backend/internal/app/models"
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
)

// ApprovalLogic contains pure business rules (no DB)
type ApprovalLogic struct{}

func NewApprovalLogic() *ApprovalLogic {
	return &ApprovalLogic{}
}

// BuildStepsFromConfigs converts approver configs into ordered ApprovalStep slice
func (l *ApprovalLogic) BuildStepsFromConfigs(cfg []models.ApproverConfig) ([]models.ApprovalStep, error) {
	if len(cfg) == 0 {
		return nil, errors.New("no approver configs")
	}
	// sort by level then priority
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
		//increment stepNumber when encountering new level
		if c.Level != currentLevel {
			currentLevel = c.Level
			stepNumber++
		}
		s := models.ApprovalStep{
			StepNumber: stepNumber,
			Mode:       c.Mode,
		}
		if c.UserID != nil {
			uid := *c.UserID
			s.UserID = &uid
		}
		if c.GroupName != "" {
			s.GroupName = c.GroupName
		}
		steps = append(steps, s)
	}
	// ensure steps have consecutive numbers staring 1..N
	for i := range steps {
		steps[i].StepNumber = i + 1
		steps[i].CreatedAt = time.Now()
	}
	return steps, nil
}

// DetermineNextStepNumber returns next step number (or 0 if none)
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

// IsLastStep checks if current step is last
func (l *ApprovalLogic) IsLastStep(flow *models.ApprovalFlow) bool {
	return flow.CurrentStep >= len(flow.Steps)
}

// ValidateApproverForStep checks if userID is allowed to act on step
// For simple mapping, allow if step.UserID == userID or step.GroupName non-empty (caller must then validate group membership)
func (l *ApprovalLogic) ValidateApproverForStep(step *models.ApprovalStep, userID uuid.UUID) bool {
	if step.UserID != nil && *step.UserID == userID {
		return true
	}
	// group based check will be handled at service layer (DB)
	if step.GroupName != "" {
		return true
	}
	return false
}
