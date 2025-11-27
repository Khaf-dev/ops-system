package handlers

import (
	"backend/internal/app/models"
	"backend/internal/app/services"
	"backend/internal/app/utils"
	"errors"
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
	SiteID        uuid.UUID  `json:"site_id" binding:"required"`
	RequestTypeID uuid.UUID  `json:"request_type_id" binding:"required"`
	ActivityID    uuid.UUID  `json:"activity_id" binding:"required"`
	LeaderName    string     `json:"leader_name"`
	RequestDate   *time.Time `json:"request_date" binding:"required"`
	Location      string     `json:"location"`
	Amount        float64    `json:"amount" binding:"required,gte=0"`
	Description   string     `json:"description"`
	Latitude      *float64   `json:"latitude"`
	Longitude     *float64   `json:"longitude"`
}

type UpdateOpsRequestInput struct {
	LeaderName    *string    `json:"leader_name"`
	RequestDate   *time.Time `json:"request_date"`
	Location      *string    `json:"location"`
	Amount        *float64   `json:"amount"`
	Description   *string    `json:"description"`
	Status        *string    `json:"status"` // only admin can set status
	SiteID        *uuid.UUID `json:"site_id"`
	RequestTypeID *uuid.UUID `json:"request_type_id"`
	ActivityID    *uuid.UUID `json:"activity_id"`
}

// POST /ops
func (h *OpsRequestHandler) CreateOpsRequest(c *gin.Context) {
	var input CreateOpsRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body: "+err.Error())
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
		RequesterID:   requesterID,
		SiteID:        input.SiteID,
		RequestTypeID: input.RequestTypeID,
		ActivityID:    input.ActivityID,
		LeaderName:    input.LeaderName,
		RequestDate:   input.RequestDate,
		Location:      input.Location,
		Amount:        input.Amount,
		Description:   input.Description,
		Latitude:      input.Latitude,
		Longitude:     input.Longitude,
	}

	dtoCreated, err := h.Svc.CreateOpsRequest(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "ops request created", dtoCreated)
}

// GET /ops/:id
func (h *OpsRequestHandler) GetOpsByRequestByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}
	dtoObj, err := h.Svc.GetByIDDTO(id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "ok", dtoObj)
}

// GET /ops -> role based: admin -> all, user -> own
func (h *OpsRequestHandler) ListOpsRequest(c *gin.Context) {
	//parse pagination
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	uidStr, _ := c.Get("user_id")
	roleStr, _ := c.Get("role")
	role := ""
	if r, ok := roleStr.(string); ok {
		role = r
	}

	var userID uuid.UUID
	if uidStr != nil {
		userID, _ = uuid.Parse(uidStr.(string))
	}

	result, err := h.Svc.List(role, userID, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

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
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}

	uidStr, ok := c.Get("user_id")
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, _ := uuid.Parse(uidStr.(string))
	role, _ := c.Get("role")

	//build updates map only with provided fields
	updates := make(map[string]interface{})
	if input.LeaderName != nil {
		updates["leader_name"] = *input.LeaderName
	}
	if input.RequestDate != nil {
		updates["request_date"] = input.RequestDate
	}
	if input.Location != nil {
		updates["location"] = *input.Location
	}
	if input.Amount != nil {
		updates["amount"] = *input.Amount
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	//admin only ypdate candidates
	if input.SiteID != nil {
		updates["site_id"] = *input.SiteID
	}
	if input.RequestTypeID != nil {
		updates["request_type_id"] = *input.RequestTypeID
	}
	if input.ActivityID != nil {
		updates["activity_id"] = *input.ActivityID
	}

	if len(updates) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "no fields to update")
		return
	}

	if err := h.Svc.UpdateOpsRequest(id, userID, role.(string), updates); err != nil {
		if errors.Is(err, utils.ErrForbidden) {
			utils.ErrorResponse(c, http.StatusForbidden, "forbidden")
			return
		}
		if errors.Is(err, utils.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "ok", nil)
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
		if errors.Is(err, utils.ErrForbidden) {
			utils.ErrorResponse(c, http.StatusForbidden, "forbidden")
			return
		}
		if errors.Is(err, utils.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "ok", nil)
}
