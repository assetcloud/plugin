package wasm

import (
	"github.com/assetcloud/plugin/plugin/dapp/wasm/commands"
	"github.com/assetcloud/plugin/plugin/dapp/wasm/executor"
	"github.com/assetcloud/plugin/plugin/dapp/wasm/rpc"
	"github.com/assetcloud/plugin/plugin/dapp/wasm/types"
	"github.com/assetcloud/chain/pluginmgr"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.WasmX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
