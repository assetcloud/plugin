package ethtxs

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/assetcloud/chain/common/log/log15"
	"github.com/assetcloud/plugin/plugin/dapp/cross2eth/contracts/contracts4eth/generated"
	"github.com/ethereum/go-ethereum/common"
)

var (
	deployLog = log15.New("contract deployer", "deployer")
)

//DeployResult ...
type DeployResult struct {
	Address common.Address
	TxHash  string
}

//X2EthContracts ...
type X2EthContracts struct {
	BridgeRegistry *generated.BridgeRegistry
	BridgeBank     *generated.BridgeBank
	ChainBridge    *generated.ChainBridge
	Valset         *generated.Valset
	Oracle         *generated.Oracle
}

//X2EthDeployResult ...
type X2EthDeployInfo struct {
	BridgeRegistry *DeployResult
	BridgeBank     *DeployResult
	ChainBridge    *DeployResult
	Valset         *DeployResult
	Oracle         *DeployResult
}

//DeployPara ...
type DeployPara struct {
	DeployPrivateKey *ecdsa.PrivateKey
	Deployer         common.Address
	Operator         common.Address
	InitValidators   []common.Address
	ValidatorPriKey  []*ecdsa.PrivateKey
	InitPowers       []*big.Int
}

//OperatorInfo ...
type OperatorInfo struct {
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
}
