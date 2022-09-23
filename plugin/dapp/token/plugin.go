// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package token 创建token
package token

import (
	_ "github.com/assetcloud/plugin/plugin/dapp/token/autotest" // register token autotest package
	"github.com/assetcloud/plugin/plugin/dapp/token/commands"
	"github.com/assetcloud/plugin/plugin/dapp/token/executor"
	"github.com/assetcloud/plugin/plugin/dapp/token/rpc"
	"github.com/assetcloud/plugin/plugin/dapp/token/types"
	"github.com/assetcloud/chain/pluginmgr"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.TokenX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.TokenCmd,
		RPC:      rpc.Init,
	})
}
