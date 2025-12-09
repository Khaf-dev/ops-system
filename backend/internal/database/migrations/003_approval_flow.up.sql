-- 003_approval_flow.up.sql
-- Approval flow, steps, logs, and approver config

CREATE TABLE IF NOT EXISTS approval_flows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID NOT NULL UNIQUE REFERENCES ops_requests(id) ON DELETE CASCADE,
    current_step INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    created_by_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_approval_flows_request_id ON approval_flows(request_id);
CREATE INDEX IF NOT EXISTS idx_approval_flows_status ON approval_flows(status);


CREATE TABLE IF NOT EXISTS approval_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flow_id UUID NOT NULL REFERENCES approval_flows(id) ON DELETE CASCADE,
    step_number INTEGER NOT NULL,
    user_id UUID,                   -- single-user approver (nullable if group)
    group_name VARCHAR(128),        -- dynamic group/role name (nullable)
    mode VARCHAR(10) NOT NULL DEFAULT 'AND', -- AND / OR
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    approved_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ensure uniqueness per flow+user (same user can't be duplicated in same flow step)
CREATE UNIQUE INDEX IF NOT EXISTS uq_approval_steps_flow_step_user
    ON approval_steps(flow_id, step_number, user_id);

CREATE INDEX IF NOT EXISTS idx_approval_steps_flow_step ON approval_steps(flow_id, step_number);
CREATE INDEX IF NOT EXISTS idx_approval_steps_status ON approval_steps(status);


CREATE TABLE IF NOT EXISTS approval_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flow_id UUID NOT NULL REFERENCES approval_flows(id) ON DELETE CASCADE,
    step_id UUID,
    action VARCHAR(64) NOT NULL, -- e.g. step_approved, step_rejected, flow_started
    by_user_id UUID,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_approval_logs_flow_id ON approval_logs(flow_id);


CREATE TABLE IF NOT EXISTS approver_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_type_id UUID NOT NULL REFERENCES request_types(id) ON DELETE CASCADE,
    level INTEGER NOT NULL, -- 1 = first step
    user_id UUID,           -- optional specific user
    group_name VARCHAR(128), -- optional group selector
    mode VARCHAR(10) NOT NULL DEFAULT 'AND', -- AND / OR
    priority INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_approver_configs_request_type ON approver_configs(request_type_id);
CREATE INDEX IF NOT EXISTS idx_approver_configs_level ON approver_configs(level);