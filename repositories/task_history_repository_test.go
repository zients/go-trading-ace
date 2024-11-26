package repositories

import (
	"testing"
	"time"

	"trading-ace/entities"
	"trading-ace/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTaskHistory(t *testing.T) {
	// 設置 mock 資料庫
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewTaskHistoryRepository(db)

	// 設定測試數據
	taskHistory := &entities.TaskHistory{
		Address:      "test_address",
		TaskID:       1,
		RewardPoints: 10.5,
		Amount:       100.0,
		CompletedAt:  nil,
	}

	// 設定查詢語句及返回結果
	mock.ExpectQuery(`INSERT INTO task_histories`).
		WithArgs(taskHistory.Address, taskHistory.TaskID, taskHistory.RewardPoints, taskHistory.Amount, taskHistory.CompletedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "task_id", "reward_points", "amount", "completed_at", "created_at", "updated_at"}).
			AddRow(1, taskHistory.Address, taskHistory.TaskID, taskHistory.RewardPoints, taskHistory.Amount, taskHistory.CompletedAt, time.Now(), time.Now()))

	// 呼叫 Create 函數
	createdTaskHistory, err := repo.Create(taskHistory)
	assert.NoError(t, err)
	assert.NotNil(t, createdTaskHistory)
	assert.Equal(t, taskHistory.Address, createdTaskHistory.Address)
}

func TestFindByID(t *testing.T) {
	// 設置 mock 資料庫
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewTaskHistoryRepository(db)

	// 設定測試數據
	taskHistoryID := int64(1)
	expectedTaskHistory := &entities.TaskHistory{
		ID:           taskHistoryID,
		Address:      "test_address",
		TaskID:       1,
		RewardPoints: 10.5,
		Amount:       100.0,
		CompletedAt:  nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 設定查詢語句及返回結果
	mock.ExpectQuery(`SELECT id, address, task_id, reward_points, amount, completed_at, created_at, updated_at`).
		WithArgs(taskHistoryID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "task_id", "reward_points", "amount", "completed_at", "created_at", "updated_at"}).
			AddRow(expectedTaskHistory.ID, expectedTaskHistory.Address, expectedTaskHistory.TaskID, expectedTaskHistory.RewardPoints, expectedTaskHistory.Amount, expectedTaskHistory.CompletedAt, expectedTaskHistory.CreatedAt, expectedTaskHistory.UpdatedAt))

	// 呼叫 FindByID 函數
	result, err := repo.FindByID(taskHistoryID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTaskHistory.ID, result.ID)
}

func TestFindByAddressAndTaskId(t *testing.T) {
	// 設置 mock 資料庫
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewTaskHistoryRepository(db)

	// 設定測試數據
	address := "test_address"
	taskId := int64(1)
	expectedTaskHistory := &entities.TaskHistory{
		ID:           1,
		Address:      address,
		TaskID:       taskId,
		RewardPoints: 10.5,
		Amount:       100.0,
		CompletedAt:  nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 設定查詢語句及返回結果
	mock.ExpectQuery(`SELECT id, address, task_id, reward_points, amount, completed_at, created_at, updated_at`).
		WithArgs(address, taskId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "task_id", "reward_points", "amount", "completed_at", "created_at", "updated_at"}).
			AddRow(expectedTaskHistory.ID, expectedTaskHistory.Address, expectedTaskHistory.TaskID, expectedTaskHistory.RewardPoints, expectedTaskHistory.Amount, expectedTaskHistory.CompletedAt, expectedTaskHistory.CreatedAt, expectedTaskHistory.UpdatedAt))

	// 呼叫 FindByAddressAndTaskId 函數
	result, err := repo.FindByAddressAndTaskId(address, taskId)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTaskHistory.ID, result.ID)
}

func TestGetByAddressIncludingTasks(t *testing.T) {
	// 設置 mock 資料庫
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewTaskHistoryRepository(db)

	// 設定測試數據
	address := "test_address"
	expectedResults := []*models.TaskTaskHistoryPair{
		{
			Task: &entities.Task{
				ID:          1,
				Name:        "Test Task",
				Description: "Test Description",
				Points:      100,
				StartedAt:   nil,
				EndAt:       nil,
				Period:      1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			TaskHistory: &entities.TaskHistory{
				ID:           1,
				Address:      address,
				TaskID:       1,
				RewardPoints: 10.5,
				Amount:       100.0,
				CompletedAt:  nil,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
		},
	}

	// 設定查詢語句及返回結果
	mock.ExpectQuery(`SELECT th.id, th.address, th.reward_points, th.amount, th.completed_at,`).
		WithArgs(address).
		WillReturnRows(sqlmock.NewRows([]string{"th.id", "th.address", "th.reward_points", "th.amount", "th.completed_at", "t.id", "t.name", "t.description", "t.points", "t.started_at", "t.end_at", "t.period", "t.created_at", "t.updated_at"}).
			AddRow(
				expectedResults[0].TaskHistory.ID,
				expectedResults[0].TaskHistory.Address,
				expectedResults[0].TaskHistory.RewardPoints,
				expectedResults[0].TaskHistory.Amount,
				expectedResults[0].TaskHistory.CompletedAt,
				expectedResults[0].Task.ID,
				expectedResults[0].Task.Name,
				expectedResults[0].Task.Description,
				expectedResults[0].Task.Points,
				expectedResults[0].Task.StartedAt,
				expectedResults[0].Task.EndAt,
				expectedResults[0].Task.Period,
				expectedResults[0].Task.CreatedAt,
				expectedResults[0].Task.UpdatedAt,
			))

	// 呼叫 GetByAddressIncludingTasks 函數
	results, err := repo.GetByAddressIncludingTasks(address)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, expectedResults[0].TaskHistory.Address, results[0].TaskHistory.Address)
}
