// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"fmt"

	log "github.com/assetcloud/chain/common/log/log15"
	drivers "github.com/assetcloud/chain/system/dapp"
	"github.com/assetcloud/chain/types"
)

var clog = log.New("module", "execs.valnode")
var driverName = "valnode"

// Init method
func Init(name string, cfg *types.ChainConfig, sub []byte) {
	clog.Debug("register valnode execer")
	drivers.Register(cfg, GetName(), newValNode, 0)
	InitExecType()
}

//InitExecType ...
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&ValNode{}))
}

// GetName method
func GetName() string {
	return newValNode().GetName()
}

// ValNode strucyt
type ValNode struct {
	drivers.DriverBase
}

func newValNode() drivers.Driver {
	n := &ValNode{}
	n.SetChild(n)
	n.SetIsFree(true)
	n.SetExecutorType(types.LoadExecutorType(driverName))
	return n
}

// GetDriverName method
func (val *ValNode) GetDriverName() string {
	return driverName
}

// CheckTx method
func (val *ValNode) CheckTx(tx *types.Transaction, index int) error {
	return nil
}

// CalcValNodeUpdateHeightIndexKey method
func CalcValNodeUpdateHeightIndexKey(height int64, index int) []byte {
	return []byte(fmt.Sprintf("LODB-valnode-Update:%18d:%18d", height, int64(index)))
}

// CalcValNodeUpdateHeightKey method
func CalcValNodeUpdateHeightKey(height int64) []byte {
	return []byte(fmt.Sprintf("LODB-valnode-Update:%18d:", height))
}

// CalcValNodeBlockInfoHeightKey method
func CalcValNodeBlockInfoHeightKey(height int64) []byte {
	return []byte(fmt.Sprintf("LODB-valnode-BlockInfo:%18d:", height))
}

// CheckReceiptExecOk return true to check if receipt ty is ok
func (val *ValNode) CheckReceiptExecOk() bool {
	return true
}
