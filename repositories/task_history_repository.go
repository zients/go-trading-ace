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
	GetByAddressIncludingTasks(address string) ([]*models.TaskTaskHistoryPair, error)
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

func (t *TaskHistoryRepository) GetByAddressIncludingTasks(address string) ([]*models.TaskTaskHistoryPair, error) {
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

	var results []*models.TaskTaskHistoryPair
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

		results = append(results, &models.TaskTaskHistoryPair{
			Task:        task,
			TaskHistory: taskHistory,
		})
	}

	return results, nil
}
