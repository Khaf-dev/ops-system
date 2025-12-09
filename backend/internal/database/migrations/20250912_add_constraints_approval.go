package migrations

import "gorm.io/gorm"

func Up_20250912_add_constraints_approval(db *gorm.DB) error {
	// 1. constraints for approval_config
	if err := db.Exec(`
		ALTER TABLE approver_configs
		ADD CONSTRAINT approver_config_user_or_group_check
		CHECK (
			(user_id IS NOT NULL) OR (group_name IS NOT NULL)
	);
	`).Error; err != nil {
		// ignore if already exist
	}

	// 2. Constraints for approval_steps
	if err := db.Exec(`
		ALTER TABLE approval_steps
		ADD CONSTRAINT approval_step_user_or_group_check
		CHECK (
			(user_id IS NOT NULL) OR (group_name IS NOT NULL)
	);
	`).Error; err != nil {
		// ignore if already exists
	}

	return nil
}

func Down_20250912_add_constraints_approval(db *gorm.DB) error {
	db.Exec(`
	ALTER TABLE approver_configs DROP CONSTRAINT IF EXISTS approver_config_user_or_group_check;`)
	db.Exec(`
	ALTER TABLE approval_steps DROP CONSTRAINT IF EXISTS approval_step_user_or_group_check;`)
	return nil
}
