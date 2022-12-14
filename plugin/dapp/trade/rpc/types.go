// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	rpctypes "github.com/assetcloud/chain/rpc/types"
	"github.com/assetcloud/chain/types"
	ptypes "github.com/assetcloud/plugin/plugin/dapp/trade/types"
)

type channelClient struct {
	rpctypes.ChannelClient
}

//Jrpc : Jrpc struct definition
type Jrpc struct {
	cli *channelClient
}

//Grpc : Grpc struct definition
type Grpc struct {
	*channelClient
}

//Init : do the init operation
func Init(name string, s rpctypes.RPCServer) {
	cli := &channelClient{}
	grpc := &Grpc{channelClient: cli}
	cli.Init(name, s, &Jrpc{cli: cli}, grpc)
	ptypes.RegisterTradeServer(s.GRPC(), grpc)
}

//GetLastMemPool : get the last memory pool
func (jrpc *Jrpc) GetLastMemPool(in types.ReqNil, result *interface{}) error {
	reply, err := jrpc.cli.GetLastMempool()
	if err != nil {
		return err
	}

	{
		var txlist rpctypes.ReplyTxList
		txs := reply.GetTxs()
		for _, tx := range txs {
			tran, err := rpctypes.DecodeTx(tx, jrpc.cli.GetConfig().GetCoinPrecision())
			if err != nil {
				continue
			}
			txlist.Txs = append(txlist.Txs, tran)
		}
		*result = &txlist
	}
	return nil
}
