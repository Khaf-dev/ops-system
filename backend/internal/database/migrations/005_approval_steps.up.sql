CREATE TABLE IF NOT EXISTS approval_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flow_id UUID NOT NULL,
    step_number INT NOT NULL,
    user_id UUID,
    group_name VARCHAR(128),
    mode VARCHAR(10) NOT NULL DEFAULT 'AND',
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    approved_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_step_flow FOREIGN KEY (flow_id)
        REFERENCES approval_flows(id) ON DELETE CASCADE,

    CONSTRAINT fk_step_user FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT uq_flow_step UNIQUE (flow_id, step_number)
);

CREATE INDEX IF NOT EXISTS idx_step_flow ON approval_steps(flow_id);
CREATE INDEX IF NOT EXISTS idx_step_user ON approval_steps(user_id);