-- 任務表
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE, 
    description TEXT,
    points BIGINT NOT NULL,
    started_at TIMESTAMP NULL,
    end_at TIMESTAMP NULL,
    is_recurring BOOLEAN DEFAULT FALSE,
    period INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 用戶完成任務記錄表
CREATE TABLE task_records (
    id SERIAL PRIMARY KEY,
    address VARCHAR(255) NOT NULL,
    task_id INT NOT NULL REFERENCES tasks(id),
    reward_points BIGINT NOT NULL,
    amount BIGINT NULL,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
