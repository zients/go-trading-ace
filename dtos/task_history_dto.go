package dtos

import (
	"time"
	"trading-ace/entities"
)

type TaskHistoryDTO struct {
	ID           int64      `json:"id"`
	Address      string     `json:"address"`
	TaskID       int64      `json:"task_id"`
	RewardPoints float64    `json:"reward_points"`
	Amount       float64    `json:"amount"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func ConvertTaskHistoryToDTO(taskHistory *entities.TaskHistory) *TaskHistoryDTO {
	return &TaskHistoryDTO{
		ID:           taskHistory.ID,
		Address:      taskHistory.Address,
		TaskID:       taskHistory.TaskID,
		RewardPoints: taskHistory.RewardPoints,
		Amount:       taskHistory.Amount,
		CompletedAt:  taskHistory.CompletedAt,
		CreatedAt:    taskHistory.CreatedAt,
		UpdatedAt:    taskHistory.UpdatedAt,
	}
}
