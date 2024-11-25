package repositories

import (
	"errors"
	"testing"
	"time"
	"trading-ace/entities"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	now := time.Now()

	task := &entities.Task{
		Name:        "Test Task",
		Description: "Test Description",
		Points:      10,
		StartedAt:   &now,
		EndAt:       &now,
		Period:      1,
	}

	// 設定 mock 查詢回傳值
	mock.ExpectQuery(`
		INSERT INTO tasks \(name, description, points, started_at, end_at, period, created_at, updated_at\)
		VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\)
		RETURNING id, name, description, points, started_at, end_at, period, created_at, updated_at
	`).
		WithArgs(task.Name, task.Description, task.Points, task.StartedAt, task.EndAt, task.Period).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "points", "started_at", "end_at", "period", "created_at", "updated_at",
		}).AddRow(1, task.Name, task.Description, task.Points, task.StartedAt, task.EndAt, task.Period, now, now))

	createdTask, err := repo.Create(task)

	assert.NoError(t, err)
	assert.NotNil(t, createdTask)
	assert.Equal(t, task.Name, createdTask.Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}

func TestFindById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	now := time.Now()

	mock.ExpectQuery(`
		SELECT id, name, description, points, started_at, end_at, period, created_at, updated_at
		FROM tasks
		WHERE id = \$1
	`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "points", "started_at", "end_at", "period", "created_at", "updated_at",
		}).AddRow(1, "Test Task", "Test Description", 10, now, now, 1, now, now))

	task, err := repo.FindById(1)

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "Test Task", task.Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}

func TestFindByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	now := time.Now()

	mock.ExpectQuery(`
		SELECT id, name, description, points, started_at, end_at, period, created_at, updated_at
		FROM tasks
		WHERE name = \$1
	`).
		WithArgs("Test Task").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "points", "started_at", "end_at", "period", "created_at", "updated_at",
		}).AddRow(1, "Test Task", "Test Description", 10, now, now, 1, now, now))

	task, err := repo.FindByName("Test Task")

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "Test Task", task.Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}

func TestGetByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	now := time.Now()

	mock.ExpectQuery(`
		SELECT id, name, description, points, started_at, end_at, period, created_at, updated_at
		FROM tasks
		WHERE name = \$1
	`).
		WithArgs("Test Task").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "points", "started_at", "end_at", "period", "created_at", "updated_at",
		}).AddRow(1, "Test Task", "Test Description", 10, now, now, 1, now, now))

	tasks, err := repo.GetByName("Test Task")

	assert.NoError(t, err)
	assert.NotNil(t, tasks)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}

func TestIsExistedByName(t *testing.T) {
	// 測試場景：Task exists
	t.Run("Task exists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(`SELECT EXISTS`).
			WithArgs("test-task").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		repo := &TaskRepository{db: db}

		result, err := repo.IsExistedByName("test-task")

		assert.NoError(t, err)
		assert.Equal(t, true, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// 測試場景：Task does not exist
	t.Run("Task does not exist", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(`SELECT EXISTS`).
			WithArgs("non-existent-task").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		repo := &TaskRepository{db: db}

		result, err := repo.IsExistedByName("non-existent-task")

		assert.NoError(t, err)
		assert.Equal(t, false, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// 測試場景：Database error
	t.Run("Database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(`SELECT EXISTS`).
			WithArgs("error-task").
			WillReturnError(errors.New("db error"))

		repo := &TaskRepository{db: db}

		result, err := repo.IsExistedByName("error-task")

		assert.Error(t, err)
		assert.Equal(t, false, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
