package dtos

import (
	"testing"
	"time"

	"trading-ace/entities"

	"github.com/stretchr/testify/assert"
)

func TestConvertTaskHistoryToDTO(t *testing.T) {
	// Arrange
	completedAt := time.Now().Add(-24 * time.Hour)
	createdAt := time.Now().Add(-48 * time.Hour)
	updatedAt := time.Now()
	taskHistory := &entities.TaskHistory{
		ID:           1,
		Address:      "0x123",
		TaskID:       10,
		RewardPoints: 100.5,
		Amount:       200.75,
		CompletedAt:  &completedAt,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	// Act
	result := ConvertTaskHistoryToDTO(taskHistory)

	// Assert
	assert.Equal(t, taskHistory.ID, result.ID, "ID should match")
	assert.Equal(t, taskHistory.Address, result.Address, "Address should match")
	assert.Equal(t, taskHistory.TaskID, result.TaskID, "TaskID should match")
	assert.Equal(t, taskHistory.RewardPoints, result.RewardPoints, "RewardPoints should match")
	assert.Equal(t, taskHistory.Amount, result.Amount, "Amount should match")
	assert.Equal(t, taskHistory.CompletedAt, result.CompletedAt, "CompletedAt should match")
	assert.Equal(t, taskHistory.CreatedAt, result.CreatedAt, "CreatedAt should match")
	assert.Equal(t, taskHistory.UpdatedAt, result.UpdatedAt, "UpdatedAt should match")
}
