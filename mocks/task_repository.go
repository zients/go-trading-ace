package mocks

import (
	"context"
	"time"
	"trading-ace/entities"
	"trading-ace/models"

	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(*entities.Task), args.Error(1)
}

func (m *MockTaskRepository) FindById(ctx context.Context, id int64) (*entities.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Task), args.Error(1)
}

func (m *MockTaskRepository) FindByName(ctx context.Context, name string) (*entities.Task, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*entities.Task), args.Error(1)
}

func (m *MockTaskRepository) GetByName(ctx context.Context, name string) ([]*entities.Task, error) {
	args := m.Called(ctx, name)
	return args.Get(0).([]*entities.Task), args.Error(1)
}

func (m *MockTaskRepository) IsExistedByName(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockTaskRepository) GetByAddressAndNamesIncludingTaskHistories(ctx context.Context, address string, names []string) ([]*models.TaskWithTaskHistory, error) {
	args := m.Called(ctx, address, names)
	return args.Get(0).([]*models.TaskWithTaskHistory), args.Error(1)
}

func (m *MockTaskRepository) ClaimDueSharePoolTask(ctx context.Context, now time.Time) (*entities.Task, error) {
	args := m.Called(ctx, now)
	return getTask(args, 0), args.Error(1)
}

func (m *MockTaskRepository) MarkSettled(ctx context.Context, id int64, claimStartedAt time.Time, settledAt time.Time) error {
	args := m.Called(ctx, id, claimStartedAt, settledAt)
	return args.Error(0)
}

func (m *MockTaskRepository) ReleaseSettlementClaim(ctx context.Context, id int64, claimStartedAt time.Time) error {
	args := m.Called(ctx, id, claimStartedAt)
	return args.Error(0)
}

func getTask(args mock.Arguments, index int) *entities.Task {
	value := args.Get(index)
	if value == nil {
		return nil
	}

	return value.(*entities.Task)
}
