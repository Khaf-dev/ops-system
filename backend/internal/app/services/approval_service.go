package services

import (
	"backend/internal/app/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApprovalService struct {
	DB *gorm.DB
}

const (
	StatusPending  = "pending"
	StatusApproved = "disetujui"
	StatusRejected = "ditolak"
)

func NewApprovalService(db *gorm.DB) *ApprovalService {
	return &ApprovalService{DB: db}
}

func (s *ApprovalService) HandleApproval(requestID, approverID uuid.UUID, decision, notes string, now time.Time) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var req models.OpsRequest
		if err := tx.First(&req, "id = ?", requestID).Error; err != nil {
			return errors.New("request not found")
		}

		if req.Status != StatusPending {
			return errors.New("request already processed")
		}

		approval := models.Approval{
			RequestID:  requestID,
			ApproverID: approverID,
			Decision:   decision,
			Notes:      notes,
			CreatedAt:  now,
		}

		if err := tx.Create(&approval).Error; err != nil {
			return err
		}

		newStatus := StatusRejected
		if decision == StatusApproved {
			newStatus = StatusApproved
		}

		req.Status = newStatus
		req.ApprovedByID = &approverID
		req.UpdatedAt = now

		if err := tx.Save(&req).Error; err != nil {
			return err
		}

		log := models.ActivityLog{
			ActorID:    approverID,
			Action:     "approval_" + decision,
			TargetType: "ops_request",
			TargetID:   &requestID,
			CreatedAt:  now,
		}
		return tx.Create(&log).Error
	})
}
