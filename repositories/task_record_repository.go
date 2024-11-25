package repositories

import (
	"database/sql"
	"fmt"
	"trading-ace/entities"
)

type ITaskRecordRepository interface {
	Create(taskRecord *entities.TaskRecord) (*entities.TaskRecord, error)
	FindByID(id int64) (*entities.TaskRecord, error)
	FindByAddressAndTaskId(address string, taskId int64) (*entities.TaskRecord, error)
	Update(record *entities.TaskRecord) (*entities.TaskRecord, error)
	Delete(id int64) error
}

type TaskRecordRepository struct {
	db *sql.DB
}

func NewTaskRecordRepository(db *sql.DB) ITaskRecordRepository {
	return &TaskRecordRepository{
		db: db,
	}
}

func (r *TaskRecordRepository) Create(taskRecord *entities.TaskRecord) (*entities.TaskRecord, error) {
	query := `
		INSERT INTO task_records (address, task_id, reward_points, amount, completed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, address, task_id, reward_points, amount, completed_at, created_at, updated_at
	`

	var result entities.TaskRecord

	err := r.db.QueryRow(
		query,
		taskRecord.Address, taskRecord.TaskID, taskRecord.RewardPoints,
		taskRecord.Amount, taskRecord.CompletedAt,
	).Scan(
		&result.ID, &result.Address, &result.TaskID, &result.RewardPoints,
		&result.Amount, &result.CompletedAt, &result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create task record: %w", err)
	}

	return &result, nil
}

func (r *TaskRecordRepository) FindByID(id int64) (*entities.TaskRecord, error) {
	query := `
		SELECT id, address, task_id, reward_points, amount, completed_at, created_at, updated_at
		FROM task_records
		WHERE id = $1
	`

	var result entities.TaskRecord
	err := r.db.QueryRow(query, id).Scan(
		&result.ID, &result.Address, &result.TaskID, &result.RewardPoints,
		&result.Amount, &result.CompletedAt, &result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task record not found: %w", err)
		}

		return nil, fmt.Errorf("failed to find task record: %w", err)
	}

	return &result, nil
}

func (r *TaskRecordRepository) FindByAddressAndTaskId(address string, taskId int64) (*entities.TaskRecord, error) {
	query := `
		SELECT id, address, task_id, reward_points, amount, completed_at, created_at, updated_at
		FROM task_records
		WHERE address = $1 AND task_id = $2
	`

	var result entities.TaskRecord
	err := r.db.QueryRow(query, address, taskId).Scan(
		&result.ID, &result.Address, &result.TaskID, &result.RewardPoints,
		&result.Amount, &result.CompletedAt, &result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task record not found: %w", err)
		}

		return nil, fmt.Errorf("failed to find task record: %w", err)
	}

	return &result, nil
}

func (r *TaskRecordRepository) Update(record *entities.TaskRecord) (*entities.TaskRecord, error) {
	query := `
		UPDATE task_records
		SET address = $1, task_id = $2, reward_points = $3, amount = $4, completed_at = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING id, address, task_id, reward_points, amount, completed_at, created_at, updated_at
	`

	var result entities.TaskRecord

	err := r.db.QueryRow(
		query,
		record.Address, record.TaskID, record.RewardPoints,
		record.Amount, record.CompletedAt, record.ID,
	).Scan(
		&result.ID, &result.Address, &result.TaskID, &result.RewardPoints,
		&result.Amount, &result.CompletedAt, &result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update task record: %w", err)
	}

	return &result, nil
}

func (r *TaskRecordRepository) Delete(id int64) error {
	query := `
		DELETE FROM task_records
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task record: %w", err)
	}

	return nil
}
