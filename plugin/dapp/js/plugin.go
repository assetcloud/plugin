package js

import (
	"github.com/assetcloud/chain/pluginmgr"
	"github.com/assetcloud/plugin/plugin/dapp/js/executor"
	ptypes "github.com/assetcloud/plugin/plugin/dapp/js/types"

	// init auto test
	_ "github.com/assetcloud/plugin/plugin/dapp/js/autotest"
	"github.com/assetcloud/plugin/plugin/dapp/js/command"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     ptypes.JsX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      command.JavaScriptCmd,
		RPC:      nil,
	})
}
