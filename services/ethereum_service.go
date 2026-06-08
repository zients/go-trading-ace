package services

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"trading-ace/config"
	"trading-ace/logger"
	"trading-ace/models"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type IEthereumClient interface {
	SubscribeFilterLogs(context.Context, ethereum.FilterQuery, chan<- types.Log) (ethereum.Subscription, error)
	TransactionByHash(context.Context, common.Hash) (*types.Transaction, bool, error)
}

type IEthereumService interface {
	SubscribeEthereumSwap(ctx context.Context) error
}

type IABI interface {
	UnpackIntoInterface(interface{}, string, []byte) error
}

type EthereumService struct {
	campaignService ICampaignService
	logger          logger.ILogger
	config          *config.Config
}

const wethDecimals int64 = 18
const usdcDecimals int64 = 6
const ethereumMainnetChainID int64 = 1

func NewEthereumService(logger logger.ILogger, config *config.Config, campaignService ICampaignService) IEthereumService {
	return &EthereumService{
		campaignService: campaignService,
		logger:          logger,
		config:          config,
	}
}

func (e *EthereumService) SubscribeEthereumSwap(ctx context.Context) error {
	client, err := e.connectToClient(ctx)
	if err != nil {
		return err
	}

	parsedABI, err := e.parseABI()
	if err != nil {
		return err
	}

	logsCh, sub, err := e.subscribeToSwapEvent(ctx, client)
	if err != nil {
		return err
	}

	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-sub.Err():
			return fmt.Errorf("ethereum subscription error: %w", err)
		case vLog, ok := <-logsCh:
			if !ok {
				return nil
			}

			event, err := e.retrieveEventData(vLog, parsedABI)
			if err != nil {
				e.logger.Error(err)
				continue
			}

			err = e.processSwapEvent(ctx, client, event)
			if err != nil {
				e.logger.Error(err)
			}
		}
	}
}

func (e *EthereumService) connectToClient(ctx context.Context) (*ethclient.Client, error) {
	url := fmt.Sprintf("wss://mainnet.infura.io/ws/v3/%s", e.config.Infura.Key)
	client, err := ethclient.DialContext(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	return client, nil
}

func (e *EthereumService) parseABI() (abi.ABI, error) {
	swapEventAbi := `
	[
		{
			"anonymous": false,
			"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount0In",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount1In",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount0Out",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount1Out",
				"type": "uint256"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "to",
				"type": "address"
			}
			],
			"name": "Swap",
			"type": "event"
		}
	]
`

	parsedABI, err := abi.JSON(strings.NewReader(swapEventAbi))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to parse ABI: %v", err)
	}

	return parsedABI, nil
}

func (e *EthereumService) subscribeToSwapEvent(ctx context.Context, client IEthereumClient) (<-chan types.Log, ethereum.Subscription, error) {
	contractAddressHex := "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc"
	contractAddress := common.HexToAddress(contractAddressHex)
	eventSignature := "Swap(address,uint256,uint256,uint256,uint256,address)"
	eventSignatureHash := crypto.Keccak256Hash([]byte(eventSignature))

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{eventSignatureHash}},
	}

	logsCh := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to subscribe: %v", err)
	}

	return logsCh, sub, nil
}

func (e *EthereumService) retrieveEventData(vLog types.Log, parsedABI IABI) (*models.SwapEvent, error) {
	event := models.SwapEvent{}

	err := parsedABI.UnpackIntoInterface(&event, "Swap", vLog.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack log: %v", err)
	}

	event.SenderAddress = vLog.Topics[1].Hex()[26:]
	event.TxHash = vLog.TxHash

	return &event, nil
}

func (e *EthereumService) processSwapEvent(ctx context.Context, client IEthereumClient, event *models.SwapEvent) error {
	senderAddress := event.SenderAddress
	e.logger.Info("Sender: %s", senderAddress)

	// Convert amounts to float for easier logging
	amountInUSDC := tokenAmountFloat(event.Amount0In, usdcDecimals)
	e.logger.Info("Amount0In (USDC): %s", amountInUSDC.String())

	amountOutUSDC := tokenAmountFloat(event.Amount0Out, usdcDecimals)
	e.logger.Info("Amount0Out (USDC): %s", amountOutUSDC.String())

	amountInWETH := tokenAmountFloat(event.Amount1In, wethDecimals)
	e.logger.Info("Amount1In (WETH): %s", amountInWETH.String())

	amountOutWETH := tokenAmountFloat(event.Amount1Out, wethDecimals)
	e.logger.Info("Amount1Out (WETH): %s", amountOutWETH.String())

	amountInUSDCFloat64, _ := amountInUSDC.Float64()
	amountOutUSDCFloat64, _ := amountOutUSDC.Float64()
	usdcAmount := amountInUSDCFloat64 + amountOutUSDCFloat64
	if usdcAmount == 0 {
		return nil
	}

	participantAddress, err := e.resolveTransactionSender(ctx, client, event.TxHash)
	if err != nil {
		return err
	}

	if _, err := e.campaignService.RecordUSDCSwapTotalAmount(ctx, participantAddress.Hex(), usdcAmount); err != nil {
		return fmt.Errorf("failed to record swap event campaign data: %w", err)
	}

	return nil
}

func (e *EthereumService) resolveTransactionSender(ctx context.Context, client IEthereumClient, txHash common.Hash) (common.Address, error) {
	tx, _, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to fetch swap transaction %s: %w", txHash.Hex(), err)
	}

	if tx == nil {
		return common.Address{}, fmt.Errorf("swap transaction not found: %s", txHash.Hex())
	}

	signer := types.LatestSignerForChainID(big.NewInt(ethereumMainnetChainID))
	sender, err := types.Sender(signer, tx)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to recover swap transaction sender %s: %w", txHash.Hex(), err)
	}

	return sender, nil
}

func tokenAmountFloat(amount *big.Int, decimals int64) *big.Float {
	if amount == nil {
		return big.NewFloat(0)
	}

	value := new(big.Float).SetInt(amount)
	denominator := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil))
	value.Quo(value, denominator)

	return value
}
