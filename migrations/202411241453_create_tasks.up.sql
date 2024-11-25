-- tasks
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL, 
    description TEXT,
    points DECIMAL NOT NULL,
    started_at TIMESTAMP NULL,
    end_at TIMESTAMP NULL,
    period INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- user's task completed records
CREATE TABLE task_records (
    id SERIAL PRIMARY KEY,
    address VARCHAR(255) NOT NULL,
    task_id INT NOT NULL REFERENCES tasks(id),
    reward_points DECIMAL NOT NULL,
    amount DECIMAL NOT NULL,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
