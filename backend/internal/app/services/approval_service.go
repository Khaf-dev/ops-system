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

// ApprovalService handles approval logic (multi-level)
type ApprovalService struct {
	DB           *gorm.DB
	ReqRepo      *repository.OpsRequestRepository
	ApprovalRepo *repository.ApprovalRepository
	approverRepo repository.ApproverConfigRepository
	Logic        *logic.ApprovalLogic
	LevelRepo    *repository.LevelRepository
	UserRepo     *repository.UserRepository
}

func NewApprovalService(
	db *gorm.DB,
	reqRepo *repository.OpsRequestRepository,
	approvalRepo *repository.ApprovalRepository,
	lr *repository.LevelRepository,
	ur *repository.UserRepository,
	approverRepo repository.ApproverConfigRepository,
	logic *logic.ApprovalLogic) *ApprovalService {
	return &ApprovalService{
		DB:           db,
		ReqRepo:      reqRepo,
		ApprovalRepo: approvalRepo,
		approverRepo: approverRepo,
		Logic:        logic,
		LevelRepo:    lr,
		UserRepo:     ur,
	}
}

// HandleApproval: transactional. decision = "approved" or "rejected"
func (s *ApprovalService) HandleApproval(requestID, approverID uuid.UUID, decision, notes string, actedAt time.Time) error {
	if decision != "approved" && decision != "rejected" {
		return errors.New("invalid decision")
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		// load request with relations needed
		req, err := s.ReqRepo.GetByID(requestID, "RequestType", "Activity", "Approvals")
		if err != nil {
			return utils.ErrNotFound
		}

		// check idempotency: an approver should not approve twice
		var existing models.Approval
		if err := tx.Where("request_id = ? AND approver_id = ?", requestID, approverID).First(&existing).Error; err == nil {
			return errors.New("already acted")
		}

		// determine required rank from request type (fallback to 1)
		requiredRank := 1
		if req.RequestType != nil {
			requiredRank = req.RequestType.RequiredLevelRank
			if requiredRank <= 0 {
				requiredRank = 1
			}
		}

		// fetch approver's highest rank (you need a user_levels table)
		type rankRow struct {
			Rank int
		}
		var rr rankRow
		// safe SQL to get the highest rank of approver
		if err := tx.Raw(`
			SELECT MAX(l.rank) as rank
			FROM levels l
			JOIN user_levels ul ON ul.level_id = l.id
			WHERE ul.user_id = ?
		`, approverID).Scan(&rr).Error; err != nil {
			return err
		}
		if rr.Rank == 0 {
			return errors.New("approver has no level assigned")
		}

		// approver must have rank >=  required minimal rank to act
		// (business rule: you can allow lower rank to approve earlier steps; adapt as needed)
		// Here we require approver rank >=  requiredRank (for simplicity)
		if rr.Rank < requiredRank {
			return errors.New("insufficient approval level")
		}

		// create approval record
		a := &models.Approval{
			ID:         uuid.New(),
			RequestID:  requestID,
			ApproverID: approverID,
			Decision:   decision,
			Notes:      notes,
			CreatedAt:  actedAt,
		}
		if err := tx.Create(a).Error; err != nil {
			return err
		}

		// if decision == rejected => immediately mark request rejected
		if decision == "rejected" {
			req.Status = "rejected"
			req.ApprovedByID = &approverID
			req.UpdatedAt = time.Now()
			if err := tx.Save(req).Error; err != nil {
				return err
			}
			// log
			if err := tx.Create(&models.ActivityLog{
				ID:         uuid.New(),
				ActorID:    approverID,
				Action:     "approval_rejected",
				TargetType: "ops_request",
				TargetID:   &requestID,
				CreatedAt:  time.Now(),
			}).Error; err != nil {
				return err
			}
			return nil
		}

		// For "approved": evaluate whether this completes the approval chain.
		// Simple rule: if approver rank >= highest required rank for this request type => finalize approved.
		// (This is conservative; adapt to your exact multi-step rules)
		// Determine max required rank from request type (we used requiredRank)
		// If approver's rank is >= requiredRank => mark approved.
		if rr.Rank >= requiredRank {
			req.Status = "approved"
			req.ApprovedByID = &approverID
			req.UpdatedAt = time.Now()
			if err := tx.Save(req).Error; err != nil {
				return err
			}
			// activity log
			if err := tx.Create(&models.ActivityLog{
				ID:         uuid.New(),
				ActorID:    approverID,
				Action:     "approval_approved",
				TargetType: "ops_request",
				TargetID:   &requestID,
				CreatedAt:  time.Now(),
			}).Error; err != nil {
				return err
			}
		} else {
			// partial approval: still pending but record inserted (status stays pending)
			if err := tx.Create(&models.ActivityLog{
				ID:         uuid.New(),
				ActorID:    approverID,
				Action:     "approval_partial",
				TargetType: "ops_request",
				TargetID:   &requestID,
				CreatedAt:  time.Now(),
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *ApprovalService) ApproveOrReject(reqID uuid.UUID, userID uuid.UUID, action string, role string) error {

	// ---1. Ambil request
	req, err := s.ReqRepo.GetByID(reqID)
	if err != nil {
		return err
	}

	// ---2. Cek apakah user authorized buat approve (level sesuai)
	if req.CurrentApproverID != userID {
		return errors.New("maaf kamu bukan approver aktif")
	}

	// ---3. Ambil config approver sesuai tipe request
	cfg, err := s.ApproverRepo.GetConfigs(req.RequestTypeID)
	if err != nil {
		return err
	}

	// ---4. Cek last approver
	isLast := s.Logic.IsLastApprover(req, cfg)

	// ---5. Tentukan final status setelah action
	finalStatus, err := s.Logic.DetermineFinalStatus(action, isLast)
	if err != nil {
		return err
	}

	// ---6. Update info approver selanjutnya (jika masih lanjut)
	if action == "approver" && !isLast {
		nextUser, err := s.Logic.DetermineNextApprovers(req, cfg)
		if err != nil {
			return err
		}
		req.CurrentApproverID = nextUser.ID
		req.CurrentApprovalLevel++
	}

	// ---7. Kalau reject -selesai
	if action == "reject" {
		req.CurrentApproverID = 0
		req.CurrentApprovalLevel = 0
	}

	// ---8. Update status final
	req.Status = finalStatus

	// ---9. Save ke DB
	return s.ReqRepo.Update(req)
}
