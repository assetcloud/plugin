// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	log "github.com/assetcloud/chain/common/log/log15"
	drivers "github.com/assetcloud/chain/system/dapp"
	"github.com/assetcloud/chain/types"
	uf "github.com/assetcloud/plugin/plugin/dapp/unfreeze/types"
)

var uflog = log.New("module", "execs.unfreeze")

var driverName = uf.UnfreezeX

// Init 重命名执行器名称
func Init(name string, cfg *types.ChainConfig, sub []byte) {
	drivers.Register(cfg, GetName(), newUnfreeze, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
}

//InitExecType ...
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Unfreeze{}))
}

// Unfreeze 执行器结构体
type Unfreeze struct {
	drivers.DriverBase
}

func newUnfreeze() drivers.Driver {
	t := &Unfreeze{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

// GetName 获得执行器名字
func GetName() string {
	return newUnfreeze().GetName()
}

// GetDriverName 获得驱动名字
func (u *Unfreeze) GetDriverName() string {
	return driverName
}
