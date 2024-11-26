package models

import "trading-ace/entities"

type TaskTaskHistoryPair struct {
	Task        *entities.Task
	TaskHistory *entities.TaskHistory
}
