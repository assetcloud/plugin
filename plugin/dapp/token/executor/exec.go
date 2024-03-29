// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/assetcloud/chain/account"
	"github.com/assetcloud/chain/types"
	tokenty "github.com/assetcloud/plugin/plugin/dapp/token/types"
)

func (t *token) Exec_Transfer(payload *types.AssetsTransfer, tx *types.Transaction, index int) (*types.Receipt, error) {
	token := payload.GetCointoken()
	cfg := t.GetAPI().GetConfig()
	db, err := account.NewAccountDB(cfg, t.GetName(), token, t.GetStateDB())
	if err != nil {
		return nil, err
	}
	tokenAction := tokenty.TokenAction{
		Ty: tokenty.ActionTransfer,
		Value: &tokenty.TokenAction_Transfer{
			Transfer: payload,
		},
	}
	return t.ExecTransWithdraw(db, tx, &tokenAction, index)
}

func (t *token) Exec_Withdraw(payload *types.AssetsWithdraw, tx *types.Transaction, index int) (*types.Receipt, error) {
	token := payload.GetCointoken()
	cfg := t.GetAPI().GetConfig()
	db, err := account.NewAccountDB(cfg, t.GetName(), token, t.GetStateDB())
	if err != nil {
		return nil, err
	}
	tokenAction := tokenty.TokenAction{
		Ty: tokenty.ActionWithdraw,
		Value: &tokenty.TokenAction_Withdraw{
			Withdraw: payload,
		},
	}
	return t.ExecTransWithdraw(db, tx, &tokenAction, index)
}

func (t *token) Exec_TokenPreCreate(payload *tokenty.TokenPreCreate, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newTokenAction(t, "", tx)
	return action.preCreate(payload)
}

func (t *token) Exec_TokenFinishCreate(payload *tokenty.TokenFinishCreate, tx *types.Transaction, index int) (*types.Receipt, error) {
	cfg := t.GetAPI().GetConfig()
	action := newTokenAction(t, cfg.MGStr("mver.consensus.fundKeyAddr", t.GetHeight()), tx)
	return action.finishCreate(payload)
}

func (t *token) Exec_TokenRevokeCreate(payload *tokenty.TokenRevokeCreate, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newTokenAction(t, "", tx)
	return action.revokeCreate(payload)
}

func (t *token) Exec_TransferToExec(payload *types.AssetsTransferToExec, tx *types.Transaction, index int) (*types.Receipt, error) {
	token := payload.GetCointoken()
	cfg := t.GetAPI().GetConfig()
	db, err := account.NewAccountDB(cfg, t.GetName(), token, t.GetStateDB())
	if err != nil {
		return nil, err
	}
	tokenAction := tokenty.TokenAction{
		Ty: tokenty.TokenActionTransferToExec,
		Value: &tokenty.TokenAction_TransferToExec{
			TransferToExec: payload,
		},
	}
	return t.ExecTransWithdraw(db, tx, &tokenAction, index)
}

func (t *token) Exec_TokenMint(payload *tokenty.TokenMint, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newTokenAction(t, "", tx)
	return action.mint(payload)
}

func (t *token) Exec_TokenBurn(payload *tokenty.TokenBurn, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newTokenAction(t, "", tx)
	return action.burn(payload)
}
