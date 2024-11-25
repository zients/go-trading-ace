package repositories

import (
	"database/sql"
	"fmt"
	"trading-ace/entities"
	"trading-ace/models"
)

type ITaskHistoryRepository interface {
	Create(taskHistory *entities.TaskHistory) (*entities.TaskHistory, error)
	FindByID(id int64) (*entities.TaskHistory, error)
	FindByAddressAndTaskId(address string, taskId int64) (*entities.TaskHistory, error)
	GetByAddressIncludingTask(address string) ([]*models.GetByAddressIncludingTask, error)
	Update(record *entities.TaskHistory) (*entities.TaskHistory, error)
	Delete(id int64) error
}

type TaskHistoryRepository struct {
	db *sql.DB
}

func NewTaskHistoryRepository(db *sql.DB) ITaskHistoryRepository {
	return &TaskHistoryRepository{
		db: db,
	}
}

func (r *TaskHistoryRepository) Create(taskHistory *entities.TaskHistory) (*entities.TaskHistory, error) {
	query := `
		INSERT INTO task_histories (address, task_id, reward_points, amount, completed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, address, task_id, reward_points, amount, completed_at, created_at, updated_at
	`

	var result entities.TaskHistory

	err := r.db.QueryRow(
		query,
		taskHistory.Address, taskHistory.TaskID, taskHistory.RewardPoints,
		taskHistory.Amount, taskHistory.CompletedAt,
	).Scan(
		&result.ID, &result.Address, &result.TaskID, &result.RewardPoints,
		&result.Amount, &result.CompletedAt, &result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create task record: %w", err)
	}

	return &result, nil
}

func (r *TaskHistoryRepository) FindByID(id int64) (*entities.TaskHistory, error) {
	query := `
		SELECT id, address, task_id, reward_points, amount, completed_at, created_at, updated_at
		FROM task_histories
		WHERE id = $1
	`

	var result entities.TaskHistory
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

func (r *TaskHistoryRepository) FindByAddressAndTaskId(address string, taskId int64) (*entities.TaskHistory, error) {
	query := `
		SELECT id, address, task_id, reward_points, amount, completed_at, created_at, updated_at
		FROM task_histories
		WHERE address = $1 AND task_id = $2
	`

	var result entities.TaskHistory
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

func (t *TaskHistoryRepository) GetByAddressIncludingTask(address string) ([]*models.GetByAddressIncludingTask, error) {
	query := `
		SELECT th.id, th.address, th.reward_points, th.amount, th.completed_at,
		       t.id, t.name, t.description, t.points, t.started_at, t.end_at, t.period, t.created_at, t.updated_at
		FROM task_histories th
		INNER JOIN tasks t ON th.task_id = t.id
		WHERE th.address = $1
	`

	rows, err := t.db.Query(query, address)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	defer rows.Close()

	var results []*models.GetByAddressIncludingTask
	for rows.Next() {
		task := &entities.Task{}
		taskHistory := &entities.TaskHistory{}

		// Scan columns
		err := rows.Scan(
			&taskHistory.ID, &taskHistory.Address, &taskHistory.RewardPoints, &taskHistory.Amount, &taskHistory.CompletedAt,

			&task.ID, &task.Name, &task.Description, &task.Points,
			&task.StartedAt, &task.EndAt, &task.Period,
			&task.CreatedAt, &task.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		results = append(results, &models.GetByAddressIncludingTask{
			Task:        task,
			TaskHistory: taskHistory,
		})
	}

	return results, nil
}

func (r *TaskHistoryRepository) Update(record *entities.TaskHistory) (*entities.TaskHistory, error) {
	query := `
		UPDATE task_histories
		SET address = $1, task_id = $2, reward_points = $3, amount = $4, completed_at = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING id, address, task_id, reward_points, amount, completed_at, created_at, updated_at
	`

	var result entities.TaskHistory

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

func (r *TaskHistoryRepository) Delete(id int64) error {
	query := `
		DELETE FROM task_histories
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task record: %w", err)
	}

	return nil
}
