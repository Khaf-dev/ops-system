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

type ApprovalActionService struct {
	DB       *gorm.DB
	FlowRepo *repository.ApprovalFlowRepository
	StepRepo *repository.ApprovalStepRepository
	ReqRepo  *repository.OpsRequestRepository
	LogRepo  *repository.ApprovalLogRepository
	UserRepo *repository.UserRepository
	Logic    *logic.ApprovalLogic
}

func NewApprovalActionService(db *gorm.DB, flowRepo *repository.ApprovalFlowRepository,
	stepRepo *repository.ApprovalStepRepository, reqRepo *repository.OpsRequestRepository,
	logRepo *repository.ApprovalLogRepository, userRepo *repository.UserRepository, logic *logic.ApprovalLogic) *ApprovalActionService {

	return &ApprovalActionService{
		DB: db, FlowRepo: flowRepo, StepRepo: stepRepo, ReqRepo: reqRepo, LogRepo: logRepo, UserRepo: userRepo, Logic: logic,
	}
}

func (s *ApprovalActionService) Approve(flowID, userID uuid.UUID, note string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		flow, err := s.FlowRepo.GetByID(flowID)
		if err != nil {
			return utils.ErrNotFound
		}
		if !(flow.Status == "pending" || flow.Status == "in_review") {
			return errors.New("flow not in approvable state")
		}

		current, err := s.StepRepo.GetCurrentStep(flow.ID, flow.CurrentStep)
		if err != nil {
			return errors.New("no current step")
		}

		// validate
		if !s.Logic.ValidateApproverForStep(current, userID) {
			if current.GroupName != "" {
				ok, err := s.UserRepo.IsUserInGroup(userID, current.GroupName)
				if err != nil {
					return err
				}
				if !ok {
					return errors.New("user not authorized")
				}
			} else {
				return errors.New("user not authorized")
			}
		}

		now := time.Now()
		current.Status = (constants.RequestApproved)
		current.ApprovedAt = &now
		current.Notes = note
		if err := tx.Save(current).Error; err != nil {
			return err
		}

		// add log
		if err := tx.Create(&models.ApprovalLog{
			ID:        uuid.New(),
			FlowID:    flow.ID,
			StepID:    &current.ID,
			Action:    "step_approved",
			ByUserID:  &userID,
			Note:      note,
			CreatedAt: time.Now(),
		}).Error; err != nil {
			return err
		}

		// decide next
		next := s.Logic.DetermineNextStepNumber(flow)
		if next == 0 {
			// finalize
			if err := tx.Model(&models.ApprovalFlow{}).Where("id = ?", flow.ID).
				Updates(map[string]interface{}{"status": "approved", "updated_at": time.Now()}).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.OpsRequest{}).Where("id = ?", flow.RequestID).
				Updates(map[string]interface{}{"status": "approved", "approved_by_id": userID, "final_approved_at": now, "updated_at": time.Now()}).Error; err != nil {
				return err
			}
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
			Updates(map[string]interface{}{"current_step": next, "status": "in_review", "updated_at": time.Now()}).Error; err != nil {
			return err
		}

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

		// validate
		if !s.Logic.ValidateApproverForStep(current, userID) {
			if current.GroupName != "" {
				ok, err := s.UserRepo.IsUserInGroup(userID, current.GroupName)
				if err != nil {
					return err
				}
				if !ok {
					return errors.New("user not authorized")
				}
			} else {
				return errors.New("user not authorized")
			}
		}

		now := time.Now()
		current.Status = (constants.RequestRejected)
		current.ApprovedAt = &now
		current.Notes = reason
		if err := tx.Save(current).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.ApprovalFlow{}).Where("id = ?", flow.ID).
			Updates(map[string]interface{}{"status": "rejected", "updated_at": time.Now()}).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.OpsRequest{}).Where("id = ?", flow.RequestID).
			Updates(map[string]interface{}{"status": "rejected", "updated_at": time.Now()}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.ApprovalLog{
			ID:        uuid.New(),
			FlowID:    flow.ID,
			StepID:    &current.ID,
			Action:    "step_rejected",
			ByUserID:  &userID,
			Note:      reason,
			CreatedAt: time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
