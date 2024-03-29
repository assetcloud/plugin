// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qbftNode

import (
	"github.com/assetcloud/chain/pluginmgr"
	"github.com/assetcloud/plugin/plugin/dapp/qbftNode/commands"
	"github.com/assetcloud/plugin/plugin/dapp/qbftNode/executor"
	"github.com/assetcloud/plugin/plugin/dapp/qbftNode/rpc"
	"github.com/assetcloud/plugin/plugin/dapp/qbftNode/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.QbftNodeX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.ValCmd,
		RPC:      rpc.Init,
	})
}
