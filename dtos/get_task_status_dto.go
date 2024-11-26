package dtos

import (
	"time"
	"trading-ace/models"
)

type TaskWithTaskHistoryDTO struct {
	TaskName               string     // Mapping to tasks.name
	TaskStartedAt          *time.Time // Mapping to tasks.started_at
	TaskEndAt              *time.Time // Mapping to tasks.end_at
	TaskPeriod             int        // Mapping to tasks.period
	Status                 string     // task status according to start and end time
	IsCompleted            bool
	RewardPoints           float64    // Mapping to task_histories.reward_points
	Amount                 *float64   // Mapping to task_histories.amount
	TaskHistoryCompletedAt *time.Time // Mapping to task_histories.completed_at
}

const NotStarted string = "Not Started"
const InProgress string = "In Progress"
const Completed string = "Completed"

func CovertTaskWithTaskHistoryToDTO(model *models.TaskWithTaskHistory) *TaskWithTaskHistoryDTO {
	var now = time.Now().UTC()
	var status string

	switch {
	case model.TaskStartedAt == nil || (*model.TaskStartedAt).After(now):
		status = NotStarted
	case model.TaskStartedAt != nil && (*model.TaskStartedAt).Before(now) && now.Before(*model.TaskEndAt):
		status = InProgress
	default:
		status = Completed
	}

	return &TaskWithTaskHistoryDTO{
		TaskName:      model.TaskName,
		TaskStartedAt: model.TaskStartedAt,
		TaskEndAt:     model.TaskEndAt,
		TaskPeriod:    model.TaskPeriod,
		Status:        status,
		IsCompleted:   model.TaskHistoryID != nil,
		RewardPoints: func() float64 {
			if model.TaskHistoryRewardPoints == nil {
				return 0
			}

			return *model.TaskHistoryRewardPoints
		}(),
		Amount: func() *float64 {
			if model.TaskHistoryAmount == nil {
				var result float64 = 0
				return &result
			}

			return model.TaskHistoryAmount
		}(),
		TaskHistoryCompletedAt: func() *time.Time {
			if model.TaskHistoryCompletedAt == nil {
				return nil
			}

			return model.TaskHistoryCompletedAt
		}(),
	}
}
