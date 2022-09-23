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
	// MsgBurn : Chain event 'Chain33Msg' type MsgBurn
	MsgBurn
	// MsgLock :  Chain event 'Chain33Msg' type MsgLock
	MsgLock
	// LogLock : Ethereum event 'LockEvent'
	LogLock
	// LogChain33TokenBurn : Ethereum event 'LogChain33TokenBurn' in contract chain33Bank
	LogChain33TokenBurn
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
	return [...]string{"unknown-x2ethereum", "Chain33ToEthBurn", "Chain33ToEthLock", "LogLock", "LogChain33TokenBurn", "LogNewProphecyClaim"}[d]
}

// Chain33MsgAttributeKey : enum containing supported attribute keys
type Chain33MsgAttributeKey int

const (
	// UnsupportedAttributeKey : unsupported attribute key
	UnsupportedAttributeKey Chain33MsgAttributeKey = iota
	// Chain33Sender : sender's address on Chain network
	Chain33Sender
	// EthereumReceiver : receiver's address on Ethereum network
	EthereumReceiver
	// Coin : coin type
	Coin
	// TokenContractAddress : coin's corresponding contract address deployed on the Ethereum network
	TokenContractAddress
)

// String : returns the event type as a string
func (d Chain33MsgAttributeKey) String() string {
	return [...]string{"unsupported", "chain33_sender", "ethereum_receiver", "amount", "token_contract_address"}[d]
}
