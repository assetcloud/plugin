package types

import (
	"github.com/assetcloud/chain/pluginmgr"
	"github.com/assetcloud/plugin/plugin/dapp/storage/commands"
	"github.com/assetcloud/plugin/plugin/dapp/storage/executor"
	"github.com/assetcloud/plugin/plugin/dapp/storage/rpc"
	storagetypes "github.com/assetcloud/plugin/plugin/dapp/storage/types"
)

/*
 * 初始化dapp相关的组件
 */

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     storagetypes.StorageX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
