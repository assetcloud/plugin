package executor

import (
	"strconv"

	"github.com/assetcloud/chain/types"
	x2eTy "github.com/assetcloud/plugin/plugin/dapp/x2ethereum/types"
)

/*
 * 实现交易相关数据本地执行，数据不上链
 * 非关键数据，本地存储(localDB), 用于辅助查询，效率高
 */

func (x *x2ethereum) ExecLocal_Eth2ChainLock(payload *x2eTy.Eth2Chain, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := x.execLocal(receiptData)
	if err != nil {
		return set, err
	}
	return x.addAutoRollBack(tx, set.KV), nil
}

func (x *x2ethereum) ExecLocal_Eth2ChainBurn(payload *x2eTy.Eth2Chain, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := x.execLocal(receiptData)
	if err != nil {
		return set, err
	}
	return x.addAutoRollBack(tx, set.KV), nil
}

func (x *x2ethereum) ExecLocal_ChainToEthBurn(payload *x2eTy.ChainToEth, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := x.execLocal(receiptData)
	if err != nil {
		return set, err
	}
	return x.addAutoRollBack(tx, set.KV), nil
}

func (x *x2ethereum) ExecLocal_ChainToEthLock(payload *x2eTy.ChainToEth, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := x.execLocal(receiptData)
	if err != nil {
		return set, err
	}
	return x.addAutoRollBack(tx, set.KV), nil
}

func (x *x2ethereum) ExecLocal_AddValidator(payload *x2eTy.MsgValidator, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return x.addAutoRollBack(tx, dbSet.KV), nil
}

func (x *x2ethereum) ExecLocal_RemoveValidator(payload *x2eTy.MsgValidator, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return x.addAutoRollBack(tx, dbSet.KV), nil
}

func (x *x2ethereum) ExecLocal_ModifyPower(payload *x2eTy.MsgValidator, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return x.addAutoRollBack(tx, dbSet.KV), nil
}

func (x *x2ethereum) ExecLocal_SetConsensusThreshold(payload *x2eTy.MsgConsensusThreshold, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return x.addAutoRollBack(tx, dbSet.KV), nil
}

//设置自动回滚
func (x *x2ethereum) addAutoRollBack(tx *types.Transaction, kv []*types.KeyValue) *types.LocalDBSet {
	dbSet := &types.LocalDBSet{}
	dbSet.KV = x.AddRollbackKV(tx, tx.Execer, kv)
	return dbSet
}

func (x *x2ethereum) execLocal(receiptData *types.ReceiptData) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	for _, log := range receiptData.Logs {
		switch log.Ty {
		case x2eTy.TyEth2ChainLog:
			var receiptEth2Chain x2eTy.ReceiptEth2Chain
			err := types.Decode(log.Log, &receiptEth2Chain)
			if err != nil {
				return nil, err
			}

			nb, err := x.GetLocalDB().Get(x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptEth2Chain.IssuerDotSymbol, receiptEth2Chain.TokenAddress, x2eTy.DirEth2Chain, "lock"))
			if err != nil && err != types.ErrNotFound {
				return nil, err
			}
			var now x2eTy.ReceiptQuerySymbolAssetsByTxType
			err = types.Decode(nb, &now)
			if err != nil {
				return nil, err
			}
			preAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(now.TotalAmount), 64)
			nowAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(receiptEth2Chain.Amount), 64)
			TokenAssetsByTxTypeBytes := types.Encode(&x2eTy.ReceiptQuerySymbolAssetsByTxType{
				TokenSymbol: receiptEth2Chain.IssuerDotSymbol,
				TxType:      "lock",
				TotalAmount: strconv.FormatFloat(preAmount+nowAmount, 'f', 4, 64),
				Direction:   1,
			})
			dbSet.KV = append(dbSet.KV, &types.KeyValue{
				Key:   x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptEth2Chain.IssuerDotSymbol, receiptEth2Chain.TokenAddress, x2eTy.DirEth2Chain, "lock"),
				Value: TokenAssetsByTxTypeBytes,
			})

			nb, err = x.GetLocalDB().Get(x2eTy.CalTokenSymbolToTokenAddress(receiptEth2Chain.IssuerDotSymbol))
			if err != nil && err != types.ErrNotFound {
				return nil, err
			}
			var t x2eTy.ReceiptTokenToTokenAddress
			err = types.Decode(nb, &t)
			if err != nil {
				return nil, err
			}
			var exist bool
			for _, addr := range t.TokenAddress {
				if addr == receiptEth2Chain.TokenAddress {
					exist = true
				}
			}
			if !exist {
				t.TokenAddress = append(t.TokenAddress, receiptEth2Chain.TokenAddress)
			}
			TokenToTokenAddressBytes := types.Encode(&x2eTy.ReceiptTokenToTokenAddress{
				TokenAddress: t.TokenAddress,
			})
			dbSet.KV = append(dbSet.KV, &types.KeyValue{
				Key:   x2eTy.CalTokenSymbolToTokenAddress(receiptEth2Chain.IssuerDotSymbol),
				Value: TokenToTokenAddressBytes,
			})
		case x2eTy.TyWithdrawEthLog:
			var receiptEth2Chain x2eTy.ReceiptEth2Chain
			err := types.Decode(log.Log, &receiptEth2Chain)
			if err != nil {
				return nil, err
			}

			nb, err := x.GetLocalDB().Get(x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptEth2Chain.IssuerDotSymbol, receiptEth2Chain.TokenAddress, x2eTy.DirEth2Chain, "withdraw"))
			if err != nil && err != types.ErrNotFound {
				return nil, err
			}
			var now x2eTy.ReceiptQuerySymbolAssetsByTxType
			err = types.Decode(nb, &now)
			if err != nil {
				return nil, err
			}

			preAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(now.TotalAmount), 64)
			nowAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(receiptEth2Chain.Amount), 64)
			TokenAssetsByTxTypeBytes := types.Encode(&x2eTy.ReceiptQuerySymbolAssetsByTxType{
				TokenSymbol: receiptEth2Chain.IssuerDotSymbol,
				TxType:      "withdraw",
				TotalAmount: strconv.FormatFloat(preAmount+nowAmount, 'f', 4, 64),
				Direction:   2,
			})
			dbSet.KV = append(dbSet.KV, &types.KeyValue{
				Key:   x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptEth2Chain.IssuerDotSymbol, receiptEth2Chain.TokenAddress, x2eTy.DirEth2Chain, "withdraw"),
				Value: TokenAssetsByTxTypeBytes,
			})
		case x2eTy.TyChainToEthLog:
			var receiptChainToEth x2eTy.ReceiptChainToEth
			err := types.Decode(log.Log, &receiptChainToEth)
			if err != nil {
				return nil, err
			}

			nb, err := x.GetLocalDB().Get(x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptChainToEth.IssuerDotSymbol, receiptChainToEth.TokenContract, x2eTy.DirChainToEth, "lock"))
			if err != nil && err != types.ErrNotFound {
				return nil, err
			}
			var now x2eTy.ReceiptQuerySymbolAssetsByTxType
			err = types.Decode(nb, &now)
			if err != nil {
				return nil, err
			}

			preAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(now.TotalAmount), 64)
			nowAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(receiptChainToEth.Amount), 64)
			TokenAssetsByTxTypeBytes := types.Encode(&x2eTy.ReceiptQuerySymbolAssetsByTxType{
				TokenSymbol: receiptChainToEth.IssuerDotSymbol,
				TxType:      "lock",
				TotalAmount: strconv.FormatFloat(preAmount+nowAmount, 'f', 4, 64),
				Direction:   1,
			})
			dbSet.KV = append(dbSet.KV, &types.KeyValue{
				Key:   x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptChainToEth.IssuerDotSymbol, receiptChainToEth.TokenContract, x2eTy.DirChainToEth, "lock"),
				Value: TokenAssetsByTxTypeBytes,
			})
		case x2eTy.TyWithdrawChainLog:
			var receiptChainToEth x2eTy.ReceiptChainToEth
			err := types.Decode(log.Log, &receiptChainToEth)
			if err != nil {
				return nil, err
			}

			nb, err := x.GetLocalDB().Get(x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptChainToEth.IssuerDotSymbol, receiptChainToEth.TokenContract, x2eTy.DirChainToEth, ""))
			if err != nil && err != types.ErrNotFound {
				return nil, err
			}
			var now x2eTy.ReceiptQuerySymbolAssetsByTxType
			err = types.Decode(nb, &now)
			if err != nil {
				return nil, err
			}

			preAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(now.TotalAmount), 64)
			nowAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(receiptChainToEth.Amount), 64)
			TokenAssetsByTxTypeBytes := types.Encode(&x2eTy.ReceiptQuerySymbolAssetsByTxType{
				TokenSymbol: receiptChainToEth.IssuerDotSymbol,
				TxType:      "withdraw",
				TotalAmount: strconv.FormatFloat(preAmount+nowAmount, 'f', 4, 64),
				Direction:   2,
			})
			dbSet.KV = append(dbSet.KV, &types.KeyValue{
				Key:   x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptChainToEth.IssuerDotSymbol, receiptChainToEth.TokenContract, x2eTy.DirChainToEth, "withdraw"),
				Value: TokenAssetsByTxTypeBytes,
			})
		default:
			continue
		}
	}
	return dbSet, nil
}
