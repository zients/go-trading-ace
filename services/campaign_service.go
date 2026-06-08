package services

import (
	"context"
	"encoding/json"
	"errors"
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
	StartCampaign(ctx context.Context) error
	GetPointHistories(ctx context.Context, address string) ([]*models.TaskTaskHistoryPair, error)
	RecordUSDCSwapTotalAmount(ctx context.Context, eventID string, senderAddress string, amount float64) (float64, error)
	GetTaskStatus(ctx context.Context, address string) ([]*models.TaskWithTaskHistory, error)
	FindOnboardingTask(ctx context.Context) (*entities.Task, error)
	FindCurrentSharePoolTask(ctx context.Context) (*entities.Task, error)
	GetLeaderboard(ctx context.Context, taskName string, period int) ([]models.LeaderboardEntry, error)
	SettleDueSharePoolTasks(ctx context.Context) error
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
const SharePoolTaskPeriods int = 4
const processedSwapEventRetention time.Duration = 30 * 24 * time.Hour

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

func (s *CampaignService) StartCampaign(ctx context.Context) error {
	if err := s.createOnboardingTask(ctx); err != nil {
		return err
	}

	if _, err := s.createSharePoolTask(ctx); err != nil {
		return err
	}

	return nil
}

func (s *CampaignService) GetPointHistories(ctx context.Context, address string) ([]*models.TaskTaskHistoryPair, error) {
	return s.taskHistoryRepo.GetByAddressIncludingTasks(ctx, address)
}

func (s *CampaignService) GetTaskStatus(ctx context.Context, address string) ([]*models.TaskWithTaskHistory, error) {
	return s.taskRepo.GetByAddressAndNamesIncludingTaskHistories(ctx, address, []string{OnboardingTaskStr, SharePoolTaskStr})
}

func (s *CampaignService) RecordUSDCSwapTotalAmount(ctx context.Context, eventID string, senderAddress string, amount float64) (float64, error) {
	// find current share task
	task, err := s.FindCurrentSharePoolTask(ctx)
	if err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%s_%d", task.Name, task.Period)
	totalKey := fmt.Sprintf("%s_total", key)
	totalAmount, recorded, err := s.redisHelper.RecordSwapVolumeOnce(ctx, eventID, key, totalKey, senderAddress, amount, swapEventDedupExpiration(task, time.Now().UTC()))
	if err != nil {
		return 0, fmt.Errorf("failed to record swap amount: %w", err)
	}

	if !recorded {
		return 0, nil
	}

	// if amount is not enough
	if totalAmount < OnboardingTaskTargetAmount {
		return totalAmount, nil
	}

	onboardingTask, err := s.FindOnboardingTask(ctx)
	if err != nil {
		return 0, err
	}

	// find existed onboarding completed task record
	_, err = s.taskHistoryRepo.FindByAddressAndTaskId(ctx, senderAddress, onboardingTask.ID)
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

	if _, err := s.taskHistoryRepo.Create(ctx, taskHistory); err != nil {
		return 0, fmt.Errorf("failed to create onboarding task history: %w", err)
	}

	return totalAmount, nil
}

func swapEventDedupExpiration(task *entities.Task, now time.Time) time.Duration {
	if task == nil || task.EndAt == nil {
		return processedSwapEventRetention
	}

	expiration := task.EndAt.Sub(now) + processedSwapEventRetention
	if expiration < time.Second {
		return time.Second
	}

	return expiration
}

func (s *CampaignService) createOnboardingTask(ctx context.Context) error {
	isExisted, err := s.taskRepo.IsExistedByName(ctx, OnboardingTaskStr)
	if err != nil {
		return err
	}

	if isExisted {
		return nil
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

	if _, err := s.taskRepo.Create(ctx, newTask); err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

func (s *CampaignService) createSharePoolTask(ctx context.Context) ([]*entities.Task, error) {
	tasks, err := s.taskRepo.GetByName(ctx, SharePoolTaskStr)
	if err != nil {
		return []*entities.Task{}, err
	}

	if len(tasks) > 0 {
		return validateSharePoolTasks(tasks)
	}

	startedAt := time.Now().UTC()
	results := []*entities.Task{}
	for i := 1; i <= SharePoolTaskPeriods; i++ {
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

		task, err := s.taskRepo.Create(ctx, newTask)
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

func validateSharePoolTasks(tasks []*entities.Task) ([]*entities.Task, error) {
	if len(tasks) != SharePoolTaskPeriods {
		return nil, fmt.Errorf("expected %d share pool tasks, found %d", SharePoolTaskPeriods, len(tasks))
	}

	seenPeriods := make(map[int]bool, SharePoolTaskPeriods)
	for _, task := range tasks {
		if task == nil {
			return nil, fmt.Errorf("share pool task list contains nil task")
		}

		if task.Name != SharePoolTaskStr {
			return nil, fmt.Errorf("unexpected share pool task name %q", task.Name)
		}

		if task.Period < 1 || task.Period > SharePoolTaskPeriods {
			return nil, fmt.Errorf("share pool task period %d is out of range", task.Period)
		}

		if seenPeriods[task.Period] {
			return nil, fmt.Errorf("duplicate share pool task period %d", task.Period)
		}

		seenPeriods[task.Period] = true
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Period < tasks[j].Period
	})

	return tasks, nil
}

func (s *CampaignService) FindCurrentSharePoolTask(ctx context.Context) (*entities.Task, error) {
	key := "curr_shared_pool_task"
	redisData, err := s.redisHelper.Get(ctx, key)
	if err == nil {
		return decodeCachedTask(redisData, "current share pool task")
	}

	tasks, err := s.taskRepo.GetByName(ctx, SharePoolTaskStr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch share pool tasks: %w", err)
	}

	now := time.Now().UTC()
	for _, task := range tasks {
		if task.StartedAt != nil && task.EndAt != nil && now.After(*task.StartedAt) && now.Before(*task.EndAt) {
			if err := s.cacheTask(ctx, key, task, "current share pool task"); err != nil {
				return nil, err
			}

			return task, nil
		}
	}

	return nil, fmt.Errorf("no active share pool task found")
}

func (s *CampaignService) FindOnboardingTask(ctx context.Context) (*entities.Task, error) {
	key := "onboarding_task"
	redisData, err := s.redisHelper.Get(ctx, key)
	if err == nil {
		return decodeCachedTask(redisData, "onboarding task")
	}

	task, err := s.taskRepo.FindByName(ctx, OnboardingTaskStr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch onboarding task: %w", err)
	}

	if err := s.cacheTask(ctx, key, task, "onboarding task"); err != nil {
		return nil, err
	}

	return task, nil
}

func decodeCachedTask(redisData string, label string) (*entities.Task, error) {
	task := &entities.Task{}
	if err := json.Unmarshal([]byte(redisData), task); err != nil {
		return nil, fmt.Errorf("failed to decode %s cache: %w", label, err)
	}

	if task.Name == "" {
		return nil, fmt.Errorf("failed to decode %s cache: task name is empty", label)
	}

	return task, nil
}

func (s *CampaignService) cacheTask(ctx context.Context, key string, task *entities.Task, label string) error {
	if task == nil {
		return fmt.Errorf("failed to cache %s: task is nil", label)
	}

	if task.EndAt == nil {
		return fmt.Errorf("failed to cache %s: task end time is missing", label)
	}

	encodedTask, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to encode %s cache: %w", label, err)
	}

	if err := s.redisHelper.Set(ctx, key, string(encodedTask), time.Until(*task.EndAt)); err != nil {
		return fmt.Errorf("failed to cache %s: %w", label, err)
	}

	return nil
}

func (s *CampaignService) calculateSharePoolPoint(ctx context.Context, task *entities.Task) error {
	if task.Name != SharePoolTaskStr {
		return fmt.Errorf("task is not shard pool task")
	}

	key := fmt.Sprintf("%s_%d", task.Name, task.Period)
	totalKey := fmt.Sprintf("%s_total", key)

	swapAmountMap, err := s.redisHelper.HGetAll(ctx, key)
	if err != nil {
		return err
	}

	if len(swapAmountMap) == 0 {
		return nil
	}

	totalStr, err := s.redisHelper.Get(ctx, totalKey)
	if err != nil {
		return err
	}
	totalAmount, err := strconv.ParseFloat(totalStr, 64)
	if err != nil {
		return fmt.Errorf("failed to parse total amount from key %s: %w", totalKey, err)
	}
	if totalAmount <= 0 {
		return fmt.Errorf("total amount for key %s must be greater than zero", totalKey)
	}

	taskPoints := task.Points

	now := time.Now().UTC()
	var settlementErrors []error
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

		if _, err := s.taskHistoryRepo.Upsert(ctx, history); err != nil {
			settlementErrors = append(settlementErrors, fmt.Errorf("upsert history failed for address %s: %w", address, err))
			continue
		}

		if err := s.redisHelper.ZAdd(ctx, fmt.Sprintf("%s_rank", key), &redis.Z{Score: rewards, Member: address}); err != nil {
			settlementErrors = append(settlementErrors, fmt.Errorf("update leaderboard failed for address %s: %w", address, err))
		}
	}

	if err := errors.Join(settlementErrors...); err != nil {
		return err
	}

	return nil
}

func (s *CampaignService) SettleDueSharePoolTasks(ctx context.Context) error {
	for {
		task, err := s.taskRepo.ClaimDueSharePoolTask(ctx, time.Now().UTC())
		if err != nil {
			return err
		}
		if task == nil {
			return nil
		}

		if task.SettlementStartedAt == nil {
			return fmt.Errorf("claimed share pool task %d is missing settlement_started_at", task.ID)
		}
		claimStartedAt := *task.SettlementStartedAt

		if err := s.calculateSharePoolPoint(ctx, task); err != nil {
			if releaseErr := s.taskRepo.ReleaseSettlementClaim(ctx, task.ID, claimStartedAt); releaseErr != nil {
				return errors.Join(
					fmt.Errorf("failed to settle share pool task %d: %w", task.ID, err),
					fmt.Errorf("failed to release settlement claim for task %d: %w", task.ID, releaseErr),
				)
			}

			return fmt.Errorf("failed to settle share pool task %d: %w", task.ID, err)
		}

		if err := s.taskRepo.MarkSettled(ctx, task.ID, claimStartedAt, time.Now().UTC()); err != nil {
			if releaseErr := s.taskRepo.ReleaseSettlementClaim(ctx, task.ID, claimStartedAt); releaseErr != nil {
				return errors.Join(
					fmt.Errorf("failed to mark share pool task %d settled: %w", task.ID, err),
					fmt.Errorf("failed to release settlement claim for task %d: %w", task.ID, releaseErr),
				)
			}

			return fmt.Errorf("failed to mark share pool task %d settled: %w", task.ID, err)
		}
	}
}

func (s *CampaignService) GetLeaderboard(ctx context.Context, taskName string, period int) ([]models.LeaderboardEntry, error) {
	key := fmt.Sprintf("%s_%d_rank", taskName, period)

	members, scores, err := s.redisHelper.ZRevRangeWithScores(ctx, key, 0, -1)
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
