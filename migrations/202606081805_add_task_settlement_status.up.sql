ALTER TABLE tasks
ADD COLUMN settlement_started_at TIMESTAMP NULL,
ADD COLUMN settled_at TIMESTAMP NULL;

CREATE INDEX tasks_unsettled_share_pool_idx
ON tasks (end_at, settlement_started_at)
WHERE name = 'SharePoolTask' AND settled_at IS NULL;
