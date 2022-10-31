// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.9
// +build go1.9

/*
每个系统的功能通过插件完成，插件分成4类：
共识 加密 dapp 存储
这个go 包提供了 官方提供的 插件。
*/
package main

import (
	"flag"
	"fmt"
	"os"

	frameVersion "github.com/assetcloud/chain/common/version"
	_ "github.com/assetcloud/chain/system"
	"github.com/assetcloud/chain/util/cli"
	_ "github.com/assetcloud/plugin/plugin"
	"github.com/assetcloud/plugin/version"
)

var (
	versionCmd = flag.Bool("version", false, "detail version")
)

func main() {
	flag.Parse()
	if *versionCmd {
		fmt.Fprintln(os.Stdout, "Build time:", version.BuildTime)
		fmt.Fprintln(os.Stdout, "System version:", version.Platform)
		fmt.Fprintln(os.Stdout, "Golang version:", version.GoVersion)
		fmt.Fprintln(os.Stdout, "plugin version:", version.GetVersion())
		fmt.Fprintln(os.Stdout, "chain frame version:", frameVersion.GetVersion())
		return
	}
	cli.RunChain("", "")
}
