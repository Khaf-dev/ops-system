CREATE TABLE IF NOT EXISTS approver_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_type_id UUID NOT NULL,
    level INT NOT NULL,
    user_id UUID,
    group_name VARCHAR(128),
    mode VARCHAR(10) NOT NULL DEFAULT 'AND',
    priority INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_cfg_request_type FOREIGN KEY (request_type_id)
        REFERENCES request_types(id),

    CONSTRAINT fk_cfg_user FOREIGN KEY (user_id)
        REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_cfg_reqtype ON approver_configs(request_type_id);
CREATE INDEX IF NOT EXISTS idx_cfg_level ON approver_configs(level);