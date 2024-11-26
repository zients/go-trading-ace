package services

import (
	"testing"
	"time"
	"trading-ace/config"
	"trading-ace/entities"
	"trading-ace/mocks"
	"trading-ace/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPointHistories(t *testing.T) {
	cfg := &config.Config{}
	loggerMock := &mocks.MockLogger{}
	taskHistoryRepoMock := &mocks.MockTaskHistoryRepository{}
	taskRepoMock := &mocks.MockTaskRepository{}
	redisHelperMock := &mocks.MockRedisHelper{}

	now := time.Now()
	// 模擬 taskHistoryRepo 的行為
	taskHistoryMock := []*models.TaskTaskHistoryPair{
		{
			TaskHistory: &entities.TaskHistory{
				ID:           1,
				Address:      "address1",
				RewardPoints: 100,
				Amount:       10.0,
				CompletedAt:  &now,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			Task: &entities.Task{
				ID:          1,
				Name:        "Onboarding",
				Description: "Onboarding task",
				Points:      10,
			},
		},
	}

	// 設置 mock 返回值
	taskHistoryRepoMock.On("GetByAddressIncludingTasks", "address1").Return(taskHistoryMock, nil)

	svc := NewCampaignService(cfg, loggerMock, taskHistoryRepoMock, taskRepoMock, redisHelperMock)
	result, err := svc.GetPointHistories("address1")

	// 驗證結果
	assert.NoError(t, err)                   // 確保沒有錯誤
	assert.Len(t, result, 1)                 // 確保返回結果長度正確
	assert.Equal(t, taskHistoryMock, result) // 確保返回的數據正確

	// 驗證 mock 方法是否被正確調用
	taskHistoryRepoMock.AssertExpectations(t)
}

func TestGetTaskStatus(t *testing.T) {
	// 初始化服務
	cfg := &config.Config{}
	loggerMock := &mocks.MockLogger{}
	taskHistoryRepoMock := &mocks.MockTaskHistoryRepository{}
	taskRepoMock := &mocks.MockTaskRepository{}
	redisHelperMock := &mocks.MockRedisHelper{}

	// 模擬 taskRepo 的行為
	taskWithHistoryMock := []*models.TaskWithTaskHistory{
		{
			TaskID:          1,
			TaskName:        OnboardingTaskStr,
			TaskDescription: "Onboarding task",
			TaskPoints:      10,
		},
	}

	// 設置 mock 返回值
	taskRepoMock.On("GetByAddressAndNamesIncludingTaskHistories", "address1", []string{OnboardingTaskStr, SharePoolTaskStr}).
		Return(taskWithHistoryMock, nil)

	svc := NewCampaignService(cfg, loggerMock, taskHistoryRepoMock, taskRepoMock, redisHelperMock)
	result, err := svc.GetTaskStatus("address1")

	// 驗證結果
	assert.NoError(t, err)                       // 確保沒有錯誤
	assert.Len(t, result, 1)                     // 確保返回結果長度正確
	assert.Equal(t, taskWithHistoryMock, result) // 確保返回的數據正確

	// 驗證 mock 方法是否被正確調用
	taskRepoMock.AssertExpectations(t)
}

func TestStartCampaign(t *testing.T) {
	cfg := &config.Config{}
	loggerMock := &mocks.MockLogger{}
	taskHistoryRepoMock := &mocks.MockTaskHistoryRepository{}
	taskRepoMock := &mocks.MockTaskRepository{}
	redisHelperMock := &mocks.MockRedisHelper{}

	// 模擬 taskRepo 的行為
	taskRepoMock.On("IsExistedByName", OnboardingTaskStr).Return(false, nil)
	taskRepoMock.On("Create", mock.Anything).Return(&entities.Task{}, nil)

	// 模擬 share pool task 的行為
	taskRepoMock.On("IsExistedByName", SharePoolTaskStr).Return(false, nil)
	taskRepoMock.On("Create", mock.Anything).Return(&entities.Task{}, nil)

	// 模擬 logger 的行為
	loggerMock.On("Info", mock.Anything).Return()

	// 呼叫 StartCampaign 方法
	svc := NewCampaignService(cfg, loggerMock, taskHistoryRepoMock, taskRepoMock, redisHelperMock)
	err := svc.StartCampaign()

	// 驗證結果
	assert.NoError(t, err) // 確保沒有錯誤

	// 驗證 createOnboardingTask 是否被調用
	taskRepoMock.AssertCalled(t, "IsExistedByName", OnboardingTaskStr)
	taskRepoMock.AssertCalled(t, "Create", mock.Anything)

	// 驗證 createSharePoolTask 是否被調用
	taskRepoMock.AssertCalled(t, "IsExistedByName", SharePoolTaskStr)
	taskRepoMock.AssertCalled(t, "Create", mock.Anything)

	// 驗證 logger 是否有記錄啟動計劃
	loggerMock.AssertCalled(t, "Info", mock.Anything)
}

func TestFindCurrentSharePoolTask(t *testing.T) {
	// 設置模擬的 RedisHelper 和 TaskRepo
	mockRedisHelper := new(mocks.MockRedisHelper)
	mockTaskRepo := new(mocks.MockTaskRepository)

	// 設置 CampaignService
	service := &CampaignService{
		redisHelper: mockRedisHelper,
		taskRepo:    mockTaskRepo,
	}

	// 測試場景：Redis 已經有資料
	mockRedisHelper.On("Get", "curr_shared_pool_task").Return(`{"id":1,"name":"share_pool","started_at":"2024-11-01T00:00:00Z","end_at":"2024-11-30T00:00:00Z"}`, nil)

	task, err := service.FindCurrentSharePoolTask()

	assert.NoError(t, err)
	assert.Equal(t, "share_pool", task.Name)
	mockRedisHelper.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestFindOnboardingTask(t *testing.T) {
	// 設置模擬的 RedisHelper 和 TaskRepo
	mockRedisHelper := new(mocks.MockRedisHelper)
	mockTaskRepo := new(mocks.MockTaskRepository)

	// 設置 CampaignService
	service := &CampaignService{
		redisHelper: mockRedisHelper,
		taskRepo:    mockTaskRepo,
	}

	// 測試場景：Redis 已經有資料
	mockRedisHelper.On("Get", "onboarding_task").Return(`{"id":1,"name":"onboarding","started_at":"2024-11-01T00:00:00Z","end_at":"2024-11-30T00:00:00Z"}`, nil)

	task, err := service.FindOnboardingTask()

	assert.NoError(t, err)
	assert.Equal(t, "onboarding", task.Name)
	mockRedisHelper.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}
