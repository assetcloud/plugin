package chain

import (
	"fmt"
	"sync/atomic"

	"github.com/assetcloud/chain/types"
	ebTypes "github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/types"
	"github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/utils"
	"github.com/ethereum/go-ethereum/common"
)

//key ...
var (
	lastSyncHeightPrefix              = []byte("lastSyncHeight:")
	chainToEthBurnLockTxHashPrefix  = "chainToEthBurnLockTxHash"
	chainToEthBurnLockTxTotalAmount = []byte("chainToEthBurnLockTxTotalAmount")
	EthTxStatusCheckedIndex           = []byte("EthTxStatusCheckedIndex")
)

func calcRelay2EthTxhash(txindex int64) []byte {
	return []byte(fmt.Sprintf("%s-%012d", chainToEthBurnLockTxHashPrefix, txindex))
}

func (chainRelayer *Relayer4Chain) updateTotalTxAmount2Eth(total int64) error {
	totalTx := &types.Int64{
		Data: atomic.LoadInt64(&chainRelayer.totalTx4ChainToEth),
	}
	//更新成功见证的交易数
	return chainRelayer.db.Set(chainToEthBurnLockTxTotalAmount, types.Encode(totalTx))
}

func (chainRelayer *Relayer4Chain) getTotalTxAmount2Eth() int64 {
	totalTx, _ := utils.LoadInt64FromDB(chainToEthBurnLockTxTotalAmount, chainRelayer.db)
	return totalTx
}

func (chainRelayer *Relayer4Chain) setLastestRelay2EthTxhash(status, txhash string, txIndex int64) error {
	key := calcRelay2EthTxhash(txIndex)
	ethTxStatus := &ebTypes.EthTxStatus{
		Status: status,
		Txhash: txhash,
	}
	data := types.Encode(ethTxStatus)
	return chainRelayer.db.Set(key, data)
}

func (chainRelayer *Relayer4Chain) getEthTxhash(txIndex int64) (common.Hash, error) {
	key := calcRelay2EthTxhash(txIndex)
	ethTxStatus := &ebTypes.EthTxStatus{}
	data, err := chainRelayer.db.Get(key)
	if nil != err {
		return common.Hash{}, err
	}
	err = types.Decode(data, ethTxStatus)
	if nil != err {
		return common.Hash{}, err
	}
	return common.HexToHash(ethTxStatus.Txhash), nil
}

func (chainRelayer *Relayer4Chain) setStatusCheckedIndex(txIndex int64) error {
	index := &types.Int64{
		Data: txIndex,
	}
	data := types.Encode(index)
	return chainRelayer.db.Set(EthTxStatusCheckedIndex, data)
}

func (chainRelayer *Relayer4Chain) getStatusCheckedIndex() int64 {
	index, _ := utils.LoadInt64FromDB(EthTxStatusCheckedIndex, chainRelayer.db)
	return index
}

//获取上次同步到app的高度
func (chainRelayer *Relayer4Chain) loadLastSyncHeight() int64 {
	height, err := utils.LoadInt64FromDB(lastSyncHeightPrefix, chainRelayer.db)
	if nil != err && err != types.ErrHeightNotExist {
		relayerLog.Error("loadLastSyncHeight", "err:", err.Error())
		return 0
	}
	return height
}

func (chainRelayer *Relayer4Chain) setLastSyncHeight(syncHeight int64) {
	bytes := types.Encode(&types.Int64{Data: syncHeight})
	_ = chainRelayer.db.Set(lastSyncHeightPrefix, bytes)
}
