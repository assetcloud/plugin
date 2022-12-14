// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	log "github.com/assetcloud/chain/common/log/log15"
	drivers "github.com/assetcloud/chain/system/dapp"
	"github.com/assetcloud/chain/types"
)

var clog = log.New("module", "execs.norm")
var driverName = "norm"

// Init norm
func Init(name string, cfg *types.ChainConfig, sub []byte) {
	clog.Debug("register norm execer")
	drivers.Register(cfg, GetName(), newNorm, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
}

//InitExecType ...
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Norm{}))
}

// GetName for norm
func GetName() string {
	return newNorm().GetName()
}

// Norm driver
type Norm struct {
	drivers.DriverBase
}

func newNorm() drivers.Driver {
	n := &Norm{}
	n.SetChild(n)
	n.SetIsFree(true)
	n.SetExecutorType(types.LoadExecutorType(driverName))
	return n
}

// GetDriverName for norm
func (n *Norm) GetDriverName() string {
	return driverName
}

// CheckTx for norm
func (n *Norm) CheckTx(tx *types.Transaction, index int) error {
	return nil
}

// Key for norm
func Key(str []byte) (key []byte) {
	key = append(key, []byte("mavl-norm-")...)
	key = append(key, str...)
	return key
}

// CheckReceiptExecOk return true to check if receipt ty is ok
func (n *Norm) CheckReceiptExecOk() bool {
	return true
}
