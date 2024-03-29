// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"golang.org/x/net/context"

	"github.com/assetcloud/chain/types"
	bw "github.com/assetcloud/plugin/plugin/dapp/blackwhite/types"
)

func (c *channelClient) Create(ctx context.Context, head *bw.BlackwhiteCreate) (*types.UnsignTx, error) {
	val := &bw.BlackwhiteAction{
		Ty:    bw.BlackwhiteActionCreate,
		Value: &bw.BlackwhiteAction_Create{Create: head},
	}
	tx := &types.Transaction{
		Payload: types.Encode(val),
	}
	cfg := c.GetConfig()
	data, err := types.FormatTxEncode(cfg, cfg.ExecName(string(bw.ExecerBlackwhite)), tx)
	if err != nil {
		return nil, err
	}
	return &types.UnsignTx{Data: data}, nil
}

func (c *channelClient) Show(ctx context.Context, head *bw.BlackwhiteShow) (*types.UnsignTx, error) {
	val := &bw.BlackwhiteAction{
		Ty:    bw.BlackwhiteActionShow,
		Value: &bw.BlackwhiteAction_Show{Show: head},
	}
	tx := &types.Transaction{
		Payload: types.Encode(val),
	}
	cfg := c.GetConfig()
	data, err := types.FormatTxEncode(cfg, cfg.ExecName(string(bw.ExecerBlackwhite)), tx)
	if err != nil {
		return nil, err
	}
	return &types.UnsignTx{Data: data}, nil
}

func (c *channelClient) Play(ctx context.Context, head *bw.BlackwhitePlay) (*types.UnsignTx, error) {
	val := &bw.BlackwhiteAction{
		Ty:    bw.BlackwhiteActionPlay,
		Value: &bw.BlackwhiteAction_Play{Play: head},
	}
	tx := &types.Transaction{
		Payload: types.Encode(val),
	}
	cfg := c.GetConfig()
	data, err := types.FormatTxEncode(cfg, cfg.ExecName(string(bw.ExecerBlackwhite)), tx)
	if err != nil {
		return nil, err
	}
	return &types.UnsignTx{Data: data}, nil
}

func (c *channelClient) TimeoutDone(ctx context.Context, head *bw.BlackwhiteTimeoutDone) (*types.UnsignTx, error) {
	val := &bw.BlackwhiteAction{
		Ty:    bw.BlackwhiteActionTimeoutDone,
		Value: &bw.BlackwhiteAction_TimeoutDone{TimeoutDone: head},
	}
	tx := &types.Transaction{
		Payload: types.Encode(val),
	}
	cfg := c.GetConfig()
	data, err := types.FormatTxEncode(cfg, cfg.ExecName(string(bw.ExecerBlackwhite)), tx)
	if err != nil {
		return nil, err
	}
	return &types.UnsignTx{Data: data}, nil
}
