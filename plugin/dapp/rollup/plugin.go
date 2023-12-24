package types

import (
	"github.com/assetcloud/chain/pluginmgr"
	"github.com/assetcloud/plugin/plugin/dapp/rollup/commands"
	"github.com/assetcloud/plugin/plugin/dapp/rollup/executor"
	"github.com/assetcloud/plugin/plugin/dapp/rollup/rpc"
	rolluptypes "github.com/assetcloud/plugin/plugin/dapp/rollup/types"
)

/*
 * 初始化dapp相关的组件
 */

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     rolluptypes.RollupX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
