package services

import (
	"fmt"
	"time"
	"trading-ace/entities"
	"trading-ace/repositories"
)

type ICampaignService interface {
	StartCampaign() error
}

type CampaignService struct {
	taskRecordRepo repositories.ITaskRecordRepository
	taskRepo       repositories.ITaskRepository
}

const OnboardingTaskStr string = "OnboardingTask"
const OnboardingTaskDescription string = "OnboardingTask"
const OnboardingTaskPoints int64 = 100

const SharePoolTaskStr string = "SharePoolTask"
const SharePoolTaskDescription string = "SharePoolTask"
const SharePoolTaskPoints int64 = 10000

func NewCampaignService(taskRecordRepo repositories.ITaskRecordRepository, taskRepo repositories.ITaskRepository) ICampaignService {
	return &CampaignService{
		taskRecordRepo: taskRecordRepo,
		taskRepo:       taskRepo,
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
		IsRecurring: false,
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
		var duration = time.Duration(7*i) * 24 * time.Hour
		endAt := startedAt.Add(duration)

		newTask := &entities.Task{
			Name:        SharePoolTaskStr,
			Description: SharePoolTaskDescription,
			Points:      SharePoolTaskPoints,
			StartedAt:   &startedAt,
			EndAt:       &endAt,
			IsRecurring: true,
			Period:      i,
		}

		if _, err := s.taskRepo.Create(newTask); err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}
	}

	return nil
}
