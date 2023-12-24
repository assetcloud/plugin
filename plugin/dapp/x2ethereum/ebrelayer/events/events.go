package events

import (
	log "github.com/assetcloud/chain/common/log/log15"
)

// Event : enum containing supported contract events
type Event int

var eventsLog = log.New("module", "ethereum_relayer")

const (
	// Unsupported : unsupported Chain or Ethereum event
	Unsupported Event = iota
	// MsgBurn : Chain event 'ChainMsg' type MsgBurn
	MsgBurn
	// MsgLock :  Chain event 'ChainMsg' type MsgLock
	MsgLock
	// LogLock : Ethereum event 'LockEvent'
	LogLock
	// LogChainTokenBurn : Ethereum event 'LogChainTokenBurn' in contract chainBank
	LogChainTokenBurn
	// LogNewProphecyClaim : Ethereum event 'NewProphecyClaimEvent'
	LogNewProphecyClaim
)

//const
const (
	ClaimTypeBurn = uint8(1)
	ClaimTypeLock = uint8(2)
)

// String : returns the event type as a string
func (d Event) String() string {
	return [...]string{"unknown-x2ethereum", "ChainToEthBurn", "ChainToEthLock", "LogLock", "LogChainTokenBurn", "LogNewProphecyClaim"}[d]
}

// ChainMsgAttributeKey : enum containing supported attribute keys
type ChainMsgAttributeKey int

const (
	// UnsupportedAttributeKey : unsupported attribute key
	UnsupportedAttributeKey ChainMsgAttributeKey = iota
	// ChainSender : sender's address on Chain network
	ChainSender
	// EthereumReceiver : receiver's address on Ethereum network
	EthereumReceiver
	// Coin : coin type
	Coin
	// TokenContractAddress : coin's corresponding contract address deployed on the Ethereum network
	TokenContractAddress
)

// String : returns the event type as a string
func (d ChainMsgAttributeKey) String() string {
	return [...]string{"unsupported", "chain_sender", "ethereum_receiver", "amount", "token_contract_address"}[d]
}
