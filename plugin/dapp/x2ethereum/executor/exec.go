package executor

import (
	"errors"

	"github.com/assetcloud/chain/system/dapp"
	manTy "github.com/assetcloud/chain/system/dapp/manage/types"

	"github.com/assetcloud/chain/types"
	x2eTy "github.com/assetcloud/plugin/plugin/dapp/x2ethereum/types"
)

/*
 * 实现交易的链上执行接口
 * 关键数据上链（statedb）并生成交易回执（log）
 */

//---------------- Ethereum(eth/erc20) --> Chain-------------------//

// 在chain上为ETH/ERC20铸币
func (x *x2ethereum) Exec_Eth2ChainLock(payload *x2eTy.Eth2Chain, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	if action == nil {
		return nil, errors.New("Create Action Error")
	}

	payload.ValidatorAddress = tx.From()

	return action.procEth2Chain_lock(payload)
}

//----------------  Chain(eth/erc20)------> Ethereum -------------------//
// 在chain端将铸的币销毁，返还给eth
func (x *x2ethereum) Exec_ChainToEthBurn(payload *x2eTy.ChainToEth, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	if action == nil {
		return nil, errors.New("Create Action Error")
	}
	return action.procChainToEth_burn(payload)
}

//---------------- Ethereum (bty) --> Chain-------------------//
// 在eth端将铸的bty币销毁，返还给chain
func (x *x2ethereum) Exec_Eth2ChainBurn(payload *x2eTy.Eth2Chain, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	if action == nil {
		return nil, errors.New("Create Action Error")
	}

	payload.ValidatorAddress = tx.From()

	return action.procEth2Chain_burn(payload)
}

//---------------- Chain --> Ethereum (bty) -------------------//
// 在 ethereum 上为 chain 铸币
func (x *x2ethereum) Exec_ChainToEthLock(payload *x2eTy.ChainToEth, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	if action == nil {
		return nil, errors.New("Create Action Error")
	}
	return action.procChainToEth_lock(payload)
}

// 转账功能
func (x *x2ethereum) Exec_Transfer(payload *types.AssetsTransfer, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	if action == nil {
		return nil, errors.New("Create Action Error")
	}
	return action.procMsgTransfer(payload)
}

func (x *x2ethereum) Exec_TransferToExec(payload *types.AssetsTransferToExec, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	if action == nil {
		return nil, errors.New("Create Action Error")
	}
	if !x2eTy.IsExecAddrMatch(payload.ExecName, tx.GetRealToAddr()) {
		return nil, types.ErrToAddrNotSameToExecAddr
	}
	return action.procMsgTransferToExec(payload)
}

func (x *x2ethereum) Exec_WithdrawFromExec(payload *types.AssetsWithdraw, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	if action == nil {
		return nil, errors.New("Create Action Error")
	}
	if dapp.IsDriverAddress(tx.GetRealToAddr(), x.GetHeight()) || x2eTy.IsExecAddrMatch(payload.ExecName, tx.GetRealToAddr()) {
		return action.procMsgWithDrawFromExec(payload)
	}
	return nil, errors.New("tx error")
}

//--------------------------合约管理员账户操作-------------------------//

// AddValidator是为了增加validator
func (x *x2ethereum) Exec_AddValidator(payload *x2eTy.MsgValidator, tx *types.Transaction, index int) (*types.Receipt, error) {
	confManager := types.ConfSub(x.GetAPI().GetConfig(), manTy.ManageX).GStrList("superManager")
	err := checkTxSignBySpecificAddr(tx, confManager)
	if err == nil {
		action := newAction(x, tx, int32(index))
		if action == nil {
			return nil, errors.New("Create Action Error")
		}
		return action.procAddValidator(payload)
	}
	return nil, err
}

// RemoveValidator是为了移除某一个validator
func (x *x2ethereum) Exec_RemoveValidator(payload *x2eTy.MsgValidator, tx *types.Transaction, index int) (*types.Receipt, error) {
	confManager := types.ConfSub(x.GetAPI().GetConfig(), manTy.ManageX).GStrList("superManager")
	err := checkTxSignBySpecificAddr(tx, confManager)
	if err == nil {
		action := newAction(x, tx, int32(index))
		if action == nil {
			return nil, errors.New("Create Action Error")
		}
		return action.procRemoveValidator(payload)
	}
	return nil, err
}

// ModifyPower是为了修改某个validator的power
func (x *x2ethereum) Exec_ModifyPower(payload *x2eTy.MsgValidator, tx *types.Transaction, index int) (*types.Receipt, error) {
	confManager := types.ConfSub(x.GetAPI().GetConfig(), manTy.ManageX).GStrList("superManager")
	err := checkTxSignBySpecificAddr(tx, confManager)
	if err == nil {
		action := newAction(x, tx, int32(index))
		if action == nil {
			return nil, errors.New("Create Action Error")
		}
		return action.procModifyValidator(payload)
	}
	return nil, err
}

// SetConsensusThreshold是为了修改对validator所提供的claim达成共识的阈值
func (x *x2ethereum) Exec_SetConsensusThreshold(payload *x2eTy.MsgConsensusThreshold, tx *types.Transaction, index int) (*types.Receipt, error) {
	confManager := types.ConfSub(x.GetAPI().GetConfig(), manTy.ManageX).GStrList("superManager")
	err := checkTxSignBySpecificAddr(tx, confManager)
	if err == nil {
		action := newAction(x, tx, int32(index))
		if action == nil {
			return nil, errors.New("Create Action Error")
		}
		return action.procMsgSetConsensusThreshold(payload)
	}
	return nil, err
}

func checkTxSignBySpecificAddr(tx *types.Transaction, addrs []string) error {
	signAddr := tx.From()
	var exist bool
	for _, addr := range addrs {
		if signAddr == addr {
			exist = true
			break
		}
	}

	if !exist {
		return x2eTy.ErrInvalidAdminAddress
	}

	if !tx.CheckSign(0) {
		return types.ErrSign
	}
	return nil
}
