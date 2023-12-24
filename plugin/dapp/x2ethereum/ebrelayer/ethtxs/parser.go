package ethtxs

// --------------------------------------------------------
//      Parser
//
//      Parses structs containing event information into
//      unsigned transactions for validators to sign, then
//      relays the data packets as transactions on the
//      chain Bridge.
// --------------------------------------------------------

import (
	"math/big"
	"strings"

	chainTypes "github.com/assetcloud/chain/types"
	"github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/events"
	ebrelayerTypes "github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/types"
	"github.com/assetcloud/plugin/plugin/dapp/x2ethereum/types"
	"github.com/ethereum/go-ethereum/common"
)

// LogLockToEthBridgeClaim : parses and packages a LockEvent struct with a validator address in an EthBridgeClaim msg
func LogLockToEthBridgeClaim(event *events.LockEvent, ethereumChainID int64, bridgeBrankAddr string, decimal int64) (*ebrelayerTypes.EthBridgeClaim, error) {
	recipient := event.To
	if 0 == len(recipient) {
		return nil, ebrelayerTypes.ErrEmptyAddress
	}
	// Symbol formatted to lowercase
	symbol := strings.ToLower(event.Symbol)
	if symbol == "eth" && event.Token != common.HexToAddress("0x0000000000000000000000000000000000000000") {
		return nil, ebrelayerTypes.ErrAddress4Eth
	}

	witnessClaim := &ebrelayerTypes.EthBridgeClaim{}
	witnessClaim.EthereumChainID = ethereumChainID
	witnessClaim.BridgeBrankAddr = bridgeBrankAddr
	witnessClaim.Nonce = event.Nonce.Int64()
	witnessClaim.TokenAddr = event.Token.String()
	witnessClaim.Symbol = event.Symbol
	witnessClaim.EthereumSender = event.From.String()
	witnessClaim.ChainReceiver = string(recipient)

	if decimal > 8 {
		event.Value = event.Value.Quo(event.Value, big.NewInt(int64(types.MultiplySpecifyTimes(1, decimal-8))))
	} else {
		event.Value = event.Value.Mul(event.Value, big.NewInt(int64(types.MultiplySpecifyTimes(1, 8-decimal))))
	}
	witnessClaim.Amount = event.Value.String()

	witnessClaim.ClaimType = types.LockClaimType
	witnessClaim.ChainName = types.LockClaim
	witnessClaim.Decimal = decimal

	return witnessClaim, nil
}

// LogBurnToEthBridgeClaim ...
func LogBurnToEthBridgeClaim(event *events.BurnEvent, ethereumChainID int64, bridgeBrankAddr string, decimal int64) (*ebrelayerTypes.EthBridgeClaim, error) {
	recipient := event.ChainReceiver
	if 0 == len(recipient) {
		return nil, ebrelayerTypes.ErrEmptyAddress
	}

	witnessClaim := &ebrelayerTypes.EthBridgeClaim{}
	witnessClaim.EthereumChainID = ethereumChainID
	witnessClaim.BridgeBrankAddr = bridgeBrankAddr
	witnessClaim.Nonce = event.Nonce.Int64()
	witnessClaim.TokenAddr = event.Token.String()
	witnessClaim.Symbol = event.Symbol
	witnessClaim.EthereumSender = event.OwnerFrom.String()
	witnessClaim.ChainReceiver = string(recipient)
	witnessClaim.Amount = event.Amount.String()
	witnessClaim.ClaimType = types.BurnClaimType
	witnessClaim.ChainName = types.BurnClaim
	witnessClaim.Decimal = decimal

	return witnessClaim, nil
}

// ParseBurnLockTxReceipt : parses data from a Burn/Lock event witnessed on chain into a ChainMsg struct
func ParseBurnLockTxReceipt(claimType events.Event, receipt *chainTypes.ReceiptData) *events.ChainMsg {
	// Set up variables
	var chainSender []byte
	var ethereumReceiver, tokenContractAddress common.Address
	var symbol string
	var amount *big.Int

	// Iterate over attributes
	for _, log := range receipt.Logs {
		if log.Ty == types.TyChainToEthLog || log.Ty == types.TyWithdrawChainLog {
			txslog.Debug("ParseBurnLockTxReceipt", "value", string(log.Log))
			var chainToEth types.ReceiptChainToEth
			err := chainTypes.Decode(log.Log, &chainToEth)
			if err != nil {
				return nil
			}
			chainSender = []byte(chainToEth.ChainSender)
			ethereumReceiver = common.HexToAddress(chainToEth.EthereumReceiver)
			tokenContractAddress = common.HexToAddress(chainToEth.TokenContract)
			symbol = chainToEth.IssuerDotSymbol
			chainToEth.Amount = types.TrimZeroAndDot(chainToEth.Amount)
			amount = big.NewInt(1)
			amount, _ = amount.SetString(chainToEth.Amount, 10)
			if chainToEth.Decimals > 8 {
				amount = amount.Mul(amount, big.NewInt(int64(types.MultiplySpecifyTimes(1, chainToEth.Decimals-8))))
			} else {
				amount = amount.Quo(amount, big.NewInt(int64(types.MultiplySpecifyTimes(1, 8-chainToEth.Decimals))))
			}

			txslog.Info("ParseBurnLockTxReceipt", "chainSender", chainSender, "ethereumReceiver", ethereumReceiver.String(), "tokenContractAddress", tokenContractAddress.String(), "symbol", symbol, "amount", amount.String())
			// Package the event data into a ChainMsg
			chainMsg := events.NewChainMsg(claimType, chainSender, ethereumReceiver, symbol, amount, tokenContractAddress)
			return &chainMsg
		}
	}
	return nil
}

// ChainMsgToProphecyClaim : parses event data from a ChainMsg, packaging it as a ProphecyClaim
func ChainMsgToProphecyClaim(event events.ChainMsg) ProphecyClaim {
	claimType := event.ClaimType
	chainSender := event.ChainSender
	ethereumReceiver := event.EthereumReceiver
	tokenContractAddress := event.TokenContractAddress
	symbol := strings.ToLower(event.Symbol)
	amount := event.Amount

	prophecyClaim := ProphecyClaim{
		ClaimType:            claimType,
		ChainSender:          chainSender,
		EthereumReceiver:     ethereumReceiver,
		TokenContractAddress: tokenContractAddress,
		Symbol:               symbol,
		Amount:               amount,
	}

	return prophecyClaim
}
