package ethtxs

import (
	"fmt"
	"testing"

	"github.com/assetcloud/chain/client/mocks"
	chainCommon "github.com/assetcloud/chain/common"
	_ "github.com/assetcloud/chain/system"
	"github.com/assetcloud/chain/system/crypto/secp256k1"
	chainTypes "github.com/assetcloud/chain/types"
	"github.com/assetcloud/chain/util/testnode"
	ebrelayerTypes "github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	chainTestCfg = chainTypes.NewChainConfig(chainTypes.GetDefaultCfgstring())
)

func Test_RelayToChain(t *testing.T) {
	var tx chainTypes.Transaction
	var ret chainTypes.Reply
	ret.IsOk = true

	mockapi := &mocks.QueueProtocolAPI{}
	// 这里对需要mock的方法打桩,Close是必须的，其它方法根据需要
	mockapi.On("Close").Return()
	mockapi.On("AddPushSubscribe", mock.Anything).Return(&ret, nil)
	mockapi.On("CreateTransaction", mock.Anything).Return(&tx, nil)
	mockapi.On("SendTx", mock.Anything).Return(&ret, nil)
	mockapi.On("SendTransaction", mock.Anything).Return(&ret, nil)
	mockapi.On("GetConfig", mock.Anything).Return(chainTestCfg, nil)

	mock33 := testnode.New("", mockapi)
	defer mock33.Close()
	rpcCfg := mock33.GetCfg().RPC
	// 这里必须设置监听端口，默认的是无效值
	rpcCfg.JrpcBindAddr = "127.0.0.1:8801"
	mock33.GetRPC().Listen()

	chainPrivateKeyStr := "0xd627968e445f2a41c92173225791bae1ba42126ae96c32f28f97ff8f226e5c68"
	var driver secp256k1.Driver
	privateKeySli, err := chainCommon.FromHex(chainPrivateKeyStr)
	require.Nil(t, err)

	priKey, err := driver.PrivKeyFromBytes(privateKeySli)
	require.Nil(t, err)

	claim := &ebrelayerTypes.EthBridgeClaim{}

	fmt.Println("======================= testRelayLockToChain =======================")
	_, err = RelayLockToChain(priKey, claim, "http://127.0.0.1:8801")
	require.Nil(t, err)

	fmt.Println("======================= testRelayBurnToChain =======================")
	_, err = RelayBurnToChain(priKey, claim, "http://127.0.0.1:8801")
	require.Nil(t, err)
}
