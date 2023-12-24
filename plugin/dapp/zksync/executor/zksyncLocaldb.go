package executor

import (
	"github.com/assetcloud/chain/types"
	zt "github.com/assetcloud/plugin/plugin/dapp/zksync/types"
)

func (z *zksync) execAutoLocalZksync(tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.Ty != types.ExecOk {
		return nil, types.ErrInvalidParam
	}
	set, err := z.execLocalZksync(tx, receiptData, index)
	if err != nil {
		return set, err
	}
	dbSet := &types.LocalDBSet{}
	dbSet.KV = z.AddRollbackKV(tx, tx.Execer, set.KV)
	return dbSet, nil
}

func (z *zksync) execLocalZksync(tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	infoTable := NewAccountTreeTable(z.GetLocalDB())

	dbSet := &types.LocalDBSet{}
	for _, log := range receiptData.Logs {
		switch log.Ty {
		case zt.TyDepositLog:
			var receipt zt.AccountTokenBalanceReceipt
			err := types.Decode(log.GetLog(), &receipt)
			if err != nil {
				return nil, err
			}
			leaf := &zt.Leaf{
				AccountId:  receipt.AccountId,
				EthAddress: receipt.EthAddress,
				ChainAddr:  receipt.ChainAddr,
			}

			err = infoTable.Replace(leaf)
			if err != nil {
				return nil, err
			}
		}
	}
	kvs, err := infoTable.Save()
	if err != nil {
		return nil, err
	}
	dbSet.KV = append(dbSet.KV, kvs...)
	return dbSet, nil
}

func (z *zksync) execAutoDelLocal(tx *types.Transaction, receiptData *types.ReceiptData) (*types.LocalDBSet, error) {
	kvs, err := z.DelRollbackKV(tx, tx.Execer)
	if err != nil {
		return nil, err
	}
	dbSet := &types.LocalDBSet{}
	dbSet.KV = append(dbSet.KV, kvs...)
	return dbSet, nil
}

func (z *zksync) execCommitProofLocal(payload *zt.ZkCommitProof, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.Ty != types.ExecOk {
		return nil, types.ErrInvalidParam
	}

	proofTable := NewCommitProofTable(z.GetLocalDB())

	set := &types.LocalDBSet{}
	payload.CommitBlockHeight = z.GetHeight()
	err := proofTable.Replace(payload)
	if err != nil {
		return nil, err
	}

	kvs, err := proofTable.Save()
	if err != nil {
		return nil, err
	}
	set.KV = append(set.KV, kvs...)

	dbSet := &types.LocalDBSet{}
	dbSet.KV = z.AddRollbackKV(tx, tx.Execer, set.KV)
	return dbSet, nil
}
