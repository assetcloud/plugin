/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package types

import (
	"reflect"

	"github.com/assetcloud/chain/types"
)

func init() {
	// init executor type
	types.AllowUserExec = append(types.AllowUserExec, []byte(OracleX))
	types.RegFork(OracleX, InitFork)
	types.RegExec(OracleX, InitExecutor)
}

//InitFork ...
func InitFork(cfg *types.ChainConfig) {
	cfg.RegisterDappFork(OracleX, "Enable", 0)
}

//InitExecutor ...
func InitExecutor(cfg *types.ChainConfig) {
	types.RegistorExecutor(OracleX, NewType(cfg))
}

// OracleType 预言机执行器类型
type OracleType struct {
	types.ExecTypeBase
}

// NewType 创建执行器类型
func NewType(cfg *types.ChainConfig) *OracleType {
	c := &OracleType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// GetName 获取执行器名称
func (o *OracleType) GetName() string {
	return OracleX
}

// GetPayload 获取oracle action
func (o *OracleType) GetPayload() types.Message {
	return &OracleAction{}
}

// GetTypeMap 获取类型map
func (o *OracleType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"EventPublish":     ActionEventPublish,
		"EventAbort":       ActionEventAbort,
		"ResultPrePublish": ActionResultPrePublish,
		"ResultAbort":      ActionResultAbort,
		"ResultPublish":    ActionResultPublish,
	}
}

// GetLogMap 获取日志map
func (o *OracleType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{
		TyLogEventPublish:     {Ty: reflect.TypeOf(ReceiptOracle{}), Name: "LogEventPublish"},
		TyLogEventAbort:       {Ty: reflect.TypeOf(ReceiptOracle{}), Name: "LogEventAbort"},
		TyLogResultPrePublish: {Ty: reflect.TypeOf(ReceiptOracle{}), Name: "LogResultPrePublish"},
		TyLogResultAbort:      {Ty: reflect.TypeOf(ReceiptOracle{}), Name: "LogResultAbort"},
		TyLogResultPublish:    {Ty: reflect.TypeOf(ReceiptOracle{}), Name: "LogResultPublish"},
	}
}
