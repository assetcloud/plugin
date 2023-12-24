package wasm

import (
	"github.com/assetcloud/chain/pluginmgr"
	"github.com/assetcloud/plugin/plugin/dapp/zksync/commands"
	"github.com/assetcloud/plugin/plugin/dapp/zksync/executor"
	"github.com/assetcloud/plugin/plugin/dapp/zksync/rpc"
	"github.com/assetcloud/plugin/plugin/dapp/zksync/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.Zksync,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.ZksyncCmd,
		RPC:      rpc.Init,
	})
}
