package commands

import (
	"testing"

	"github.com/assetcloud/chain/rpc/jsonclient"
	rpctypes "github.com/assetcloud/chain/rpc/types"
	"github.com/assetcloud/chain/types"
	"github.com/assetcloud/chain/util/testnode"
	"github.com/stretchr/testify/assert"

	// 因为测试程序在外层，而合约类型的初始化在里面，所以需要显示引用，否则不会加载合约插件
	evm "github.com/assetcloud/plugin/plugin/dapp/evm/executor"
	evmtypes "github.com/assetcloud/plugin/plugin/dapp/evm/types"

	// 需要显示引用系统插件，以加载系统内置合约
	"github.com/assetcloud/chain/client/mocks"
	_ "github.com/assetcloud/chain/system"
	"github.com/stretchr/testify/mock"
)

// TestQueryDebug 测试命令行调用rpc接口
func TestQueryDebug(t *testing.T) {
	var cfg = types.NewChainConfig(types.GetDefaultCfgstring())
	evm.Init(evmtypes.ExecutorName, cfg, nil)
	var debugReq = evmtypes.EvmDebugReq{Optype: 1}
	js, err := types.PBToJSON(&debugReq)
	assert.Nil(t, err)
	in := &rpctypes.Query4Jrpc{
		Execer:   "evm",
		FuncName: "EvmDebug",
		Payload:  js,
	}

	var mockResp = evmtypes.EvmDebugResp{DebugStatus: "on"}

	mockapi := &mocks.QueueProtocolAPI{}
	// 这里对需要mock的方法打桩,Close是必须的，其它方法根据需要
	mockapi.On("Close").Return()
	mockapi.On("Query", "evm", "EvmDebug", &debugReq).Return(&mockResp, nil)
	mockapi.On("GetConfig", mock.Anything).Return(cfg, nil)

	mock33 := testnode.New("", mockapi)
	defer mock33.Close()
	rpcCfg := mock33.GetCfg().RPC
	// 这里必须设置监听端口，默认的是无效值
	rpcCfg.JrpcBindAddr = "127.0.0.1:8899"
	mock33.GetRPC().Listen()

	jsonClient, err := jsonclient.NewJSONClient("http://" + rpcCfg.JrpcBindAddr + "/")
	assert.Nil(t, err)
	assert.NotNil(t, jsonClient)

	var debugResp evmtypes.EvmDebugResp
	err = jsonClient.Call("Chain.Query", in, &debugResp)
	assert.Nil(t, err)
	assert.Equal(t, "on", debugResp.DebugStatus)
}
