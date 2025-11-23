package handlers

import (
	"backend/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttachmentHandler struct {
	Svc *services.AttachmentService
}

func NewAttachmentHandler(svc *services.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{Svc: svc}
}

func (h *AttachmentHandler) Upload(c *gin.Context) {
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "no file"})
		return
	}

	att, err := h.Svc.Upload(c, requestID, file)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, att)
}
