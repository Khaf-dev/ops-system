-- 004_fix_user_levels_and_configs.up.sql
-- Fix user_levels foreign key (ensure level_id references levels.id)
ALTER TABLE user_levels
  DROP CONSTRAINT IF EXISTS user_levels_level_id_fkey;

ALTER TABLE user_levels
  ADD CONSTRAINT fk_user_levels_level FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE CASCADE;

-- If there was an old approval_configs table using BIGINT, migrate data if exists
-- Drop legacy approval_configs_int if present (manual step recommended)