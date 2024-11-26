package dtos

import (
	"testing"
	"time"
	"trading-ace/entities"

	"github.com/stretchr/testify/assert"
)

func TestConvertTaskToDTO(t *testing.T) {
	createdAt := time.Now()
	updatedAt := time.Now().Add(time.Hour)
	startedAt := time.Now().Add(-time.Hour)
	endAt := time.Now().Add(time.Hour * 2)

	task := &entities.Task{
		ID:          1,
		Name:        "Sample Task",
		Description: "This is a sample task",
		Points:      10.5,
		StartedAt:   &startedAt,
		EndAt:       &endAt,
		Period:      30,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	taskDTO := ConvertTaskToDTO(task)

	assert.NotNil(t, taskDTO)
	assert.Equal(t, task.ID, taskDTO.ID)
	assert.Equal(t, task.Name, taskDTO.Name)
	assert.Equal(t, task.Description, taskDTO.Description)
	assert.Equal(t, task.Points, taskDTO.Points)
	assert.Equal(t, task.StartedAt, taskDTO.StartedAt)
	assert.Equal(t, task.EndAt, taskDTO.EndAt)
	assert.Equal(t, task.Period, taskDTO.Period)
	assert.Equal(t, task.CreatedAt, taskDTO.CreatedAt)
	assert.Equal(t, task.UpdatedAt, taskDTO.UpdatedAt)
}
