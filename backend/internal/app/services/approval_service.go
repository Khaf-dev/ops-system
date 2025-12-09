package services

import (
	"backend/internal/app/logic"
	"backend/internal/app/models"
	"backend/internal/app/repository"
	"backend/internal/app/utils"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ApprovalService orchestrates flow + steps + logs
// Approval Service ini fungsinya adalah User action : approve/reject
type ApprovalService struct {
	DB         *gorm.DB
	ReqRepo    *repository.OpsRequestRepository
	FlowRepo   *repository.ApprovalFlowRepository
	ConfigRepo *repository.ApproverConfigRepository
	Logic      *logic.ApprovalLogic
	UserRepo   *repository.UserRepository // used for group checks if needed
	LogRepo    *repository.ApprovalLogRepository
	StepRepo   *repository.ApprovalStepRepository
}

func NewApprovalService(
	db *gorm.DB,
	reqRepo *repository.OpsRequestRepository,
	flowRepo *repository.ApprovalFlowRepository,
	cfgRepo *repository.ApproverConfigRepository,
	stepRepo *repository.ApprovalStepRepository,
	logRepo *repository.ApprovalLogRepository,
	userRepo *repository.UserRepository,
	logic *logic.ApprovalLogic,
) *ApprovalService {
	return &ApprovalService{
		DB:         db,
		ReqRepo:    reqRepo,
		FlowRepo:   flowRepo,
		ConfigRepo: cfgRepo,
		StepRepo:   stepRepo,
		LogRepo:    logRepo,
		UserRepo:   userRepo,
		Logic:      logic,
	}
}

// StartFlow build approval-steps from configs and persist flow + steps
func (s *ApprovalService) StartFlow(requestID uuid.UUID, startedBy uuid.UUID) (*models.ApprovalFlow, error) {
	// load request
	req, err := s.ReqRepo.GetByID(requestID, "RequestType")
	if err != nil {
		return nil, utils.ErrNotFound
	}
	if req == nil {
		return nil, utils.ErrNotFound
	}
	if req.RequestTypeID == uuid.Nil {
		return nil, errors.New("request has no request_type")
	}

	// get configs
	cfgs, err := s.ConfigRepo.ListByRequestType(req.RequestTypeID)
	if err != nil {
		return nil, err
	}
	if len(cfgs) == 0 {
		return nil, errors.New("no approver configured for this request type")
	}

	steps, err := s.Logic.BuildStepsFromConfigs(cfgs)
	if err != nil {
		return nil, err
	}

	flow := &models.ApprovalFlow{
		RequestID:   requestID,
		CurrentStep: 0, // not started yet
		Status:      "pending",
		CreatedByID: &startedBy,
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(flow).Error; err != nil {
			return err
		}
		// assign FlowID to steps and save
		for i := range steps {
			steps[i].FlowID = flow.ID
			if err := tx.Create(&steps[i]).Error; err != nil {
				return err
			}
		}
		// set flow current = 1 to start first step
		flow.CurrentStep = 1
		if err := tx.Save(flow).Error; err != nil {
			return err
		}
		// log
		if err := tx.Create(&models.ApprovalLog{
			FlowID:   flow.ID,
			Action:   "flow_started",
			ByUserID: &startedBy,
		}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// load flow with steps
	return s.FlowRepo.GetByID(flow.ID)
}

// ApproveStep: user approves current step
func (s *ApprovalService) ApproveStep(flowID, userID uuid.UUID, note string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		flow, err := s.FlowRepo.GetByID(flowID)
		if err != nil {
			return utils.ErrNotFound
		}
		if flow.Status != "pending" && flow.Status != "in_review" {
			return errors.New("flow not in approvable state")
		}
		// find current step
		var current *models.ApprovalStep
		for i := range flow.Steps {
			if flow.Steps[i].StepNumber == flow.CurrentStep {
				current = &flow.Steps[i]
				break
			}
		}
		if current == nil {
			return errors.New("no current step")
		}
		// validate approver (user matches or group membership -- here we check user mmatch, group check done  via UserRepo if needed)
		if !s.Logic.ValidateApproverForStep(*current, userID) {
			// if step is group-based we should validate membership; try group validation (UserRepo must implement HasGroup)
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
		// update step
		now := time.Now()
		current.Status = "approved"
		current.ApprovedAt = &now
		current.Notes = note
		if err := tx.Save(current).Error; err != nil {
			return err
		}
		// create log
		if err := tx.Create(&models.ApprovalLog{
			FlowID:   flow.ID,
			StepID:   &current.ID,
			Action:   "step_approved",
			ByUserID: &userID,
			Note:     note,
		}).Error; err != nil {
			return err
		}

		// decide next step
		next := s.Logic.DetermineNextStepNumber(flow)
		if next == 0 {
			// last -> finalize
			flow.Status = "approved"
			flow.CurrentStep = len(flow.Steps)
			if err := tx.Save(flow).Error; err != nil {
				return err
			}
			// write ops_request: mark approved + aproved_by
			req, err := s.ReqRepo.GetByID(flow.RequestID, "")
			if err != nil {
				return err
			}
			req.Status = "approved"
			req.ApprovedByID = &userID
			req.UpdatedAt = time.Now()
			if err := tx.Save(req).Error; err != nil {
				return err
			}
			// log
			if err := tx.Create(&models.ApprovalLog{
				FlowID:   flow.ID,
				Action:   "flow_approved",
				ByUserID: &userID,
				Note:     "",
			}).Error; err != nil {
				return err
			}
			return nil
		}
		// move to next step
		flow.CurrentStep = next
		flow.Status = "in_review"
		if err := tx.Save(flow).Error; err != nil {
			return err
		}
		// log transition
		if err := tx.Create(&models.ApprovalLog{
			FlowID:   flow.ID,
			Action:   "moved_to_next_step",
			ByUserID: &userID,
			Note:     "",
		}).Error; err != nil {
			return err
		}
		return nil
	})
}

// RejectStep: user rejects flow (terminal)
func (s *ApprovalService) RejectStep(flowID, userID uuid.UUID, reason string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		flow, err := s.FlowRepo.GetByID(flowID)
		if err != nil {
			return utils.ErrNotFound
		}
		// find current step
		var current *models.ApprovalStep
		for i := range flow.Steps {
			if flow.Steps[i].StepNumber == flow.CurrentStep {
				current = &flow.Steps[i]
				break
			}
		}
		if current == nil {
			return errors.New("no current step")
		}
		// valodate approver
		if !s.Logic.ValidateApproverForStep(*current, userID) {
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
		current.Status = "rejected"
		current.ApprovedAt = &now
		current.Notes = reason
		if err := tx.Save(current).Error; err != nil {
			return err
		}
		flow.Status = "rejected"
		if err := tx.Save(flow).Error; err != nil {
			return err
		}
		// update ops_request
		req, err := s.ReqRepo.GetByID(flow.RequestID, "")
		if err != nil {
			return err
		}
		req.Status = "rejected"
		req.UpdatedAt = time.Now()
		if err := tx.Save(req).Error; err != nil {
			return err
		}
		// log
		if err := tx.Create(&models.ApprovalLog{
			FlowID:   flow.ID,
			StepID:   &current.ID,
			Action:   "step_rejected",
			ByUserID: &userID,
			Note:     reason,
		}).Error; err != nil {
			return err
		}
		return nil
	})
}
