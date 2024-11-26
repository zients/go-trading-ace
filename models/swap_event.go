package models

import "math/big"

type SwapEvent struct {
	SenderAddress string
	Amount0In     *big.Int
	Amount1In     *big.Int
	Amount0Out    *big.Int
	Amount1Out    *big.Int
}
