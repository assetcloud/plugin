// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collateralize

import (
	"github.com/assetcloud/chain/pluginmgr"
	"github.com/assetcloud/plugin/plugin/dapp/collateralize/commands"
	"github.com/assetcloud/plugin/plugin/dapp/collateralize/executor"
	"github.com/assetcloud/plugin/plugin/dapp/collateralize/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.CollateralizeX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.CollateralizeCmd,
	})
}
