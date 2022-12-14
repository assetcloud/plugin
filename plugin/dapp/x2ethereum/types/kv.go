package types

import "strings"

/*
 * 用户合约存取kv数据时，key值前缀需要满足一定规范
 * 即key = keyPrefix + userKey
 * 需要字段前缀查询时，使用’-‘作为分割符号
 */

var (
	//KeyPrefixStateDB state db key必须前缀
	KeyPrefixStateDB = "mavl-x2ethereum-"
	//KeyPrefixLocalDB local db的key必须前缀
	KeyPrefixLocalDB = "LODB-x2ethereum-"
)

//CalProphecyPrefix ...
func CalProphecyPrefix(id string) []byte {
	return []byte(KeyPrefixStateDB + string(ProphecyKey) + id)
}

//CalEth2ChainPrefix ...
func CalEth2ChainPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(Eth2ChainKey))
}

//CalWithdrawEthPrefix ...
func CalWithdrawEthPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(WithdrawEthKey))
}

//CalChainToEthPrefix ...
func CalChainToEthPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(ChainToEthKey))
}

//CalWithdrawChainPrefix ...
func CalWithdrawChainPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(WithdrawChainKey))
}

//CalValidatorMapsPrefix ...
func CalValidatorMapsPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(ValidatorMapsKey))
}

//CalLastTotalPowerPrefix ...
func CalLastTotalPowerPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(LastTotalPowerKey))
}

//CalConsensusThresholdPrefix ...
func CalConsensusThresholdPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(ConsensusThresholdKey))
}

//CalTokenSymbolTotalLockOrBurnAmount ...
func CalTokenSymbolTotalLockOrBurnAmount(symbol, tokenAddress, direction, txType string) []byte {
	return []byte(KeyPrefixLocalDB + string(TokenSymbolTotalLockOrBurnAmountKey) + direction + "-" + txType + "-" + strings.ToLower(symbol) + "-" + strings.ToLower(tokenAddress))
}

//CalTokenSymbolToTokenAddress ...
func CalTokenSymbolToTokenAddress(symbol string) []byte {
	return []byte(KeyPrefixLocalDB + string(TokenSymbolToTokenAddressKey) + strings.ToLower(symbol))
}
