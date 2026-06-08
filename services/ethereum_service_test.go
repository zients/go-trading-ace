package services

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"trading-ace/mocks"
	"trading-ace/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
	logsCh, sub, err := e.subscribeToSwapEvent(context.Background(), mockClient)

	// 驗證結果
	assert.NoError(t, err, "expected no error")
	assert.NotNil(t, logsCh, "expected logs channel to be returned")
	assert.NotNil(t, sub, "expected subscription to be returned")

	// 驗證 Mock 方法被正確呼叫
	mockClient.AssertExpectations(t)
}

func TestSubscribeToSwapEventUsesProvidedContext(t *testing.T) {
	mockClient := new(mocks.MockEthereumClient)
	mockSubscription := new(mocks.MockEthereumSubscription)
	ctx := context.WithValue(context.Background(), struct{}{}, "ethereum-listener")

	mockClient.On("SubscribeFilterLogs", mock.MatchedBy(func(actual context.Context) bool {
		return actual == ctx
	}), mock.Anything, mock.Anything).Return(mockSubscription, nil).Once()

	e := &EthereumService{}

	logsCh, sub, err := e.subscribeToSwapEvent(ctx, mockClient)

	assert.NoError(t, err)
	assert.NotNil(t, logsCh)
	assert.NotNil(t, sub)
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
	assert.Equal(t, vLog.TxHash, event.TxHash)

	// Verify that UnpackIntoInterface was called with correct parameters
	mockABI.AssertExpectations(t)

	// Test error case for UnpackIntoInterface
	mockABI.On("UnpackIntoInterface", &models.SwapEvent{}, "Swap", vLog.Data).Return(fmt.Errorf("unpacking error"))
}

func TestProcessSwapEvent(t *testing.T) {
	// Set up mocks
	mockLogger := new(mocks.MockLogger)
	mockCampaignService := new(mocks.MockCampaignService)
	mockClient := new(mocks.MockEthereumClient)
	ctx := context.WithValue(context.Background(), struct{}{}, "swap-context")
	txHash := common.HexToHash("0x1234")
	signedTx, transactionSender := signedMainnetTransaction(t)

	// Mock logger behavior
	mockLogger.On("Info", mock.Anything).Return()

	// Mock RecordUSDCSwapTotalAmount behavior
	mockClient.On("TransactionByHash", sameContext(ctx), txHash).Return(signedTx, false, nil).Once()
	mockCampaignService.On("RecordUSDCSwapTotalAmount", sameContext(ctx), transactionSender.Hex(), 0.00002).Return(100.0, nil).Once()

	// Create EthereumService instance
	e := &EthereumService{
		logger:          mockLogger,
		campaignService: mockCampaignService,
	}

	event := &models.SwapEvent{
		SenderAddress: "0xRouterAddress",
		TxHash:        txHash,
		Amount0In:     big.NewInt(10),
		Amount0Out:    big.NewInt(10),
		Amount1In:     big.NewInt(10),
		Amount1Out:    big.NewInt(10),
	}

	// Test processSwapEvent
	err := e.processSwapEvent(ctx, mockClient, event)
	assert.NoError(t, err, "expected no error")

	// Verify expectations
	mockLogger.AssertExpectations(t)
	mockClient.AssertExpectations(t)
	mockCampaignService.AssertExpectations(t)

	// Additional assertions for verifying specific behaviors
	mockCampaignService.AssertNumberOfCalls(t, "RecordUSDCSwapTotalAmount", 1)
}

func TestProcessSwapEventReturnsErrorWhenCampaignRecordingFails(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	mockCampaignService := new(mocks.MockCampaignService)
	mockClient := new(mocks.MockEthereumClient)
	ctx := context.Background()
	txHash := common.HexToHash("0xabcd")
	signedTx, transactionSender := signedMainnetTransaction(t)

	mockLogger.On("Info", mock.Anything).Return()
	mockClient.On("TransactionByHash", mock.Anything, txHash).Return(signedTx, false, nil).Once()
	mockCampaignService.On("RecordUSDCSwapTotalAmount", mock.Anything, transactionSender.Hex(), 0.00002).
		Return(0.0, fmt.Errorf("record failed")).Once()

	e := &EthereumService{
		logger:          mockLogger,
		campaignService: mockCampaignService,
	}

	event := &models.SwapEvent{
		SenderAddress: "0xRouterAddress",
		TxHash:        txHash,
		Amount0In:     big.NewInt(10),
		Amount0Out:    big.NewInt(10),
		Amount1In:     big.NewInt(10),
		Amount1Out:    big.NewInt(10),
	}

	err := e.processSwapEvent(ctx, mockClient, event)

	assert.ErrorContains(t, err, "record failed")
	mockLogger.AssertExpectations(t)
	mockClient.AssertExpectations(t)
	mockCampaignService.AssertExpectations(t)
}

func TestProcessSwapEventRecordsTransactionSenderOnceWithCombinedUSDCAmount(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	mockCampaignService := new(mocks.MockCampaignService)
	mockClient := new(mocks.MockEthereumClient)
	ctx := context.WithValue(context.Background(), struct{}{}, "participant-context")
	txHash := common.HexToHash("0xbeef")
	signedTx, transactionSender := signedMainnetTransaction(t)

	mockLogger.On("Info", mock.Anything).Return()
	mockClient.On("TransactionByHash", sameContext(ctx), txHash).Return(signedTx, false, nil).Once()
	mockCampaignService.On("RecordUSDCSwapTotalAmount", sameContext(ctx), transactionSender.Hex(), 1000.0).
		Return(1000.0, nil).Once()

	e := &EthereumService{
		logger:          mockLogger,
		campaignService: mockCampaignService,
	}

	event := &models.SwapEvent{
		SenderAddress: "0x000000000000000000000000000000000000dEaD",
		TxHash:        txHash,
		Amount0In:     big.NewInt(250_000000),
		Amount0Out:    big.NewInt(750_000000),
		Amount1In:     big.NewInt(0),
		Amount1Out:    big.NewInt(0),
	}

	err := e.processSwapEvent(ctx, mockClient, event)

	assert.NoError(t, err)
	mockCampaignService.AssertExpectations(t)
	mockCampaignService.AssertNumberOfCalls(t, "RecordUSDCSwapTotalAmount", 1)
	mockClient.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestProcessSwapEventSkipsZeroUSDCAmount(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	mockCampaignService := new(mocks.MockCampaignService)
	mockClient := new(mocks.MockEthereumClient)
	txHash := common.HexToHash("0xcafe")

	mockLogger.On("Info", mock.Anything).Return()

	e := &EthereumService{
		logger:          mockLogger,
		campaignService: mockCampaignService,
	}

	event := &models.SwapEvent{
		SenderAddress: "0x000000000000000000000000000000000000dEaD",
		TxHash:        txHash,
		Amount0In:     big.NewInt(0),
		Amount0Out:    big.NewInt(0),
		Amount1In:     big.NewInt(10),
		Amount1Out:    big.NewInt(10),
	}

	err := e.processSwapEvent(context.Background(), mockClient, event)

	assert.NoError(t, err)
	mockCampaignService.AssertNotCalled(t, "RecordUSDCSwapTotalAmount", mock.Anything, mock.Anything, mock.Anything)
	mockClient.AssertNotCalled(t, "TransactionByHash", mock.Anything, mock.Anything)
	mockLogger.AssertExpectations(t)
}

func TestProcessSwapEventReturnsErrorWhenTransactionIsMissing(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	mockCampaignService := new(mocks.MockCampaignService)
	mockClient := new(mocks.MockEthereumClient)
	txHash := common.HexToHash("0xfeed")

	mockLogger.On("Info", mock.Anything).Return()
	mockClient.On("TransactionByHash", mock.Anything, txHash).Return(nil, false, nil).Once()

	e := &EthereumService{
		logger:          mockLogger,
		campaignService: mockCampaignService,
	}

	event := &models.SwapEvent{
		SenderAddress: "0x000000000000000000000000000000000000dEaD",
		TxHash:        txHash,
		Amount0In:     big.NewInt(1),
		Amount0Out:    big.NewInt(0),
		Amount1In:     big.NewInt(0),
		Amount1Out:    big.NewInt(0),
	}

	err := e.processSwapEvent(context.Background(), mockClient, event)

	assert.ErrorContains(t, err, "swap transaction not found")
	mockCampaignService.AssertNotCalled(t, "RecordUSDCSwapTotalAmount", mock.Anything, mock.Anything, mock.Anything)
	mockClient.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func signedMainnetTransaction(t *testing.T) (*types.Transaction, common.Address) {
	t.Helper()

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	transactionSender := crypto.PubkeyToAddress(privateKey.PublicKey)
	tx := types.NewTransaction(
		0,
		common.HexToAddress("0x0000000000000000000000000000000000000001"),
		big.NewInt(0),
		21_000,
		big.NewInt(1),
		nil,
	)

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(big.NewInt(1)), privateKey)
	require.NoError(t, err)

	return signedTx, transactionSender
}
