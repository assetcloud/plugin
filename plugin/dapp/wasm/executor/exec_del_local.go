package executor

import (
	types2 "github.com/assetcloud/plugin/plugin/dapp/wasm/types"
	"github.com/assetcloud/chain/types"
)

func (w *Wasm) ExecDelLocal_Create(payload *types2.WasmCreate, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return &types.LocalDBSet{}, nil
}

func (w *Wasm) ExecDelLocal_Update(payload *types2.WasmUpdate, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return &types.LocalDBSet{}, nil
}

func (w *Wasm) ExecDelLocal_Call(payload *types2.WasmCall, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	localExecer := w.userExecName(payload.Contract, true)
	kvs, err := w.DelRollbackKV(tx, []byte(localExecer))
	if err != nil {
		return nil, err
	}
	return &types.LocalDBSet{KV: kvs}, nil

}
