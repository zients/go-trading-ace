package dtos

import (
	"time"
	"trading-ace/entities"
)

type TaskDTO struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Points      float64    `json:"points"`
	StartedAt   *time.Time `json:"started_at"`
	EndAt       *time.Time `json:"end_at"`
	Period      int        `json:"period"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ConvertTaskToDTO(task *entities.Task) *TaskDTO {
	return &TaskDTO{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Points:      task.Points,
		StartedAt:   task.StartedAt,
		EndAt:       task.EndAt,
		Period:      task.Period,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
