package services

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"trading-ace/config"
	"trading-ace/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type IEthereumService interface {
	SubscribeEthereumSwap() error
}

type EthereumService struct {
	logger logger.ILogger
	config *config.Config
}

const wethDecimals int64 = 18
const usdcDecimals int64 = 6

func NewEthereumService(logger logger.ILogger, config *config.Config) IEthereumService {
	return &EthereumService{
		logger: logger,
		config: config,
	}
}

func (e *EthereumService) SubscribeEthereumSwap() error {
	url := fmt.Sprintf("wss://mainnet.infura.io/ws/v3/%s", e.config.Infura.Key)
	contractAddressHex := "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc"

	client, err := ethclient.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

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

	contractAddress := common.HexToAddress(contractAddressHex)

	parsedABI, err := abi.JSON(strings.NewReader(swapEventAbi))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %v", err)
	}

	eventSignature := "Swap(address,uint256,uint256,uint256,uint256,address)"
	eventSignatureHash := crypto.Keccak256Hash([]byte(eventSignature))

	// 計算事件簽名的哈希
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{eventSignatureHash}},
	}

	logsCh := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logsCh)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %v", err)
	}

	defer sub.Unsubscribe()

	for vLog := range logsCh {
		event := struct {
			Sender     common.Address
			Amount0In  *big.Int
			Amount1In  *big.Int
			Amount0Out *big.Int
			Amount1Out *big.Int
			To         common.Address
		}{}

		err := parsedABI.UnpackIntoInterface(&event, "Swap", vLog.Data)
		if err != nil {
			e.logger.Error("Failed to unpack log: %v", err)
			continue
		}

		e.logger.Info("Sender: %s", vLog.Topics[1].Hex()[26:])
		e.logger.Info("To: %s", vLog.Topics[2].Hex()[26:])

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
	}

	return nil
}
