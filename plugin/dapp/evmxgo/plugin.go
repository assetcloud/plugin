package types

import (
	"github.com/assetcloud/plugin/plugin/dapp/evmxgo/commands"
	"github.com/assetcloud/plugin/plugin/dapp/evmxgo/executor"
	"github.com/assetcloud/plugin/plugin/dapp/evmxgo/rpc"
	evmxgotypes "github.com/assetcloud/plugin/plugin/dapp/evmxgo/types"
	"github.com/assetcloud/chain/pluginmgr"
)

/*
 * 初始化dapp相关的组件
 */

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     evmxgotypes.EvmxgoX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
