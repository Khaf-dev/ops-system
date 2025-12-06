package database

import (
	"backend/internal/app/models"

	"gorm.io/gorm"
)

func SeedApprovalConfig(db *gorm.DB) error {
	data := []models.ApproverConfig{
		{TypeID: 1, Level: 1, UserID: 2}, // Admin
		{TypeID: 1, Level: 2, UserID: 3}, // Manager
		{TypeID: 1, Level: 3, UserID: 4}, // Finance
	}

	for _, cfg := range data {
		var exists int64
		db.Model(&models.ApproverConfig{}).
			Where("type_id = ? AND level = ?", cfg.TypeID, cfg.Level).
			Count(&exists)

		if exists == 0 {
			if err := db.Create(&cfg).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
