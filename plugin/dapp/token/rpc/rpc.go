// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"encoding/hex"

	"github.com/assetcloud/chain/account"
	"github.com/assetcloud/chain/common/address"
	rpctypes "github.com/assetcloud/chain/rpc/types"
	"github.com/assetcloud/chain/types"
	tokenty "github.com/assetcloud/plugin/plugin/dapp/token/types"
	context "golang.org/x/net/context"
)

//TODO:和GetBalance进行泛化处理，同时LoadAccounts和LoadExecAccountQueue也需要进行泛化处理, added by hzj
func (c *channelClient) getTokenBalance(in *tokenty.ReqTokenBalance) ([]*types.Account, error) {
	cfg := c.GetConfig()
	accountTokendb, err := account.NewAccountDB(cfg, tokenty.TokenX, in.GetTokenSymbol(), nil)
	if err != nil {
		return nil, err
	}

	switch in.GetExecer() {
	case cfg.ExecName(tokenty.TokenX):
		addrs := in.GetAddresses()
		var queryAddrs []string
		queryAddrs = append(queryAddrs, addrs...)

		accounts, err := accountTokendb.LoadAccounts(c.QueueProtocolAPI, queryAddrs)
		if err != nil {
			log.Error("GetTokenBalance", "err", err.Error(), "token symbol", in.GetTokenSymbol(), "address", queryAddrs)
			return nil, err
		}
		return accounts, nil

	default: //trade
		execaddress := address.ExecAddress(in.GetExecer())
		addrs := in.GetAddresses()
		var accounts []*types.Account
		for _, addr := range addrs {
			acc, err := accountTokendb.LoadExecAccountQueue(c.QueueProtocolAPI, addr, execaddress)
			if err != nil {
				log.Error("GetTokenBalance for exector", "err", err.Error(), "token symbol", in.GetTokenSymbol(),
					"address", addr)
				continue
			}
			accounts = append(accounts, acc)
		}

		return accounts, nil
	}
}

// GetTokenBalance 获取token金额（channelClient）
func (c *channelClient) GetTokenBalance(ctx context.Context, in *tokenty.ReqTokenBalance) (*types.Accounts, error) {
	reply, err := c.getTokenBalance(in)
	if err != nil {
		return nil, err
	}
	return &types.Accounts{Acc: reply}, nil
}

// GetTokenBalance 获取token金额 (Jrpc)
func (c *Jrpc) GetTokenBalance(in tokenty.ReqTokenBalance, result *interface{}) error {
	balances, err := c.cli.GetTokenBalance(context.Background(), &in)
	if err != nil {
		return err
	}
	var accounts []*rpctypes.Account
	for _, balance := range balances.Acc {
		accounts = append(accounts, &rpctypes.Account{Addr: balance.GetAddr(),
			Balance:  balance.GetBalance(),
			Currency: balance.GetCurrency(),
			Frozen:   balance.GetFrozen()})
	}
	*result = accounts
	return nil
}

// CreateRawTokenPreCreateTx 创建未签名的创建Token交易
func (c *Jrpc) CreateRawTokenPreCreateTx(param *tokenty.TokenPreCreate, result *interface{}) error {
	if param == nil || param.Symbol == "" {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(tokenty.TokenX), "TokenPreCreate", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// CreateRawTokenFinishTx 创建未签名的结束Token交易
func (c *Jrpc) CreateRawTokenFinishTx(param *tokenty.TokenFinishCreate, result *interface{}) error {
	if param == nil || param.Symbol == "" {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(tokenty.TokenX), "TokenFinishCreate", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// CreateRawTokenRevokeTx 创建未签名的撤销Token交易
func (c *Jrpc) CreateRawTokenRevokeTx(param *tokenty.TokenRevokeCreate, result *interface{}) error {
	if param == nil || param.Symbol == "" {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(tokenty.TokenX), "TokenRevokeCreate", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// CreateRawTokenMintTx 创建未签名的mint Token交易
func (c *Jrpc) CreateRawTokenMintTx(param *tokenty.TokenMint, result *interface{}) error {
	if param == nil || param.Symbol == "" || param.Amount <= 0 {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(tokenty.TokenX), "TokenMint", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// CreateRawTokenBurnTx 创建未签名的 burn Token交易
func (c *Jrpc) CreateRawTokenBurnTx(param *tokenty.TokenBurn, result *interface{}) error {
	if param == nil || param.Symbol == "" || param.Amount <= 0 {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(tokenty.TokenX), "TokenBurn", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}
