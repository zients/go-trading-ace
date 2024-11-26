package services

import (
	"fmt"
	"math/big"
	"testing"
	"trading-ace/mocks"
	"trading-ace/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
	mockClient := new(mocks.MockEthereumClient)
	mockSubscription := new(mocks.MockEthereumSubscription)

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

func TestRetrieveEventData(t *testing.T) {
	// Set up mocks
	mockABI := new(mocks.MockABI)

	// Sample event data
	vLog := types.Log{
		Topics: []common.Hash{
			{},                                  // Event signature
			common.HexToHash("0xSenderAddress"), // Sender Address (mocked)
		},
		Data: []byte{0x01, 0x02}, // Mock event data (this will be used by ABI unpacking)
	}

	// Create a mock SwapEvent struct (models.SwapEvent)
	expectedEvent := models.SwapEvent{}

	// Mock UnpackIntoInterface behavior to simulate successful unpacking
	mockABI.On("UnpackIntoInterface", &expectedEvent, "Swap", vLog.Data).Return(nil)

	// Create EthereumService instance
	e := &EthereumService{
		logger:          new(mocks.MockLogger),          // assuming mockLogger is set up elsewhere
		campaignService: new(mocks.MockCampaignService), // assuming mockCampaignService is set up elsewhere
	}

	// Test retrieveEventData
	event, err := e.retrieveEventData(vLog, mockABI)
	assert.NoError(t, err, "expected no error")

	// Verify that the returned event matches expectations
	assert.Equal(t, event.SenderAddress, vLog.Topics[1].Hex()[26:])
	assert.Equal(t, event.Amount0In, expectedEvent.Amount0In)
	assert.Equal(t, event.Amount1In, expectedEvent.Amount1In)
	assert.Equal(t, event.Amount0Out, expectedEvent.Amount0Out)
	assert.Equal(t, event.Amount1Out, expectedEvent.Amount1Out)

	// Verify that UnpackIntoInterface was called with correct parameters
	mockABI.AssertExpectations(t)

	// Test error case for UnpackIntoInterface
	mockABI.On("UnpackIntoInterface", &models.SwapEvent{}, "Swap", vLog.Data).Return(fmt.Errorf("unpacking error"))
}

func TestProcessSwapEvent(t *testing.T) {
	// Set up mocks
	mockLogger := new(mocks.MockLogger)
	mockCampaignService := new(mocks.MockCampaignService)

	// Mock logger behavior
	mockLogger.On("Info", mock.Anything).Return()

	// Mock RecordUSDCSwapTotalAmount behavior
	mockCampaignService.On("RecordUSDCSwapTotalAmount", "0xSenderAddress", mock.Anything).Return(100.0, nil)

	// Create EthereumService instance
	e := &EthereumService{
		logger:          mockLogger,
		campaignService: mockCampaignService,
	}

	event := &models.SwapEvent{
		SenderAddress: "0xSenderAddress",
		Amount0In:     big.NewInt(10),
		Amount0Out:    big.NewInt(10),
		Amount1In:     big.NewInt(10),
		Amount1Out:    big.NewInt(10),
	}

	// Test processSwapEvent
	err := e.processSwapEvent(event)
	assert.NoError(t, err, "expected no error")

	// Verify expectations
	mockLogger.AssertExpectations(t)
	mockCampaignService.AssertExpectations(t)

	// Additional assertions for verifying specific behaviors
	mockCampaignService.AssertCalled(t, "RecordUSDCSwapTotalAmount", "0xSenderAddress", mock.Anything)
}
