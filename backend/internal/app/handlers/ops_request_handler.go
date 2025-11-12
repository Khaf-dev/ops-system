package handlers

import (
	"backend/internal/app/dto"
	"backend/internal/app/models"
	"backend/internal/app/services"
	"backend/internal/app/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OpsRequestHandler struct {
	Svc *services.OpsRequestService
}

func NewOpsRequestHandler(svc *services.OpsRequestService) *OpsRequestHandler {
	return &OpsRequestHandler{Svc: svc}
}

type CreateOpsRequestInput struct {
	SiteID       *uuid.UUID `json:"site_id"`
	RequestType  string     `json:"request_type" binding:"required"`
	ActivityName string     `json:"activity_name" binding:"required"`
	LeaderName   string     `json:"leader_name"`
	RequestDate  *time.Time `json:"request_date"`
	Location     string     `json:"location"`
	Amount       float64    `json:"amount" binding:"required,gt=0"`
	Description  string     `json:"description"`
	Latitude     *float64   `json:"latitude"`
	Longitude    *float64   `json:"longitude"`
}

type UpdateOpsRequestInput struct {
	RequestType  string     `json:"request_type"`
	ActivityName string     `json:"activity_name"`
	LeaderName   string     `json:"leader_name"`
	RequestDate  *time.Time `json:"request_date"`
	Location     string     `json:"location"`
	Amount       *float64   `json:"amount"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
}

// POST /ops
func (h *OpsRequestHandler) CreateOpsRequest(c *gin.Context) {
	var input CreateOpsRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	uidStr, ok := c.Get("user_id")
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	requesterID, err := uuid.Parse(uidStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	req := &models.OpsRequest{
		RequesterID:  requesterID,
		SiteID:       input.SiteID,
		RequestType:  input.RequestType,
		ActivityName: input.ActivityName,
		LeaderName:   input.LeaderName,
		RequestDate:  input.RequestDate,
		Location:     input.Location,
		Amount:       input.Amount,
		Description:  input.Description,
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
		Status:       "pending",
	}

	if err := h.Svc.CreateOpsRequest(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	dto := mapToDTOForResponse(req)
	utils.SuccessResponse(c, http.StatusCreated, "ops request created", dto)
}

func mapToDTOForResponse(m *models.OpsRequest) dto.OpsRequestDTO {
	// reuse service mapping logic if u want -- duplication small for direct create response
	var requesterName, siteName string
	if m.Requester != nil {
		requesterName = m.Requester.Name
	}
	if m.Site != nil {
		siteName = m.Site.Name
	}
	return dto.OpsRequestDTO{
		ID:            m.ID,
		RequesterID:   m.RequesterID,
		RequesterName: requesterName,
		SiteID:        m.SiteID,
		SiteName:      siteName,
		RequestType:   m.RequestType,
		ActivityName:  m.ActivityName,
		LeaderName:    m.LeaderName,
		RequestDate:   m.RequestDate,
		Location:      m.Location,
		Amount:        m.Amount,
		Description:   m.Description,
		Status:        m.Status,
		Latitude:      m.Latitude,
		Longitude:     m.Longitude,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

// GET /ops/:id
func (h *OpsRequestHandler) GetOpsByRequestByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}
	dto, err := h.Svc.GetByIDDTO(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "not found")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "ok", dto)
}

// GET /ops -> role based: admin -> all, user -> own
func (h *OpsRequestHandler) ListOpsRequests(c *gin.Context) {
	//parse pagination
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	uidStr, _ := c.Get("user_id")
	role, _ := c.Get("role")
	var userID uuid.UUID
	if uidStr != nil {
		userID, _ = uuid.Parse(uidStr.(string))
	}

	result, err := h.Svc.List(role.(string), userID, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//build paginated response
	resp := map[string]interface{}{
		"items":  result.Items,
		"total":  result.Total,
		"limit":  limit,
		"offset": offset,
	}
	utils.SuccessResponse(c, http.StatusOK, "ok", resp)
}

// PUT /ops/:id
func (h *OpsRequestHandler) UpdateOpsRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	var input UpdateOpsRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid body")
		return
	}

	uidStr, ok := c.Get("user_id")
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, _ := uuid.Parse(uidStr.(string))
	role, _ := c.Get("role")

	updated := &models.OpsRequest{}
	if input.RequestType != "" {
		updated.RequestType = input.RequestType
	}
	if input.ActivityName != "" {
		updated.ActivityName = input.ActivityName
	}
	if input.LeaderName != "" {
		updated.LeaderName = input.LeaderName
	}
	if input.RequestDate != nil {
		updated.RequestDate = input.RequestDate
	}
	if input.Location != "" {
		updated.Location = input.Location
	}
	if input.Amount != nil {
		updated.Amount = *input.Amount
	}
	if input.Description != "" {
		updated.Description = input.Description
	}
	if input.Status != "" {
		updated.Status = input.Status
	}

	if err := h.Svc.UpdateOpsRequest(id, userID, role.(string), updated); err != nil {
		if err.Error() == "Forbidden" {
			utils.ErrorResponse(c, http.StatusForbidden, "Forbidden")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "updated", nil)
}

// DELETE /ops/:id
func (h *OpsRequestHandler) DeleteOpsRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	uidStr, ok := c.Get("user_id")
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, _ := uuid.Parse(uidStr.(string))
	role, _ := c.Get("role")

	if err := h.Svc.DeleteOpsRequest(id, userID, role.(string)); err != nil {
		if err.Error() == "Forbidden" {
			utils.ErrorResponse(c, http.StatusForbidden, "Forbidden")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "deleted", nil)
}
