package mocks

import (
	"trading-ace/models"

	"github.com/stretchr/testify/mock"
)

type MockCampaignService struct {
	mock.Mock
}

func (m *MockCampaignService) StartCampaign() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCampaignService) GetPointHistories(address string) ([]*models.TaskTaskHistoryPair, error) {
	args := m.Called(address)
	return args.Get(0).([]*models.TaskTaskHistoryPair), args.Error(1)
}

func (m *MockCampaignService) RecordUSDCSwapTotalAmount(senderAddress string, amount float64) (float64, error) {
	args := m.Called(senderAddress, amount)
	return args.Get(0).(float64), args.Error(1)
}
