package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type SwapEvent struct {
	SenderAddress string
	TxHash        common.Hash
	LogIndex      uint
	Amount0In     *big.Int
	Amount1In     *big.Int
	Amount0Out    *big.Int
	Amount1Out    *big.Int
}
