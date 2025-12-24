package services

import (
	"backend/internal/app/constants"
	"backend/internal/app/logic"
	"backend/internal/app/models"
	"backend/internal/app/repository"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupApprovalActionTest(t *testing.T) (*ApprovalActionService, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&models.ApprovalFlow{},
		&models.ApprovalStep{},
		&models.ApprovalLog{},
		&models.OpsRequest{},
	)
	require.NoError(t, err)

	svc := NewApprovalActionService(
		db,
		repository.NewApprovalFlowRepository(db),
		repository.NewApprovalStepRepository(db),
		repository.NewOpsRequestRepository(db),
		repository.NewApprovalLogRepository(db),
		repository.NewUserRepository(db),
		logic.NewApprovalLogic(),
	)
	return svc, db
}

// ==================== ACTION SERVICE TESTING ==================== //
func TestApprove_OR_FinalizesFlow(t *testing.T) {
	svc, db := setupApprovalActionTest(t)

	userID := uuid.New()
	flowID := uuid.New()
	reqID := uuid.New()

	flow := models.ApprovalFlow{
		ID:          flowID,
		RequestID:   reqID,
		Status:      constants.RequestInReview,
		CurrentStep: 1,
	}
	require.NoError(t, db.Create(&flow).Error)

	steps := []models.ApprovalStep{
		{
			ID:         uuid.New(),
			FlowID:     flowID,
			StepNumber: 1,
			UserID:     &userID,
			Mode:       constants.ModeOR,
			Status:     constants.RequestPending,
		},
		{
			ID:         uuid.New(),
			FlowID:     flowID,
			StepNumber: 1,
			GroupName:  "finance",
			Mode:       constants.ModeOR,
			Status:     constants.RequestPending,
		},
	}
	require.NoError(t, db.Create(&steps).Error)

	err := svc.Approve(context.Background(), flowID, userID, "ok")
	require.NoError(t, err)

	var updated models.ApprovalFlow
	db.First(&updated, "id = ?", flowID)

	assert.Equal(t, constants.RequestApproved.String(), updated.Status)
}
