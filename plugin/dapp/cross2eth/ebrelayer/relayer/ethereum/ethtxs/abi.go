package ethtxs

import (
	"strings"

	"github.com/assetcloud/plugin/plugin/dapp/cross2eth/contracts/contracts4eth/generated"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

// const
const (
	BridgeBankABI   = "BridgeBankABI"
	ChainBankABI    = "ChainBankABI"
	ChainBridgeABI  = "ChainBridgeABI"
	EthereumBankABI = "EthereumBankABI"
	OracleABI       = "OracleABI"
)

// LoadABI ...
func LoadABI(contractName string) abi.ABI {
	var abiJSON string
	switch contractName {
	case BridgeBankABI:
		abiJSON = generated.BridgeBankABI
	case ChainBankABI:
		abiJSON = generated.ChainBankABI
	case ChainBridgeABI:
		abiJSON = generated.ChainBridgeABI
	case EthereumBankABI:
		abiJSON = generated.EthereumBankABI
	case OracleABI:
		abiJSON = generated.OracleABI
	default:
		panic("No abi matched")
	}

	// Convert the raw abi into a usable format
	contractABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		panic(err)
	}

	return contractABI
}
