// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package commands

import (
	"strings"

	pt "github.com/assetcloud/plugin/plugin/dapp/paracross/types"

	"github.com/assetcloud/chain/common/address"
	"github.com/assetcloud/chain/types"
)

// GetExecAddr 获取执行器地址
func GetExecAddr(exec string) (string, error) {
	if ok := types.IsAllowExecName([]byte(exec), []byte(exec)); !ok {
		return "", types.ErrExecNameNotAllow
	}

	addrResult := address.ExecAddress(exec)
	result := addrResult
	return result, nil
}

func getRealExecName(paraName string, name string) string {
	if strings.HasPrefix(name, pt.ParaPrefix) {
		return name
	}
	return paraName + name
}
