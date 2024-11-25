package services

import (
	"encoding/json"
	"fmt"
	"time"
	"trading-ace/config"
	"trading-ace/entities"
	"trading-ace/helpers"
	"trading-ace/logger"
	"trading-ace/repositories"
)

type ICampaignService interface {
	StartCampaign() error
}

type CampaignService struct {
	config         *config.Config
	logger         logger.ILogger
	taskRecordRepo repositories.ITaskRecordRepository
	taskRepo       repositories.ITaskRepository
	redisHelper    helpers.IRedisHelper
}

const OnboardingTaskStr string = "OnboardingTask"
const OnboardingTaskDescription string = "OnboardingTask"
const OnboardingTaskPoints int64 = 100

const SharePoolTaskStr string = "SharePoolTask"
const SharePoolTaskDescription string = "SharePoolTask"
const SharePoolTaskPoints int64 = 10000

func NewCampaignService(
	config *config.Config,
	logger logger.ILogger,
	taskRecordRepo repositories.ITaskRecordRepository,
	taskRepo repositories.ITaskRepository,
	redisHelper helpers.IRedisHelper,
) ICampaignService {
	return &CampaignService{
		config:         config,
		logger:         logger,
		taskRecordRepo: taskRecordRepo,
		taskRepo:       taskRepo,
		redisHelper:    redisHelper,
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

func (s *CampaignService) getCurrentSharePoolTask() (*entities.Task, error) {
	key := s.config.Redis.Prefix + "curr_shared_pool_task"
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
