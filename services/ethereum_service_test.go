package services

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEthereumClient struct {
	mock.Mock
}

func (m *MockEthereumClient) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, logsCh chan<- types.Log) (ethereum.Subscription, error) {
	args := m.Called(ctx, query, logsCh)
	return args.Get(0).(ethereum.Subscription), args.Error(1)
}

type MockSubscription struct {
	mock.Mock
}

func (m *MockSubscription) Unsubscribe() {
	m.Called()
}

func (m *MockSubscription) Err() <-chan error {
	args := m.Called()
	return args.Get(0).(<-chan error)
}

func TestParseABI(t *testing.T) {
	e := &EthereumService{}

	parsedABI, err := e.parseABI()

	assert.NoError(t, err, "Expected no error when parsing ABI")

	assert.NotNil(t, parsedABI, "Expected ABI object to be valid")

	parsedEvent, exists := parsedABI.Events["Swap"]
	assert.True(t, exists, "Expected 'Swap' event to be found in the ABI")
	assert.Equal(t, 6, len(parsedEvent.Inputs), "Expected 6 inputs in 'Swap' event")
}

func TestSubscribeToSwapEvent(t *testing.T) {
	// 創建 Mock 物件
	mockClient := new(MockEthereumClient)
	mockSubscription := new(MockSubscription)

	// 設定期望：我們希望 SubscribeFilterLogs 被呼叫，並且返回一個 Subscription 和 nil 錯誤
	mockClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Return(mockSubscription, nil)

	// 創建 EthereumService 實例
	e := &EthereumService{}

	// 呼叫 subscribeToSwapEvent 方法
	logsCh, sub, err := e.subscribeToSwapEvent(mockClient)

	// 驗證結果
	assert.NoError(t, err, "expected no error")
	assert.NotNil(t, logsCh, "expected logs channel to be returned")
	assert.NotNil(t, sub, "expected subscription to be returned")

	// 驗證 Mock 方法被正確呼叫
	mockClient.AssertExpectations(t)
}
