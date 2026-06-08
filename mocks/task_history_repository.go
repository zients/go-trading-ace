package mocks

import (
	"context"
	"trading-ace/entities"
	"trading-ace/models"

	"github.com/stretchr/testify/mock"
)

type MockTaskHistoryRepository struct {
	mock.Mock
}

func (m *MockTaskHistoryRepository) Create(ctx context.Context, taskHistory *entities.TaskHistory) (*entities.TaskHistory, error) {
	args := m.Called(ctx, taskHistory)
	return getTaskHistory(args, 0), args.Error(1)
}

func (m *MockTaskHistoryRepository) Upsert(ctx context.Context, taskHistory *entities.TaskHistory) (*entities.TaskHistory, error) {
	args := m.Called(ctx, taskHistory)
	return getTaskHistory(args, 0), args.Error(1)
}

func (m *MockTaskHistoryRepository) FindByID(ctx context.Context, id int64) (*entities.TaskHistory, error) {
	args := m.Called(ctx, id)
	return getTaskHistory(args, 0), args.Error(1)
}

func (m *MockTaskHistoryRepository) FindByAddressAndTaskId(ctx context.Context, address string, taskId int64) (*entities.TaskHistory, error) {
	args := m.Called(ctx, address, taskId)
	return getTaskHistory(args, 0), args.Error(1)
}

func (m *MockTaskHistoryRepository) GetByAddressIncludingTasks(ctx context.Context, address string) ([]*models.TaskTaskHistoryPair, error) {
	args := m.Called(ctx, address)
	return args.Get(0).([]*models.TaskTaskHistoryPair), args.Error(1)
}

func getTaskHistory(args mock.Arguments, index int) *entities.TaskHistory {
	value := args.Get(index)
	if value == nil {
		return nil
	}

	return value.(*entities.TaskHistory)
}
