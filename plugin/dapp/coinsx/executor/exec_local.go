// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/assetcloud/chain/types"
)

// ExecLocal_Transfer  transfer of local exec
func (c *Coinsx) ExecLocal_Transfer(transfer *types.AssetsTransfer, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	kv, err := updateAddrReciver(c.GetLocalDB(), tx.GetRealToAddr(), transfer.Amount, true)
	if err != nil {
		return nil, err
	}
	return &types.LocalDBSet{KV: []*types.KeyValue{kv}}, nil
}

// ExecLocal_TransferToExec  transfer of local exec to exec
func (c *Coinsx) ExecLocal_TransferToExec(transfer *types.AssetsTransferToExec, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	kv, err := updateAddrReciver(c.GetLocalDB(), tx.GetRealToAddr(), transfer.Amount, true)
	if err != nil {
		return nil, err
	}
	return &types.LocalDBSet{KV: []*types.KeyValue{kv}}, nil
}

// ExecLocal_Withdraw  withdraw local exec
func (c *Coinsx) ExecLocal_Withdraw(withdraw *types.AssetsWithdraw, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	from := tx.From()
	kv, err := updateAddrReciver(c.GetLocalDB(), from, withdraw.Amount, true)
	if err != nil {
		return nil, err
	}
	return &types.LocalDBSet{KV: []*types.KeyValue{kv}}, nil
}

// ExecLocal_Genesis Genesis of local exec
func (c *Coinsx) ExecLocal_Genesis(gen *types.AssetsGenesis, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	kv, err := updateAddrReciver(c.GetLocalDB(), tx.GetRealToAddr(), gen.Amount, true)
	if err != nil {
		return nil, err
	}
	return &types.LocalDBSet{KV: []*types.KeyValue{kv}}, nil
}
