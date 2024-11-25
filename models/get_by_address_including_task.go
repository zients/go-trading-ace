package models

import "trading-ace/entities"

type GetByAddressIncludingTask struct {
	Task        *entities.Task
	TaskHistory *entities.TaskHistory
}
