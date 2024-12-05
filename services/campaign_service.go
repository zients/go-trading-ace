package services

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"
	"trading-ace/config"
	"trading-ace/entities"
	"trading-ace/helpers"
	"trading-ace/logger"
	"trading-ace/models"
	"trading-ace/repositories"

	"github.com/go-redis/redis/v8"
)

type ICampaignService interface {
	StartCampaign() error
	GetPointHistories(address string) ([]*models.TaskTaskHistoryPair, error)
	RecordUSDCSwapTotalAmount(senderAddress string, amount float64) (float64, error)
	GetTaskStatus(address string) ([]*models.TaskWithTaskHistory, error)
	FindOnboardingTask() (*entities.Task, error)
	FindCurrentSharePoolTask() (*entities.Task, error)
	GetLeaderboard(taskName string, period int) ([]models.LeaderboardEntry, error)
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
	shareTasks, err := s.createSharePoolTask()
	if err != nil {
		return err
	}

	s.startLimitedWeeklySettlementScheduler(shareTasks)

	return nil
}

func (s *CampaignService) GetPointHistories(address string) ([]*models.TaskTaskHistoryPair, error) {
	return s.taskHistoryRepo.GetByAddressIncludingTasks(address)
}

func (s *CampaignService) GetTaskStatus(address string) ([]*models.TaskWithTaskHistory, error) {
	return s.taskRepo.GetByAddressAndNamesIncludingTaskHistories(address, []string{OnboardingTaskStr, SharePoolTaskStr})
}

func (s *CampaignService) RecordUSDCSwapTotalAmount(senderAddress string, amount float64) (float64, error) {
	// find current share task
	task, err := s.FindCurrentSharePoolTask()
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

	onboardingTask, err := s.FindOnboardingTask()
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

func (s *CampaignService) createSharePoolTask() ([]*entities.Task, error) {
	isExisted, err := s.taskRepo.IsExistedByName(SharePoolTaskStr)
	if err != nil {
		return []*entities.Task{}, err
	}

	if isExisted {
		return []*entities.Task{}, fmt.Errorf("share pool task is existed")
	}

	startedAt := time.Now().UTC()
	results := []*entities.Task{}
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

		task, err := s.taskRepo.Create(newTask)
		if err != nil {
			return []*entities.Task{}, fmt.Errorf("failed to create task: %w", err)
		}

		results = append(results, task)

		startedAt = endAt
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Period < results[j].Period
	})

	return results, nil
}

func (s *CampaignService) FindCurrentSharePoolTask() (*entities.Task, error) {
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

func (s *CampaignService) FindOnboardingTask() (*entities.Task, error) {
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

func (s *CampaignService) startLimitedWeeklySettlementScheduler(tasks []*entities.Task) {
	ticker := time.NewTicker(7 * 24 * time.Hour)
	maxRuns := 4
	runCount := 0

	go func() {
		for range ticker.C {
			if runCount >= maxRuns {
				s.logger.Info("Weekly settlement scheduler reached its limit, stopping...")
				ticker.Stop()
				return
			}

			if err := s.calculateSharePoolPoint(tasks[runCount]); err != nil {
				s.logger.Error("Failed to perform weekly settlement: %v", err)
			}

			runCount++
		}
	}()

	s.logger.Info("Limited weekly settlement scheduler started")
}

func (s *CampaignService) calculateSharePoolPoint(task *entities.Task) error {
	if task.Name != SharePoolTaskStr {
		return fmt.Errorf("task is not shard pool task")
	}

	key := fmt.Sprintf("%s_%d", task.Name, task.Period)
	totalKey := fmt.Sprintf("%s_total", key)

	totalStr, err := s.redisHelper.Get(totalKey)
	if err != nil {
		return err
	}

	totalAmount, err := strconv.ParseFloat(totalStr, 64)
	if err != nil {
		return fmt.Errorf("failed to parse total amount from key %s: %w", totalKey, err)
	}

	swapAmountMap, err := s.redisHelper.HGetAll(key)
	if err != nil {
		return err
	}

	taskPoints := task.Points

	now := time.Now().UTC()
	for address, v := range swapAmountMap {
		amount, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}

		quotes := amount / totalAmount
		rewards := taskPoints * quotes

		history := &entities.TaskHistory{
			Address:      address,
			TaskID:       task.ID,
			RewardPoints: rewards,
			Amount:       amount,
			CompletedAt:  &now,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		if _, err := s.taskHistoryRepo.Create(history); err != nil {
			s.logger.Error("create history failed for address %s, %v", address, err)
			continue
		}

		s.redisHelper.ZAdd(fmt.Sprintf("%s_rank", key), &redis.Z{Score: rewards, Member: address})
	}

	return nil
}

func (s *CampaignService) GetLeaderboard(taskName string, period int) ([]models.LeaderboardEntry, error) {
	key := fmt.Sprintf("%s_%d_rank", taskName, period)

	members, scores, err := s.redisHelper.ZRevRangeWithScores(key, 0, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch leaderboard for key %s: %w", key, err)
	}

	if len(members) != len(scores) {
		return nil, fmt.Errorf("mismatched lengths: %d members, %d scores", len(members), len(scores))
	}

	entries := make([]models.LeaderboardEntry, len(members))
	for i := range members {
		entries[i] = models.LeaderboardEntry{
			Address: members[i],
			Score:   scores[i],
		}
	}

	return entries, nil
}
