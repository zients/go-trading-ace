package services

import (
	"context"
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
	taskHistoryRepoMock.On("GetByAddressIncludingTasks", mock.Anything, "address1").Return(taskHistoryMock, nil)

	svc := NewCampaignService(cfg, loggerMock, taskHistoryRepoMock, taskRepoMock, redisHelperMock)
	result, err := svc.GetPointHistories(context.Background(), "address1")

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
	taskRepoMock.On("GetByAddressAndNamesIncludingTaskHistories", mock.Anything, "address1", []string{OnboardingTaskStr, SharePoolTaskStr}).
		Return(taskWithHistoryMock, nil)

	svc := NewCampaignService(cfg, loggerMock, taskHistoryRepoMock, taskRepoMock, redisHelperMock)
	result, err := svc.GetTaskStatus(context.Background(), "address1")

	// 驗證結果
	assert.NoError(t, err)                       // 確保沒有錯誤
	assert.Len(t, result, 1)                     // 確保返回結果長度正確
	assert.Equal(t, taskWithHistoryMock, result) // 確保返回的數據正確

	// 驗證 mock 方法是否被正確調用
	taskRepoMock.AssertExpectations(t)
}

func TestGetTaskStatusPropagatesContextToRepository(t *testing.T) {
	cfg := &config.Config{}
	loggerMock := &mocks.MockLogger{}
	taskHistoryRepoMock := &mocks.MockTaskHistoryRepository{}
	taskRepoMock := &mocks.MockTaskRepository{}
	redisHelperMock := &mocks.MockRedisHelper{}
	ctx := context.WithValue(context.Background(), struct{}{}, "request-context")

	taskRepoMock.On(
		"GetByAddressAndNamesIncludingTaskHistories",
		sameContext(ctx),
		"address1",
		[]string{OnboardingTaskStr, SharePoolTaskStr},
	).Return([]*models.TaskWithTaskHistory{}, nil).Once()

	svc := NewCampaignService(cfg, loggerMock, taskHistoryRepoMock, taskRepoMock, redisHelperMock)
	result, err := svc.GetTaskStatus(ctx, "address1")

	assert.NoError(t, err)
	assert.Empty(t, result)
	taskRepoMock.AssertExpectations(t)
}

func TestStartCampaign(t *testing.T) {
	cfg := &config.Config{}
	loggerMock := &mocks.MockLogger{}
	taskHistoryRepoMock := &mocks.MockTaskHistoryRepository{}
	taskRepoMock := &mocks.MockTaskRepository{}
	redisHelperMock := &mocks.MockRedisHelper{}

	// 模擬 taskRepo 的行為
	taskRepoMock.On("IsExistedByName", mock.Anything, OnboardingTaskStr).Return(false, nil)
	taskRepoMock.On("Create", mock.Anything, mock.Anything).Return(&entities.Task{}, nil)

	// 模擬 share pool task 的行為
	taskRepoMock.On("GetByName", mock.Anything, SharePoolTaskStr).Return([]*entities.Task{}, nil)
	taskRepoMock.On("Create", mock.Anything, mock.Anything).Return(&entities.Task{}, nil)

	// 呼叫 StartCampaign 方法
	svc := NewCampaignService(cfg, loggerMock, taskHistoryRepoMock, taskRepoMock, redisHelperMock)
	err := svc.StartCampaign(context.Background())

	// 驗證結果
	assert.NoError(t, err) // 確保沒有錯誤

	// 驗證 createOnboardingTask 是否被調用
	taskRepoMock.AssertCalled(t, "IsExistedByName", mock.Anything, OnboardingTaskStr)
	taskRepoMock.AssertCalled(t, "Create", mock.Anything, mock.Anything)

	// 驗證 createSharePoolTask 是否被調用
	taskRepoMock.AssertCalled(t, "GetByName", mock.Anything, SharePoolTaskStr)
	taskRepoMock.AssertCalled(t, "Create", mock.Anything, mock.Anything)

	loggerMock.AssertNotCalled(t, "Info", mock.Anything)
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

	taskRepoMock.On("IsExistedByName", mock.Anything, OnboardingTaskStr).Return(true, nil).Once()
	taskRepoMock.On("GetByName", mock.Anything, SharePoolTaskStr).Return(existingSharePoolTasks, nil).Once()

	svc := NewCampaignService(cfg, loggerMock, taskHistoryRepoMock, taskRepoMock, redisHelperMock)
	err := svc.StartCampaign(context.Background())

	assert.NoError(t, err)
	taskRepoMock.AssertNotCalled(t, "Create", mock.Anything)
	taskRepoMock.AssertExpectations(t)
	loggerMock.AssertNotCalled(t, "Info", mock.Anything)
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

	taskRepoMock.On("IsExistedByName", mock.Anything, OnboardingTaskStr).Return(true, nil).Once()
	taskRepoMock.On("GetByName", mock.Anything, SharePoolTaskStr).Return(incompleteSharePoolTasks, nil).Once()

	svc := NewCampaignService(cfg, loggerMock, taskHistoryRepoMock, taskRepoMock, redisHelperMock)
	err := svc.StartCampaign(context.Background())

	assert.Error(t, err)
	assert.ErrorContains(t, err, "expected 4 share pool tasks")
	taskRepoMock.AssertNotCalled(t, "Create", mock.Anything)
	taskRepoMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)
}

func TestStartCampaignDoesNotStartSettlementWorker(t *testing.T) {
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

	taskRepoMock.On("IsExistedByName", mock.Anything, OnboardingTaskStr).Return(false, nil).Once()
	taskRepoMock.On("Create", mock.Anything, mock.Anything).Return(&entities.Task{}, nil).Times(5)
	taskRepoMock.On("GetByName", mock.Anything, SharePoolTaskStr).Return([]*entities.Task{}, nil).Once()

	taskRepoMock.On("IsExistedByName", mock.Anything, OnboardingTaskStr).Return(true, nil).Once()
	taskRepoMock.On("GetByName", mock.Anything, SharePoolTaskStr).Return(existingSharePoolTasks, nil).Once()

	svc := NewCampaignService(cfg, loggerMock, taskHistoryRepoMock, taskRepoMock, redisHelperMock)

	assert.NoError(t, svc.StartCampaign(context.Background()))
	assert.NoError(t, svc.StartCampaign(context.Background()))

	taskRepoMock.AssertNumberOfCalls(t, "Create", 5)
	loggerMock.AssertNotCalled(t, "Info", mock.Anything)
	taskRepoMock.AssertExpectations(t)
}

func TestSettleDueSharePoolTasksClaimsDueTaskAndMarksSettled(t *testing.T) {
	mockRedisHelper := new(mocks.MockRedisHelper)
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockTaskHistoryRepo := new(mocks.MockTaskHistoryRepository)
	loggerMock := new(mocks.MockLogger)
	ctx := context.WithValue(context.Background(), struct{}{}, "settlement-context")
	claimStartedAt := time.Now().UTC()
	dueTask := &entities.Task{
		ID:                  7,
		Name:                SharePoolTaskStr,
		Period:              2,
		Points:              SharePoolTaskPoints,
		SettlementStartedAt: &claimStartedAt,
	}
	key := "SharePoolTask_2"

	service := &CampaignService{
		logger:          loggerMock,
		redisHelper:     mockRedisHelper,
		taskRepo:        mockTaskRepo,
		taskHistoryRepo: mockTaskHistoryRepo,
	}

	mockTaskRepo.On("ClaimDueSharePoolTask", sameContext(ctx), mock.AnythingOfType("time.Time")).
		Return(dueTask, nil).Once()
	mockRedisHelper.On("HGetAll", sameContext(ctx), key).Return(map[string]string{
		"0xA": "1000",
		"0xB": "3000",
	}, nil).Once()
	mockRedisHelper.On("Get", sameContext(ctx), key+"_total").Return("4000", nil).Once()
	mockTaskHistoryRepo.On("Upsert", sameContext(ctx), mock.AnythingOfType("*entities.TaskHistory")).
		Return(&entities.TaskHistory{}, nil).Twice()
	mockRedisHelper.On("ZAdd", sameContext(ctx), key+"_rank", mock.Anything).Return(nil).Twice()
	mockTaskRepo.On("MarkSettled", sameContext(ctx), dueTask.ID, claimStartedAt, mock.AnythingOfType("time.Time")).
		Return(nil).Once()
	mockTaskRepo.On("ClaimDueSharePoolTask", sameContext(ctx), mock.AnythingOfType("time.Time")).
		Return(nil, nil).Once()

	err := service.SettleDueSharePoolTasks(ctx)

	assert.NoError(t, err)
	mockTaskRepo.AssertExpectations(t)
	mockRedisHelper.AssertExpectations(t)
	mockTaskHistoryRepo.AssertExpectations(t)
	loggerMock.AssertExpectations(t)
}

func TestSettleDueSharePoolTasksReleasesClaimWhenCalculationFails(t *testing.T) {
	mockRedisHelper := new(mocks.MockRedisHelper)
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockTaskHistoryRepo := new(mocks.MockTaskHistoryRepository)
	loggerMock := new(mocks.MockLogger)
	ctx := context.Background()
	claimStartedAt := time.Now().UTC()
	dueTask := &entities.Task{
		ID:                  7,
		Name:                SharePoolTaskStr,
		Period:              2,
		Points:              SharePoolTaskPoints,
		SettlementStartedAt: &claimStartedAt,
	}
	key := "SharePoolTask_2"

	service := &CampaignService{
		logger:          loggerMock,
		redisHelper:     mockRedisHelper,
		taskRepo:        mockTaskRepo,
		taskHistoryRepo: mockTaskHistoryRepo,
	}

	mockTaskRepo.On("ClaimDueSharePoolTask", mock.Anything, mock.AnythingOfType("time.Time")).
		Return(dueTask, nil).Once()
	mockRedisHelper.On("HGetAll", mock.Anything, key).Return(map[string]string{"0xA": "1000"}, nil).Once()
	mockRedisHelper.On("Get", mock.Anything, key+"_total").Return("", errors.New("missing total")).Once()
	mockTaskRepo.On("ReleaseSettlementClaim", mock.Anything, dueTask.ID, claimStartedAt).Return(nil).Once()

	err := service.SettleDueSharePoolTasks(ctx)

	assert.ErrorContains(t, err, "missing total")
	mockTaskRepo.AssertExpectations(t)
	mockRedisHelper.AssertExpectations(t)
	mockTaskHistoryRepo.AssertExpectations(t)
	loggerMock.AssertExpectations(t)
}

func TestSettleDueSharePoolTasksMarksNoVolumeTaskSettled(t *testing.T) {
	mockRedisHelper := new(mocks.MockRedisHelper)
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockTaskHistoryRepo := new(mocks.MockTaskHistoryRepository)
	loggerMock := new(mocks.MockLogger)
	ctx := context.Background()
	claimStartedAt := time.Now().UTC()
	dueTask := &entities.Task{
		ID:                  7,
		Name:                SharePoolTaskStr,
		Period:              2,
		Points:              SharePoolTaskPoints,
		SettlementStartedAt: &claimStartedAt,
	}
	key := "SharePoolTask_2"

	service := &CampaignService{
		logger:          loggerMock,
		redisHelper:     mockRedisHelper,
		taskRepo:        mockTaskRepo,
		taskHistoryRepo: mockTaskHistoryRepo,
	}

	mockTaskRepo.On("ClaimDueSharePoolTask", mock.Anything, mock.AnythingOfType("time.Time")).
		Return(dueTask, nil).Once()
	mockRedisHelper.On("HGetAll", mock.Anything, key).Return(map[string]string{}, nil).Once()
	mockTaskRepo.On("MarkSettled", mock.Anything, dueTask.ID, claimStartedAt, mock.AnythingOfType("time.Time")).
		Return(nil).Once()
	mockTaskRepo.On("ClaimDueSharePoolTask", mock.Anything, mock.AnythingOfType("time.Time")).
		Return(nil, nil).Once()

	err := service.SettleDueSharePoolTasks(ctx)

	assert.NoError(t, err)
	mockTaskRepo.AssertExpectations(t)
	mockRedisHelper.AssertExpectations(t)
	mockRedisHelper.AssertNotCalled(t, "Get", mock.Anything, key+"_total")
	mockTaskHistoryRepo.AssertExpectations(t)
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
	mockRedisHelper.On("Get", mock.Anything, "curr_shared_pool_task").Return(`{"id":1,"name":"share_pool","started_at":"2024-11-01T00:00:00Z","end_at":"2024-11-30T00:00:00Z"}`, nil)

	task, err := service.FindCurrentSharePoolTask(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, "share_pool", task.Name)
	mockRedisHelper.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestFindCurrentSharePoolTaskReturnsErrorWhenCachedTaskCannotDecode(t *testing.T) {
	mockRedisHelper := new(mocks.MockRedisHelper)
	mockTaskRepo := new(mocks.MockTaskRepository)
	service := &CampaignService{
		redisHelper: mockRedisHelper,
		taskRepo:    mockTaskRepo,
	}

	mockRedisHelper.On("Get", mock.Anything, "curr_shared_pool_task").Return(`{bad json`, nil)

	task, err := service.FindCurrentSharePoolTask(context.Background())

	assert.Nil(t, task)
	assert.ErrorContains(t, err, "failed to decode current share pool task cache")
	mockRedisHelper.AssertExpectations(t)
	mockTaskRepo.AssertNotCalled(t, "GetByName", mock.Anything, mock.Anything)
}

func TestFindCurrentSharePoolTaskReturnsErrorWhenCachingActiveTaskFails(t *testing.T) {
	mockRedisHelper := new(mocks.MockRedisHelper)
	mockTaskRepo := new(mocks.MockTaskRepository)
	now := time.Now()
	activeTask := &entities.Task{
		ID:        1,
		Name:      SharePoolTaskStr,
		StartedAt: ptrTime(now.Add(-time.Hour)),
		EndAt:     ptrTime(now.Add(time.Hour)),
		Period:    1,
	}
	service := &CampaignService{
		redisHelper: mockRedisHelper,
		taskRepo:    mockTaskRepo,
	}

	mockRedisHelper.On("Get", mock.Anything, "curr_shared_pool_task").Return("", errors.New("cache miss"))
	mockTaskRepo.On("GetByName", mock.Anything, SharePoolTaskStr).Return([]*entities.Task{activeTask}, nil)
	mockRedisHelper.On("Set", mock.Anything, "curr_shared_pool_task", mock.Anything, mock.Anything).Return(errors.New("redis set failed"))

	task, err := service.FindCurrentSharePoolTask(context.Background())

	assert.Nil(t, task)
	assert.ErrorContains(t, err, "failed to cache current share pool task")
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
	mockRedisHelper.On("Get", mock.Anything, "onboarding_task").Return(`{"id":1,"name":"onboarding","started_at":"2024-11-01T00:00:00Z","end_at":"2024-11-30T00:00:00Z"}`, nil)

	task, err := service.FindOnboardingTask(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, "onboarding", task.Name)
	mockRedisHelper.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestFindOnboardingTaskReturnsErrorWhenCachedTaskCannotDecode(t *testing.T) {
	mockRedisHelper := new(mocks.MockRedisHelper)
	mockTaskRepo := new(mocks.MockTaskRepository)
	service := &CampaignService{
		redisHelper: mockRedisHelper,
		taskRepo:    mockTaskRepo,
	}

	mockRedisHelper.On("Get", mock.Anything, "onboarding_task").Return(`{bad json`, nil)

	task, err := service.FindOnboardingTask(context.Background())

	assert.Nil(t, task)
	assert.ErrorContains(t, err, "failed to decode onboarding task cache")
	mockRedisHelper.AssertExpectations(t)
	mockTaskRepo.AssertNotCalled(t, "FindByName", mock.Anything, mock.Anything)
}

func TestFindOnboardingTaskReturnsErrorWhenCachingTaskFails(t *testing.T) {
	mockRedisHelper := new(mocks.MockRedisHelper)
	mockTaskRepo := new(mocks.MockTaskRepository)
	now := time.Now()
	task := &entities.Task{
		ID:    1,
		Name:  OnboardingTaskStr,
		EndAt: ptrTime(now.Add(time.Hour)),
	}
	service := &CampaignService{
		redisHelper: mockRedisHelper,
		taskRepo:    mockTaskRepo,
	}

	mockRedisHelper.On("Get", mock.Anything, "onboarding_task").Return("", errors.New("cache miss"))
	mockTaskRepo.On("FindByName", mock.Anything, OnboardingTaskStr).Return(task, nil)
	mockRedisHelper.On("Set", mock.Anything, "onboarding_task", mock.Anything, mock.Anything).Return(errors.New("redis set failed"))

	result, err := service.FindOnboardingTask(context.Background())

	assert.Nil(t, result)
	assert.ErrorContains(t, err, "failed to cache onboarding task")
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
	now := time.Now()
	currentTask := &entities.Task{
		Name:      SharePoolTaskStr,
		StartedAt: ptrTime(now.Add(-time.Hour)),
		EndAt:     ptrTime(now.Add(time.Hour)),
		Period:    1,
	}
	totalAmountStr := "100.0"
	key := "SharePoolTask_1"

	// Mock Redis responses
	mockRedisHelper.On("Get", mock.Anything, "curr_shared_pool_task").Return("", errors.New("cache miss"))
	mockRedisHelper.On("Set", mock.Anything, "curr_shared_pool_task", mock.Anything, mock.Anything).Return(nil)
	mockRedisHelper.On("HIncrFloat", mock.Anything, key, senderAddress, amount).Return(nil)
	mockRedisHelper.On("HGet", mock.Anything, key, senderAddress).Return(totalAmountStr, nil)
	mockRedisHelper.On("IncrFloat", mock.Anything, key+"_total", amount).Return(nil)

	// Mock FindCurrentSharePoolTask response
	mockTaskRepo.On("GetByName", mock.Anything, SharePoolTaskStr).Return([]*entities.Task{currentTask}, nil)

	// Call the method under test
	totalAmountReturned, err := campaignService.RecordUSDCSwapTotalAmount(context.Background(), senderAddress, amount)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, amount, totalAmountReturned)

	// Assert that the Redis helper and task history repo methods were called
	mockRedisHelper.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
	mockTaskHistoryRepo.AssertExpectations(t)
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

	mockRedisHelper.On("Get", mock.Anything, "curr_shared_pool_task").Return("", errors.New("cache miss"))
	mockTaskRepo.On("GetByName", mock.Anything, SharePoolTaskStr).Return([]*entities.Task{currentTask}, nil)
	mockRedisHelper.On("Set", mock.Anything, "curr_shared_pool_task", mock.Anything, mock.Anything).Return(nil)
	mockRedisHelper.On("HIncrFloat", mock.Anything, key, "0x123", 100.0).Return(errors.New("redis hincr failed"))
	mockRedisHelper.On("HGet", mock.Anything, key, "0x123").Return("100", nil).Maybe()
	mockRedisHelper.On("IncrFloat", mock.Anything, key+"_total", 100.0).Return(nil).Maybe()

	_, err := campaignService.RecordUSDCSwapTotalAmount(context.Background(), "0x123", 100.0)

	assert.ErrorContains(t, err, "redis hincr failed")
	mockRedisHelper.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestRecordUSDCSwapTotalAmountPropagatesContextToDependencies(t *testing.T) {
	mockRedisHelper := new(mocks.MockRedisHelper)
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockTaskHistoryRepo := new(mocks.MockTaskHistoryRepository)
	ctx := context.WithValue(context.Background(), struct{}{}, "record-context")
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

	mockRedisHelper.On("Get", sameContext(ctx), "curr_shared_pool_task").Return("", errors.New("cache miss")).Once()
	mockTaskRepo.On("GetByName", sameContext(ctx), SharePoolTaskStr).Return([]*entities.Task{currentTask}, nil).Once()
	mockRedisHelper.On("Set", sameContext(ctx), "curr_shared_pool_task", mock.Anything, mock.Anything).Return(nil).Once()
	mockRedisHelper.On("HIncrFloat", sameContext(ctx), key, "0x123", 100.0).Return(nil).Once()
	mockRedisHelper.On("HGet", sameContext(ctx), key, "0x123").Return("100", nil).Once()
	mockRedisHelper.On("IncrFloat", sameContext(ctx), key+"_total", 100.0).Return(nil).Once()

	totalAmount, err := campaignService.RecordUSDCSwapTotalAmount(ctx, "0x123", 100.0)

	assert.NoError(t, err)
	assert.Equal(t, 100.0, totalAmount)
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

	mockRedisHelper.On("Get", mock.Anything, "curr_shared_pool_task").Return("", errors.New("cache miss"))
	mockTaskRepo.On("GetByName", mock.Anything, SharePoolTaskStr).Return([]*entities.Task{currentTask}, nil)
	mockRedisHelper.On("Set", mock.Anything, "curr_shared_pool_task", mock.Anything, mock.Anything).Return(nil)
	mockRedisHelper.On("HIncrFloat", mock.Anything, key, "0x123", 100.0).Return(nil)
	mockRedisHelper.On("HGet", mock.Anything, key, "0x123").Return("100", nil)
	mockRedisHelper.On("IncrFloat", mock.Anything, key+"_total", 100.0).Return(errors.New("redis incr failed"))

	_, err := campaignService.RecordUSDCSwapTotalAmount(context.Background(), "0x123", 100.0)

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

	mockRedisHelper.On("Get", mock.Anything, "curr_shared_pool_task").Return("", errors.New("cache miss"))
	mockTaskRepo.On("GetByName", mock.Anything, SharePoolTaskStr).Return([]*entities.Task{currentTask}, nil)
	mockRedisHelper.On("Set", mock.Anything, "curr_shared_pool_task", mock.Anything, mock.Anything).Return(nil)
	mockRedisHelper.On("HIncrFloat", mock.Anything, key, "0x123", 1000.0).Return(nil)
	mockRedisHelper.On("HGet", mock.Anything, key, "0x123").Return("1000", nil)
	mockRedisHelper.On("IncrFloat", mock.Anything, key+"_total", 1000.0).Return(nil)
	mockRedisHelper.On("Get", mock.Anything, "onboarding_task").Return("", errors.New("cache miss"))
	mockTaskRepo.On("FindByName", mock.Anything, OnboardingTaskStr).Return(onboardingTask, nil)
	mockRedisHelper.On("Set", mock.Anything, "onboarding_task", mock.Anything, mock.Anything).Return(nil)
	mockTaskHistoryRepo.On("FindByAddressAndTaskId", mock.Anything, "0x123", onboardingTask.ID).Return(nil, errors.New("not found"))
	mockTaskHistoryRepo.On("Create", mock.Anything, mock.Anything).Return(&entities.TaskHistory{}, errors.New("db create failed"))

	_, err := campaignService.RecordUSDCSwapTotalAmount(context.Background(), "0x123", 1000.0)

	assert.ErrorContains(t, err, "db create failed")
	mockRedisHelper.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
	mockTaskHistoryRepo.AssertExpectations(t)
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func sameContext(expected context.Context) interface{} {
	return mock.MatchedBy(func(actual context.Context) bool {
		return actual == expected
	})
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
		mockRedisHelper.On("ZRevRangeWithScores", mock.Anything, key, int64(0), int64(-1)).
			Return([]string{"address1", "address2"}, []float64{100.5, 75.3}, nil)

		// Expected result
		expectedLeaderboard := []models.LeaderboardEntry{
			{Address: "address1", Score: 100.5},
			{Address: "address2", Score: 75.3},
		}

		// Call the method
		result, err := campaignService.GetLeaderboard(context.Background(), taskName, period)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedLeaderboard, result)

		// Assert that Redis helper was called as expected
		mockRedisHelper.AssertExpectations(t)
	})
}
