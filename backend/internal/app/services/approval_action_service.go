package services

import (
	"backend/internal/app/constants"
	"backend/internal/app/logic"
	"backend/internal/app/models"
	"backend/internal/app/repository"
	"backend/internal/app/utils"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ApprovalActionService: approve/reject steps (transactional)
type ApprovalActionService struct {
	DB       *gorm.DB
	FlowRepo *repository.ApprovalFlowRepository
	StepRepo *repository.ApprovalStepRepository
	ReqRepo  *repository.OpsRequestRepository
	LogRepo  *repository.ApprovalLogRepository // optional - we also use tx.Create directly
	UserRepo *repository.UserRepository
	Logic    *logic.ApprovalLogic
}

func NewApprovalActionService(
	db *gorm.DB,
	flowRepo *repository.ApprovalFlowRepository,
	stepRepo *repository.ApprovalStepRepository,
	reqRepo *repository.OpsRequestRepository,
	logRepo *repository.ApprovalLogRepository,
	userRepo *repository.UserRepository,
	logic *logic.ApprovalLogic,
) *ApprovalActionService {
	return &ApprovalActionService{
		DB:       db,
		FlowRepo: flowRepo,
		StepRepo: stepRepo,
		ReqRepo:  reqRepo,
		LogRepo:  logRepo,
		UserRepo: userRepo,
		Logic:    logic,
	}
}

// Approve: user approves current step. If last -> finalize flow and ops_request.
func (s *ApprovalActionService) Approve(
	ctx context.Context,
	flowID uuid.UUID,
	userID uuid.UUID,
	note string,
) error {

	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// ================== load flow ================== //
		flow, err := s.FlowRepo.GetByID(flowID)
		if err != nil || flow == nil {
			return utils.ErrNotFound
		}

		if flow.Status != constants.RequestInReview {
			return errors.New("flow not in approvable state")
		}

		// ================== load steps (current level) ================== //
		steps, err := s.StepRepo.GetStepsByStepNumber(flow.ID, flow.CurrentStep)
		if err != nil || len(steps) == 0 {
			return errors.New("approval steps not found")
		}

		// ================== find user's step ================== //
		var current *models.ApprovalStep
		for i := range steps {
			if s.Logic.ValidateApproverForStep(&steps[i], userID) {
				current = &steps[i]
				break
			}
		}
		if current == nil {
			return errors.New("user not authorized for this step")
		}

		// ================== group validation ================== //
		if current.GroupName != "" {
			ok, err := s.UserRepo.IsUserInGroup(userID, current.GroupName)
			if err != nil {
				return err
			}
			if !ok {
				return errors.New("user not authorized for this step")
			}
		}

		now := time.Now()

		// ================== APPROVE CURRENT ================== //
		current.Status = constants.RequestApproved
		current.ApprovedAt = &now
		current.Notes = note

		if err := tx.Save(current).Error; err != nil {
			return err
		}

		// ================== LOG APPROVED ================== //
		if err := tx.Create(&models.ApprovalLog{
			ID:        uuid.New(),
			FlowID:    flow.ID,
			StepID:    &current.ID,
			Action:    "step_approved",
			ByUserID:  &userID,
			Note:      note,
			CreatedAt: now,
		}).Error; err != nil {
			return err
		}

		// ================== MODE HANDLING ================== //
		switch current.Mode {

		// AND MODE //
		case constants.ModeAND:
			if !s.Logic.IsStepCompleted(steps, constants.ModeAND) {
				return nil // masih nunggu approver lain
			}

			// OR MODE //
		case constants.ModeOR:
			if !s.Logic.IsStepCompleted(steps, constants.ModeOR) {
				return nil
			}

			for i := range steps {
				if steps[i].ID == current.ID {
					continue
				}
				if steps[i].Status != constants.RequestPending {
					continue
				}

				steps[i].Status = constants.RequestCanceled
				steps[i].ApprovedAt = &now
				steps[i].Notes = "Auto-cancelled due to OR approval"

				if err := tx.Save(&steps[i]).Error; err != nil {
					return err
				}

				if err := tx.Create(&models.ApprovalLog{
					ID:        uuid.New(),
					FlowID:    flow.ID,
					StepID:    &steps[i].ID,
					Action:    "step_auto_cancelled_or",
					ByUserID:  &userID,
					Note:      steps[i].Notes,
					CreatedAt: now,
				}).Error; err != nil {
					return err
				}
			}
		}

		// ================== NEXT STEP ================== //
		next := s.Logic.DetermineNextStepNumber(flow)

		if next == 0 {
			// ====================== FINAL APPROVAL ====================== //
			if err := tx.Model(&models.ApprovalFlow{}).Where("id = ?", flow.ID).Updates(map[string]interface{}{
				"status":            constants.RequestApproved.String(),
				"approved_by_id":    userID,
				"final_approved_at": now,
				"updated_at":        now,
			}).Error; err != nil {
				return err
			}

			if err := tx.Create(&models.ApprovalLog{
				ID:        uuid.New(),
				FlowID:    flow.ID,
				Action:    "flow_approved",
				ByUserID:  &userID,
				CreatedAt: now,
			}).Error; err != nil {
				return err
			}

			return nil
		}

		// ================== MOVE TO NEXT LEVEL ================== //
		if err := tx.Model(&models.ApprovalFlow{}).Where("id = ?", flow.ID).Updates(map[string]interface{}{
			"current_step": next,
			"status":       constants.RequestInReview.String(),
			"updated_at":   now,
		}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.ApprovalLog{
			ID:        uuid.New(),
			FlowID:    flow.ID,
			Action:    "moved_to_next_step",
			ByUserID:  &userID,
			CreatedAt: now,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// Reject: user rejects current step (terminal)
func (s *ApprovalActionService) Reject(
	ctx context.Context,
	flowID uuid.UUID,
	userID uuid.UUID,
	reason string,
) error {

	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// ================== load flow ================== //
		flow, err := s.FlowRepo.GetByID(flowID)
		if err != nil || flow == nil {
			return utils.ErrNotFound
		}

		if flow.Status != constants.RequestInReview {
			return errors.New("flow not in rejectable state")
		}

		// ================== load steps (current level) ================== //
		steps, err := s.StepRepo.GetStepsByStepNumber(flow.ID, flow.CurrentStep)
		if err != nil || len(steps) == 0 {
			return errors.New("approval steps not found")
		}

		// ================== find user's step ================== //
		var current *models.ApprovalStep
		for i := range steps {
			if s.Logic.ValidateApproverForStep(&steps[i], userID) {
				current = &steps[i]
				break
			}
		}
		if current == nil {
			return errors.New("user not authorized for this step")
		}

		// ================== group validation ================== //
		if current.GroupName != "" {
			ok, err := s.UserRepo.IsUserInGroup(userID, current.GroupName)
			if err != nil {
				return err
			}
			if !ok {
				return errors.New("user not authorized for this step")
			}
		}

		now := time.Now()

		// ================== REJECT CURRENT ================== //
		current.Status = constants.RequestRejected
		current.ApprovedAt = &now
		current.Notes = reason

		if err := tx.Save(current).Error; err != nil {
			return err
		}

		// ================== CANCEL OTHER STEPS (same level) ================== //
		for i := range steps {
			if steps[i].ID == current.ID {
				continue
			}
			if steps[i].Status != constants.RequestPending {
				continue
			}

			steps[i].Status = constants.RequestCanceled
			steps[i].Notes = "Cancelled due to rejection in same approval step"

			if err := tx.Save(&steps[i]).Error; err != nil {
				return err
			}

			if err := tx.Create(&models.ApprovalLog{
				ID:        uuid.New(),
				FlowID:    flow.ID,
				StepID:    &steps[i].ID,
				Action:    "step_cancelled",
				ByUserID:  &userID,
				Note:      "Cancelled due to rejection in same approval step",
				CreatedAt: now,
			}).Error; err != nil {
				return err
			}
		}

		// ================== UPDATE OPS REQUEST ================== //
		if err := tx.Model(&models.OpsRequest{}).Where("id = ?", flow.RequestID).Updates(map[string]interface{}{
			"status":     constants.RequestRejected.String(),
			"updated_at": now,
		}).Error; err != nil {
			return err
		}

		// ================== LOG REJECT ================== //
		if err := tx.Create(&models.ApprovalLog{
			ID:        uuid.New(),
			FlowID:    flow.ID,
			StepID:    &current.ID,
			Action:    "step_rejected",
			ByUserID:  &userID,
			Note:      reason,
			CreatedAt: now,
		}).Error; err != nil {
			return err
		}

		// ================== LOG FLOW REJECT ================== //
		if err := tx.Create(&models.ApprovalLog{
			ID:        uuid.New(),
			FlowID:    flow.ID,
			Action:    "flow_rejected",
			ByUserID:  &userID,
			Note:      reason,
			CreatedAt: now,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
