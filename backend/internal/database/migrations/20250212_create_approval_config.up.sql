CREATE TABLE approval_configs (
    id UUID PRIMARY KEY,
    request_type_id UUID NOT NULL,
    user_id UUID NOT NULL,
    level INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT fk_approval_config_request_type FOREIGN KEY (request_type_id) REFERENCES requested_types(id) ON DELETE CASCADE,
    CONSTRAINT fk_approval_config_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- biar query makin ngebut
CREATE INDEX idx_approval_config_request_type ON approval_configs (request_type_id);
CREATE INDEX idx_approval_config_level ON approval_configs (level);