package chain

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	chainEvmCommon "github.com/assetcloud/plugin/plugin/dapp/evm/executor/vm/common"

	evmtypes "github.com/assetcloud/plugin/plugin/dapp/evm/types"

	"github.com/assetcloud/chain/common"
	chainCrypto "github.com/assetcloud/chain/common/crypto"
	dbm "github.com/assetcloud/chain/common/db"
	log "github.com/assetcloud/chain/common/log/log15"
	"github.com/assetcloud/chain/rpc/jsonclient"
	rpctypes "github.com/assetcloud/chain/rpc/types"
	chainTypes "github.com/assetcloud/chain/types"
	syncTx "github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/relayer/chain/transceiver/sync"
	"github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/relayer/events"
	ebTypes "github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/types"
	"github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
)

var relayerLog = log.New("module", "chain_relayer")

//Relayer4Chain ...
type Relayer4Chain struct {
	syncEvmTxLogs       *syncTx.EVMTxLogs
	rpcLaddr            string //用户向指定的blockchain节点进行rpc调用
	chainRpcUrls      []string
	chainName           string //用来区别主链中继还是平行链，主链为空，平行链则是user.p.xxx.
	chainID             int32
	fetchHeightPeriodMs int64
	db                  dbm.DB
	lastHeight4Tx       int64 //等待被处理的具有相应的交易回执的高度
	matDegree           int32 //成熟度         heightSync2App    matDegress   height

	privateKey4Chain         chainCrypto.PrivKey
	privateKey4Chain_ecdsa   *ecdsa.PrivateKey
	ctx                        context.Context
	rwLock                     sync.RWMutex
	unlockChan                 chan int
	bridgeBankEventLockSig     string
	bridgeBankEventBurnSig     string
	bridgeBankEventWithdrawSig string
	bridgeBankAbi              abi.ABI
	totalTx4RelayEth2chai33    int64
	//新增//
	ethBridgeClaimChan        <-chan *ebTypes.EthBridgeClaim
	txRelayAckRecvChan        <-chan *ebTypes.TxRelayAck
	txRelayAckSendChan        map[string]chan<- *ebTypes.TxRelayAck
	chainMsgChan            map[string]chan<- *events.ChainMsg
	bridgeRegistryAddr        string
	oracleAddr                string
	bridgeBankAddr            string
	mulSignAddr               string
	deployResult              *X2EthDeployResult
	symbol2Addr               map[string]string
	bridgeSymbol2EthChainName map[string]string //在chain上发行的跨链token的名称到以太坊链的名称映射
	processWithDraw           bool
	delayedSend               bool
	delayedSendTime           int64
}

type ChainStartPara struct {
	ChainName          string
	Ctx                context.Context
	SyncTxConfig       *ebTypes.SyncTxConfig
	BridgeRegistryAddr string
	DBHandle           dbm.DB
	EthBridgeClaimChan <-chan *ebTypes.EthBridgeClaim
	TxRelayAckRecvChan <-chan *ebTypes.TxRelayAck
	TxRelayAckSendChan map[string]chan<- *ebTypes.TxRelayAck
	ChainMsgChan     map[string]chan<- *events.ChainMsg
	ChainID            int32
	ProcessWithDraw    bool
	DelayedSend        bool
	DelayedSendTime    int64
}

// StartChainRelayer : initializes a relayer which witnesses events on the chain network and relays them to Ethereum
func StartChainRelayer(startPara *ChainStartPara) *Relayer4Chain {
	chainRelayer := &Relayer4Chain{
		rpcLaddr:                startPara.SyncTxConfig.ChainHost,
		chainRpcUrls:          startPara.SyncTxConfig.ChainRpcUrls,
		chainName:               startPara.ChainName,
		chainID:                 startPara.ChainID,
		fetchHeightPeriodMs:     startPara.SyncTxConfig.FetchHeightPeriodMs,
		unlockChan:              make(chan int),
		db:                      startPara.DBHandle,
		ctx:                     startPara.Ctx,
		bridgeRegistryAddr:      startPara.BridgeRegistryAddr,
		ethBridgeClaimChan:      startPara.EthBridgeClaimChan,
		txRelayAckRecvChan:      startPara.TxRelayAckRecvChan,
		txRelayAckSendChan:      startPara.TxRelayAckSendChan,
		chainMsgChan:          startPara.ChainMsgChan,
		totalTx4RelayEth2chai33: 0,
		symbol2Addr:             make(map[string]string),
		processWithDraw:         startPara.ProcessWithDraw,
		delayedSend:             startPara.DelayedSend,
		delayedSendTime:         startPara.DelayedSendTime,
	}

	syncCfg := &ebTypes.SyncTxReceiptConfig{
		ChainHost:       startPara.SyncTxConfig.ChainHost,
		PushHost:          startPara.SyncTxConfig.PushHost,
		PushName:          startPara.SyncTxConfig.PushName,
		PushBind:          startPara.SyncTxConfig.PushBind,
		StartSyncHeight:   startPara.SyncTxConfig.StartSyncHeight,
		StartSyncSequence: startPara.SyncTxConfig.StartSyncSequence,
		StartSyncHash:     startPara.SyncTxConfig.StartSyncHash,
		KeepAliveDuration: startPara.SyncTxConfig.KeepAliveDuration,
	}

	registrAddrInDB, err := chainRelayer.getBridgeRegistryAddr()
	//如果输入的registry地址非空，且和数据库保存地址不一致，则直接使用输入注册地址
	if chainRelayer.bridgeRegistryAddr != "" && nil == err && registrAddrInDB != chainRelayer.bridgeRegistryAddr {
		relayerLog.Error("StartChainRelayer", "BridgeRegistry is setted already with value", registrAddrInDB,
			"but now setting to", startPara.BridgeRegistryAddr)
		_ = chainRelayer.setBridgeRegistryAddr(startPara.BridgeRegistryAddr)
	} else if startPara.BridgeRegistryAddr == "" && registrAddrInDB != "" {
		//输入地址为空，且数据库中保存地址不为空，则直接使用数据库中的地址
		chainRelayer.bridgeRegistryAddr = registrAddrInDB
	}
	chainRelayer.totalTx4RelayEth2chai33 = chainRelayer.getTotalTxAmount()
	if 0 == chainRelayer.totalTx4RelayEth2chai33 {
		statics := &ebTypes.Ethereum2ChainStatics{}
		data := chainTypes.Encode(statics)
		err := chainRelayer.setLastestRelay2ChainTxStatics(0, int32(events.ClaimTypeLock), data)
		if err != nil {
			relayerLog.Error("StartChainRelayer", "setLastestRelay2ChainTxStatics ClaimTypeLock error", err.Error())
		}
		err = chainRelayer.setLastestRelay2ChainTxStatics(0, int32(events.ClaimTypeBurn), data)
		if err != nil {
			relayerLog.Error("StartChainRelayer", "setLastestRelay2ChainTxStatics ClaimTypeBurn error", err.Error())
		}
	}

	go chainRelayer.syncProc(syncCfg)
	return chainRelayer
}

func (chainRelayer *Relayer4Chain) syncProc(syncCfg *ebTypes.SyncTxReceiptConfig) {
	_, _ = fmt.Fprintln(os.Stdout, "Pls unlock or import private key for Chain relayer")
	<-chainRelayer.unlockChan
	_, _ = fmt.Fprintln(os.Stdout, "Chain relayer starts to run...")
	if err := chainRelayer.RestoreTokenAddress(); nil != err {
		relayerLog.Info("Failed to RestoreTokenAddress")
		return
	}
	setChainID(chainRelayer.chainID)
	//如果该中继器的bridgeRegistryAddr为空，就说明合约未部署，需要等待部署成功之后再继续
	if "" == chainRelayer.bridgeRegistryAddr {
		chaintxLog.Debug("bridgeRegistryAddr empty")
		<-chainRelayer.unlockChan
	}
	//如果oracleAddr为空，则通过bridgeRegistry合约进行查询
	if "" != chainRelayer.bridgeRegistryAddr && "" == chainRelayer.oracleAddr {
		oracleAddr, bridgeBankAddr := recoverContractAddrFromRegistry(chainRelayer.bridgeRegistryAddr, chainRelayer.rpcLaddr)
		if "" == oracleAddr || "" == bridgeBankAddr {
			panic("Failed to recoverContractAddrFromRegistry")
		}
		chainRelayer.oracleAddr = oracleAddr
		chainRelayer.bridgeBankAddr = bridgeBankAddr
		chaintxLog.Debug("recoverContractAddrFromRegistry", "bridgeRegistryAddr", chainRelayer.bridgeRegistryAddr,
			"oracleAddr", chainRelayer.oracleAddr, "bridgeBankAddr", chainRelayer.bridgeBankAddr)
	}

	syncCfg.Contracts = append(syncCfg.Contracts, chainRelayer.bridgeBankAddr)
	chainRelayer.syncEvmTxLogs = syncTx.StartSyncEvmTxLogs(syncCfg, chainRelayer.db)
	chainRelayer.lastHeight4Tx = chainRelayer.loadLastSyncHeight()
	chainRelayer.mulSignAddr = chainRelayer.getMultiSignAddress()
	chainRelayer.bridgeSymbol2EthChainName = chainRelayer.restoreSymbol2chainName()
	chainRelayer.prePareSubscribeEvent()
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

		case ethBridgeClaim := <-chainRelayer.ethBridgeClaimChan:
			chainRelayer.relayLockBurnToChain(ethBridgeClaim)

		case txRelayAck := <-chainRelayer.txRelayAckRecvChan:
			chainRelayer.procTxRelayAck(txRelayAck)
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
	chainRelayer.updateTxStatus()
	chainRelayer.checkTxRelay2Ethereum()

	//未达到足够的成熟度，不进行处理
	//  +++++++++||++++++++++++||++++++++++||
	//           ^             ^           ^
	// lastHeight4Tx    matDegress   currentHeight
	for chainRelayer.lastHeight4Tx+int64(chainRelayer.matDegree)+1 <= currentHeight {
		relayerLog.Info("onNewHeightProc", "currHeight", currentHeight, "lastHeight4Tx", chainRelayer.lastHeight4Tx)

		lastHeight4Tx := chainRelayer.lastHeight4Tx
		txLogs, err := chainRelayer.syncEvmTxLogs.GetNextValidEvmTxLogs(lastHeight4Tx)
		if nil == txLogs || nil != err {
			if err != nil {
				relayerLog.Error("onNewHeightProc", "Failed to GetNextValidTxReceipts due to:", err.Error())
			}
			break
		}
		relayerLog.Debug("onNewHeightProc", "currHeight", currentHeight, "valid tx receipt with height:", txLogs.Height)

		txAndLogs := txLogs.TxAndLogs
		for _, txAndLog := range txAndLogs {
			tx := txAndLog.Tx

			//确认订阅的evm交易类型和合约地址
			if !strings.Contains(string(tx.Execer), "evm") {
				relayerLog.Error("onNewHeightProc received logs not from evm tx", "tx.Execer", string(tx.Execer))
				continue
			}

			var evmAction evmtypes.EVMContractAction
			err := chainTypes.Decode(tx.Payload, &evmAction)
			if nil != err {
				relayerLog.Error("onNewHeightProc", "Failed to decode action for tx with hash", common.ToHex(tx.Hash()))
				continue
			}

			//确认监听的合约地址
			if evmAction.ContractAddr != chainRelayer.bridgeBankAddr {
				relayerLog.Error("onNewHeightProc received logs not from bridgeBank", "evmAction.ContractAddr", evmAction.ContractAddr)
				continue
			}

			for _, evmlog := range txAndLog.LogsPerTx.Logs {
				var evmEventType events.ChainEvmEvent
				if chainRelayer.bridgeBankEventBurnSig == common.ToHex(evmlog.Topic[0]) {
					evmEventType = events.ChainEventLogBurn
				} else if chainRelayer.bridgeBankEventLockSig == common.ToHex(evmlog.Topic[0]) {
					evmEventType = events.ChainEventLogLock
				} else if chainRelayer.bridgeBankEventWithdrawSig == common.ToHex(evmlog.Topic[0]) {
					evmEventType = events.ChainEventLogWithdraw
				} else {
					continue
				}

				if evmEventType == events.ChainEventLogWithdraw && !chainRelayer.processWithDraw {
					//代理提币消息只由代理提币节点处理
					continue
				}
				if evmEventType != events.ChainEventLogWithdraw && chainRelayer.processWithDraw {
					//lock和burn消息消息只由普通中继节点处理
					continue
				}

				if err := chainRelayer.handleBurnLockWithdrawEvent(evmEventType, evmlog.Data, tx.Hash()); nil != err {
					relayerLog.Error("onNewHeightProc", "Failed to handleBurnLockWithdrawEvent due to:%s", err.Error())
				}

			}
		}
		chainRelayer.lastHeight4Tx = txLogs.Height
		chainRelayer.setLastSyncHeight(chainRelayer.lastHeight4Tx)
	}
}

// handleBurnLockMsg : parse event data as a ChainMsg, package it into a ProphecyClaim, then relay tx to the Ethereum Network
func (chainRelayer *Relayer4Chain) handleBurnLockWithdrawEvent(evmEventType events.ChainEvmEvent, data []byte, chainTxHash []byte) error {
	txHashStr := common.ToHex(chainTxHash)
	relayerLog.Info("handleBurnLockWithdrawEvent", "Received tx with hash", txHashStr)

	// 删除已发送校验, 如果ethereum端发生交易后没有打包, 可重新再发生
	//if chainRelayer.checkTxProcessed(txHashStr) {
	//	relayerLog.Info("handleBurnLockWithdrawEvent", "Tx has been already Processed with hash:", txHashStr)
	//	return nil
	//}

	// Parse the witnessed event's data into a new ChainMsg
	chainMsg, err := events.ParseBurnLock4chain(evmEventType, data, chainRelayer.bridgeBankAbi, chainTxHash)
	if nil != err {
		return err
	}
	fdIndex := chainRelayer.getFdTx2EthTotalAmount() + 1
	chainMsg.ForwardTimes = 1
	chainMsg.ForwardIndex = fdIndex

	relayerLog.Info("handleBurnLockWithdrawEvent", "Going to send chainMsg.ClaimType", chainMsg.ClaimType.String())

	var chainName string
	//specical process: withdraw YCC　only to bsc
	if events.ChainEventLogWithdraw == evmEventType && "YCC" == chainMsg.Symbol {
		chainName = ebTypes.BinanceChainName
	} else {
		ok := false
		chainName, ok = chainRelayer.bridgeSymbol2EthChainName[chainMsg.Symbol]
		if !ok {
			relayerLog.Error("handleBurnLockWithdrawEvent", "No bridgeSymbol2EthChainName", chainMsg.Symbol)
			return errors.New("ErrNoEthChainName4BridgeSymbol")
		}
	}

	channel, ok := chainRelayer.chainMsgChan[chainName]
	if !ok {
		relayerLog.Error("handleBurnLockWithdrawEvent", "No bridgeSymbol2EthChainName", chainName)
		return errors.New("ErrNoChainMsgChan4EthChainName")
	}

	_ = chainRelayer.updateFdTx2EthTotalAmount(fdIndex)
	txRelayConfirm4Chain := &ebTypes.TxRelayConfirm4Chain{
		EventType:   int32(evmEventType),
		Data:        data,
		FdTimes:     1,
		FdIndex:     fdIndex,
		ToChainName: chainName,
		TxHash:      chainTxHash,
		Resend:      false,
	}

	if chainRelayer.delayedSend {
		go chainRelayer.delayedSendTxs(chainName, chainMsg, chainTxHash, txRelayConfirm4Chain)
	} else {
		channel <- chainMsg
		//relaychainToEthereumCheckPonit 1:send chainMsg to ethereum relay service
		relayerLog.Info("handleBurnLockWithdrawEvent::relaychainToEthereumCheckPonit_1", "chainTxHash", txHashStr, "ForwardIndex", chainMsg.ForwardIndex, "FdTimes", 1)
		err = chainRelayer.setChainTxIsRelayedUnconfirm(txHashStr, fdIndex, txRelayConfirm4Chain)
	}

	return err
}

func (chainRelayer *Relayer4Chain) delayedSendTxs(chainName string, chainMsg *events.ChainMsg, chainTxHash []byte, txRelayConfirm4Chain *ebTypes.TxRelayConfirm4Chain) {
	delayedSendTime := time.Duration(chainRelayer.delayedSendTime) * time.Millisecond
	relayerLog.Debug("delayedSendTxs", "setEthTxWaitingForSend chainTxHash", common.ToHex(chainTxHash))
	time.Sleep(delayedSendTime)
	channel, ok := chainRelayer.chainMsgChan[chainName]
	if !ok {
		relayerLog.Error("handleBurnLockWithdrawEvent", "No bridgeSymbol2EthChainName", chainName)
		return
	}

	channel <- chainMsg

	//relaychainToEthereumCheckPonit 1:send chainMsg to ethereum relay service
	relayerLog.Info("handleBurnLockWithdrawEvent::relaychainToEthereumCheckPonit_1", "chainTxHash", common.ToHex(chainTxHash), "ForwardIndex", chainMsg.ForwardIndex, "FdTimes", 1)
	_ = chainRelayer.setChainTxIsRelayedUnconfirm(common.ToHex(chainTxHash), txRelayConfirm4Chain.FdIndex, txRelayConfirm4Chain)
}

func (chainRelayer *Relayer4Chain) ResendChainEvent(height int64) (err error) {
	txLogs, err := chainRelayer.syncEvmTxLogs.GetNextValidEvmTxLogs(height)
	if nil == txLogs || nil != err {
		if err != nil {
			relayerLog.Error("ResendChainEvent", "Failed to GetNextValidTxReceipts due to:", err.Error())
			return err
		}
		return nil
	}
	relayerLog.Debug("ResendChainEvent", "lastHeight4Tx", chainRelayer.lastHeight4Tx, "valid tx receipt with height:", txLogs.Height)

	txAndLogs := txLogs.TxAndLogs
	for _, txAndLog := range txAndLogs {
		tx := txAndLog.Tx

		//确认订阅的evm交易类型和合约地址
		if !strings.Contains(string(tx.Execer), "evm") {
			relayerLog.Error("ResendChainEvent received logs not from evm tx", "tx.Execer", string(tx.Execer))
			continue
		}

		var evmAction evmtypes.EVMContractAction
		err := chainTypes.Decode(tx.Payload, &evmAction)
		if nil != err {
			relayerLog.Error("ResendChainEvent", "Failed to decode action for tx with hash", common.ToHex(tx.Hash()))
			continue
		}

		//确认监听的合约地址
		if evmAction.ContractAddr != chainRelayer.bridgeBankAddr {
			relayerLog.Error("ResendChainEvent received logs not from bridgeBank", "evmAction.ContractAddr", evmAction.ContractAddr)
			continue
		}

		for _, evmlog := range txAndLog.LogsPerTx.Logs {
			var evmEventType events.ChainEvmEvent
			if chainRelayer.bridgeBankEventBurnSig == common.ToHex(evmlog.Topic[0]) {
				evmEventType = events.ChainEventLogBurn
			} else if chainRelayer.bridgeBankEventLockSig == common.ToHex(evmlog.Topic[0]) {
				evmEventType = events.ChainEventLogLock
			} else if chainRelayer.bridgeBankEventWithdrawSig == common.ToHex(evmlog.Topic[0]) {
				evmEventType = events.ChainEventLogWithdraw
			} else {
				continue
			}

			if evmEventType == events.ChainEventLogWithdraw && !chainRelayer.processWithDraw {
				//代理提币消息只由代理提币节点处理
				continue
			}
			if evmEventType != events.ChainEventLogWithdraw && chainRelayer.processWithDraw {
				//lock和burn消息消息只由普通中继节点处理
				continue
			}

			if err := chainRelayer.handleBurnLockWithdrawEvent(evmEventType, evmlog.Data, tx.Hash()); nil != err {
				return err
			}
		}
	}

	return nil
}

func (chainRelayer *Relayer4Chain) checkIsResendEthClaim(claim *ebTypes.EthBridgeClaim) bool {
	if claim.ForwardTimes <= 1 {
		return false
	}
	ethTxHash := claim.EthTxHash
	relayerLog.Info("checkIsResendEthClaim", "Received the same EthBridgeClaim more than once with times", claim.ForwardTimes, "tx hash string", ethTxHash)
	relayTxDetail, _ := chainRelayer.getEthTxRelayAlreadyInfo(ethTxHash)
	if nil == relayTxDetail {
		relayerLog.Info("checkIsResendEthClaim::haven't relay yet")
		return false
	}

	//if relay already, just ack it
	chainRelayer.txRelayAckSendChan[claim.ChainName] <- &ebTypes.TxRelayAck{
		TxHash:  ethTxHash,
		FdIndex: claim.ForwardIndex,
	}
	relayerLog.Info("checkIsResendEthClaim", "have relay already with tx hash:", relayTxDetail.Txhash)
	return true
}

func (chainRelayer *Relayer4Chain) relayLockBurnToChain(claim *ebTypes.EthBridgeClaim) {
	relayerLog.Debug("relayLockBurnToChain", "new EthBridgeClaim received", claim)
	if chainRelayer.checkIsResendEthClaim(claim) {
		return
	}

	nonceBytes := big.NewInt(claim.Nonce).Bytes()
	bigAmount := big.NewInt(0)
	bigAmount.SetString(claim.Amount, 10)
	amountBytes := bigAmount.Bytes()
	claimID := crypto.Keccak256Hash(nonceBytes, []byte(claim.EthereumSender), []byte(claim.ChainReceiver), []byte(claim.Symbol), amountBytes)

	// Sign the hash using the active validator's private key
	signature, err := utils.SignClaim4Evm(claimID, chainRelayer.privateKey4Chain_ecdsa)
	if nil != err {
		panic("SignClaim4Evm due to" + err.Error())
	}

	var tokenAddr string
	operationType := events.ClaimType(claim.ClaimType).String()
	if int32(events.ClaimTypeBurn) == claim.ClaimType {
		//burn 分支
		if ebTypes.SYMBOL_BTY == claim.Symbol {
			tokenAddr = ebTypes.BTYAddrChain
		} else {
			tokenAddr = getLockedTokenAddress(chainRelayer.bridgeBankAddr, claim.Symbol, chainRelayer.rpcLaddr)
			if "" == tokenAddr {
				relayerLog.Error("relayLockBurnToChain", "No locked token address created for symbol", claim.Symbol)
				return
			}
		}
	} else {
		//lock 分支
		if _, ok := chainRelayer.bridgeSymbol2EthChainName[claim.Symbol]; !ok {
			chainRelayer.bridgeSymbol2EthChainName[claim.Symbol] = claim.ChainName
			chainRelayer.storeSymbol2chainName(chainRelayer.bridgeSymbol2EthChainName)
		}
		//如果是代理打币节点，则只收集symbol和chain name相关信息
		if chainRelayer.processWithDraw {
			return
		}

		var exist bool
		tokenAddr, exist = chainRelayer.symbol2Addr[claim.Symbol]
		if !exist {
			tokenAddr = getBridgeToken2address(chainRelayer.bridgeBankAddr, claim.Symbol, chainRelayer.rpcLaddr)
			if "" == tokenAddr {
				relayerLog.Error("relayLockBurnToChain", "No bridge token address created for symbol", claim.Symbol)
				return
			}
			relayerLog.Info("relayLockBurnToChain", "Succeed to get bridge token address for symbol", claim.Symbol,
				"address", tokenAddr)

			token2set := &ebTypes.TokenAddress{
				Address:   tokenAddr,
				Symbol:    claim.Symbol,
				ChainName: ebTypes.ChainBlockChainName,
			}
			if err := chainRelayer.SetTokenAddress(token2set); nil != err {
				relayerLog.Info("relayLockBurnToChain", "Failed to SetTokenAddress due to", err.Error())
			}
		}
	}

	//因为发行的合约的精度为8，所以需要进行相应的缩放
	if 8 != claim.Decimal {
		if claim.Decimal > 8 {
			dist := claim.Decimal - 8
			value, exist := utils.Decimal2value[int(dist)]
			if !exist {
				panic(fmt.Sprintf("does support for decimal, %d", claim.Decimal))
			}
			bigAmount.Div(bigAmount, big.NewInt(value))
			claim.Amount = bigAmount.String()
		} else {
			dist := 8 - claim.Decimal
			value, exist := utils.Decimal2value[int(dist)]
			if !exist {
				panic(fmt.Sprintf("does support for decimal, %d", claim.Decimal))
			}
			bigAmount.Mul(bigAmount, big.NewInt(value))
			claim.Amount = bigAmount.String()
		}
	}

	parameter := fmt.Sprintf("newOracleClaim(%d, %s, %s, %s, %s, %s, %s, %s)",
		claim.ClaimType,
		claim.EthereumSender,
		claim.ChainReceiver,
		tokenAddr,
		claim.Symbol,
		claim.Amount,
		claimID.String(),
		common.ToHex(signature))
	relayerLog.Info("relayLockBurnToChain", "parameter", parameter)

	txhash, err := relayEvmTx2Chain(chainRelayer.privateKey4Chain, claim, parameter, chainRelayer.oracleAddr, chainRelayer.chainName, chainRelayer.chainRpcUrls)
	if err != nil {
		relayerLog.Error("relayLockBurnToChain", "Failed to RelayEvmTx2Chain due to:", err.Error(), "EthereumTxhash", claim.EthTxHash)
		return
	}

	chainRelayer.txRelayAckSendChan[claim.ChainName] <- &ebTypes.TxRelayAck{
		TxHash:  claim.EthTxHash,
		FdIndex: claim.ForwardIndex,
	}
	//relayEthereum2chainCheckPonit 2:send ack
	relayerLog.Info("relayLockBurnToChain::relayEthereum2chainCheckPonit_2::sendAck", "ethTxhash", claim.EthTxHash, "ForwardIndex", claim.ForwardIndex, "FdTimes", claim.ForwardTimes)

	relayTxDetail := &ebTypes.RelayTxDetail{
		ClaimType:      claim.ClaimType,
		TxIndexRelayed: claim.ForwardIndex,
		Txhash:         txhash,
	}

	//set flag to indicate that the eth tx has been relayed to chain
	if err = chainRelayer.setEthTxRelayAlreadyInfo(claim.EthTxHash, relayTxDetail); nil != err {
		relayerLog.Error("relayLockBurnToChain", "Failed to setTxRelayAlreadyInfo due to:", err.Error())
		return
	}
	//relayEthereum2chainCheckPonit 3:setFalgRelayFinish
	relayerLog.Info("relayLockBurnToChain::relayEthereum2chainCheckPonit_3::setFalgRelayFinish", "ethTxhash", claim.EthTxHash, "ForwardIndex", claim.ForwardIndex, "FdTimes", claim.ForwardTimes)

	//第一个有效的index从１开始，方便list
	txIndex := atomic.AddInt64(&chainRelayer.totalTx4RelayEth2chai33, 1)
	if err = chainRelayer.updateTotalTxAmount2Eth(txIndex); nil != err {
		relayerLog.Error("relayLockBurnToChain", "Failed to updateTotalTxAmount2Eth due to:", err.Error())
		return
	}

	statics := &ebTypes.Ethereum2ChainStatics{
		ChainTxstatus: ebTypes.Tx_Status_Pending,
		ChainTxhash:   txhash,
		EthereumTxhash:  claim.EthTxHash,
		BurnLock:        claim.ClaimType,
		EthereumSender:  claim.EthereumSender,
		ChainReceiver: claim.ChainReceiver,
		Symbol:          claim.Symbol,
		Amount:          claim.Amount,
		Nonce:           claim.Nonce,
		TxIndex:         txIndex,
		OperationType:   operationType,
	}
	data := chainTypes.Encode(statics)
	if err = chainRelayer.setLastestRelay2ChainTxStatics(txIndex, claim.ClaimType, data); nil != err {
		relayerLog.Error("relayLockBurnToChain", "Failed to setLastestRelay2ChainTxStatics due to:", err.Error())
		return
	}
	relayerLog.Info("relayLockBurnToChain::successful",
		"txIndex", txIndex,
		"ChainTxhash", txhash,
		"EthereumTxhash", claim.EthTxHash,
		"type", operationType,
		"Symbol", claim.Symbol,
		"Amount", claim.Amount,
		"EthereumSender", claim.EthereumSender,
		"ChainReceiver", claim.ChainReceiver)
}

func (chainRelayer *Relayer4Chain) BurnAsyncFromChain(ownerPrivateKey, tokenAddr, ethereumReceiver, amount string) (string, error) {
	bn := big.NewInt(1)
	bn, _ = bn.SetString(utils.TrimZeroAndDot(amount), 10)
	return burnAsync(ownerPrivateKey, tokenAddr, ethereumReceiver, bn.Int64(), chainRelayer.bridgeBankAddr, chainRelayer.chainName, chainRelayer.rpcLaddr)
}

func (chainRelayer *Relayer4Chain) LockBTYAssetAsync(ownerPrivateKey, ethereumReceiver, amount string) (string, error) {
	bn := big.NewInt(1)
	bn, _ = bn.SetString(utils.TrimZeroAndDot(amount), 10)
	return lockAsync(ownerPrivateKey, ethereumReceiver, bn.Int64(), chainRelayer.bridgeBankAddr, chainRelayer.chainName, chainRelayer.rpcLaddr)
}

//ShowBridgeRegistryAddr ...
func (chainRelayer *Relayer4Chain) ShowBridgeRegistryAddr() (string, error) {
	if "" == chainRelayer.bridgeRegistryAddr {
		return "", errors.New("the relayer is not started yet")
	}

	return chainRelayer.bridgeRegistryAddr, nil
}

func (chainRelayer *Relayer4Chain) ShowStatics(request *ebTypes.TokenStaticsRequest) (*ebTypes.TokenStaticsResponse, error) {
	res := &ebTypes.TokenStaticsResponse{}

	datas, err := chainRelayer.getStatics(request.Operation, request.TxIndex, request.Count)
	if nil != err {
		return nil, err
	}
	//todo:完善分页显示功能
	for _, data := range datas {
		var statics ebTypes.Ethereum2ChainStatics
		_ = chainTypes.Decode(data, &statics)
		if request.Status != 0 && ebTypes.Tx_Status_Map[request.Status] != statics.ChainTxstatus {
			continue
		}
		if len(request.Symbol) > 0 && request.Symbol != statics.Symbol {
			continue
		}
		res.E2Cstatics = append(res.E2Cstatics, &statics)
	}
	return res, nil
}

func (chainRelayer *Relayer4Chain) updateTxStatus() {
	chainRelayer.updateSingleTxStatus(events.ClaimTypeBurn)
	chainRelayer.updateSingleTxStatus(events.ClaimTypeLock)
}

// 该函数用于定期检查是否有需要重新发送给以太坊协成的chain事件信息,用于产生relay event
func (chainRelayer *Relayer4Chain) checkTxRelay2Ethereum() {
	txInfos, err := chainRelayer.getAllTxsUnconfirm()
	if err != nil {
		relayerLog.Error("chainRelayer::checkTxRelay2Ethereum", "Failed to getAllTxsUnconfirm due to", err.Error())
		return
	}
	if 0 == len(txInfos) {
		return
	}
	for _, txInfo := range txInfos {
		txHashStr := chainEvmCommon.Bytes2Hex(txInfo.TxHash)

		if !txInfo.Resend {
			//为了防止转发出去的消息之后，下一个区块时间马上到来，首次转发的消息需要至少等一个区块间隔之后才会进行转发
			txInfo.Resend = true
			err = chainRelayer.setChainTxIsRelayedUnconfirm(txHashStr, txInfo.FdIndex, txInfo)
			if nil != err {
				relayerLog.Error("chainRelayer::checkTxRelay2Ethereum", "Failed to SetTxIsRelayedconfirm due to", err.Error())
				return
			}
			continue
		}

		chainMsg, err := events.ParseBurnLock4chain(events.ChainEvmEvent(txInfo.EventType), txInfo.Data, chainRelayer.bridgeBankAbi, txInfo.TxHash)
		if nil != err {
			relayerLog.Error("chainRelayer::checkTxRelay2Ethereum", "Failed to ParseBurnLock4chain due to", err.Error())
			return
		}
		txInfo.FdTimes = txInfo.FdTimes + 1
		chainMsg.ForwardTimes = txInfo.FdTimes
		chainMsg.ForwardIndex = txInfo.FdIndex

		channel, ok := chainRelayer.chainMsgChan[txInfo.ToChainName]
		if !ok {
			relayerLog.Error("chainRelayer::checkTxRelay2Ethereum", "No chainMsgChan for ethereum chain with name", txInfo.ToChainName)
			return
		}
		channel <- chainMsg

		//relaychainToEthereumCheckPonit 5: checkTxRelay2Ethereum
		relayerLog.Info("chainRelayer::relaychainToEthereumCheckPonit_5::checkTxRelay2Ethereum", "chainTxHash", txHashStr, "ForwardIndex", chainMsg.ForwardIndex, "FdTimes", chainMsg.ForwardTimes)
		err = chainRelayer.setChainTxIsRelayedUnconfirm(txHashStr, txInfo.FdIndex, txInfo)
		if nil != err {
			relayerLog.Error("chainRelayer::checkTxRelay2Ethereum", "Failed to SetTxIsRelayedconfirm due to", err.Error())
			return
		}
	}
}

//用于chain的事件信息被中继之后的ack信息，重置标志位
func (chainRelayer *Relayer4Chain) procTxRelayAck(ack *ebTypes.TxRelayAck) {
	//reset with another key to exclude from the check list to resend the same message
	if err := chainRelayer.resetKeyChainTxRelayedAlready(ack.TxHash); nil != err {
		relayerLog.Error("chainRelayer::procTxRelayAck", "Failed to resetKeyTxRelayedAlready due to:", err.Error())
		return
	}
	//relaychainToEthereumCheckPonit 4: recv ack from ethereum relay service
	relayerLog.Info("chainRelayer::procTxRelayAck::relaychainToEthereumCheckPonit_4", "chainTxHash", ack.TxHash, "ForwardIndex", ack.FdIndex)
}

func (chainRelayer *Relayer4Chain) updateSingleTxStatus(claimType events.ClaimType) {
	txIndex := chainRelayer.getChainUpdateTxIndex(claimType)
	datas, _ := chainRelayer.getStatics(int32(claimType), txIndex, 0)
	if nil == datas {
		return
	}
	for _, data := range datas {
		var statics ebTypes.Ethereum2ChainStatics
		_ = chainTypes.Decode(data, &statics)
		result := GetTxStatusByHashesRpc(statics.ChainTxhash, chainRelayer.rpcLaddr)
		//当前处理机制比较简单，如果发现该笔交易未执行，就不再产寻后续交易的回执
		if ebTypes.Invalid_ChainTx_Status == result {
			relayerLog.Debug("chainRelayer::updateSingleTxStatus", "no receipt for tx index", statics.TxIndex)
			break
		}
		status := ebTypes.Tx_Status_Success
		if result != chainTypes.ExecOk {
			status = ebTypes.Tx_Status_Failed
		}
		statics.ChainTxstatus = status
		dataNew := chainTypes.Encode(&statics)
		_ = chainRelayer.setLastestRelay2ChainTxStatics(statics.TxIndex, int32(claimType), dataNew)
		_ = chainRelayer.setChainUpdateTxIndex(statics.TxIndex, claimType)
		relayerLog.Debug("updateSingleTxStatus", "TxIndex", statics.TxIndex, "operationType", statics.OperationType, "txHash", statics.ChainTxhash, "updated status", status)
	}
}

func (chainRelayer *Relayer4Chain) SetupMulSign(setupMulSign *ebTypes.SetupMulSign) (string, error) {
	if "" == chainRelayer.mulSignAddr {
		return "", ebTypes.ErrMulSignNotDeployed
	}

	return setupMultiSign(setupMulSign.OperatorPrivateKey, chainRelayer.mulSignAddr, chainRelayer.chainName, chainRelayer.rpcLaddr, setupMulSign.Owners)
}

func (chainRelayer *Relayer4Chain) SafeTransfer(para *ebTypes.SafeTransfer) (string, error) {
	if "" == chainRelayer.mulSignAddr {
		return "", ebTypes.ErrMulSignNotDeployed
	}

	return safeTransfer(para.OwnerPrivateKeys[0], chainRelayer.mulSignAddr, chainRelayer.chainName,
		chainRelayer.rpcLaddr, para.To, para.Token, para.OwnerPrivateKeys, para.Amount)
}

func (chainRelayer *Relayer4Chain) SetMultiSignAddr(address string) {
	chainRelayer.rwLock.Lock()
	chainRelayer.mulSignAddr = address
	chainRelayer.rwLock.Unlock()

	chainRelayer.setMultiSignAddress(address)
}

func (chainRelayer *Relayer4Chain) GetMultiSignAddr() string {
	return chainRelayer.getMultiSignAddress()
}

func (chainRelayer *Relayer4Chain) WithdrawFromChain(ownerPrivateKey, tokenAddr, ethereumReceiver, amount string) (string, error) {
	bn := big.NewInt(1)
	bn, _ = bn.SetString(utils.TrimZeroAndDot(amount), 10)
	return withdrawAsync(ownerPrivateKey, tokenAddr, ethereumReceiver, bn.Int64(), chainRelayer.bridgeBankAddr, chainRelayer.chainName, chainRelayer.rpcLaddr)
}

func (chainRelayer *Relayer4Chain) BurnWithIncreaseAsyncFromChain(ownerPrivateKey, tokenAddr, ethereumReceiver, amount string) (string, error) {
	bn := big.NewInt(1)
	bn, _ = bn.SetString(utils.TrimZeroAndDot(amount), 10)
	return burnWithIncreaseAsync(ownerPrivateKey, tokenAddr, ethereumReceiver, bn.Int64(), chainRelayer.bridgeBankAddr, chainRelayer.chainName, chainRelayer.rpcLaddr)
}
