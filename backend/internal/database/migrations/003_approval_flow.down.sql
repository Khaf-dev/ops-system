-- 003_approval_flow.down.sql

DROP INDEX IF EXISTS idx_approver_configs_level;
DROP INDEX IF EXISTS idx_approver_configs_request_type;
DROP TABLE IF EXISTS approver_configs;

DROP INDEX IF EXISTS idx_approval_logs_flow_id;
DROP TABLE IF EXISTS approval_logs;

DROP INDEX IF EXISTS idx_approval_steps_status;
DROP INDEX IF EXISTS idx_approval_steps_flow_step;
DROP INDEX IF EXISTS uq_approval_steps_flow_step_user;
DROP TABLE IF EXISTS approval_steps;

DROP INDEX IF EXISTS idx_approval_flows_status;
DROP INDEX IF EXISTS idx_approval_flows_request_id;
DROP TABLE IF EXISTS approval_flows;