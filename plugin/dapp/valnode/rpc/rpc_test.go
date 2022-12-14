/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package rpc

//only load all plugin and system
import (
	"testing"

	"strings"

	"github.com/assetcloud/chain/client"
	"github.com/assetcloud/chain/client/mocks"
	rpctypes "github.com/assetcloud/chain/rpc/types"
	"github.com/assetcloud/chain/types"
	vt "github.com/assetcloud/plugin/plugin/dapp/valnode/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func newGrpc(api client.QueueProtocolAPI) *channelClient {
	return &channelClient{
		ChannelClient: rpctypes.ChannelClient{QueueProtocolAPI: api},
	}
}

func newJrpc(api client.QueueProtocolAPI) *Jrpc {
	return &Jrpc{cli: newGrpc(api)}
}

func TestChannelClient_IsSync(t *testing.T) {
	cfg := types.NewChainConfig(strings.Replace(types.GetDefaultCfgstring(), "Title=\"local\"", "Title=\"chain\"", 1))
	api := new(mocks.QueueProtocolAPI)
	api.On("GetConfig", mock.Anything).Return(cfg, nil)
	client := newGrpc(api)
	client.Init("valnode", nil, nil, nil)
	req := &types.ReqNil{}
	api.On("QueryConsensusFunc", "tendermint", "IsHealthy", req).Return(&vt.IsHealthy{IsHealthy: true}, nil)
	result, err := client.IsSync(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, true, result.IsHealthy)
}

func TestJrpc_IsSync(t *testing.T) {
	api := new(mocks.QueueProtocolAPI)
	J := newJrpc(api)
	req := &types.ReqNil{}
	var result interface{}
	api.On("QueryConsensusFunc", "tendermint", "IsHealthy", req).Return(&vt.IsHealthy{IsHealthy: true}, nil)
	err := J.IsSync(req, &result)
	assert.Nil(t, err)
	assert.Equal(t, true, result)
}

func TestChannelClient_GetNodeInfo(t *testing.T) {
	cfg := types.NewChainConfig(strings.Replace(types.GetDefaultCfgstring(), "Title=\"local\"", "Title=\"chain\"", 1))
	api := new(mocks.QueueProtocolAPI)
	api.On("GetConfig", mock.Anything).Return(cfg, nil)
	client := newGrpc(api)
	client.Init("valnode", nil, nil, nil)
	req := &types.ReqNil{}
	node := &vt.ValNodeInfo{
		NodeIP:      "127.0.0.1",
		NodeID:      "001",
		Address:     "aaa",
		PubKey:      "bbb",
		VotingPower: 10,
		Accum:       -1,
	}
	set := &vt.ValNodeInfoSet{
		Nodes: []*vt.ValNodeInfo{node},
	}
	api.On("QueryConsensusFunc", "tendermint", "NodeInfo", req).Return(set, nil)
	result, err := client.GetNodeInfo(context.Background(), req)
	assert.Nil(t, err)
	assert.EqualValues(t, set, result)
}

func TestJrpc_GetNodeInfo(t *testing.T) {
	api := new(mocks.QueueProtocolAPI)
	J := newJrpc(api)
	req := &types.ReqNil{}
	var result interface{}
	node := &vt.ValNodeInfo{
		NodeIP:      "127.0.0.1",
		NodeID:      "001",
		Address:     "aaa",
		PubKey:      "bbb",
		VotingPower: 10,
		Accum:       -1,
	}
	set := &vt.ValNodeInfoSet{
		Nodes: []*vt.ValNodeInfo{node},
	}
	api.On("QueryConsensusFunc", "tendermint", "NodeInfo", req).Return(set, nil)
	err := J.GetNodeInfo(req, &result)
	assert.Nil(t, err)
	assert.EqualValues(t, set, result)
}
