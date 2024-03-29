// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ticket

import (
	"github.com/assetcloud/chain/pluginmgr"
	"github.com/assetcloud/plugin/plugin/dapp/ticket/commands"
	"github.com/assetcloud/plugin/plugin/dapp/ticket/executor"
	"github.com/assetcloud/plugin/plugin/dapp/ticket/rpc"
	"github.com/assetcloud/plugin/plugin/dapp/ticket/types"

	// init wallet
	_ "github.com/assetcloud/plugin/plugin/dapp/ticket/wallet"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.TicketX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.TicketCmd,
		RPC:      rpc.Init,
	})
}
