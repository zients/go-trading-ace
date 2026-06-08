package mocks

import (
	"context"
	"trading-ace/entities"
	"trading-ace/models"

	"github.com/stretchr/testify/mock"
)

type MockCampaignService struct {
	mock.Mock
}

func (m *MockCampaignService) StartCampaign(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCampaignService) GetPointHistories(ctx context.Context, address string) ([]*models.TaskTaskHistoryPair, error) {
	args := m.Called(ctx, address)
	return args.Get(0).([]*models.TaskTaskHistoryPair), args.Error(1)
}

func (m *MockCampaignService) RecordUSDCSwapTotalAmount(ctx context.Context, senderAddress string, amount float64) (float64, error) {
	args := m.Called(ctx, senderAddress, amount)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockCampaignService) GetTaskStatus(ctx context.Context, address string) ([]*models.TaskWithTaskHistory, error) {
	args := m.Called(ctx, address)
	return args.Get(0).([]*models.TaskWithTaskHistory), args.Error(1)
}

func (m *MockCampaignService) FindOnboardingTask(ctx context.Context) (*entities.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).(*entities.Task), args.Error(1)
}

func (m *MockCampaignService) FindCurrentSharePoolTask(ctx context.Context) (*entities.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).(*entities.Task), args.Error(1)
}

func (m *MockCampaignService) GetLeaderboard(ctx context.Context, taskName string, period int) ([]models.LeaderboardEntry, error) {
	args := m.Called(ctx, taskName, period)
	return args.Get(0).([]models.LeaderboardEntry), args.Error(1)
}
