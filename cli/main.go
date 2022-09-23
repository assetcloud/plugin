// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.8
// +build go1.8

package main

import (
	_ "github.com/assetcloud/chain/system"
	"github.com/assetcloud/plugin/cli/buildflags"
	_ "github.com/assetcloud/plugin/plugin"

	"github.com/assetcloud/chain/util/cli"
)

func main() {
	if buildflags.RPCAddr == "" {
		buildflags.RPCAddr = "http://localhost:8801"
	}
	cli.Run(buildflags.RPCAddr, buildflags.ParaName, "")
}
