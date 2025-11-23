package handlers

import (
	"backend/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminRequestTypeHandler struct {
	svc services.RequestTypeService
}

func NewAdminRequestTypeHandler(svc services.RequestTypeService) *AdminRequestTypeHandler {
	return &AdminRequestTypeHandler{svc: svc}
}

func (h *AdminRequestTypeHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/", h.GetAll)
	rg.GET("/:id", h.GetByID)
	rg.POST("/", h.Create)
	rg.PUT("/:id", h.Update)
	rg.PUT("/:id/activate", h.Activate)
	rg.PUT("/:id/deactivate", h.Deactivate)
}

func (h *AdminRequestTypeHandler) GetAll(c *gin.Context) {
	items, err := h.svc.GetAll(c, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *AdminRequestTypeHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	result, err := h.svc.GetByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AdminRequestTypeHandler) Create(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rt, err := h.svc.Create(c, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rt)
}

func (h *AdminRequestTypeHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	var req struct {
		Name   string `json:"name"`
		Active bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.svc.Update(c, id, req.Name, req.Active)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"updated": true})
}

func (h *AdminRequestTypeHandler) Activate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	err = h.svc.SetActive(c, id, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"active": true})
}

func (h *AdminRequestTypeHandler) Deactivate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	err = h.svc.SetActive(c, id, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"active": false})
}
