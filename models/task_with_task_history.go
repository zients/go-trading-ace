package models

import "time"

type TaskWithTaskHistory struct {
	TaskID          int64      // Mapping to tasks.id
	TaskName        string     // Mapping to tasks.name
	TaskDescription string     // Mapping to tasks.description
	TaskPoints      float64    // Mapping to tasks.points
	TaskStartedAt   *time.Time // Mapping to tasks.started_at
	TaskEndAt       *time.Time // Mapping to tasks.end_at
	TaskPeriod      int        // Mapping to tasks.period
	TaskCreatedAt   time.Time  // Mapping to tasks.created_at
	TaskUpdatedAt   time.Time  // Mapping to tasks.updated_at

	TaskHistoryID           *int64     // Mapping to task_histories.id
	TaskHistoryAddress      *string    // Mapping to task_histories.address
	TaskHistoryRewardPoints *float64   // Mapping to task_histories.reward_points
	TaskHistoryAmount       *float64   // Mapping to task_histories.amount
	TaskHistoryCompletedAt  *time.Time // Mapping to task_histories.completed_at
	TaskHistoryCreatedAt    *time.Time // Mapping to task_histories.created_at
	TaskHistoryUpdatedAt    *time.Time // Mapping to task_histories.updated_at
}
