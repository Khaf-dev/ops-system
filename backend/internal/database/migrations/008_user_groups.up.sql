CREATE TABLE IF NOT EXISTS user_groups (
    user_id UUID NOT NULL,
    group_name VARCHAR(128) NOT NULL,

    CONSTRAINT fk_ug_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_ug_user ON user_groups(user_id);