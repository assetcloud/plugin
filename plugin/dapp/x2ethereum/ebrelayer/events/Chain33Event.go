package events

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// ChainMsg : contains data from MsgBurn and MsgLock events
type ChainMsg struct {
	ClaimType            Event
	ChainSender        []byte
	EthereumReceiver     common.Address
	TokenContractAddress common.Address
	Symbol               string
	Amount               *big.Int
}

// NewChainMsg : creates a new ChainMsg
func NewChainMsg(
	claimType Event,
	chainSender []byte,
	ethereumReceiver common.Address,
	symbol string,
	amount *big.Int,
	tokenContractAddress common.Address,
) ChainMsg {
	// Package data into a ChainMsg
	chainMsg := ChainMsg{
		ClaimType:            claimType,
		ChainSender:        chainSender,
		EthereumReceiver:     ethereumReceiver,
		Symbol:               symbol,
		Amount:               amount,
		TokenContractAddress: tokenContractAddress,
	}

	return chainMsg
}
