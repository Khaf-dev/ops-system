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

type ApprovalFlowService struct {
	DB         *gorm.DB
	ReqRepo    *repository.OpsRequestRepository
	FlowRepo   *repository.ApprovalFlowRepository
	StepRepo   *repository.ApprovalStepRepository
	ConfigRepo *repository.ApproverConfigRepository
	Logic      *logic.ApprovalLogic
}

func NewApprovalFlowService(
	db *gorm.DB,
	reqRepo *repository.OpsRequestRepository,
	flowRepo *repository.ApprovalFlowRepository,
	stepRepo *repository.ApprovalStepRepository,
	configRepo *repository.ApproverConfigRepository,
	logic *logic.ApprovalLogic,
) *ApprovalFlowService {
	return &ApprovalFlowService{
		DB:         db,
		ReqRepo:    reqRepo,
		FlowRepo:   flowRepo,
		StepRepo:   stepRepo,
		ConfigRepo: configRepo,
		Logic:      logic,
	}
}

// StartFlow: build steps from approver configs and persist flow+steps
// Returns created flow with steps preloaded
func (s *ApprovalFlowService) StartFlow(requestID, startedBy uuid.UUID) (*models.ApprovalFlow, error) {
	// validate request exists
	req, err := s.ReqRepo.GetByID(requestID, "RequestType")
	if err != nil {
		return nil, utils.ErrNotFound
	}

	// get configs
	cfgs, err := s.ConfigRepo.ListByRequestType(req.RequestTypeID)
	if err != nil {
		return nil, err
	}
	if len(cfgs) == 0 {
		return nil, errors.New("no approver configured for this request type")
	}

	// build steps
	steps, err := s.Logic.BuildStepsFromConfigs(cfgs)
	if err != nil {
		return nil, err
	}

	flow := &models.ApprovalFlow{
		RequestID:   requestID,
		CurrentStep: 0, // will set to 1 inside tx
		Status:      constants.RequestStatus(constants.RequestPending.String()),
		CreatedByID: &startedBy,
	}

	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		// create flow
		if err := tx.Create(flow).Error; err != nil {
			return err
		}
		// create steps and attach flow id
		for i := range steps {
			steps[i].FlowID = flow.ID
			// ensure created_at is set
			if steps[i].CreatedAt.IsZero() {
				steps[i].CreatedAt = time.Now()
			}
			if err := tx.Create(&steps[i]).Error; err != nil {
				return err
			}
		}

		// set currentstep to 1 and status to in_review
		flow.CurrentStep = 1
		flow.Status = constants.RequestStatus(constants.RequestInReview.String())
		flow.UpdatedAt = time.Now()
		if err := tx.Save(flow).Error; err != nil {
			return err
		}

		// update ops_request current approver if first step has user
		if len(steps) > 0 && steps[0].UserID != nil {
			if err := tx.Model(&models.OpsRequest{}).Where("id = ?", requestID).
				Updates(map[string]interface{}{
					"current_approver_level": 1,
					"status":                 constants.RequestInReview.String(),
					"updated_at":             time.Now(),
				}).Error; err != nil {
				return err
			}
		}
		// log start
		if err := tx.Create(&models.ApprovalLog{
			ID:        uuid.New(),
			FlowID:    flow.ID,
			Action:    "flow_started",
			ByUserID:  &startedBy,
			CreatedAt: time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return s.FlowRepo.GetByRequestID(requestID)
}

// cancel flow : admin or owner can cancel an active flow
func (s *ApprovalFlowService) CancelFlow(flowID, actor uuid.UUID) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		flow, err := s.FlowRepo.GetByID(flowID)
		if err != nil {
			return utils.ErrNotFound
		}
		if flow.Status == constants.RequestStatus(constants.RequestApproved.String()) || flow.Status == constants.RequestStatus(constants.RequestRejected.String()) {
			return errors.New("flow already finalized")
		}
		if err := tx.Model(&models.ApprovalFlow{}).Where("id = ?", flow.ID).
			Updates(map[string]interface{}{
				"status":     constants.RequestCanceled.String(),
				"updated_at": time.Now(),
			}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.ApprovalLog{
			ID:        uuid.New(),
			FlowID:    flow.ID,
			Action:    "flow_cancelled",
			ByUserID:  &actor,
			CreatedAt: time.Now(),
		}).Error; err != nil {
			return err
		}
		// set ops_request status to cancelled
		if err := tx.Model(&models.OpsRequest{}).Where("id = ?", flow.RequestID).
			Updates(map[string]interface{}{
				"status":     constants.RequestCanceled.String(),
				"updated_at": time.Now(),
			}).Error; err != nil {
			return err
		}
		return nil
	})
}
