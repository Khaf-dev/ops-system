-- LEVELS (master)
CREATE TABLE IF NOT EXISTS levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    rank INTEGER NOT NULL UNIQUE, -- higher = higher authority
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_levels_rank ON levels(rank);

--USER_LEVELS (mapping)
CREATE TABLE IF NOT EXISTS user_levels (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    level_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, level_id),
    assigned_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_user_levels_user_id ON user_levels(user_id);
CREATE INDEX IF NOT EXISTS idx_user_levels_level_id ON user_levels(level_id);

-- Extend request_types: required approval rank (nullable => default to 1)
ALTER TABLE request_types
    ADD COLUMN IF NOT EXISTS required_level_rank INTEGER DEFAULT 1;

-- Extend attachments: file metadata
ALTER TABLE attachments 
    ADD COLUMN IF NOT EXISTS file_size BIGINT,
    ADD COLUMN IF NOT EXISTS file_name VARCHAR(512),
    ADD COLUMN IF NOT EXISTS file_ext VARCHAR(32),
    ADD COLUMN IF NOT EXISTS mime_type VARCHAR(128),
    ADD COLUMN IF NOT EXISTS checksum VARCHAR(128);

-- Ensure aja untuk ops_requests.request_date is timestamptz (if you wanted timestamp instead of date)
ALTER TABLE ops_requests
    ALTER COLUMN request_date TYPE TIMESTAMPTZ USING request_date::timestamptz;
