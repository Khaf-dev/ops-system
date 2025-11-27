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

}

// DetermineNextApprovers
func (l *ApprovalLogic) DetermineNextApprovers() error {

}
