// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"reflect"

	"github.com/assetcloud/chain/types"
)

func init() {
	// init executor type
	types.AllowUserExec = append(types.AllowUserExec, ExecerPokerBull)
	types.RegFork(PokerBullX, InitFork)
	types.RegExec(PokerBullX, InitExecutor)
}

//InitFork ...
func InitFork(cfg *types.ChainConfig) {
	cfg.RegisterDappFork(PokerBullX, "Enable", 0)
}

//InitExecutor ...
func InitExecutor(cfg *types.ChainConfig) {
	types.RegistorExecutor(PokerBullX, NewType(cfg))
}

// PokerBullType 斗牛执行器类型
type PokerBullType struct {
	types.ExecTypeBase
}

// NewType 创建pokerbull执行器类型
func NewType(cfg *types.ChainConfig) *PokerBullType {
	c := &PokerBullType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// GetName 获取执行器名称
func (t *PokerBullType) GetName() string {
	return PokerBullX
}

// GetPayload 获取payload
func (t *PokerBullType) GetPayload() types.Message {
	return &PBGameAction{}
}

// GetTypeMap 获取类型map
func (t *PokerBullType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"Start":    PBGameActionStart,
		"Continue": PBGameActionContinue,
		"Quit":     PBGameActionQuit,
		"Query":    PBGameActionQuery,
		"Play":     PBGameActionPlay,
	}
}

// GetLogMap 获取日志map
func (t *PokerBullType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{
		TyLogPBGameStart:    {Ty: reflect.TypeOf(ReceiptPBGame{}), Name: "TyLogPBGameStart"},
		TyLogPBGameContinue: {Ty: reflect.TypeOf(ReceiptPBGame{}), Name: "TyLogPBGameContinue"},
		TyLogPBGameQuit:     {Ty: reflect.TypeOf(ReceiptPBGame{}), Name: "TyLogPBGameQuit"},
		TyLogPBGameQuery:    {Ty: reflect.TypeOf(ReceiptPBGame{}), Name: "TyLogPBGameQuery"},
		TyLogPBGamePlay:     {Ty: reflect.TypeOf(ReceiptPBGame{}), Name: "TyLogPBGamePlay"},
	}
}
