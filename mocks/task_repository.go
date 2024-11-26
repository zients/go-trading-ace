package mocks

import (
	"trading-ace/entities"
	"trading-ace/models"

	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(task *entities.Task) (*entities.Task, error) {
	args := m.Called(task)
	return args.Get(0).(*entities.Task), args.Error(1)
}

func (m *MockTaskRepository) FindById(id int64) (*entities.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Task), args.Error(1)
}

func (m *MockTaskRepository) FindByName(name string) (*entities.Task, error) {
	args := m.Called(name)
	return args.Get(0).(*entities.Task), args.Error(1)
}

func (m *MockTaskRepository) GetByName(name string) ([]*entities.Task, error) {
	args := m.Called(name)
	return args.Get(0).([]*entities.Task), args.Error(1)
}

func (m *MockTaskRepository) IsExistedByName(name string) (bool, error) {
	args := m.Called(name)
	return args.Bool(0), args.Error(1)
}

func (m *MockTaskRepository) GetByAddressAndNamesIncludingTaskHistories(address string, names []string) ([]*models.TaskWithTaskHistory, error) {
	args := m.Called(address, names)
	return args.Get(0).([]*models.TaskWithTaskHistory), args.Error(1)
}
