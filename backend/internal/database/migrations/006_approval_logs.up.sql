CREATE TABLE IF NOT EXISTS approval_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flow_id UUID,
    step_id UUID,
    action VARCHAR(64) NOT NULL,
    by_user_id UUID,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_log_flow FOREIGN KEY (flow_id)
        REFERENCES approval_flows(id) ON DELETE CASCADE,

    CONSTRAINT fk_log_step FOREIGN KEY (step_id)
        REFERENCES approval_steps(id),

    CONSTRAINT fk_log_user FOREIGN KEY (by_user_id)
        REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_log_flow ON approval_logs(flow_id);
CREATE INDEX IF NOT EXISTS idx_log_step ON approval_logs(step_id);