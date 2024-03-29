package ethtxs

// ------------------------------------------------------------
//	Relay : Builds and encodes EthBridgeClaim Msgs with the
//  	specified variables, before presenting the unsigned
//      transaction to validators for optional signing.
//      Once signed, the data packets are sent as transactions
//      on the chain Bridge.
// ------------------------------------------------------------

import (
	"github.com/assetcloud/chain/common"
	chainCrypto "github.com/assetcloud/chain/common/crypto"
	"github.com/assetcloud/chain/rpc/jsonclient"
	rpctypes "github.com/assetcloud/chain/rpc/types"
	chainTypes "github.com/assetcloud/chain/types"
	ebrelayerTypes "github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/types"
	"github.com/assetcloud/plugin/plugin/dapp/x2ethereum/types"
)

// RelayLockToChain : RelayLockToChain applies validator's signature to an EthBridgeClaim message
//		containing information about an event on the Ethereum blockchain before relaying to the Bridge
func RelayLockToChain(privateKey chainCrypto.PrivKey, claim *ebrelayerTypes.EthBridgeClaim, rpcURL string) (string, error) {
	var res string

	params := &types.Eth2Chain{
		EthereumChainID:       claim.EthereumChainID,
		BridgeContractAddress: claim.BridgeBrankAddr,
		Nonce:                 claim.Nonce,
		IssuerDotSymbol:       claim.Symbol,
		TokenContractAddress:  claim.TokenAddr,
		EthereumSender:        claim.EthereumSender,
		ChainReceiver:       claim.ChainReceiver,
		Amount:                claim.Amount,
		ClaimType:             int64(claim.ClaimType),
		Decimals:              claim.Decimal,
	}

	pm := rpctypes.CreateTxIn{
		Execer:     X2Eth,
		ActionName: types.NameEth2ChainAction,
		Payload:    chainTypes.MustPBToJSON(params),
	}
	ctx := jsonclient.NewRPCCtx(rpcURL, "Chain.CreateTransaction", pm, &res)
	_, _ = ctx.RunResult()

	data, err := common.FromHex(res)
	if err != nil {
		return "", err
	}
	var tx chainTypes.Transaction
	err = chainTypes.Decode(data, &tx)
	if err != nil {
		return "", err
	}

	if tx.Fee == 0 {
		tx.Fee, err = tx.GetRealFee(1e5)
		if err != nil {
			return "", err
		}
	}
	//构建交易，验证人validator用来向chain合约证明自己验证了该笔从以太坊向chain跨链转账的交易
	tx.Sign(chainTypes.SECP256K1, privateKey)

	txData := chainTypes.Encode(&tx)
	dataStr := common.ToHex(txData)
	pms := rpctypes.RawParm{
		Token: "BTY",
		Data:  dataStr,
	}
	var txhash string

	ctx = jsonclient.NewRPCCtx(rpcURL, "Chain.SendTransaction", pms, &txhash)
	_, err = ctx.RunResult()
	return txhash, err
}

//RelayBurnToChain ...
func RelayBurnToChain(privateKey chainCrypto.PrivKey, claim *ebrelayerTypes.EthBridgeClaim, rpcURL string) (string, error) {
	var res string

	params := &types.Eth2Chain{
		EthereumChainID:       claim.EthereumChainID,
		BridgeContractAddress: claim.BridgeBrankAddr,
		Nonce:                 claim.Nonce,
		IssuerDotSymbol:       claim.Symbol,
		TokenContractAddress:  claim.TokenAddr,
		EthereumSender:        claim.EthereumSender,
		ChainReceiver:       claim.ChainReceiver,
		Amount:                claim.Amount,
		ClaimType:             int64(claim.ClaimType),
		Decimals:              claim.Decimal,
	}

	pm := rpctypes.CreateTxIn{
		Execer:     X2Eth,
		ActionName: types.NameWithdrawEthAction,
		Payload:    chainTypes.MustPBToJSON(params),
	}
	ctx := jsonclient.NewRPCCtx(rpcURL, "Chain.CreateTransaction", pm, &res)
	_, _ = ctx.RunResult()

	data, err := common.FromHex(res)
	if err != nil {
		return "", err
	}
	var tx chainTypes.Transaction
	err = chainTypes.Decode(data, &tx)
	if err != nil {
		return "", err
	}

	if tx.Fee == 0 {
		tx.Fee, err = tx.GetRealFee(1e5)
		if err != nil {
			return "", err
		}
	}
	//构建交易，验证人validator用来向chain合约证明自己验证了该笔从以太坊向chain跨链转账的交易
	tx.Sign(chainTypes.SECP256K1, privateKey)

	txData := chainTypes.Encode(&tx)
	dataStr := common.ToHex(txData)
	pms := rpctypes.RawParm{
		Token: "BTY",
		Data:  dataStr,
	}
	var txhash string

	ctx = jsonclient.NewRPCCtx(rpcURL, "Chain.SendTransaction", pms, &txhash)
	_, err = ctx.RunResult()
	return txhash, err
}
