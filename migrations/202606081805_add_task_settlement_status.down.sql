DROP INDEX IF EXISTS tasks_unsettled_share_pool_idx;

ALTER TABLE tasks
DROP COLUMN IF EXISTS settlement_started_at,
DROP COLUMN IF EXISTS settled_at;
