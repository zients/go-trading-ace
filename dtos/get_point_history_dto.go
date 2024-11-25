package dtos

type GetPointHistoryDTO struct {
	Task        *TaskDTO        `json:"task"`
	TaskHistory *TaskHistoryDTO `json:"task_history"`
}
