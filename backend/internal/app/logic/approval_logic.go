package logic

import (
	"backend/internal/app/models"
	"errors"
	"sort"

	"github.com/google/uuid"
)

type ApprovalLogic struct{}

func NewApprovalLogic() *ApprovalLogic {
	return &ApprovalLogic{}
}

// MAIN LOGIC menentukan approver selanjutnya
func (l *ApprovalLogic) DetermineNextApprover(
	req *models.OpsRequest,
	cfg []models.ApproverConfig,
) (uuid.UUID, bool, error) {

	if len(cfg) == 0 {
		return uuid.Nil, false, errors.New("approved config is null")
	}

	// sort ascending berdasarkan level
	sort.Slice(cfg, func(i, j int) bool {
		return cfg[i].Level < cfg[j].Level
	})

	// cari posisi idx sekarang
	currentIdx := -1
	for i, c := range cfg {
		if c.Level == req.CurrentApprovalLevel {
			currentIdx = i
		}
	}

	// belum mulai approval -> cari level paling kecil
	if currentIdx == -1 {
		first := cfg[0] // karena sudah sorted
		if first.UserID == uuid.Nil {
			return uuid.Nil, false, errors.New("config level pertama tidak punya UserID")
		}
		return first.UserID, len(cfg) == 1, nil
	}

	// sudah paling akhir? berarti selesai
	if currentIdx == len(cfg)-1 {
		return uuid.Nil, true, nil
	}

	next := cfg[currentIdx+1]
	if next.UserID == uuid.Nil {
		return uuid.Nil, false, errors.New("config level berikut tanpa UserID")
	}

	isLast := currentIdx+1 == len(cfg)-1

	return next.UserID, isLast, nil
}

// FINAL STATUS
func (l *ApprovalLogic) DetermineFinalStatus(action string, isLast bool) (string, error) {
	switch action {
	case "approver":
		if isLast {
			return "APPROVED", nil
		}
		return "IN_REVIEW", nil
	case "reject":
		return "REJECTED", nil
	default:
		return "", errors.New("invalid action")
	}
}
