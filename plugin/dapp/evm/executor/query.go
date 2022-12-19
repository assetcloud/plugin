// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync/atomic"

	log "github.com/assetcloud/chain/common/log/log15"

	"github.com/assetcloud/plugin/plugin/dapp/evm/executor/vm/runtime"

	"github.com/assetcloud/chain/common"
	"github.com/assetcloud/chain/types"
	evmAbi "github.com/assetcloud/plugin/plugin/dapp/evm/executor/abi"
	evmCommon "github.com/assetcloud/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/assetcloud/plugin/plugin/dapp/evm/executor/vm/model"
	evmtypes "github.com/assetcloud/plugin/plugin/dapp/evm/types"
)

// Query_CheckAddrExists 检查合约地址是否存在，此操作不会改变任何状态，所以可以直接从statedb查询
func (evm *EVMExecutor) Query_CheckAddrExists(in *evmtypes.CheckEVMAddrReq) (types.Message, error) {
	evm.CheckInit()
	addrStr := in.Addr
	if len(addrStr) == 0 {
		return nil, model.ErrAddrNotExists
	}

	var addr evmCommon.Address
	// 合约名称
	cfg := evm.GetAPI().GetConfig()
	if strings.HasPrefix(addrStr, cfg.ExecName(evmtypes.EvmPrefix)) {
		addr = evmCommon.ExecAddress(addrStr)
	} else {
		// 合约地址
		nAddr := evmCommon.StringToAddress(addrStr)
		if nAddr == nil {
			return nil, model.ErrAddrNotExists
		}
		addr = *nAddr
	}

	exists := evm.GetMStateDB().Exist(addr.String())
	ret := &evmtypes.CheckEVMAddrResp{Contract: exists}
	if exists {
		account := evm.GetMStateDB().GetAccount(addr.String())
		if account != nil {
			ret.ContractAddr = account.Addr
			ret.ContractName = account.GetExecName()
			ret.AliasName = account.GetAliasName()
		}
	}
	return ret, nil
}

// Query_EstimateGas 此方法用来估算合约消耗的Gas，不能修改原有执行器的状态数据
// Query_EstimateGas 此方法用来估算合约消耗的Gas，不能修改原有执行器的状态数据
func (evm *EVMExecutor) Query_EstimateGas(req *evmtypes.EstimateEVMGasReq) (types.Message, error) {
	evm.CheckInit()
	txBytes, err := hex.DecodeString(req.Tx)
	if nil != err {
		return nil, err
	}
	var tx types.Transaction
	err = types.Decode(txBytes, &tx)
	if nil != err {
		return nil, err
	}
	var index int
	from := evmCommon.StringToAddress(req.From)
	msg, err := evm.GetMessage(&tx, index, from)
	if err != nil {
		return nil, err
	}
	//初始状态设置
	var lo uint64 = 21000
	var hi uint64 = evmtypes.MaxGasLimit
	var cap = hi
	//get coins balance
	if evm.mStateDB != nil && !evm.GetAPI().GetConfig().IsPara() {
		fromBalance := evm.mStateDB.GetBalance(from.String())
		if fromBalance-msg.Value() > 0 {
			hi = fromBalance - msg.Value()
			if hi > evmtypes.MaxGasLimit {
				cap = hi
			}

		}
	}

	// 创建EVM运行时对象
	env := runtime.NewEVM(evm.NewEVMContext(msg, tx.Hash()), evm.mStateDB, *evm.vmCfg, evm.GetAPI().GetConfig())
	isTransferOnly := strings.Compare(msg.To().String(), EvmAddress) == 0 && 0 == len(msg.Data())
	//coins转账，para数据作为备注交易
	isTransferNote := strings.Compare(msg.To().String(), EvmAddress) != 0 && !env.StateDB.Exist(msg.To().String()) && len(msg.Para()) > 0 && msg.Value() != 0
	//如果是普通转账或者带有备注的Coins 转账 则直接返回
	if isTransferOnly || isTransferNote {

		result := &evmtypes.EstimateEVMGasResp{}
		result.Gas = lo
		log.Info("Query_EstimateGas", "gas:", result.Gas, "isTransferOnly:", isTransferOnly, "isTransferNote:", isTransferNote)
		return result, nil

	}

	executable := func(evm *EVMExecutor, msg *evmCommon.Message, gas uint64) (bool, *evmtypes.EstimateEVMGasResp, error) {
		msg.SetGasLimit(gas)
		receipt, err := evm.innerExec(msg, tx.Hash(), tx.GetSignature().GetTy(), index, evmtypes.MaxGasLimit, true)
		if err != nil {
			if strings.Contains(err.Error(), "out of gas") {
				return false, nil, nil
			}
			return false, nil, err
		}

		if receipt.Ty != types.ExecOk {
			return false, nil, errors.New("contract call error")
		}
		callData := getCallReceipt(receipt.GetLogs())
		if callData == nil {
			return false, nil, errors.New("nil receipt")
		}
		log.Info("executable", "evm usedGas:", callData.UsedGas)
		result := &evmtypes.EstimateEVMGasResp{}
		result.Gas = callData.UsedGas
		return true, result, nil
	}

	var count int
	//通过二分查找确定可执行的gaslimit.
	for lo+1 < hi {
		count++
		evm.mStateDB.Snapshot()
		snapID := evm.mStateDB.GetLastSnapshot().GetID()
		mid := (hi + lo) / 2
		log.Info("Query_EstimateGas", "[executable  count]:", count, "the last low gas:", lo, "the last high gas:", hi, "the mid gas:", mid)
		// ok 设置的gas可以执行
		ok, _, err := executable(evm, msg, mid)
		evm.mStateDB.RevertToSnapshot(snapID)
		if err != nil {
			return nil, err
		}
		if !ok { //如果!ok 说明GaS 不够用，则把上一轮计算的mid gas 赋值给low gas, 进而提高mid gas 的值
			lo = mid
		} else { //如果ok,说明mid Gas 有较多余量，则把mid gas 赋值给hi 降低high gas 的值,进而压缩mid gas 的值
			hi = mid
		}

	}
	log.Info("Query_EstimateGas", "[complete,executable count]:", count, "the last low gas:", lo, "the last high gas:", hi)

	if hi == cap {
		ok, result, err := executable(evm, msg, hi)
		if err != nil || !ok {
			return nil, err
		}
		return result, nil
	}

	result := &evmtypes.EstimateEVMGasResp{}
	result.Gas = quickFixGas(hi)
	log.Info("Query_EstimateGas", "gas:", result.Gas)
	return result, nil

}

func quickFixGas(gas uint64) uint64 {
  gas = gas * 12 / 10
  if gas > evmtypes.MaxGasLimit {
    gas = evmtypes.MaxGasLimit
  }
  return gas
}

func (evm *EVMExecutor) Query1_EstimateGas(req *evmtypes.EstimateEVMGasReq) (types.Message, error) {
	evm.CheckInit()

	txBytes, err := hex.DecodeString(req.Tx)
	if nil != err {
		return nil, err
	}
	var tx types.Transaction
	err = types.Decode(txBytes, &tx)
	if nil != err {
		return nil, err
	}

	index := 0
	from := evmCommon.StringToAddress(req.From)
	msg, err := evm.GetMessage(&tx, index, from)
	if err != nil {
		return nil, err
	}

	msg.SetGasLimit(evmtypes.MaxGasLimit)
	receipt, err := evm.innerExec(msg, tx.Hash(), tx.GetSignature().GetTy(), index, evmtypes.MaxGasLimit, true)
	if err != nil {
		return nil, err
	}

	if receipt.Ty != types.ExecOk {
		return nil, errors.New("contract call error")
	}
	callData := getCallReceipt(receipt.GetLogs())
	if callData == nil {
		return nil, errors.New("nil receipt")
	}
	result := &evmtypes.EstimateEVMGasResp{}
	result.Gas = callData.UsedGas
	return result, nil
}

// 从日志中查找调用结果
func getCallReceipt(logs []*types.ReceiptLog) *evmtypes.ReceiptEVMContract {
	if len(logs) == 0 {
		return nil
	}
	for _, v := range logs {
		if v.Ty == evmtypes.TyLogCallContract {
			var res evmtypes.ReceiptEVMContract
			err := types.Decode(v.Log, &res)
			if err != nil {
				return nil
			}
			return &res
		}
	}
	return nil
}

// Query_EvmDebug 此方法用来控制evm调试打印开关
func (evm *EVMExecutor) Query_EvmDebug(in *evmtypes.EvmDebugReq) (types.Message, error) {
	evm.CheckInit()
	optype := in.Optype

	if optype < 0 {
		atomic.StoreInt32(&evm.vmCfg.Debug, runtime.EVMDebugOff)
	} else if optype > 0 {
		atomic.StoreInt32(&evm.vmCfg.Debug, runtime.EVMDebugOn)
	}

	ret := &evmtypes.EvmDebugResp{DebugStatus: fmt.Sprintf("%v", evm.vmCfg.Debug)}
	return ret, nil
}

// Query_Query 此方法用来调用合约的只读接口，不修改原有执行器的状态数据
func (evm *EVMExecutor) Query_Query(in *evmtypes.EvmQueryReq) (types.Message, error) {
	evm.CheckInit()

	ret := &evmtypes.EvmQueryResp{}
	ret.Address = in.Address
	ret.Input = in.Input
	ret.Caller = in.Caller

	var (
		caller evmCommon.Address
	)

	to := evmCommon.StringToAddress(in.Address)
	if to == nil {
		ret.JsonData = fmt.Sprintf("invalid address:%v", in.Address)
		return ret, nil
	}

	// 如果未指定调用地址，则直接使用一个虚拟的地址发起调用
	cfg := evm.GetAPI().GetConfig()
	if len(in.Caller) > 0 {
		callAddr := evmCommon.StringToAddress(in.Caller)
		if callAddr != nil {
			caller = *callAddr
		}
	} else {
		caller = evmCommon.ExecAddress(cfg.ExecName(evmtypes.ExecutorName))
	}

	msg := evmCommon.NewMessage(caller, evmCommon.StringToAddress(in.Address), 0, 0, evmtypes.MaxGasLimit, 1, nil, evmCommon.FromHex(in.Input), "estimateGas")
	txHash := evmCommon.BigToHash(big.NewInt(evmtypes.MaxGasLimit)).Bytes()

	receipt, err := evm.innerExec(msg, txHash, 0, 1, evmtypes.MaxGasLimit, true)
	if err != nil {
		ret.JsonData = fmt.Sprintf("%v", err)
		return ret, nil
	}
	if receipt.Ty == types.ExecOk {
		callData := getCallReceipt(receipt.GetLogs())
		if callData != nil {
			ret.RawData = evmCommon.Bytes2Hex(callData.Ret)
			ret.JsonData = callData.JsonRet
			return ret, nil
		}
	}
	return ret, nil
}

//Query_GetNonce 获取普通账户的Nonce
func (evm *EVMExecutor) Query_GetNonce(in *evmtypes.EvmGetNonceReq) (types.Message, error) {
	evm.CheckInit()
	nonce := evm.mStateDB.GetAccountNonce(in.Address)
	return &evmtypes.EvmGetNonceRespose{Nonce: int64(nonce)}, nil
}

//Query_GetPackData ...
func (evm *EVMExecutor) Query_GetPackData(in *evmtypes.EvmGetPackDataReq) (types.Message, error) {
	evm.CheckInit()
	_, packData, err := evmAbi.Pack(in.Parameter, in.Abi, false)
	if nil != err {
		return nil, errors.New("Failed to do evmAbi.Pack" + err.Error())
	}
	packStr := common.ToHex(packData)

	return &evmtypes.EvmGetPackDataRespose{PackData: packStr}, nil
}

//Query_GetUnpackData ...
func (evm *EVMExecutor) Query_GetUnpackData(in *evmtypes.EvmGetUnpackDataReq) (types.Message, error) {
	evm.CheckInit()
	data, err := common.FromHex(in.Data)
	if nil != err {
		return nil, errors.New("common.FromHex failed due to:" + err.Error())
	}

	outputs, err := evmAbi.UnpackAllTypes(data, in.Name, in.Abi)
	if err != nil {
		return nil, errors.New("unpack evm return error" + err.Error())
	}

	ret := evmtypes.EvmGetUnpackDataRespose{}

	for _, v := range outputs {
		ret.UnpackData = append(ret.UnpackData, fmt.Sprintf("%v", v.Value))
	}
	return &ret, nil
}

//Query_GetCode 获取合约地址下的code
func (evm *EVMExecutor) Query_GetCode(in *evmtypes.CheckEVMAddrReq) (types.Message, error) {
	evm.CheckInit()
	addrStr := in.Addr
	if len(addrStr) == 0 {
		return nil, model.ErrAddrNotExists
	}

	addr := evmCommon.StringToAddress(in.GetAddr())
	log.Debug("Query_GetCode", "addr", in.GetAddr(), "addrstring", addr.String())
	codeData := evm.mStateDB.GetCode(addr.String())
	abiData := evm.mStateDB.GetAbi(addr.String())
	account := evm.mStateDB.GetAccount(addr.String())
	var ret evmtypes.EVMContractData
	ret.Code = codeData
	ret.Abi = abiData
	if account != nil {
		ret.Creator = account.GetCreator()
		ret.Alias = account.GetAliasName()
	}
	return &ret, nil

}
