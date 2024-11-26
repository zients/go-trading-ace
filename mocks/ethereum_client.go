package mocks

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
)

type MockEthereumClient struct {
	mock.Mock
}

func (m *MockEthereumClient) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, logsCh chan<- types.Log) (ethereum.Subscription, error) {
	args := m.Called(ctx, query, logsCh)
	return args.Get(0).(ethereum.Subscription), args.Error(1)
}
