package ethereum

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	dbm "github.com/assetcloud/chain/common/db"
	chainTypes "github.com/assetcloud/chain/types"
	ebTypes "github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/types"
	"github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/utils"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	eth2chainTxHashPrefix        = "Eth2chainTxHash"
	eth2chainTxTotalAmount       = []byte("Eth2chainTxTotalAmount")
	chainToEthTxHashPrefix       = "chainToEthTxHash"
	bridgeRegistryAddrPrefix       = []byte("x2EthBridgeRegistryAddr")
	bridgeBankLogProcessedAt       = []byte("bridgeBankLogProcessedAt")
	ethTxEventPrefix               = []byte("ethTxEventPrefix")
	lastBridgeBankHeightProcPrefix = []byte("lastBridgeBankHeight")
)

func ethTxEventKey4Height(height uint64, index uint32) []byte {
	return append(ethTxEventPrefix, []byte(fmt.Sprintf("%020d-%d", height, index))...)
}

func calcRelay2ChainTxhash(txindex int64) []byte {
	return []byte(fmt.Sprintf("%s-%012d", eth2chainTxHashPrefix, txindex))
}

func (ethRelayer *Relayer4Ethereum) setBridgeRegistryAddr(bridgeRegistryAddr string) error {
	return ethRelayer.db.Set(bridgeRegistryAddrPrefix, []byte(bridgeRegistryAddr))
}

func (ethRelayer *Relayer4Ethereum) getBridgeRegistryAddr() (string, error) {
	addr, err := ethRelayer.db.Get(bridgeRegistryAddrPrefix)
	if nil != err {
		return "", err
	}
	return string(addr), nil
}

func (ethRelayer *Relayer4Ethereum) updateTotalTxAmount2chain(total int64) error {
	totalTx := &chainTypes.Int64{
		Data: atomic.LoadInt64(&ethRelayer.totalTx4Eth2Chain),
	}
	//更新成功见证的交易数
	return ethRelayer.db.Set(eth2chainTxTotalAmount, chainTypes.Encode(totalTx))
}

func (ethRelayer *Relayer4Ethereum) setLastestRelay2ChainTxhash(txhash string, txIndex int64) error {
	key := calcRelay2ChainTxhash(txIndex)
	return ethRelayer.db.Set(key, []byte(txhash))
}

func (ethRelayer *Relayer4Ethereum) queryTxhashes(prefix []byte) []string {
	return utils.QueryTxhashes(prefix, ethRelayer.db)
}

func (ethRelayer *Relayer4Ethereum) setHeight4BridgeBankLogAt(height uint64) error {
	return ethRelayer.setLogProcHeight(bridgeBankLogProcessedAt, height)
}

func (ethRelayer *Relayer4Ethereum) getHeight4BridgeBankLogAt() uint64 {
	return ethRelayer.getLogProcHeight(bridgeBankLogProcessedAt)
}

func (ethRelayer *Relayer4Ethereum) setLogProcHeight(key []byte, height uint64) error {
	data := &ebTypes.Uint64{
		Data: height,
	}
	return ethRelayer.db.Set(key, chainTypes.Encode(data))
}

func (ethRelayer *Relayer4Ethereum) getLogProcHeight(key []byte) uint64 {
	value, err := ethRelayer.db.Get(key)
	if nil != err {
		return 0
	}
	var height ebTypes.Uint64
	err = chainTypes.Decode(value, &height)
	if nil != err {
		return 0
	}
	return height.Data
}

//保存处理过的交易
func (ethRelayer *Relayer4Ethereum) setTxProcessed(txhash []byte) error {
	return ethRelayer.db.Set(txhash, []byte("1"))
}

//判断是否已经被处理，如果能够在数据库中找到该笔交易，则认为已经被处理
func (ethRelayer *Relayer4Ethereum) checkTxProcessed(txhash []byte) bool {
	_, err := ethRelayer.db.Get(txhash)
	return nil == err
}

func (ethRelayer *Relayer4Ethereum) setEthTxEvent(vLog types.Log) error {
	key := ethTxEventKey4Height(vLog.BlockNumber, uint32(vLog.TxIndex))
	value, err := json.Marshal(vLog)
	if nil != err {
		return err
	}
	return ethRelayer.db.Set(key, value)
}

func (ethRelayer *Relayer4Ethereum) getEthTxEvent(blockNumber uint64, txIndex uint32) (*types.Log, error) {
	key := ethTxEventKey4Height(blockNumber, txIndex)
	data, err := ethRelayer.db.Get(key)
	if nil != err {
		return nil, err
	}
	log := types.Log{}
	err = json.Unmarshal(data, &log)
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (ethRelayer *Relayer4Ethereum) getNextValidEthTxEventLogs(height uint64, index uint32, fetchCnt int32) ([]*types.Log, error) {
	key := ethTxEventKey4Height(height, index)
	helper := dbm.NewListHelper(ethRelayer.db)
	datas := helper.List(ethTxEventPrefix, key, fetchCnt, dbm.ListASC)
	if nil == datas {
		return nil, nil
	}
	var logs []*types.Log
	for _, data := range datas {
		log := &types.Log{}
		err := json.Unmarshal(data, log)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (ethRelayer *Relayer4Ethereum) setBridgeBankProcessedHeight(height uint64, index uint32) {
	bytes := chainTypes.Encode(&ebTypes.EventLogIndex{
		Height: height,
		Index:  index})
	_ = ethRelayer.db.Set(lastBridgeBankHeightProcPrefix, bytes)
}

func (ethRelayer *Relayer4Ethereum) getLastBridgeBankProcessedHeight() ebTypes.EventLogIndex {
	data, err := ethRelayer.db.Get(lastBridgeBankHeightProcPrefix)
	if nil != err {
		return ebTypes.EventLogIndex{}
	}
	logIndex := ebTypes.EventLogIndex{}
	_ = chainTypes.Decode(data, &logIndex)
	return logIndex
}

//构建一个引导查询使用的bridgeBankTx
func (ethRelayer *Relayer4Ethereum) initBridgeBankTx() {
	log, _ := ethRelayer.getEthTxEvent(0, 0)
	if nil != log {
		return
	}
	_ = ethRelayer.setEthTxEvent(types.Log{})
}
