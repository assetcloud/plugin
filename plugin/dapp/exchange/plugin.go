package types

import (
	"github.com/assetcloud/plugin/plugin/dapp/exchange/commands"
	"github.com/assetcloud/plugin/plugin/dapp/exchange/executor"
	"github.com/assetcloud/plugin/plugin/dapp/exchange/rpc"
	exchangetypes "github.com/assetcloud/plugin/plugin/dapp/exchange/types"
	"github.com/assetcloud/chain/pluginmgr"
)

/*
 * 初始化dapp相关的组件
 */

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     exchangetypes.ExchangeX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
