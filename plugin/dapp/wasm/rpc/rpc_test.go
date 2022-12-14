package rpc

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/assetcloud/chain/client/mocks"
	rpctypes "github.com/assetcloud/chain/rpc/types"
	"github.com/assetcloud/chain/types"
	types2 "github.com/assetcloud/plugin/plugin/dapp/wasm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var cfg *types.ChainConfig

func init() {
	cfg = types.NewChainConfig(strings.Replace(types.GetDefaultCfgstring(), "Title=\"local\"", "Title=\"chain\"", 1))
}

func TestJrpc_CheckContract(t *testing.T) {
	api := new(mocks.QueueProtocolAPI)
	api.On("Query", types2.WasmX, "Check", mock.Anything).Return(&types.Reply{}, nil)
	jrpc := &Jrpc{
		cli: &channelClient{
			rpctypes.ChannelClient{
				QueueProtocolAPI: api,
			},
		},
	}
	var result interface{}
	err := jrpc.CheckContract(&types2.QueryCheckContract{Name: "dice"}, &result)
	assert.Nil(t, err, "CheckContract error not nil")
	assert.Equal(t, false, result.(bool))
}

func TestJrpc_CreateContract(t *testing.T) {
	jrpc := &Jrpc{}
	code, err := ioutil.ReadFile("../contracts/dice/dice.wasm")
	assert.Nil(t, err, "read wasm file error")
	var result interface{}
	err = jrpc.CreateContract(&types2.WasmCreate{Name: "dice", Code: code}, &result)
	assert.Nil(t, err, "create contract error")
	t.Log(result)
}

func TestJrpc_CallContract(t *testing.T) {
	jrpc := &Jrpc{}
	var result interface{}
	err := jrpc.CallContract(&types2.WasmCall{Contract: "dice", Method: "play"}, &result)
	assert.Nil(t, err, "call contract error")
	t.Log(result)
}
