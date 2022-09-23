package types

import (
	"github.com/assetcloud/plugin/plugin/dapp/accountmanager/commands"
	"github.com/assetcloud/plugin/plugin/dapp/accountmanager/executor"
	"github.com/assetcloud/plugin/plugin/dapp/accountmanager/rpc"
	accountmanagertypes "github.com/assetcloud/plugin/plugin/dapp/accountmanager/types"
	"github.com/assetcloud/chain/pluginmgr"
)

/*
 * 初始化dapp相关的组件
 */

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     accountmanagertypes.AccountmanagerX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
