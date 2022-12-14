// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pokerbull

import (
	"github.com/assetcloud/chain/pluginmgr"
	"github.com/assetcloud/plugin/plugin/dapp/pokerbull/cmd"
	"github.com/assetcloud/plugin/plugin/dapp/pokerbull/executor"
	"github.com/assetcloud/plugin/plugin/dapp/pokerbull/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.PokerBullX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      cmd.PokerBullCmd,
	})
}
