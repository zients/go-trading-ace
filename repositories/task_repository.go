package repositories

import (
	"database/sql"
	"fmt"
	"trading-ace/entities"
)

type ITaskRepository interface {
	Create(task *entities.Task) (*entities.Task, error)
	FindById(id int64) (*entities.Task, error)
	FindByName(name string) (*entities.Task, error)
	GetByName(name string) ([]*entities.Task, error)
	IsExistedByName(name string) (bool, error)
}

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) ITaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (t *TaskRepository) Create(task *entities.Task) (*entities.Task, error) {
	query := `
		INSERT INTO tasks (name, description, points, started_at, end_at, period, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, name, description, points, started_at, end_at, period, created_at, updated_at
	`

	var createdTask entities.Task
	err := t.db.QueryRow(
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

func (t *TaskRepository) FindById(id int64) (*entities.Task, error) {
	query := `
		SELECT id, name, description, points, started_at, end_at, period, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	var task entities.Task
	err := t.db.QueryRow(query, id).Scan(
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

func (t *TaskRepository) FindByName(name string) (*entities.Task, error) {
	query := `
		SELECT id, name, description, points, started_at, end_at, period, created_at, updated_at
		FROM tasks
		WHERE name = $1
	`

	var task entities.Task
	err := t.db.QueryRow(query, name).Scan(
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

func (t *TaskRepository) GetByName(name string) ([]*entities.Task, error) {
	query := `
		SELECT id, name, description, points, started_at, end_at, period, created_at, updated_at
		FROM tasks
		WHERE name = $1
	`

	rows, err := t.db.Query(query, name)
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

	return tasks, nil
}

func (t *TaskRepository) IsExistedByName(name string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM tasks WHERE name = $1 LIMIT 1
		)
	`

	var exists bool
	err := t.db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if task exists: %w", err)
	}

	return exists, nil
}
