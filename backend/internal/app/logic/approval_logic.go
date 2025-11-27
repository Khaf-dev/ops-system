package logic

import (
	"backend/internal/app/models"
	"errors"
)

type ApprovalLogic struct{}

/*

	Determine Approvers :
	Mendefisinikan terkait dengan logika bisnis
	Menggerakkan sistem yang tidak terkait dengan Database dan Modul lainnya

	Determine Approvers:
	requester : ngambil request dari user
	reqType : object dari table request_type
	userLevel : assign level (manager, admin, requester)
	Levels : detect object ke setiap level
	users : semua users

*/

func NewApprovalLogic() *ApprovalLogic {
	return &ApprovalLogic{}
}

func (l *ApprovalLogic) DetermineApprovers(
	requester *models.User,
	reqType *models.RequestType,
	userLevel *models.UserLevel,
	allLevels []models.Level,
	allUsers []models.User,
) ([]models.User, error) {

	if requester == nil || reqType == nil || userLevel == nil {
		return nil, errors.New("invalid request input for approval logic")
	}

	// == 1 : ambil minimum approval lecel dari request type
	minLevel := reqType.MinApprovalLevel // <== TODO (Liat di models/request_type.go)
	if minLevel == 0 {
		minLevel = userLevel.Order + 1
	}

	// == 2 : tentukan level approver
	var approvalLevels []models.Level
	for _, lvl := range allLevels {
		if lvl.Order >= minLevel {
			approvalLevels = append(approvalLevels, lvl)
		}
	}

	if len(approvalLevels) == 0 {
		return nil, errors.New("no valid approval level found")
	}

	// == 3 : Cari user yang levelnya match approval level
	var approvers []models.User
	for _, apprLvl := range approvalLevels {
		for _, usr := range allUsers {
			if usr.LevelID == apprLvl.ID {
				approvers = append(approvers, usr)
			}
		}
	}

	if len(approvers) == 0 {
		return nil, errors.New("no approver found for levels")
	}
	return approvers, nil

}

// DetermineNextApprovers
func (l *ApprovalLogic) DetermineNextApprovers(req *models.OpsRequest, approvers []models.ApproverConfig) (*models.User, error) {
	// Untuk cari position approver
	currentIdx := 1
	for i, a := range approvers {
		if a.Level == req.CurrentApprovalLevel {
			currentIdx = i
			break
		}
	}

	// Klo currentIndex -1 artinya request belum mulai approval
	if currentIdx == -1 {
		//ambil approver level 1
		for _, a := range approvers {
			if a.Level == 1 {
				return &a.User, nil
			}
		}
		return nil, errors.New("tidak menemukan approver level pertama")
	}

}
