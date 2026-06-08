package mocks

import (
	"trading-ace/entities"
	"trading-ace/models"

	"github.com/stretchr/testify/mock"
)

type MockTaskHistoryRepository struct {
	mock.Mock
}

func (m *MockTaskHistoryRepository) Create(taskHistory *entities.TaskHistory) (*entities.TaskHistory, error) {
	args := m.Called(taskHistory)
	return getTaskHistory(args, 0), args.Error(1)
}

func (m *MockTaskHistoryRepository) FindByID(id int64) (*entities.TaskHistory, error) {
	args := m.Called(id)
	return getTaskHistory(args, 0), args.Error(1)
}

func (m *MockTaskHistoryRepository) FindByAddressAndTaskId(address string, taskId int64) (*entities.TaskHistory, error) {
	args := m.Called(address, taskId)
	return getTaskHistory(args, 0), args.Error(1)
}

func (m *MockTaskHistoryRepository) GetByAddressIncludingTasks(address string) ([]*models.TaskTaskHistoryPair, error) {
	args := m.Called(address)
	return args.Get(0).([]*models.TaskTaskHistoryPair), args.Error(1)
}

func getTaskHistory(args mock.Arguments, index int) *entities.TaskHistory {
	value := args.Get(index)
	if value == nil {
		return nil
	}

	return value.(*entities.TaskHistory)
}
