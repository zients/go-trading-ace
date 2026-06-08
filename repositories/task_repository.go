package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
	"trading-ace/entities"
	"trading-ace/models"
)

type ITaskRepository interface {
	Create(ctx context.Context, task *entities.Task) (*entities.Task, error)
	FindById(ctx context.Context, id int64) (*entities.Task, error)
	FindByName(ctx context.Context, name string) (*entities.Task, error)
	GetByName(ctx context.Context, name string) ([]*entities.Task, error)
	IsExistedByName(ctx context.Context, name string) (bool, error)
	GetByAddressAndNamesIncludingTaskHistories(ctx context.Context, address string, names []string) ([]*models.TaskWithTaskHistory, error)
	ClaimDueSharePoolTask(ctx context.Context, now time.Time) (*entities.Task, error)
	MarkSettled(ctx context.Context, id int64, claimStartedAt time.Time, settledAt time.Time) error
	ReleaseSettlementClaim(ctx context.Context, id int64, claimStartedAt time.Time) error
}

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) ITaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (t *TaskRepository) Create(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	query := `
		INSERT INTO tasks (name, description, points, started_at, end_at, period, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, name, description, points, started_at, end_at, period, created_at, updated_at
	`

	var createdTask entities.Task
	err := t.db.QueryRowContext(
		ctx,
		query,
		task.Name, task.Description, task.Points,
		task.StartedAt, task.EndAt, task.Period,
	).Scan(
		&createdTask.ID, &createdTask.Name, &createdTask.Description,
		&createdTask.Points, &createdTask.StartedAt, &createdTask.EndAt,
		&createdTask.Period, &createdTask.CreatedAt, &createdTask.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return &createdTask, nil
}

func (t *TaskRepository) FindById(ctx context.Context, id int64) (*entities.Task, error) {
	query := `
		SELECT id, name, description, points, started_at, end_at, period, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	var task entities.Task
	err := t.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.Name, &task.Description, &task.Points,
		&task.StartedAt, &task.EndAt, &task.Period,
		&task.CreatedAt, &task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found: %w", err)
		}

		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func (t *TaskRepository) FindByName(ctx context.Context, name string) (*entities.Task, error) {
	query := `
		SELECT id, name, description, points, started_at, end_at, period, created_at, updated_at
		FROM tasks
		WHERE name = $1
	`

	var task entities.Task
	err := t.db.QueryRowContext(ctx, query, name).Scan(
		&task.ID, &task.Name, &task.Description, &task.Points,
		&task.StartedAt, &task.EndAt, &task.Period,
		&task.CreatedAt, &task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found: %w", err)
		}

		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func (t *TaskRepository) GetByName(ctx context.Context, name string) ([]*entities.Task, error) {
	query := `
		SELECT id, name, description, points, started_at, end_at, period, created_at, updated_at
		FROM tasks
		WHERE name = $1
	`

	rows, err := t.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	defer rows.Close()

	var tasks []*entities.Task
	for rows.Next() {
		task := &entities.Task{}
		err := rows.Scan(
			&task.ID, &task.Name, &task.Description, &task.Points,
			&task.StartedAt, &task.EndAt, &task.Period,
			&task.CreatedAt, &task.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}

	return tasks, nil
}

func (t *TaskRepository) IsExistedByName(ctx context.Context, name string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM tasks WHERE name = $1 LIMIT 1
		)
	`

	var exists bool
	err := t.db.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if task exists: %w", err)
	}

	return exists, nil
}

func (t *TaskRepository) GetByAddressAndNamesIncludingTaskHistories(ctx context.Context, address string, names []string) ([]*models.TaskWithTaskHistory, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("task names cannot be empty")
	}

	placeholders := make([]string, len(names))
	for i := range names {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
	}

	query := fmt.Sprintf(`
		SELECT t.id, t.name, t.description, t.points, t.started_at, t.end_at, t.period, t.created_at, t.updated_at,
			th.id, th.address, th.reward_points, th.amount, th.completed_at, th.created_at, th.updated_at
		FROM tasks t
		LEFT JOIN task_histories th ON t.id = th.task_id AND th.address = $1
		WHERE t.name IN (%s)
		ORDER BY t.name, t.period
	`, strings.Join(placeholders, ","))

	args := make([]interface{}, len(names)+1)
	args[0] = address
	for i, name := range names {
		args[i+1] = name
	}

	rows, err := t.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	defer rows.Close()

	var results []*models.TaskWithTaskHistory
	for rows.Next() {
		taskWithHistory := &models.TaskWithTaskHistory{}

		err := rows.Scan(
			&taskWithHistory.TaskID, &taskWithHistory.TaskName, &taskWithHistory.TaskDescription, &taskWithHistory.TaskPoints,
			&taskWithHistory.TaskStartedAt, &taskWithHistory.TaskEndAt, &taskWithHistory.TaskPeriod,
			&taskWithHistory.TaskCreatedAt, &taskWithHistory.TaskUpdatedAt,
			&taskWithHistory.TaskHistoryID, &taskWithHistory.TaskHistoryAddress, &taskWithHistory.TaskHistoryRewardPoints,
			&taskWithHistory.TaskHistoryAmount, &taskWithHistory.TaskHistoryCompletedAt, &taskWithHistory.TaskHistoryCreatedAt, &taskWithHistory.TaskHistoryUpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		results = append(results, taskWithHistory)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}

	return results, nil
}

func (t *TaskRepository) ClaimDueSharePoolTask(ctx context.Context, now time.Time) (*entities.Task, error) {
	query := `
		UPDATE tasks
		SET settlement_started_at = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = (
			SELECT id
			FROM tasks
			WHERE name = $1
				AND end_at <= $2
				AND settled_at IS NULL
				AND (
					settlement_started_at IS NULL
					OR settlement_started_at <= $2 - INTERVAL '2 hours'
				)
			ORDER BY end_at, period
			LIMIT 1
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, name, description, points, started_at, end_at, period, created_at, updated_at,
			settlement_started_at, settled_at
	`

	task := &entities.Task{}
	err := t.db.QueryRowContext(ctx, query, "SharePoolTask", now).Scan(
		&task.ID, &task.Name, &task.Description, &task.Points,
		&task.StartedAt, &task.EndAt, &task.Period,
		&task.CreatedAt, &task.UpdatedAt,
		&task.SettlementStartedAt, &task.SettledAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to claim due share pool task: %w", err)
	}

	return task, nil
}

func (t *TaskRepository) MarkSettled(ctx context.Context, id int64, claimStartedAt time.Time, settledAt time.Time) error {
	query := `
		UPDATE tasks
		SET settled_at = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND settlement_started_at = $2 AND settled_at IS NULL
	`

	result, err := t.db.ExecContext(ctx, query, id, claimStartedAt, settledAt)
	if err != nil {
		return fmt.Errorf("failed to mark task settled: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check settled task update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("settlement claim does not match task %d", id)
	}

	return nil
}

func (t *TaskRepository) ReleaseSettlementClaim(ctx context.Context, id int64, claimStartedAt time.Time) error {
	query := `
		UPDATE tasks
		SET settlement_started_at = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND settlement_started_at = $2 AND settled_at IS NULL
	`

	result, err := t.db.ExecContext(ctx, query, id, claimStartedAt)
	if err != nil {
		return fmt.Errorf("failed to release settlement claim: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check settlement claim release result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("settlement claim does not match task %d", id)
	}

	return nil
}
