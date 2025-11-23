package services

import (
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
}

type ApprovalLogic struct {
	UserRepo        *repository.UserRepository
	RequestTypeRepo *repository.RequestTypeRepository
	LevelRepo       *repository.LevelRepository
}

func NewApprovalLogic(u *repository.UserRepository, rt *repository.RequestTypeRepository, l *repository.LevelRepository) *ApprovalLogic {
	return &ApprovalLogic{
		UserRepo:        u,
		RequestTypeRepo: rt,
		LevelRepo:       l,
	}
}

func NewApprovalService(db *gorm.DB, reqRepo *repository.OpsRequestRepository, approvalRepo *repository.ApprovalRepository) *ApprovalService {
	return &ApprovalService{
		DB:           db,
		ReqRepo:      reqRepo,
		ApprovalRepo: approvalRepo,
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

func (s *ApprovalLogic) DetermineApprovers(request *models.OpsRequest) ([]models.User, error) {
	if request == nil {
		return nil, errors.New("request nil")
	}
	rt, err := a.RequestTypeRepo.GetByID(request.RequestTypeID)
	if err != nil {
		return nil, err
	}
	requiredRank := rt.RequiredLevelRank
	// find levels with rank >= requiredRank
	levels, err := a.LevelRepo.FindByMinRank(requiredRank)
	if err != nil {
		return nil, err
	}
	if len(levels) == 0 {
		return nil, err
	}
	userIDs := make([]uuid.UUID, 0)
	for _, l := range levels {
		uls, err := a.UserRepo.FindUsersByLevel(l.ID)
		if err != nil {
			return nil, err
		}
		for _, u := range uls {
			userIDs = append(userIDs, u.ID)
		}
	}
	// dedupe & load users
	uniq := map[uuid.UUID]bool{}
	users := make([]models.User, 0)
	for _, uid := range userIDs {
		if uniq[uid] {
			continue
		}
		uniq[uid] = true
		u, err := a.UserRepo.FindByID(uid)
		if err == nil {
			users = append(users, *u)
		}
	}
	return users, nil
}
