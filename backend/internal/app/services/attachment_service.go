package services

import (
	"backend/internal/app/models"
	"backend/internal/app/repository"
	"backend/internal/app/utils"
	"errors"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttachmentService struct {
	Repo *repository.AttachmentRepository
}

func NewAttachmentService(repo *repository.AttachmentRepository) *AttachmentService {
	return &AttachmentService{Repo: repo}
}

func (s *AttachmentService) Upload(c *gin.Context, requestID uuid.UUID, file *multipart.FileHeader) (*models.Attachment, error) {
	if file.Size > 15*1024*1024 {
		return nil, errors.New("file to large")
	}

	ext := filepath.Ext(file.Filename)
	mime := file.Header.Get("Content-Type")

	checksum, bytes, err := utils.CalculateCheckSum(file)
	if err != nil {
		return nil, err
	}

	// upload to s3 (fake path, nanti bisa diganti dengan real s3 client)
	url, err := utils.UploadToS3(bytes, file.Filename, mime)
	if err != nil {
		return nil, err
	}

	a := &models.Attachment{
		RequestID:  requestID,
		FileURL:    url,
		FileType:   mime,
		FileExt:    ext,
		MimeType:   mime,
		FileSize:   file.Size,
		Checksum:   checksum,
		UploadedAt: time.Now(),
	}

	if err := s.Repo.Create(a); err != nil {
		return nil, err
	}

	return a, nil
}
