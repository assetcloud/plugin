// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"reflect"

	log "github.com/assetcloud/chain/common/log/log15"
	coins "github.com/assetcloud/chain/system/dapp/coins/types"
	"github.com/assetcloud/chain/types"
	token "github.com/assetcloud/plugin/plugin/dapp/token/types"
)

var (
	// ParaX paracross exec name
	MixX = "mix"
	glog = log.New("module", MixX)
)

func init() {
	// init executor type
	types.AllowUserExec = append(types.AllowUserExec, []byte(MixX))
	types.RegFork(MixX, InitFork)
	types.RegExec(MixX, InitExecutor)

}

//InitFork ...
func InitFork(cfg *types.ChainConfig) {
	cfg.RegisterDappFork(MixX, "Enable", 0)

}

//InitExecutor ...
func InitExecutor(cfg *types.ChainConfig) {
	types.RegistorExecutor(MixX, NewType(cfg))
}

// GetExecName get para exec name
func GetExecName(cfg *types.ChainConfig) string {
	return cfg.ExecName(MixX)
}

// ParacrossType base paracross type
type MixType struct {
	types.ExecTypeBase
}

// NewType get paracross type
func NewType(cfg *types.ChainConfig) *MixType {
	c := &MixType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// GetName 获取执行器名称
func (p *MixType) GetName() string {
	return MixX
}

// GetLogMap get receipt log map
func (p *MixType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{
		TyLogMixConfigVk:           {Ty: reflect.TypeOf(ZkVerifyKeys{}), Name: "LogMixConfigVk"},
		TyLogMixConfigAuth:         {Ty: reflect.TypeOf(AuthKeys{}), Name: "LogMixConfigAuthPubKey"},
		TyLogSubLeaves:             {Ty: reflect.TypeOf(ReceiptCommitSubLeaves{}), Name: "LogSubLeaves"},
		TyLogCommitTreeStatus:      {Ty: reflect.TypeOf(ReceiptCommitTreeStatus{}), Name: "LogCommitTreeStatus"},
		TyLogSubRoots:              {Ty: reflect.TypeOf(ReceiptCommitSubRoots{}), Name: "LogSubRoots"},
		TyLogArchiveRootLeaves:     {Ty: reflect.TypeOf(ReceiptArchiveLeaves{}), Name: "LogArchiveRootLeaves"},
		TyLogCommitTreeArchiveRoot: {Ty: reflect.TypeOf(ReceiptArchiveTreeRoot{}), Name: "LogCommitTreeArchiveRoot"},
		TyLogMixConfigPaymentKey:   {Ty: reflect.TypeOf(NoteAccountKey{}), Name: "LogConfigReceivingKey"},
		TyLogNulliferSet:           {Ty: reflect.TypeOf(ExistValue{}), Name: "LogNullifierSet"},
	}
}

// GetTypeMap get action type
func (p *MixType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"Config":    MixActionConfig,
		"Deposit":   MixActionDeposit,
		"Withdraw":  MixActionWithdraw,
		"Transfer":  MixActionTransfer,
		"Authorize": MixActionAuth,
	}
}

// GetPayload mix get action payload
func (p *MixType) GetPayload() types.Message {
	return &MixAction{}
}

func GetAssetExecSymbol(cfg *types.ChainConfig, execer, symbol string) (string, string) {
	if symbol == "" {
		return coins.CoinsX, cfg.GetCoinSymbol()
	}
	if execer == "" {
		return token.TokenX, symbol
	}
	return execer, symbol
}

func GetTransferTxFee(cfg *types.ChainConfig, assetExecer string) int64 {
	conf := types.ConfSub(cfg, MixX)
	txFee := conf.GInt("txFee")
	tokenFee := conf.IsEnable("tokenFee")
	//一切非coins的token资产 在tokenFee=false都不收txfee,特殊地址代扣
	if assetExecer != coins.CoinsX && !tokenFee {
		return 0
	}
	//tokenFee=true或者coins都按txfee数量收txFee
	return txFee
}
