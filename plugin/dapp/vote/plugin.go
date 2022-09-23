package types

import (
	"github.com/assetcloud/plugin/plugin/dapp/vote/commands"
	"github.com/assetcloud/plugin/plugin/dapp/vote/executor"
	"github.com/assetcloud/plugin/plugin/dapp/vote/rpc"
	votetypes "github.com/assetcloud/plugin/plugin/dapp/vote/types"
	"github.com/assetcloud/chain/pluginmgr"
)

/*
 * 初始化dapp相关的组件
 */

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     votetypes.VoteX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
