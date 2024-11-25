package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"trading-ace/config"
	"trading-ace/entities"
	"trading-ace/helpers"
	"trading-ace/logger"
	"trading-ace/repositories"
)

type ICampaignService interface {
	StartCampaign() error
	RecordUSDCSwapTotalAmount(senderAddress string, amount float64) (float64, error)
}

type CampaignService struct {
	config          *config.Config
	logger          logger.ILogger
	taskHistoryRepo repositories.ITaskHistoryRepository
	taskRepo        repositories.ITaskRepository
	redisHelper     helpers.IRedisHelper
}

const OnboardingTaskStr string = "OnboardingTask"
const OnboardingTaskDescription string = "OnboardingTask"
const OnboardingTaskPoints float64 = 100
const OnboardingTaskTargetAmount float64 = 1000

const SharePoolTaskStr string = "SharePoolTask"
const SharePoolTaskDescription string = "SharePoolTask"
const SharePoolTaskPoints float64 = 10000

func NewCampaignService(
	config *config.Config,
	logger logger.ILogger,
	taskHistoryRepo repositories.ITaskHistoryRepository,
	taskRepo repositories.ITaskRepository,
	redisHelper helpers.IRedisHelper,
) ICampaignService {
	return &CampaignService{
		config:          config,
		logger:          logger,
		taskHistoryRepo: taskHistoryRepo,
		taskRepo:        taskRepo,
		redisHelper:     redisHelper,
	}
}

func (s *CampaignService) StartCampaign() error {
	// if exists
	if err := s.createOnboardingTask(); err != nil {
		return err
	}

	//share pool task
	if err := s.createSharePoolTask(); err != nil {
		return err
	}

	return nil
}

func (s *CampaignService) RecordUSDCSwapTotalAmount(senderAddress string, amount float64) (float64, error) {
	// find current share task
	task, err := s.findCurrentSharePoolTask()
	if err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%s_%d", task.Name, task.Period)
	s.redisHelper.HIncrFloat(key, senderAddress, amount)
	totalAmountStr, err := s.redisHelper.HGet(key, senderAddress)
	if err != nil {
		return 0, err
	}

	totalKey := fmt.Sprintf("%s_total", key)
	s.redisHelper.IncrFloat(totalKey, amount)

	totalAmount, err := strconv.ParseFloat(totalAmountStr, 64)
	if err != nil {
		return 0, err
	}

	// if amount is not enough
	if totalAmount < OnboardingTaskTargetAmount {
		return totalAmount, nil
	}

	onboardingTask, err := s.findOnboardingTask()
	if err != nil {
		return 0, err
	}

	// find existed onboarding completed task record
	_, err = s.taskHistoryRepo.FindByAddressAndTaskId(senderAddress, onboardingTask.ID)
	if err == nil {
		return totalAmount, nil
	}

	// create onboarding task record
	now := time.Now().UTC()
	taskHistory := &entities.TaskHistory{
		Address:      senderAddress,
		TaskID:       onboardingTask.ID,
		RewardPoints: OnboardingTaskPoints,
		Amount:       totalAmount,
		CompletedAt:  &now,
	}

	s.taskHistoryRepo.Create(taskHistory)

	return totalAmount, nil
}

func (s *CampaignService) createOnboardingTask() error {
	isExisted, err := s.taskRepo.IsExistedByName(OnboardingTaskStr)
	if err != nil {
		return err
	}

	if isExisted {
		return fmt.Errorf("onboarding task is existed")
	}

	startedAt := time.Now().UTC()
	endAt := startedAt.Add(28 * 24 * time.Hour)

	newTask := &entities.Task{
		Name:        OnboardingTaskStr,
		Description: OnboardingTaskDescription,
		Points:      OnboardingTaskPoints,
		StartedAt:   &startedAt,
		EndAt:       &endAt,
		Period:      1,
	}

	if _, err := s.taskRepo.Create(newTask); err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

func (s *CampaignService) createSharePoolTask() error {
	isExisted, err := s.taskRepo.IsExistedByName(SharePoolTaskStr)
	if err != nil {
		return err
	}

	if isExisted {
		return fmt.Errorf("share pool task is existed")
	}

	startedAt := time.Now().UTC()
	for i := 1; i <= 4; i++ {
		var duration = 7 * 24 * time.Hour
		endAt := startedAt.Add(duration)

		newTask := &entities.Task{
			Name:        SharePoolTaskStr,
			Description: SharePoolTaskDescription,
			Points:      SharePoolTaskPoints,
			StartedAt:   &startedAt,
			EndAt:       &endAt,
			Period:      i,
		}

		if _, err := s.taskRepo.Create(newTask); err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}

		startedAt = endAt
	}

	return nil
}

func (s *CampaignService) findCurrentSharePoolTask() (*entities.Task, error) {
	key := "curr_shared_pool_task"
	redisData, err := s.redisHelper.Get(key)
	if err == nil {
		task := &entities.Task{}
		json.Unmarshal([]byte(redisData), &task)

		return task, nil
	}

	tasks, err := s.taskRepo.GetByName(SharePoolTaskStr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch share pool tasks: %w", err)
	}

	now := time.Now().UTC()
	for _, task := range tasks {
		if task.StartedAt != nil && task.EndAt != nil && now.After(*task.StartedAt) && now.Before(*task.EndAt) {
			encodedTask, _ := json.Marshal(task)
			s.redisHelper.Set(key, string(encodedTask), time.Until(*task.EndAt))

			return task, nil
		}
	}

	return nil, fmt.Errorf("no active share pool task found")
}

func (s *CampaignService) findOnboardingTask() (*entities.Task, error) {
	key := "onboarding_task"
	redisData, err := s.redisHelper.Get(key)
	if err == nil {
		task := &entities.Task{}
		json.Unmarshal([]byte(redisData), &task)

		return task, nil
	}

	task, err := s.taskRepo.FindByName(OnboardingTaskStr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch onboarding task: %w", err)
	}

	encodedTask, _ := json.Marshal(task)
	s.redisHelper.Set(key, string(encodedTask), time.Until(*task.EndAt))

	return task, nil
}
