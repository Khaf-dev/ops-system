package services

import (
	"backend/internal/app/constants"
	"backend/internal/app/logic"
	"backend/internal/app/models"
	"backend/internal/app/repository"
	"backend/internal/app/utils"
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
func (s *ApprovalActionService) Approve(flowID, userID uuid.UUID, note string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		flow, err := s.FlowRepo.GetByID(flowID)
		if err != nil {
			return utils.ErrNotFound
		}
		if !(flow.Status == constants.RequestStatus(constants.RequestPending.String()) || flow.Status == constants.RequestStatus(constants.RequestInReview.String())) {
			return errors.New("flow not in approvable state")
		}

		// current step
		current, err := s.StepRepo.GetCurrentStep(flow.ID, flow.CurrentStep)
		if err != nil {
			return errors.New("no current step")
		}

		// validate approver
		if !s.Logic.ValidateApproverForStep(current, userID) {
			if current.GroupName != "" {
				ok, err := s.UserRepo.IsUserInGroup(userID, current.GroupName)
				if err != nil {
					return err
				}
				if !ok {
					return errors.New("user not authorized for this step")
				}
			} else {
				return errors.New("user not authorized for this step")
			}
		}

		now := time.Now()
		current.Status = constants.RequestApproved
		current.ApprovedAt = &now
		current.Notes = note
		if err := tx.Save(current).Error; err != nil {
			return err
		}

		// log step approved
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

		// determine next step
		next := s.Logic.DetermineNextStepNumber(flow)
		if next == 0 {
			// finalize
			if err := tx.Model(&models.ApprovalFlow{}).Where("id = ?", flow.ID).
				Updates(map[string]interface{}{"status": constants.RequestApproved.String(), "updated_at": time.Now()}).Error; err != nil {
				return err
			}
			// update ops_request
			if err := tx.Model(&models.OpsRequest{}).Where("id = ?", flow.RequestID).
				Updates(map[string]interface{}{
					"status":            constants.RequestApproved.String(),
					"approved_by_id":    userID,
					"final_approved_at": now,
					"updated_at":        time.Now(),
				}).Error; err != nil {
				return err
			}
			// log flow approved
			if err := tx.Create(&models.ApprovalLog{
				ID:        uuid.New(),
				FlowID:    flow.ID,
				Action:    "flow_approved",
				ByUserID:  &userID,
				CreatedAt: time.Now(),
			}).Error; err != nil {
				return err
			}
			return nil
		}

		// move to next step
		if err := tx.Model(&models.ApprovalFlow{}).Where("id = ?", flow.ID).
			Updates(map[string]interface{}{"current_step": next, "status": constants.RequestInReview.String(), "updated_at": time.Now()}).Error; err != nil {
			return err
		}

		// set ops_request.current_approver_id if next step has explicit user
		nextStep, err := s.StepRepo.GetCurrentStep(flow.ID, next)
		if err == nil && nextStep != nil && nextStep.UserID != nil {
			if err := tx.Model(&models.OpsRequest{}).Where("id = ?", flow.RequestID).
				Updates(map[string]interface{}{"current_approver_id": *nextStep.UserID, "current_approval_level": next, "updated_at": time.Now()}).Error; err != nil {
				return err
			}
		} else {
			// group-based next step: clear current_approver_id, still increase level
			if err := tx.Model(&models.OpsRequest{}).Where("id = ?", flow.RequestID).
				Updates(map[string]interface{}{"current_approver_id": nil, "current_approval_level": next, "updated_at": time.Now()}).Error; err != nil {
				return err
			}
		}

		// log move
		if err := tx.Create(&models.ApprovalLog{
			ID:        uuid.New(),
			FlowID:    flow.ID,
			Action:    "moved_to_next_step",
			ByUserID:  &userID,
			Note:      "",
			CreatedAt: time.Now(),
		}).Error; err != nil {
			return err
		}
		return nil
	})
}

// Reject: user rejects current step (terminal)
func (s *ApprovalActionService) Reject(flowID, userID uuid.UUID, reason string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		flow, err := s.FlowRepo.GetByID(flowID)
		if err != nil {
			return utils.ErrNotFound
		}

		current, err := s.StepRepo.GetCurrentStep(flow.ID, flow.CurrentStep)
		if err != nil {
			return errors.New("no current step")
		}

		// validate approver
		if !s.Logic.ValidateApproverForStep(current, userID) {
			if current.GroupName != "" {
				ok, err := s.UserRepo.IsUserInGroup(userID, current.GroupName)
				if err != nil {
					return err
				}
				if !ok {
					return errors.New("user not authorized for this step")
				}
			} else {
				return errors.New("user not authorized for this step")
			}
		}

		now := time.Now()
		current.Status = constants.RequestRejected
		current.ApprovedAt = &now
		current.Notes = reason
		if err := tx.Save(current).Error; err != nil {
			return err
		}

		// update flow + ops_request
		if err := tx.Model(&models.ApprovalFlow{}).Where("id = ?", flow.ID).
			Updates(map[string]interface{}{"status": constants.RequestRejected.String(), "updated_at": time.Now()}).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.OpsRequest{}).Where("id = ?", flow.RequestID).
			Updates(map[string]interface{}{"status": constants.RequestRejected.String(), "updated_at": time.Now()}).Error; err != nil {
			return err
		}

		// log reject
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

		return nil
	})
}
