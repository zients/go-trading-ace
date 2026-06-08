package services

import (
	"errors"
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
	taskRepoMock.On("GetByName", SharePoolTaskStr).Return([]*entities.Task{}, nil)
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
	taskRepoMock.AssertCalled(t, "GetByName", SharePoolTaskStr)
	taskRepoMock.AssertCalled(t, "Create", mock.Anything)

	// 驗證 logger 是否有記錄啟動計劃
	loggerMock.AssertCalled(t, "Info", mock.Anything)
}

func TestStartCampaignReturnsNoErrorWhenCampaignTasksAlreadyExist(t *testing.T) {
	cfg := &config.Config{}
	loggerMock := &mocks.MockLogger{}
	taskHistoryRepoMock := &mocks.MockTaskHistoryRepository{}
	taskRepoMock := &mocks.MockTaskRepository{}
	redisHelperMock := &mocks.MockRedisHelper{}

	existingSharePoolTasks := []*entities.Task{
		{Name: SharePoolTaskStr, Period: 1},
		{Name: SharePoolTaskStr, Period: 2},
		{Name: SharePoolTaskStr, Period: 3},
		{Name: SharePoolTaskStr, Period: 4},
	}

	taskRepoMock.On("IsExistedByName", OnboardingTaskStr).Return(true, nil).Once()
	taskRepoMock.On("GetByName", SharePoolTaskStr).Return(existingSharePoolTasks, nil).Once()
	loggerMock.On("Info", mock.Anything).Return().Once()

	svc := NewCampaignService(cfg, loggerMock, taskHistoryRepoMock, taskRepoMock, redisHelperMock)
	err := svc.StartCampaign()

	assert.NoError(t, err)
	taskRepoMock.AssertNotCalled(t, "Create", mock.Anything)
	taskRepoMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)
}

func TestStartCampaignReturnsErrorWhenExistingSharePoolTasksAreIncomplete(t *testing.T) {
	cfg := &config.Config{}
	loggerMock := &mocks.MockLogger{}
	taskHistoryRepoMock := &mocks.MockTaskHistoryRepository{}
	taskRepoMock := &mocks.MockTaskRepository{}
	redisHelperMock := &mocks.MockRedisHelper{}

	incompleteSharePoolTasks := []*entities.Task{
		{Name: SharePoolTaskStr, Period: 1},
		{Name: SharePoolTaskStr, Period: 2},
	}

	taskRepoMock.On("IsExistedByName", OnboardingTaskStr).Return(true, nil).Once()
	taskRepoMock.On("GetByName", SharePoolTaskStr).Return(incompleteSharePoolTasks, nil).Once()

	svc := NewCampaignService(cfg, loggerMock, taskHistoryRepoMock, taskRepoMock, redisHelperMock)
	err := svc.StartCampaign()

	assert.Error(t, err)
	assert.ErrorContains(t, err, "expected 4 share pool tasks")
	taskRepoMock.AssertNotCalled(t, "Create", mock.Anything)
	taskRepoMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)
}

func TestStartCampaignDoesNotStartDuplicateScheduler(t *testing.T) {
	cfg := &config.Config{}
	loggerMock := &mocks.MockLogger{}
	taskHistoryRepoMock := &mocks.MockTaskHistoryRepository{}
	taskRepoMock := &mocks.MockTaskRepository{}
	redisHelperMock := &mocks.MockRedisHelper{}

	existingSharePoolTasks := []*entities.Task{
		{Name: SharePoolTaskStr, Period: 1},
		{Name: SharePoolTaskStr, Period: 2},
		{Name: SharePoolTaskStr, Period: 3},
		{Name: SharePoolTaskStr, Period: 4},
	}

	taskRepoMock.On("IsExistedByName", OnboardingTaskStr).Return(false, nil).Once()
	taskRepoMock.On("Create", mock.Anything).Return(&entities.Task{}, nil).Times(5)
	taskRepoMock.On("GetByName", SharePoolTaskStr).Return([]*entities.Task{}, nil).Once()
	loggerMock.On("Info", mock.Anything).Return().Once()

	taskRepoMock.On("IsExistedByName", OnboardingTaskStr).Return(true, nil).Once()
	taskRepoMock.On("GetByName", SharePoolTaskStr).Return(existingSharePoolTasks, nil).Once()

	svc := NewCampaignService(cfg, loggerMock, taskHistoryRepoMock, taskRepoMock, redisHelperMock)

	assert.NoError(t, svc.StartCampaign())
	assert.NoError(t, svc.StartCampaign())

	taskRepoMock.AssertNumberOfCalls(t, "Create", 5)
	loggerMock.AssertNumberOfCalls(t, "Info", 1)
	taskRepoMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)
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

func TestRecordUSDCSwapTotalAmount(t *testing.T) {
	// Mock dependencies
	mockRedisHelper := new(mocks.MockRedisHelper)
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockTaskHistoryRepo := new(mocks.MockTaskHistoryRepository)

	campaignService := &CampaignService{
		redisHelper:     mockRedisHelper,
		taskRepo:        mockTaskRepo,
		taskHistoryRepo: mockTaskHistoryRepo,
	}

	// Mock data
	senderAddress := "0x123"
	amount := 100.0
	currentTask := &entities.Task{
		Name:   "share_pool_task",
		Period: 1,
	}
	onboardingTask := &entities.Task{
		ID:     1,
		Name:   "onboarding_task",
		Period: 1,
	}
	totalAmountStr := "100.0"

	// Mock Redis responses
	mockRedisHelper.On("HIncrFloat", mock.Anything, senderAddress, amount).Return(nil)
	mockRedisHelper.On("HGet", mock.Anything, senderAddress).Return(totalAmountStr, nil)
	mockRedisHelper.On("IncrFloat", mock.Anything, amount).Return(nil)
	mockRedisHelper.On("Get", mock.Anything).Return(totalAmountStr, nil)

	// Mock FindCurrentSharePoolTask response
	mockTaskRepo.On("GetByName", "share_pool_task").Return([]*entities.Task{currentTask}, nil)
	mockTaskRepo.On("FindByName", "onboarding_task").Return(onboardingTask, nil)

	// Mock task history repo to simulate no existing onboarding task history
	mockTaskHistoryRepo.On("FindByAddressAndTaskId", senderAddress, onboardingTask.ID).Return(nil, errors.New("not found"))

	// Simulate successful creation of task history
	mockTaskHistoryRepo.On("Create", mock.Anything).Return(nil)

	// Call the method under test
	totalAmountReturned, err := campaignService.RecordUSDCSwapTotalAmount(senderAddress, amount)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, amount, totalAmountReturned)

	// Assert that the Redis helper and task history repo methods were called
	mockRedisHelper.AssertExpectations(t)
}

func TestRecordUSDCSwapTotalAmountReturnsErrorWhenUserAmountIncrementFails(t *testing.T) {
	mockRedisHelper := new(mocks.MockRedisHelper)
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockTaskHistoryRepo := new(mocks.MockTaskHistoryRepository)
	campaignService := &CampaignService{
		redisHelper:     mockRedisHelper,
		taskRepo:        mockTaskRepo,
		taskHistoryRepo: mockTaskHistoryRepo,
	}

	now := time.Now()
	currentTask := &entities.Task{
		Name:      SharePoolTaskStr,
		Period:    1,
		StartedAt: &now,
		EndAt:     ptrTime(now.Add(time.Hour)),
	}
	key := "SharePoolTask_1"

	mockRedisHelper.On("Get", "curr_shared_pool_task").Return("", errors.New("cache miss"))
	mockTaskRepo.On("GetByName", SharePoolTaskStr).Return([]*entities.Task{currentTask}, nil)
	mockRedisHelper.On("Set", "curr_shared_pool_task", mock.Anything, mock.Anything).Return(nil)
	mockRedisHelper.On("HIncrFloat", key, "0x123", 100.0).Return(errors.New("redis hincr failed"))
	mockRedisHelper.On("HGet", key, "0x123").Return("100", nil).Maybe()
	mockRedisHelper.On("IncrFloat", key+"_total", 100.0).Return(nil).Maybe()

	_, err := campaignService.RecordUSDCSwapTotalAmount("0x123", 100.0)

	assert.ErrorContains(t, err, "redis hincr failed")
	mockRedisHelper.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestRecordUSDCSwapTotalAmountReturnsErrorWhenTotalAmountIncrementFails(t *testing.T) {
	mockRedisHelper := new(mocks.MockRedisHelper)
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockTaskHistoryRepo := new(mocks.MockTaskHistoryRepository)
	campaignService := &CampaignService{
		redisHelper:     mockRedisHelper,
		taskRepo:        mockTaskRepo,
		taskHistoryRepo: mockTaskHistoryRepo,
	}

	now := time.Now()
	currentTask := &entities.Task{
		Name:      SharePoolTaskStr,
		Period:    1,
		StartedAt: &now,
		EndAt:     ptrTime(now.Add(time.Hour)),
	}
	key := "SharePoolTask_1"

	mockRedisHelper.On("Get", "curr_shared_pool_task").Return("", errors.New("cache miss"))
	mockTaskRepo.On("GetByName", SharePoolTaskStr).Return([]*entities.Task{currentTask}, nil)
	mockRedisHelper.On("Set", "curr_shared_pool_task", mock.Anything, mock.Anything).Return(nil)
	mockRedisHelper.On("HIncrFloat", key, "0x123", 100.0).Return(nil)
	mockRedisHelper.On("HGet", key, "0x123").Return("100", nil)
	mockRedisHelper.On("IncrFloat", key+"_total", 100.0).Return(errors.New("redis incr failed"))

	_, err := campaignService.RecordUSDCSwapTotalAmount("0x123", 100.0)

	assert.ErrorContains(t, err, "redis incr failed")
	mockRedisHelper.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestRecordUSDCSwapTotalAmountReturnsErrorWhenOnboardingHistoryCreateFails(t *testing.T) {
	mockRedisHelper := new(mocks.MockRedisHelper)
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockTaskHistoryRepo := new(mocks.MockTaskHistoryRepository)
	campaignService := &CampaignService{
		redisHelper:     mockRedisHelper,
		taskRepo:        mockTaskRepo,
		taskHistoryRepo: mockTaskHistoryRepo,
	}

	now := time.Now()
	currentTask := &entities.Task{
		Name:      SharePoolTaskStr,
		Period:    1,
		StartedAt: &now,
		EndAt:     ptrTime(now.Add(time.Hour)),
	}
	onboardingTask := &entities.Task{
		ID:    10,
		Name:  OnboardingTaskStr,
		EndAt: ptrTime(now.Add(24 * time.Hour)),
	}
	key := "SharePoolTask_1"

	mockRedisHelper.On("Get", "curr_shared_pool_task").Return("", errors.New("cache miss"))
	mockTaskRepo.On("GetByName", SharePoolTaskStr).Return([]*entities.Task{currentTask}, nil)
	mockRedisHelper.On("Set", "curr_shared_pool_task", mock.Anything, mock.Anything).Return(nil)
	mockRedisHelper.On("HIncrFloat", key, "0x123", 1000.0).Return(nil)
	mockRedisHelper.On("HGet", key, "0x123").Return("1000", nil)
	mockRedisHelper.On("IncrFloat", key+"_total", 1000.0).Return(nil)
	mockRedisHelper.On("Get", "onboarding_task").Return("", errors.New("cache miss"))
	mockTaskRepo.On("FindByName", OnboardingTaskStr).Return(onboardingTask, nil)
	mockRedisHelper.On("Set", "onboarding_task", mock.Anything, mock.Anything).Return(nil)
	mockTaskHistoryRepo.On("FindByAddressAndTaskId", "0x123", onboardingTask.ID).Return(nil, errors.New("not found"))
	mockTaskHistoryRepo.On("Create", mock.Anything).Return(&entities.TaskHistory{}, errors.New("db create failed"))

	_, err := campaignService.RecordUSDCSwapTotalAmount("0x123", 1000.0)

	assert.ErrorContains(t, err, "db create failed")
	mockRedisHelper.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
	mockTaskHistoryRepo.AssertExpectations(t)
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func TestCampaignService_GetLeaderboard(t *testing.T) {
	// Mock dependencies
	mockRedisHelper := new(mocks.MockRedisHelper)
	mockTaskRepo := new(mocks.MockTaskRepository)

	// Initialize the service with mocked dependencies
	campaignService := &CampaignService{
		redisHelper: mockRedisHelper,
		taskRepo:    mockTaskRepo,
	}

	// Test case variables
	taskName := "SharePoolTask"
	period := 7
	key := "SharePoolTask_7_rank"

	t.Run("Success", func(t *testing.T) {
		// Mock Redis response
		mockRedisHelper.On("ZRevRangeWithScores", key, int64(0), int64(-1)).
			Return([]string{"address1", "address2"}, []float64{100.5, 75.3}, nil)

		// Expected result
		expectedLeaderboard := []models.LeaderboardEntry{
			{Address: "address1", Score: 100.5},
			{Address: "address2", Score: 75.3},
		}

		// Call the method
		result, err := campaignService.GetLeaderboard(taskName, period)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedLeaderboard, result)

		// Assert that Redis helper was called as expected
		mockRedisHelper.AssertExpectations(t)
	})
}
