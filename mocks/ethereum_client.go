package mocks

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
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

func (m *MockEthereumClient) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	args := m.Called(ctx, hash)

	var tx *types.Transaction
	if args.Get(0) != nil {
		tx = args.Get(0).(*types.Transaction)
	}

	return tx, args.Bool(1), args.Error(2)
}
