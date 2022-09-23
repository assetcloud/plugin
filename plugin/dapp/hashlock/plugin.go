// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashlock

import (
	"github.com/assetcloud/plugin/plugin/dapp/hashlock/commands"
	"github.com/assetcloud/plugin/plugin/dapp/hashlock/executor"
	"github.com/assetcloud/plugin/plugin/dapp/hashlock/types"
	"github.com/assetcloud/chain/pluginmgr"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.HashlockX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.HashlockCmd,
		RPC:      nil,
	})
}
