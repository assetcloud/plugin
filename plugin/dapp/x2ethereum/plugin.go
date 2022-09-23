package x2ethereum

import (
	"github.com/assetcloud/plugin/plugin/dapp/x2ethereum/commands"
	"github.com/assetcloud/plugin/plugin/dapp/x2ethereum/executor"
	"github.com/assetcloud/plugin/plugin/dapp/x2ethereum/rpc"
	x2ethereumtypes "github.com/assetcloud/plugin/plugin/dapp/x2ethereum/types"
	"github.com/assetcloud/chain/pluginmgr"
)

/*
 * 初始化dapp相关的组件
 */

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     x2ethereumtypes.X2ethereumX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
