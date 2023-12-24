package chain

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	dbm "github.com/assetcloud/chain/common/db"
	log "github.com/assetcloud/chain/common/log/log15"
	"github.com/assetcloud/chain/rpc/jsonclient"
	rpctypes "github.com/assetcloud/chain/rpc/types"
	chainTypes "github.com/assetcloud/chain/types"
	"github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/ethcontract/generated"
	"github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/ethinterface"
	relayerTx "github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/ethtxs"
	"github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/events"
	syncTx "github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/relayer/chain/transceiver/sync"
	ebTypes "github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/types"
	"github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/utils"
	"github.com/assetcloud/plugin/plugin/dapp/x2ethereum/types"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

var relayerLog = log.New("module", "chain_relayer")

//Relayer4Chain ...
type Relayer4Chain struct {
	syncTxReceipts      *syncTx.TxReceipts
	ethClient           ethinterface.EthClientSpec
	rpcLaddr            string //用户向指定的blockchain节点进行rpc调用
	fetchHeightPeriodMs int64
	db                  dbm.DB
	lastHeight4Tx       int64 //等待被处理的具有相应的交易回执的高度
	matDegree           int32 //成熟度         heightSync2App    matDegress   height
	//passphase            string
	privateKey4Ethereum  *ecdsa.PrivateKey
	ethSender            ethCommon.Address
	bridgeRegistryAddr   ethCommon.Address
	oracleInstance       *generated.Oracle
	totalTx4ChainToEth int64
	statusCheckedIndex   int64
	ctx                  context.Context
	rwLock               sync.RWMutex
	unlock               chan int
}

// StartChainRelayer : initializes a relayer which witnesses events on the chain network and relays them to Ethereum
func StartChainRelayer(ctx context.Context, syncTxConfig *ebTypes.SyncTxConfig, registryAddr, provider string, db dbm.DB) *Relayer4Chain {
	chian33Relayer := &Relayer4Chain{
		rpcLaddr:            syncTxConfig.ChainHost,
		fetchHeightPeriodMs: syncTxConfig.FetchHeightPeriodMs,
		unlock:              make(chan int),
		db:                  db,
		ctx:                 ctx,
		bridgeRegistryAddr:  ethCommon.HexToAddress(registryAddr),
	}

	syncCfg := &ebTypes.SyncTxReceiptConfig{
		ChainHost:       syncTxConfig.ChainHost,
		PushHost:          syncTxConfig.PushHost,
		PushName:          syncTxConfig.PushName,
		PushBind:          syncTxConfig.PushBind,
		StartSyncHeight:   syncTxConfig.StartSyncHeight,
		StartSyncSequence: syncTxConfig.StartSyncSequence,
		StartSyncHash:     syncTxConfig.StartSyncHash,
	}

	client, err := relayerTx.SetupWebsocketEthClient(provider)
	if err != nil {
		panic(err)
	}
	chian33Relayer.ethClient = client
	chian33Relayer.totalTx4ChainToEth = chian33Relayer.getTotalTxAmount2Eth()
	chian33Relayer.statusCheckedIndex = chian33Relayer.getStatusCheckedIndex()

	go chian33Relayer.syncProc(syncCfg)
	return chian33Relayer
}

//QueryTxhashRelay2Eth ...
func (chainRelayer *Relayer4Chain) QueryTxhashRelay2Eth() ebTypes.Txhashes {
	txhashs := utils.QueryTxhashes([]byte(chainToEthBurnLockTxHashPrefix), chainRelayer.db)
	return ebTypes.Txhashes{Txhash: txhashs}
}

func (chainRelayer *Relayer4Chain) syncProc(syncCfg *ebTypes.SyncTxReceiptConfig) {
	_, _ = fmt.Fprintln(os.Stdout, "Pls unlock or import private key for Chain relayer")
	<-chainRelayer.unlock
	_, _ = fmt.Fprintln(os.Stdout, "Chain relayer starts to run...")

	chainRelayer.syncTxReceipts = syncTx.StartSyncTxReceipt(syncCfg, chainRelayer.db)
	chainRelayer.lastHeight4Tx = chainRelayer.loadLastSyncHeight()

	oracleInstance, err := relayerTx.RecoverOracleInstance(chainRelayer.ethClient, chainRelayer.bridgeRegistryAddr, chainRelayer.bridgeRegistryAddr)
	if err != nil {
		panic(err.Error())
	}
	chainRelayer.oracleInstance = oracleInstance

	timer := time.NewTicker(time.Duration(chainRelayer.fetchHeightPeriodMs) * time.Millisecond)
	for {
		select {
		case <-timer.C:
			height := chainRelayer.getCurrentHeight()
			relayerLog.Debug("syncProc", "getCurrentHeight", height)
			chainRelayer.onNewHeightProc(height)

		case <-chainRelayer.ctx.Done():
			timer.Stop()
			return
		}
	}
}

func (chainRelayer *Relayer4Chain) getCurrentHeight() int64 {
	var res rpctypes.Header
	ctx := jsonclient.NewRPCCtx(chainRelayer.rpcLaddr, "Chain.GetLastHeader", nil, &res)
	_, err := ctx.RunResult()
	if nil != err {
		relayerLog.Error("getCurrentHeight", "Failede due to:", err.Error())
	}
	return res.Height
}

func (chainRelayer *Relayer4Chain) onNewHeightProc(currentHeight int64) {
	//检查已经提交的交易结果
	chainRelayer.rwLock.Lock()
	for chainRelayer.statusCheckedIndex < chainRelayer.totalTx4ChainToEth {
		index := chainRelayer.statusCheckedIndex + 1
		txhash, err := chainRelayer.getEthTxhash(index)
		if nil != err {
			relayerLog.Error("onNewHeightProc", "getEthTxhash for index ", index, "error", err.Error())
			break
		}
		status := relayerTx.GetEthTxStatus(chainRelayer.ethClient, txhash)
		//按照提交交易的先后顺序检查交易，只要出现当前交易还在pending状态，就不再检查后续交易，等到下个区块再从该交易进行检查
		//TODO:可能会由于网络和打包挖矿的原因，使得交易执行顺序和提交顺序有差别，后续完善该检查逻辑
		if status == relayerTx.EthTxPending.String() {
			break
		}
		_ = chainRelayer.setLastestRelay2EthTxhash(status, txhash.Hex(), index)
		atomic.AddInt64(&chainRelayer.statusCheckedIndex, 1)
		_ = chainRelayer.setStatusCheckedIndex(chainRelayer.statusCheckedIndex)
	}
	chainRelayer.rwLock.Unlock()
	//未达到足够的成熟度，不进行处理
	//  +++++++++||++++++++++++||++++++++++||
	//           ^             ^           ^
	// lastHeight4Tx    matDegress   currentHeight
	for chainRelayer.lastHeight4Tx+int64(chainRelayer.matDegree)+1 <= currentHeight {
		relayerLog.Info("onNewHeightProc", "currHeight", currentHeight, "lastHeight4Tx", chainRelayer.lastHeight4Tx)

		lastHeight4Tx := chainRelayer.lastHeight4Tx
		TxReceipts, err := chainRelayer.syncTxReceipts.GetNextValidTxReceipts(lastHeight4Tx)
		if nil == TxReceipts || nil != err {
			if err != nil {
				relayerLog.Error("onNewHeightProc", "Failed to GetNextValidTxReceipts due to:", err.Error())
			}
			break
		}
		relayerLog.Debug("onNewHeightProc", "currHeight", currentHeight, "valid tx receipt with height:", TxReceipts.Height)

		txs := TxReceipts.Tx
		for i, tx := range txs {
			//检查是否为lns的交易(包括平行链：user.p.xxx.lns)，将闪电网络交易进行收集
			if 0 != bytes.Compare(tx.Execer, []byte(relayerTx.X2Eth)) &&
				(len(tx.Execer) > 4 && string(tx.Execer[(len(tx.Execer)-4):]) != "."+relayerTx.X2Eth) {
				relayerLog.Debug("onNewHeightProc, the tx is not x2ethereum", "Execer", string(tx.Execer), "height:", TxReceipts.Height)
				continue
			}
			var ss types.X2EthereumAction
			_ = chainTypes.Decode(tx.Payload, &ss)
			actionName := ss.GetActionName()
			if relayerTx.BurnAction == actionName || relayerTx.LockAction == actionName {
				relayerLog.Debug("^_^ ^_^ Processing chain tx receipt", "ActionName", actionName, "fromAddr", tx.From(), "exec", string(tx.Execer))
				actionEvent := getOracleClaimType(actionName)
				if err := chainRelayer.handleBurnLockMsg(actionEvent, TxReceipts.ReceiptData[i], tx.Hash()); nil != err {
					errInfo := fmt.Sprintf("Failed to handleBurnLockMsg due to:%s", err.Error())
					panic(errInfo)
				}
			}
		}
		chainRelayer.lastHeight4Tx = TxReceipts.Height
		chainRelayer.setLastSyncHeight(chainRelayer.lastHeight4Tx)
	}
}

// getOracleClaimType : sets the OracleClaim's claim type based upon the witnessed event type
func getOracleClaimType(eventType string) events.Event {
	var claimType events.Event

	switch eventType {
	case events.MsgBurn.String():
		claimType = events.Event(events.ClaimTypeBurn)
	case events.MsgLock.String():
		claimType = events.Event(events.ClaimTypeLock)
	default:
		panic(errors.New("eventType invalid"))
	}

	return claimType
}

// handleBurnLockMsg : parse event data as a ChainMsg, package it into a ProphecyClaim, then relay tx to the Ethereum Network
func (chainRelayer *Relayer4Chain) handleBurnLockMsg(claimEvent events.Event, receipt *chainTypes.ReceiptData, chainTxHash []byte) error {
	relayerLog.Info("handleBurnLockMsg", "Received tx with hash", ethCommon.Bytes2Hex(chainTxHash))

	// Parse the witnessed event's data into a new ChainMsg
	chainMsg := relayerTx.ParseBurnLockTxReceipt(claimEvent, receipt)
	if nil == chainMsg {
		//收到执行失败的交易，直接跳过
		relayerLog.Error("handleBurnLockMsg", "Received failed tx with hash", ethCommon.Bytes2Hex(chainTxHash))
		return nil
	}

	// Parse the ChainMsg into a ProphecyClaim for relay to Ethereum
	prophecyClaim := relayerTx.ChainMsgToProphecyClaim(*chainMsg)

	// Relay the ChainMsg to the Ethereum network
	txhash, err := relayerTx.RelayOracleClaimToEthereum(chainRelayer.oracleInstance, chainRelayer.ethClient, chainRelayer.ethSender, claimEvent, prophecyClaim, chainRelayer.privateKey4Ethereum, chainTxHash)
	if nil != err {
		return err
	}

	//保存交易hash，方便查询
	atomic.AddInt64(&chainRelayer.totalTx4ChainToEth, 1)
	txIndex := atomic.LoadInt64(&chainRelayer.totalTx4ChainToEth)
	if err = chainRelayer.updateTotalTxAmount2Eth(txIndex); nil != err {
		relayerLog.Error("handleLogNewProphecyClaimEvent", "Failed to RelayLockToChain due to:", err.Error())
		return err
	}
	if err = chainRelayer.setLastestRelay2EthTxhash(relayerTx.EthTxPending.String(), txhash, txIndex); nil != err {
		relayerLog.Error("handleLogNewProphecyClaimEvent", "Failed to RelayLockToChain due to:", err.Error())
		return err
	}
	return nil
}
