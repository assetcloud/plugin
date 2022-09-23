package chain

import (
	"errors"
	"fmt"

	dbm "github.com/assetcloud/chain/common/db"
	chainTypes "github.com/assetcloud/chain/types"
	"github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/relayer/events"
	ebTypes "github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/types"
	"github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/utils"
)

//key ...
var (
	lastSyncHeightPrefix             = []byte("chain-lastSyncHeight:")
	eth2ChainBurnLockTxStaticsPrefix = "chain-eth2chainBurnLockStatics"
	eth2ChainBurnLockTxFinished      = "chain-eth2ChainBurnLockTxFinished"
	relayEthBurnLockTxTotalAmount    = []byte("chain-relayEthBurnLockTxTotalAmount")
	chainBurnTxUpdateTxIndex         = []byte("chain-chainBurnTxUpdateTxIndx")
	chainLockTxUpdateTxIndex         = []byte("chain-chainLockTxUpdateTxIndex")
	bridgeRegistryAddrOnChain        = []byte("chain-x2EthBridgeRegistryAddrOnChain")
	tokenSymbol2AddrPrefix           = []byte("chain-chainTokenSymbol2AddrPrefix")
	multiSignAddressPrefix           = []byte("chain-multiSignAddress")
	symbol2Ethchain                  = []byte("chain-symbol2Ethchain")
	txIsRelayedUnconfirm             = []byte("chain-txIsRelayedUnconfirm")
	chainTxRelayedAlready            = []byte("chain-txRelayedAlready")
	fdTx2EthTotalAmount              = []byte("chain-fdTx2EthTotalAmount")
	ethTxRelayAlreadyPrefix          = []byte("chain-ethTxRelayAlready")
)

func ethTxRelayAlreadyKey(chainTxhash string) []byte {
	return append(ethTxRelayAlreadyPrefix, []byte(fmt.Sprintf("-txHash-%s", chainTxhash))...)
}

func chainTxIsRelayedUnconfirmKey(txHash string) []byte {
	return append(txIsRelayedUnconfirm, []byte(fmt.Sprintf("-txHash-%s", txHash))...)
}

func chainTxRelayedAlreadyKey(txHash string) []byte {
	return append(chainTxRelayedAlready, []byte(fmt.Sprintf("-txHash-%s", txHash))...)
}

func tokenSymbol2AddrKey(symbol string) []byte {
	return append(tokenSymbol2AddrPrefix, []byte(fmt.Sprintf("-symbol-%s", symbol))...)
}

func calcRelayFromEthStaticsKey(txindex int64, claimType int32) []byte {
	return []byte(fmt.Sprintf("%s-%d-%012d", eth2ChainBurnLockTxStaticsPrefix, claimType, txindex))
}

//未完成，处在pending状态
func calcRelayFromEthStaticsList(claimType int32) []byte {
	return []byte(fmt.Sprintf("%s-%d-", eth2ChainBurnLockTxStaticsPrefix, claimType))
}

func calcFromEthFinishedStaticsKey(txindex int64, claimType int32) []byte {
	return []byte(fmt.Sprintf("%s-%d-%012d", eth2ChainBurnLockTxFinished, claimType, txindex))
}

func calcFromEthFinishedStaticsList(claimType int32) []byte {
	return []byte(fmt.Sprintf("%s-%d-", eth2ChainBurnLockTxFinished, claimType))
}

func (chainRelayer *Relayer4Chain) updateFdTx2EthTotalAmount(index int64) error {
	totalTx := &chainTypes.Int64{
		Data: index,
	}
	//更新成功见证的交易数
	return chainRelayer.db.SetSync(fdTx2EthTotalAmount, chainTypes.Encode(totalTx))
}

func (chainRelayer *Relayer4Chain) getFdTx2EthTotalAmount() int64 {
	totalTx, _ := utils.LoadInt64FromDB(fdTx2EthTotalAmount, chainRelayer.db)
	return totalTx
}

func (chainRelayer *Relayer4Chain) getAllTxsUnconfirm() (txInfos []*ebTypes.TxRelayConfirm4Chain, err error) {
	helper := dbm.NewListHelper(chainRelayer.db)
	datas := helper.List(txIsRelayedUnconfirm, nil, 0, dbm.ListASC)
	cnt := len(datas)
	if 0 == cnt {
		return nil, nil
	}

	txInfos = make([]*ebTypes.TxRelayConfirm4Chain, cnt)
	for i, data := range datas {
		txInfo := &ebTypes.TxRelayConfirm4Chain{}
		if err := chainTypes.Decode(data, txInfo); nil != err {
			return nil, err
		}

		txInfos[i] = txInfo
	}
	return
}

func (chainRelayer *Relayer4Chain) resetKeyChainTxRelayedAlready(txHash string) error {
	key := chainTxIsRelayedUnconfirmKey(txHash)
	data, err := chainRelayer.db.Get(key)
	if nil != err {
		relayerLog.Info("resetKeyTxRelayedAlready", "No data for tx", txHash)
		return err
	}
	_ = chainRelayer.db.DeleteSync(key)
	setkey := chainTxRelayedAlreadyKey(txHash)

	return chainRelayer.db.SetSync(setkey, data)
}

func (chainRelayer *Relayer4Chain) setChainTxIsRelayedUnconfirm(txHash string, index int64, txRelayConfirm4Chain *ebTypes.TxRelayConfirm4Chain) error {
	key := chainTxIsRelayedUnconfirmKey(txHash)
	data := chainTypes.Encode(txRelayConfirm4Chain)
	relayerLog.Info("setChainTxIsRelayedUnconfirm", "TxHash", txHash, "index", index, "ForwardTimes", txRelayConfirm4Chain.FdTimes)
	return chainRelayer.db.SetSync(key, data)
}

func (chainRelayer *Relayer4Chain) setEthTxRelayAlreadyInfo(ethTxhash string, relayTxDetail *ebTypes.RelayTxDetail) error {
	key := ethTxRelayAlreadyKey(ethTxhash)
	data := chainTypes.Encode(relayTxDetail)
	return chainRelayer.db.SetSync(key, data)
}

func (chainRelayer *Relayer4Chain) getEthTxRelayAlreadyInfo(ethTxhash string) (*ebTypes.RelayTxDetail, error) {
	key := ethTxRelayAlreadyKey(ethTxhash)
	data, err := chainRelayer.db.Get(key)
	if nil != err {
		return nil, err
	}
	var relayTxDetail ebTypes.RelayTxDetail
	err = chainTypes.Decode(data, &relayTxDetail)
	return &relayTxDetail, err
}

func (chainRelayer *Relayer4Chain) updateTotalTxAmount2Eth(txIndex int64) error {
	totalTx := &chainTypes.Int64{
		Data: txIndex,
	}
	//更新成功见证的交易数
	return chainRelayer.db.SetSync(relayEthBurnLockTxTotalAmount, chainTypes.Encode(totalTx))
}

func (chainRelayer *Relayer4Chain) getTotalTxAmount() int64 {
	totalTx, _ := utils.LoadInt64FromDB(relayEthBurnLockTxTotalAmount, chainRelayer.db)
	return totalTx
}

func (chainRelayer *Relayer4Chain) setLastestRelay2ChainTxStatics(txIndex int64, claimType int32, data []byte) error {
	key := calcRelayFromEthStaticsKey(txIndex, claimType)
	return chainRelayer.db.SetSync(key, data)
}

func (chainRelayer *Relayer4Chain) getStatics(claimType int32, txIndex int64, count int32) ([][]byte, error) {
	//第一步：获取处在pending状态的
	keyPrefix := calcRelayFromEthStaticsList(claimType)
	keyFrom := calcRelayFromEthStaticsKey(txIndex, claimType)
	helper := dbm.NewListHelper(chainRelayer.db)
	datas := helper.List(keyPrefix, keyFrom, count, dbm.ListASC)
	if nil == datas {
		return nil, errors.New("Not found")
	}

	return datas, nil
}

func (chainRelayer *Relayer4Chain) setChainUpdateTxIndex(txindex int64, claimType events.ClaimType) error {
	txIndexWrapper := &chainTypes.Int64{
		Data: txindex,
	}

	if events.ClaimTypeBurn == claimType {
		return chainRelayer.db.SetSync(chainBurnTxUpdateTxIndex, chainTypes.Encode(txIndexWrapper))
	}
	return chainRelayer.db.SetSync(chainLockTxUpdateTxIndex, chainTypes.Encode(txIndexWrapper))
}

func (chainRelayer *Relayer4Chain) getChainUpdateTxIndex(claimType events.ClaimType) int64 {
	var key []byte
	if events.ClaimTypeBurn == claimType {
		key = chainBurnTxUpdateTxIndex
	} else {
		key = chainLockTxUpdateTxIndex
	}
	data, err := chainRelayer.db.Get(key)
	if nil != err {
		return ebTypes.Invalid_Tx_Index
	}

	var txIndexWrapper chainTypes.Int64
	err = chainTypes.Decode(data, &txIndexWrapper)
	if nil != err {
		return ebTypes.Invalid_Tx_Index
	}
	return txIndexWrapper.Data
}

//获取上次同步到app的高度
func (chainRelayer *Relayer4Chain) loadLastSyncHeight() int64 {
	height, err := utils.LoadInt64FromDB(lastSyncHeightPrefix, chainRelayer.db)
	if nil != err && err != chainTypes.ErrHeightNotExist {
		relayerLog.Error("loadLastSyncHeight", "err:", err.Error())
		return 0
	}
	return height
}

func (chainRelayer *Relayer4Chain) setLastSyncHeight(syncHeight int64) {
	bytes := chainTypes.Encode(&chainTypes.Int64{Data: syncHeight})
	_ = chainRelayer.db.SetSync(lastSyncHeightPrefix, bytes)
}

func (chainRelayer *Relayer4Chain) setBridgeRegistryAddr(bridgeRegistryAddr string) error {
	return chainRelayer.db.SetSync(bridgeRegistryAddrOnChain, []byte(bridgeRegistryAddr))
}

func (chainRelayer *Relayer4Chain) getBridgeRegistryAddr() (string, error) {
	addr, err := chainRelayer.db.Get(bridgeRegistryAddrOnChain)
	if nil != err {
		return "", err
	}
	return string(addr), nil
}

func (chainRelayer *Relayer4Chain) SetTokenAddress(token2set *ebTypes.TokenAddress) error {
	bytes := chainTypes.Encode(token2set)
	chainRelayer.rwLock.Lock()
	chainRelayer.symbol2Addr[token2set.Symbol] = token2set.Address
	chainRelayer.rwLock.Unlock()
	return chainRelayer.db.SetSync(tokenSymbol2AddrKey(token2set.Symbol), bytes)
}

func (chainRelayer *Relayer4Chain) RestoreTokenAddress() error {
	chainRelayer.rwLock.Lock()
	defer chainRelayer.rwLock.Unlock()
	chainRelayer.symbol2Addr[ebTypes.SYMBOL_BTY] = ebTypes.BTYAddrChain

	helper := dbm.NewListHelper(chainRelayer.db)
	datas := helper.List(tokenSymbol2AddrPrefix, nil, 100, dbm.ListASC)
	if nil == datas {
		return nil
	}

	for _, data := range datas {
		var token2set ebTypes.TokenAddress
		err := chainTypes.Decode(data, &token2set)
		if nil != err {
			return err
		}
		relayerLog.Info("RestoreTokenAddress", "symbol", token2set.Symbol, "address", token2set.Address)
		chainRelayer.symbol2Addr[token2set.Symbol] = token2set.Address
	}
	return nil
}

func (chainRelayer *Relayer4Chain) ShowTokenAddress(token2show *ebTypes.TokenAddress) (*ebTypes.TokenAddressArray, error) {
	res := &ebTypes.TokenAddressArray{}

	if len(token2show.Symbol) > 0 {
		data, err := chainRelayer.db.Get(tokenSymbol2AddrKey(token2show.Symbol))
		if err != nil {
			return nil, err
		}
		var token2set ebTypes.TokenAddress
		err = chainTypes.Decode(data, &token2set)
		if nil != err {
			return nil, err
		}
		res.TokenAddress = append(res.TokenAddress, &token2set)
		return res, nil
	}
	helper := dbm.NewListHelper(chainRelayer.db)
	datas := helper.List(tokenSymbol2AddrPrefix, nil, 100, dbm.ListASC)
	if nil == datas {
		return nil, errors.New("Not found")
	}

	for _, data := range datas {

		var token2set ebTypes.TokenAddress
		err := chainTypes.Decode(data, &token2set)
		if nil != err {
			return nil, err
		}
		res.TokenAddress = append(res.TokenAddress, &token2set)

	}
	return res, nil
}

func (chainRelayer *Relayer4Chain) setMultiSignAddress(address string) {
	bytes := []byte(address)
	_ = chainRelayer.db.SetSync(multiSignAddressPrefix, bytes)
}

func (chainRelayer *Relayer4Chain) getMultiSignAddress() string {
	bytes, _ := chainRelayer.db.Get(multiSignAddressPrefix)
	if 0 == len(bytes) {
		return ""
	}
	return string(bytes)
}

func (chainRelayer *Relayer4Chain) storeSymbol2chainName(symbol2Name map[string]string) {
	Symbol2EthChain := &ebTypes.Symbol2EthChain{
		Symbol2Name: symbol2Name,
	}
	data := chainTypes.Encode(Symbol2EthChain)
	_ = chainRelayer.db.SetSync(symbol2Ethchain, data)
}

func (chainRelayer *Relayer4Chain) restoreSymbol2chainName() map[string]string {
	data, _ := chainRelayer.db.Get(symbol2Ethchain)
	if 0 == len(data) {
		return make(map[string]string)
	}

	symbol2EthChain := &ebTypes.Symbol2EthChain{}
	if err := chainTypes.Decode(data, symbol2EthChain); nil != err {
		return make(map[string]string)
	}
	return symbol2EthChain.Symbol2Name
}

//判断是否已经被处理，如果能够在数据库中找到该笔交易，则认为已经被处理
func (chainRelayer *Relayer4Chain) checkTxProcessed(txhash string) bool {
	key1 := chainTxIsRelayedUnconfirmKey(txhash)
	data, err := chainRelayer.db.Get(key1)
	if 0 != len(data) && nil == err {
		return true
	}

	key2 := chainTxRelayedAlreadyKey(txhash)
	data, err = chainRelayer.db.Get(key2)
	if 0 != len(data) && nil == err {
		return true
	}

	return false
}
