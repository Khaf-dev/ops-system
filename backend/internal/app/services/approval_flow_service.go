package services

// Approval Flow Service ini berfungsi sebagai Initialize (mMove to the next step)

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

type ApprovalFlowService struct {
	DB         *gorm.DB
	ReqRepo    *repository.OpsRequestRepository
	FlowRepo   *repository.ApprovalFlowRepository
	StepRepo   *repository.ApprovalStepRepository
	ConfigRepo *repository.ApproverConfigRepository
	LogRepo    *repository.ApprovalLogRepository
	Logic      *logic.ApprovalLogic
}

func NewApprovalFlowService(db *gorm.DB, reqRepo *repository.OpsRequestRepository,
	flowRepo *repository.ApprovalFlowRepository, stepRepo *repository.ApprovalStepRepository,
	cfgRepo *repository.ApproverConfigRepository, logRepo *repository.ApprovalLogRepository, logic *logic.ApprovalLogic) *ApprovalFlowService {

	return &ApprovalFlowService{
		DB: db, ReqRepo: reqRepo, FlowRepo: flowRepo, StepRepo: stepRepo,
		ConfigRepo: cfgRepo, LogRepo: logRepo, Logic: logic,
	}
}

func (s *ApprovalFlowService) StartFlow(requestID uuid.UUID, startedBy uuid.UUID) (*models.ApprovalFlow, error) {
	req, err := s.ReqRepo.GetByID(requestID)
	if err != nil {
		return nil, utils.ErrNotFound
	}
	cfgs, err := s.ConfigRepo.ListByRequestType(req.RequestTypeID)
	if err != nil {
		return nil, err
	}
	if len(cfgs) == 0 {
		return nil, errors.New("no approver configured")
	}

	steps, err := s.Logic.BuildStepsFromConfigs(cfgs)
	if err != nil {
		return nil, err
	}

	flow := &models.ApprovalFlow{
		RequestID:   requestID,
		CurrentStep: 1,
		Status:      (constants.RequestPending), // or constants.RequestPending
		CreatedByID: &startedBy,
	}

	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(flow).Error; err != nil {
			return err
		}
		for i := range steps {
			steps[i].FlowID = flow.ID
			if err := tx.Create(&steps[i]).Error; err != nil {
				return err
			}
		}
		// log
		if err := tx.Create(&models.ApprovalLog{
			ID:        uuid.New(),
			FlowID:    flow.ID,
			Action:    "flow_started",
			ByUserID:  &startedBy,
			CreatedAt: time.Now(),
		}).Error; err != nil {
			return err
		}
		// update ops_request: set current approver if first step user-based
		if len(steps) > 0 && steps[0].UserID != nil {
			if err := tx.Model(&models.OpsRequest{}).Where("id = ?", requestID).
				Updates(map[string]interface{}{"current_approver_id": *steps[0].UserID, "current_approval_level": 1, "status": "in_review", "updated_at": time.Now()}).Error; err != nil {
				return err
			}
		} else {
			// if first step is group-based: current_approver_id stays null; UI must resolve group members
			if err := tx.Model(&models.OpsRequest{}).Where("id = ?", requestID).
				Updates(map[string]interface{}{"current_approval_level": 1, "status": "in_review", "updated_at": time.Now()}).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return s.FlowRepo.GetByRequestID(requestID)
}
