package repositories

import (
	"database/sql"
	"fmt"
	"trading-ace/entities"
)

type ITaskRepository interface {
	Create(task *entities.Task) (*entities.Task, error)
	FindById(id int64) (*entities.Task, error)
	IsExistedByName(name string) (bool, error)
	Update(task *entities.Task) (*entities.Task, error)
	Delete(id int64) error
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
		INSERT INTO tasks (name, description, points, started_at, end_at, is_recurring, period, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, name, description, points, started_at, end_at, is_recurring, period, created_at, updated_at
	`

	var createdTask entities.Task
	err := t.db.QueryRow(
		query,
		task.Name, task.Description, task.Points,
		task.StartedAt, task.EndAt, task.IsRecurring, task.Period,
	).Scan(
		&createdTask.ID, &createdTask.Name, &createdTask.Description,
		&createdTask.Points, &createdTask.StartedAt, &createdTask.EndAt,
		&createdTask.IsRecurring, &createdTask.Period, &createdTask.CreatedAt, &createdTask.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return &createdTask, nil
}

func (t *TaskRepository) FindById(id int64) (*entities.Task, error) {
	query := `
		SELECT id, name, description, points, started_at, end_at, is_recurring, period, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	var task entities.Task
	err := t.db.QueryRow(query, id).Scan(
		&task.ID, &task.Name, &task.Description, &task.Points,
		&task.StartedAt, &task.EndAt, &task.IsRecurring, &task.Period,
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

func (t *TaskRepository) Update(task *entities.Task) (*entities.Task, error) {
	query := `
		UPDATE tasks
		SET name = $1, description = $2, points = $3, started_at = $4, end_at = $5, 
			is_recurring = $6, period = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING id, name, description, points, started_at, end_at, is_recurring, created_at, updated_at
	`

	var updatedTask entities.Task
	err := t.db.QueryRow(
		query,
		task.Name, task.Description, task.Points,
		task.StartedAt, task.EndAt, task.IsRecurring, task.Period,
		task.ID,
	).Scan(
		&updatedTask.ID, &updatedTask.Name, &updatedTask.Description,
		&updatedTask.Points, &updatedTask.StartedAt, &updatedTask.EndAt,
		&updatedTask.IsRecurring, &updatedTask.Period, &updatedTask.CreatedAt, &updatedTask.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return &updatedTask, nil
}

func (t *TaskRepository) Delete(id int64) error {
	query := `
		DELETE FROM tasks
		WHERE id = $1
	`
	_, err := t.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
