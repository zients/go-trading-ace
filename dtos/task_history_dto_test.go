package dtos

import (
	"testing"
	"time"
	"trading-ace/entities"

	"github.com/stretchr/testify/assert"
)

func TestConvertTaskHistoryToDTO(t *testing.T) {
	createdAt := time.Now()
	updatedAt := time.Now().Add(time.Hour)
	completedAt := time.Now().Add(-time.Hour)

	taskHistory := &entities.TaskHistory{
		ID:           1,
		Address:      "Sample Address",
		TaskID:       101,
		RewardPoints: 50.0,
		Amount:       100.0,
		CompletedAt:  &completedAt,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	taskHistoryDTO := ConvertTaskHistoryToDTO(taskHistory)

	assert.NotNil(t, taskHistoryDTO)
	assert.Equal(t, taskHistory.ID, taskHistoryDTO.ID)
	assert.Equal(t, taskHistory.Address, taskHistoryDTO.Address)
	assert.Equal(t, taskHistory.TaskID, taskHistoryDTO.TaskID)
	assert.Equal(t, taskHistory.RewardPoints, taskHistoryDTO.RewardPoints)
	assert.Equal(t, taskHistory.Amount, taskHistoryDTO.Amount)
	assert.Equal(t, taskHistory.CompletedAt, taskHistoryDTO.CompletedAt)
	assert.Equal(t, taskHistory.CreatedAt, taskHistoryDTO.CreatedAt)
	assert.Equal(t, taskHistory.UpdatedAt, taskHistoryDTO.UpdatedAt)
}
