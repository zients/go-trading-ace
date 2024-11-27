package dtos

import (
	"testing"
	"time"

	"trading-ace/entities"

	"github.com/stretchr/testify/assert"
)

func TestConvertTaskToDTO(t *testing.T) {
	// Arrange
	startedAt := time.Now().Add(-72 * time.Hour) // 3 days ago
	endAt := time.Now().Add(24 * time.Hour)     // 1 day in the future
	createdAt := time.Now().Add(-96 * time.Hour) // 4 days ago
	updatedAt := time.Now()
	task := &entities.Task{
		ID:          1,
		Name:        "Test Task",
		Description: "This is a test task",
		Points:      50.5,
		StartedAt:   &startedAt,
		EndAt:       &endAt,
		Period:      7,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	// Act
	result := ConvertTaskToDTO(task)

	// Assert
	assert.Equal(t, task.ID, result.ID, "ID should match")
	assert.Equal(t, task.Name, result.Name, "Name should match")
	assert.Equal(t, task.Description, result.Description, "Description should match")
	assert.Equal(t, task.Points, result.Points, "Points should match")
	assert.Equal(t, task.StartedAt, result.StartedAt, "StartedAt should match")
	assert.Equal(t, task.EndAt, result.EndAt, "EndAt should match")
	assert.Equal(t, task.Period, result.Period, "Period should match")
	assert.Equal(t, task.CreatedAt, result.CreatedAt, "CreatedAt should match")
	assert.Equal(t, task.UpdatedAt, result.UpdatedAt, "UpdatedAt should match")
}
