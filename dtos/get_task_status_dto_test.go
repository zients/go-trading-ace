package dtos

import (
	"testing"
	"time"

	"trading-ace/models"

	"github.com/stretchr/testify/assert"
)

func TestCovertTaskWithTaskHistoryToDTO(t *testing.T) {
	// Arrange
	startedAt := time.Now().Add(-72 * time.Hour)  // 3 days ago
	endAt := time.Now().Add(96 * time.Hour)       // 1 day in the future
	createdAt := time.Now().Add(-96 * time.Hour)  // 4 days ago
	completedAt := time.Now().Add(48 * time.Hour) // 2 days in the future
	updatedAt := time.Now()
	taskWithHistory := &models.TaskWithTaskHistory{
		TaskID:          1,
		TaskName:        "Test Task",
		TaskDescription: "This is a test task",
		TaskPoints:      50.5,
		TaskStartedAt:   &startedAt,
		TaskEndAt:       &endAt,
		TaskPeriod:      7,
		TaskCreatedAt:   createdAt,
		TaskUpdatedAt:   updatedAt,

		TaskHistoryID:           newInt64Ptr(1),
		TaskHistoryAddress:      newStringPtr("123"),
		TaskHistoryRewardPoints: newFloat64Ptr(1000),
		TaskHistoryAmount:       newFloat64Ptr(10),
		TaskHistoryCompletedAt:  &completedAt,
		TaskHistoryCreatedAt:    &createdAt,
		TaskHistoryUpdatedAt:    &updatedAt,
	}

	expectedTaskWithHistory := &TaskWithTaskHistoryDTO{
		taskWithHistory.TaskName,
		taskWithHistory.TaskStartedAt,
		taskWithHistory.TaskEndAt,
		taskWithHistory.TaskPeriod,
		InProgress,
		true,
		*taskWithHistory.TaskHistoryRewardPoints,
		taskWithHistory.TaskHistoryAmount,
		&completedAt,
	}

	// Act
	result := CovertTaskWithTaskHistoryToDTO(taskWithHistory)

	// Assert
	assert.Equal(t, expectedTaskWithHistory.TaskName, result.TaskName, "TaskName should match")
	assert.Equal(t, expectedTaskWithHistory.TaskStartedAt, result.TaskStartedAt, "TaskStartedAt should match")
	assert.Equal(t, expectedTaskWithHistory.TaskEndAt, result.TaskEndAt, "TaskEndAt should match")
	assert.Equal(t, expectedTaskWithHistory.TaskPeriod, result.TaskPeriod, "TaskPeriod should match")
	assert.Equal(t, expectedTaskWithHistory.Status, result.Status, "Status should match")
	assert.Equal(t, expectedTaskWithHistory.IsCompleted, result.IsCompleted, "IsCompleted should match")
	assert.Equal(t, expectedTaskWithHistory.RewardPoints, result.RewardPoints, "RewardPoints should match")
	assert.Equal(t, expectedTaskWithHistory.Amount, result.Amount, "Amount should match")
	assert.Equal(t, expectedTaskWithHistory.TaskHistoryCompletedAt, result.TaskHistoryCompletedAt, "CompletedAt should match")

}

func newInt64Ptr(a int64) *int64 {
	return &a
}

func newFloat64Ptr(a float64) *float64 {
	return &a
}

func newStringPtr(a string) *string {
	return &a
}
