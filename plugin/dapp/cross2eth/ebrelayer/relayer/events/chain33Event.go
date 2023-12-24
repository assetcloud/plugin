package events

import (
	"errors"
	"math/big"

	ebrelayerTypes "github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/types"
	chainEvmCommon "github.com/assetcloud/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type ChainEvmEvent int

const (
	UnsupportedEvent ChainEvmEvent = iota
	//在chain的evm合约中产生了lock事件
	ChainEventLogLock
	//在chain的evm合约中产生了burn事件
	ChainEventLogBurn
	//在chain的evm合约中产生了withdraw事件
	ChainEventLogWithdraw
)

// String : returns the event type as a string
func (d ChainEvmEvent) String() string {
	return [...]string{"unknown-event", "LogLock", "LogEthereumTokenBurn", "LogEthereumTokenWithdraw"}[d]
}

// ChainMsg : contains data from MsgBurn and MsgLock events
type ChainMsg struct {
	ClaimType            ClaimType
	ChainSender        chainEvmCommon.Address
	EthereumReceiver     common.Address
	TokenContractAddress chainEvmCommon.Address
	Symbol               string
	Amount               *big.Int
	TxHash               []byte
	Nonce                int64
	ForwardTimes         int32
	ForwardIndex         int64
}

// 发生在chainevm上的lock事件，当bty跨链转移到eth时会发生该种事件
type LockEventOnChain struct {
	From   chainEvmCommon.Hash160Address
	To     []byte
	Token  chainEvmCommon.Hash160Address
	Symbol string
	Value  *big.Int
	Nonce  *big.Int
}

// 发生在chain evm上的withdraw事件，当用户发起通过代理人提币交易时，则弹射出该事件信息
type WithdrawEventOnChain struct {
	BridgeToken      chainEvmCommon.Hash160Address
	Symbol           string
	Amount           *big.Int
	OwnerFrom        chainEvmCommon.Hash160Address
	EthereumReceiver []byte
	ProxyReceiver    chainEvmCommon.Hash160Address
	Nonce            *big.Int
}

// 发生在chainevm上的burn事件，当eth/erc20资产需要提币回到以太坊链上时，会发生该种事件
type BurnEventOnChain struct {
	Token            chainEvmCommon.Hash160Address
	Symbol           string
	Amount           *big.Int
	OwnerFrom        chainEvmCommon.Hash160Address
	EthereumReceiver []byte
	Nonce            *big.Int
}

func UnpackChainLogLock(contractAbi abi.ABI, eventName string, eventData []byte) (lockEvent *LockEventOnChain, err error) {
	lockEvent = &LockEventOnChain{}
	// Parse the event's attributes as Ethereum network variables
	err = contractAbi.UnpackIntoInterface(lockEvent, eventName, eventData)
	if err != nil {
		eventsLog.Error("UnpackLogLock", "Failed to unpack abi due to:", err.Error())
		return nil, ebrelayerTypes.ErrUnpack
	}

	eventsLog.Info("UnpackLogLock", "value", lockEvent.Value.String(),
		"symbol", lockEvent.Symbol,
		"token addr on chain evm", lockEvent.Token.ToAddress().String(),
		"chain sender", lockEvent.From.ToAddress().String(),
		"ethereum recipient", common.BytesToAddress(lockEvent.To).String(),
		"nonce", lockEvent.Nonce.String())

	return lockEvent, nil
}

func UnpackChainLogBurn(contractAbi abi.ABI, eventName string, eventData []byte) (burnEvent *BurnEventOnChain, err error) {
	burnEvent = &BurnEventOnChain{}
	// Parse the event's attributes as Ethereum network variables
	err = contractAbi.UnpackIntoInterface(burnEvent, eventName, eventData)
	if err != nil {
		eventsLog.Error("UnpackLogBurn", "Failed to unpack abi due to:", err.Error())
		return nil, ebrelayerTypes.ErrUnpack
	}

	eventsLog.Info("UnpackLogBurn", "token addr on chain evm", burnEvent.Token.ToAddress().String(),
		"symbol", burnEvent.Symbol,
		"Amount", burnEvent.Amount.String(),
		"Owner address from chain", burnEvent.OwnerFrom.ToAddress().String(),
		"EthereumReceiver", common.BytesToAddress(burnEvent.EthereumReceiver).String(),
		"nonce", burnEvent.Nonce.String())
	return burnEvent, nil
}

func UnpackLogWithdraw(contractAbi abi.ABI, eventName string, eventData []byte) (withdrawEvent *WithdrawEventOnChain, err error) {
	withdrawEvent = &WithdrawEventOnChain{}
	err = contractAbi.UnpackIntoInterface(withdrawEvent, eventName, eventData)
	if err != nil {
		eventsLog.Error("UnpackLogWithdraw", "Failed to unpack abi due to:", err.Error())
		return nil, err
	}

	eventsLog.Info("UnpackLogWithdraw", "bridge token addr on chain evm", withdrawEvent.BridgeToken.ToAddress().String(),
		"symbol", withdrawEvent.Symbol,
		"Amount", withdrawEvent.Amount.String(),
		"Owner address from chain", withdrawEvent.OwnerFrom.ToAddress().String(),
		"EthereumReceiver", common.BytesToAddress(withdrawEvent.EthereumReceiver).String(),
		"ProxyReceiver", withdrawEvent.ProxyReceiver.ToAddress().String(),
		"nonce", withdrawEvent.Nonce.String())
	return withdrawEvent, nil
}

// ParseBurnLock4chain ParseBurnLockTxReceipt : parses data from a Burn/Lock/Withdraw event witnessed on chain into a ChainMsg struct
func ParseBurnLock4chain(evmEventType ChainEvmEvent, data []byte, bridgeBankAbi abi.ABI, chainTxHash []byte) (*ChainMsg, error) {
	if ChainEventLogLock == evmEventType {
		lockEvent, err := UnpackChainLogLock(bridgeBankAbi, evmEventType.String(), data)
		if nil != err {
			return nil, err
		}

		chainMsg := &ChainMsg{
			ClaimType:            ClaimTypeLock,
			ChainSender:        lockEvent.From.ToAddress(),
			EthereumReceiver:     common.BytesToAddress(lockEvent.To),
			TokenContractAddress: lockEvent.Token.ToAddress(),
			Symbol:               lockEvent.Symbol,
			Amount:               lockEvent.Value,
			TxHash:               chainTxHash,
			Nonce:                lockEvent.Nonce.Int64(),
		}
		return chainMsg, nil

	} else if ChainEventLogBurn == evmEventType {
		burnEvent, err := UnpackChainLogBurn(bridgeBankAbi, evmEventType.String(), data)
		if nil != err {
			return nil, err
		}

		chainMsg := &ChainMsg{
			ClaimType:            ClaimTypeBurn,
			ChainSender:        burnEvent.OwnerFrom.ToAddress(),
			EthereumReceiver:     common.BytesToAddress(burnEvent.EthereumReceiver),
			TokenContractAddress: burnEvent.Token.ToAddress(),
			Symbol:               burnEvent.Symbol,
			Amount:               burnEvent.Amount,
			TxHash:               chainTxHash,
			Nonce:                burnEvent.Nonce.Int64(),
		}
		return chainMsg, nil
	} else if ChainEventLogWithdraw == evmEventType {
		burnEvent, err := UnpackLogWithdraw(bridgeBankAbi, evmEventType.String(), data)
		if nil != err {
			return nil, err
		}

		chainMsg := &ChainMsg{
			ClaimType:            ClaimTypeWithdraw,
			ChainSender:        burnEvent.OwnerFrom.ToAddress(),
			EthereumReceiver:     common.BytesToAddress(burnEvent.EthereumReceiver),
			TokenContractAddress: burnEvent.BridgeToken.ToAddress(),
			Symbol:               burnEvent.Symbol,
			Amount:               burnEvent.Amount,
			TxHash:               chainTxHash,
			Nonce:                burnEvent.Nonce.Int64(),
		}
		return chainMsg, nil
	}

	return nil, errors.New("unknown-event")
}
