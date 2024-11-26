package services

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
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
}

type IEthereumService interface {
	SubscribeEthereumSwap() error
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

func NewEthereumService(logger logger.ILogger, config *config.Config, campaignService ICampaignService) IEthereumService {
	return &EthereumService{
		campaignService: campaignService,
		logger:          logger,
		config:          config,
	}
}

func (e *EthereumService) SubscribeEthereumSwap() error {
	client, err := e.connectToClient()
	if err != nil {
		return err
	}

	parsedABI, err := e.parseABI()
	if err != nil {
		return err
	}

	logsCh, sub, err := e.subscribeToSwapEvent(client)
	if err != nil {
		return err
	}

	defer sub.Unsubscribe()

	for vLog := range logsCh {
		event, err := e.retrieveEventData(vLog, parsedABI)
		if err != nil {
			e.logger.Error(err)
			continue
		}

		err = e.processSwapEvent(event)
		if err != nil {
			e.logger.Error(err)
		}
	}

	return nil
}

func (e *EthereumService) connectToClient() (*ethclient.Client, error) {
	url := fmt.Sprintf("wss://mainnet.infura.io/ws/v3/%s", e.config.Infura.Key)
	client, err := ethclient.Dial(url)
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

func (e *EthereumService) subscribeToSwapEvent(client IEthereumClient) (<-chan types.Log, ethereum.Subscription, error) {
	contractAddressHex := "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc"
	contractAddress := common.HexToAddress(contractAddressHex)
	eventSignature := "Swap(address,uint256,uint256,uint256,uint256,address)"
	eventSignatureHash := crypto.Keccak256Hash([]byte(eventSignature))

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{eventSignatureHash}},
	}

	logsCh := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logsCh)
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

	return &event, nil
}

func (e *EthereumService) processSwapEvent(event *models.SwapEvent) error {
	senderAddress := event.SenderAddress
	e.logger.Info("Sender: %s", senderAddress)

	// Convert amounts to float for easier logging
	amountInUSDC := new(big.Float).SetInt(event.Amount0In)
	amountInUSDC.Quo(amountInUSDC, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(usdcDecimals), nil)))
	e.logger.Info("Amount0In (USDC): %s", amountInUSDC.String())

	amountOutUSDC := new(big.Float).SetInt(event.Amount0Out)
	amountOutUSDC.Quo(amountOutUSDC, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(usdcDecimals), nil)))
	e.logger.Info("Amount0Out (USDC): %s", amountOutUSDC.String())

	amountInWETH := new(big.Float).SetInt(event.Amount1In)
	amountInWETH.Quo(amountInWETH, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(wethDecimals), nil)))
	e.logger.Info("Amount1In (WETH): %s", amountInWETH.String())

	amountOutWETH := new(big.Float).SetInt(event.Amount1Out)
	amountOutWETH.Quo(amountOutWETH, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(wethDecimals), nil)))
	e.logger.Info("Amount1Out (WETH): %s", amountOutWETH.String())

	// Record the campaign data asynchronously
	amountInUSDCFloat64, _ := amountInUSDC.Float64()
	amountOutUSDCFloat64, _ := amountOutUSDC.Float64()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		e.campaignService.RecordUSDCSwapTotalAmount(event.SenderAddress, amountInUSDCFloat64)
	}()

	go func() {
		defer wg.Done()
		e.campaignService.RecordUSDCSwapTotalAmount(event.SenderAddress, amountOutUSDCFloat64)
	}()

	wg.Wait()

	return nil
}
