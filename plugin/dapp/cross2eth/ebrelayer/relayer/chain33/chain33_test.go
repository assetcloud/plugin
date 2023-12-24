package chain

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	chainCommon "github.com/assetcloud/chain/common"
	dbm "github.com/assetcloud/chain/common/db"
	"github.com/assetcloud/chain/types"
	"github.com/assetcloud/chain/util/testnode"
	"github.com/assetcloud/plugin/plugin/dapp/cross2eth/contracts/test/setup"
	"github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/relayer/ethereum/ethtxs"
	"github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/relayer/events"
	ebTypes "github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/types"
	ebrelayerTypes "github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/types"
	relayerTypes "github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/types"
	tml "github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/assetcloud/chain/system"

	// 需要显示引用系统插件，以加载系统内置合约
	"github.com/assetcloud/chain/client/mocks"
	"github.com/stretchr/testify/mock"
)

var (
	configPath    = flag.String("f", "./../../relayer.toml", "configfile")
	privateKeyStr = "0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944"
	accountAddr   = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
	passphrase    = "123456hzj"
)

func Test_ImportRestorePrivateKey(t *testing.T) {
	var ret = types.ReplySubscribePush{IsOk: true, Msg: ""}
	var he = types.Header{Height: 10000}

	mockapi := &mocks.QueueProtocolAPI{}
	// 这里对需要mock的方法打桩,Close是必须的，其它方法根据需要
	mockapi.On("Close").Return()
	mockapi.On("AddPushSubscribe", mock.Anything).Return(&ret, nil)
	mockapi.On("GetLastHeader", mock.Anything).Return(&he, nil)

	mock33 := testnode.New("", mockapi)
	defer mock33.Close()
	rpcCfg := mock33.GetCfg().RPC
	// 这里必须设置监听端口，默认的是无效值
	rpcCfg.JrpcBindAddr = "127.0.0.1:8801"
	mock33.GetRPC().Listen()

	_, _, _, x2EthDeployInfo, err := setup.DeployContracts()
	require.NoError(t, err)
	chainRelayer := newChainRelayer(x2EthDeployInfo, "127.0.0.1:60000")

	err = chainRelayer.ImportPrivateKey(passphrase, privateKeyStr)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	addr, err := chainRelayer.GetAccountAddr()
	assert.NoError(t, err)
	assert.Equal(t, addr, accountAddr)

	key, _, _ := chainRelayer.GetAccount("123")
	assert.NotEqual(t, key, privateKeyStr)

	key, _, _ = chainRelayer.GetAccount(passphrase)
	assert.Equal(t, key, privateKeyStr)

	//////////////restore part//////////
	go func() {
		for range chainRelayer.unlockChan {
		}
	}()

	err = chainRelayer.RestorePrivateKeys("123")
	assert.NotEqual(t, chainCommon.ToHex(chainRelayer.privateKey4Chain.Bytes()), privateKeyStr)
	fmt.Println("err", err)
	assert.NoError(t, err)

	err = chainRelayer.RestorePrivateKeys(passphrase)
	assert.Equal(t, chainCommon.ToHex(chainRelayer.privateKey4Chain.Bytes()), privateKeyStr)
	assert.NoError(t, err)

	err = chainRelayer.StoreAccountWithNewPassphase("new123", passphrase)
	assert.NoError(t, err)

	err = chainRelayer.RestorePrivateKeys("new123")
	assert.Equal(t, chainCommon.ToHex(chainRelayer.privateKey4Chain.Bytes()), privateKeyStr)
	assert.NoError(t, err)

	time.Sleep(20 * time.Millisecond)
}

func newChainRelayer(x2EthDeployInfo *ethtxs.X2EthDeployInfo, pushBind string) *Relayer4Chain {
	cfg := initCfg(*configPath)

	chainMsgChan2Eths := make(map[string]chan<- *events.ChainMsg)
	ethBridgeClaimChan := make(chan *ebrelayerTypes.EthBridgeClaim, 100)

	for i := range cfg.EthRelayerCfg {
		cfg.EthRelayerCfg[i].BridgeRegistry = x2EthDeployInfo.BridgeRegistry.Address.String()
		chainMsgChan := make(chan *events.ChainMsg, 100)
		chainMsgChan2Eths[cfg.EthRelayerCfg[i].EthChainName] = chainMsgChan
	}
	cfg.ChainRelayerCfg.SyncTxConfig.PushBind = pushBind
	cfg.ChainRelayerCfg.SyncTxConfig.FetchHeightPeriodMs = 50
	cfg.ChainRelayerCfg.SyncTxConfig.ChainHost = "http://127.0.0.1:8801"
	cfg.Dbdriver = "memdb"

	db := dbm.NewDB("relayer_db_service", cfg.Dbdriver, cfg.DbPath, cfg.DbCache)
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	startPara := &ChainStartPara{
		ChainName:          cfg.ChainRelayerCfg.ChainName,
		Ctx:                ctx,
		SyncTxConfig:       cfg.ChainRelayerCfg.SyncTxConfig,
		BridgeRegistryAddr: cfg.ChainRelayerCfg.BridgeRegistryOnChain,
		DBHandle:           db,
		EthBridgeClaimChan: ethBridgeClaimChan,
		ChainMsgChan:     chainMsgChan2Eths,
		ChainID:            cfg.ChainRelayerCfg.ChainID4Chain,
	}

	relayer := &Relayer4Chain{
		rpcLaddr:                startPara.SyncTxConfig.ChainHost,
		chainName:               startPara.ChainName,
		chainID:                 startPara.ChainID,
		fetchHeightPeriodMs:     startPara.SyncTxConfig.FetchHeightPeriodMs,
		unlockChan:              make(chan int),
		db:                      startPara.DBHandle,
		ctx:                     startPara.Ctx,
		bridgeRegistryAddr:      startPara.BridgeRegistryAddr,
		ethBridgeClaimChan:      startPara.EthBridgeClaimChan,
		chainMsgChan:          startPara.ChainMsgChan,
		totalTx4RelayEth2chai33: 0,
		symbol2Addr:             make(map[string]string),
	}

	syncCfg := &ebTypes.SyncTxReceiptConfig{
		ChainHost:       startPara.SyncTxConfig.ChainHost,
		PushHost:          startPara.SyncTxConfig.PushHost,
		PushName:          startPara.SyncTxConfig.PushName,
		PushBind:          startPara.SyncTxConfig.PushBind,
		StartSyncHeight:   startPara.SyncTxConfig.StartSyncHeight,
		StartSyncSequence: startPara.SyncTxConfig.StartSyncSequence,
		StartSyncHash:     startPara.SyncTxConfig.StartSyncHash,
	}
	go relayer.syncProc(syncCfg)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		<-ch
		cancel()
		wg.Wait()
		os.Exit(0)
	}()
	return relayer
}

func initCfg(path string) *relayerTypes.RelayerConfig {
	var cfg relayerTypes.RelayerConfig
	if _, err := tml.DecodeFile(path, &cfg); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return &cfg
}

func Test_getExecerName(t *testing.T) {
	assert.Equal(t, getExecerName(""), "evm")
	assert.Equal(t, getExecerName("user.p.para."), "user.p.para.evm")
	assert.Equal(t, getExecerName("user.p.para.."), "user.p.para.evm")
	assert.Equal(t, getExecerName("user...p.para.."), "user.p.para.evm")
	assert.Equal(t, getExecerName("user.p...para.."), "user.p.para.evm")
	assert.Equal(t, getExecerName("user.p.para"), "user.p.para.evm")
	assert.Equal(t, getExecerName("user"), "user.evm")
}
