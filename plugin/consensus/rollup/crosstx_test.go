package rollup

import (
	"context"
	"fmt"
	"testing"
	"time"

	pt "github.com/assetcloud/plugin/plugin/dapp/paracross/types"

	"github.com/assetcloud/chain/client"
	"github.com/assetcloud/chain/queue"
	"github.com/assetcloud/chain/system/consensus"

	"github.com/assetcloud/chain/rpc/grpcclient"
	_ "github.com/assetcloud/chain/system/consensus/init"
	_ "github.com/assetcloud/chain/system/dapp/init"
	_ "github.com/assetcloud/chain/system/mempool/init"
	_ "github.com/assetcloud/chain/system/store/init"
	"github.com/assetcloud/chain/types"
	"github.com/assetcloud/chain/util"
	"github.com/assetcloud/chain/util/testnode"
	rtypes "github.com/assetcloud/plugin/plugin/dapp/rollup/types"
	"github.com/stretchr/testify/require"
)

func newTestHandler() *crossTxHandler {
	ru := &RollUp{base: &consensus.BaseClient{}}
	h := &crossTxHandler{}
	h.init(ru, &rtypes.RollupStatus{})
	return h
}

func TestCrossTxHandler(t *testing.T) {

	h := newTestHandler()

	tx := &types.Transaction{Payload: []byte("test")}
	tx1 := &types.Transaction{Execer: []byte("user.p.test.paracross")}
	tx1.Payload = types.Encode(&pt.ParacrossAction{Ty: pt.ParacrossActionCrossAssetTransfer})
	h.addMainChainCrossTx(2, nil)
	require.Equal(t, 0, len(h.txIdxCache))
	h.addMainChainCrossTx(2, []*types.Transaction{tx, tx, tx1})
	require.Equal(t, 1, len(h.txIdxCache))
	idxArr := h.removePackedCrossTx([][]byte{tx1.Hash()})
	require.Equal(t, 0, len(h.txIdxCache))
	require.Equal(t, 1, len(idxArr))
	require.Equal(t, int64(2), idxArr[0].BlockHeight)
	require.Equal(t, int32(0), idxArr[0].FilterIndex)
	h.removePackedCrossTx(nil)
	require.Equal(t, 0, len(h.txIdxCache))
	idxArr = h.removePackedCrossTx([][]byte{tx.Hash()})
	require.Equal(t, 1, len(idxArr))
	require.Equal(t, int64(0), idxArr[0].BlockHeight)
	require.Equal(t, int32(0), idxArr[0].FilterIndex)
	require.Equal(t, tx.Hash(), idxArr[0].TxHash)
}

func TestRefreshSyncedHeight(t *testing.T) {

	h := newTestHandler()
	tx := &types.Transaction{Execer: []byte("user.p.test.paracross")}
	tx.Payload = types.Encode(&pt.ParacrossAction{Ty: pt.ParacrossActionCrossAssetTransfer})
	h.addMainChainCrossTx(2, []*types.Transaction{tx})
	require.Equal(t, 1, len(h.txIdxCache))
	info := h.txIdxCache[shortHash(tx.Hash())]
	require.Equal(t, int64(1), h.refreshSyncedHeight())
	info.enterTimestamp = types.Now().Unix() - 600
	require.Equal(t, int64(2), h.refreshSyncedHeight())
	require.Equal(t, 0, len(h.txIdxCache))
}

func TestRemoveErrTx(t *testing.T) {

	h := newTestHandler()
	tx := &types.Transaction{Execer: []byte("user.p.test.paracross")}
	tx.Payload = types.Encode(&pt.ParacrossAction{Ty: pt.ParacrossActionCrossAssetTransfer})
	h.addMainChainCrossTx(2, []*types.Transaction{tx})
	require.Equal(t, 1, len(h.txIdxCache))

	h.removeErrTxs([]*types.Transaction{tx})
	require.Equal(t, 0, len(h.txIdxCache))
}

func TestPullCrossTx(t *testing.T) {

	cfg := types.NewChainConfig(types.GetDefaultCfgstring())
	cfg.GetModuleConfig().RPC.GrpcBindAddr = fmt.Sprintf("localhost:%d", 9965)
	node := testnode.NewWithRPC(cfg, nil)
	defer node.Close()
	h := newTestHandler()
	grpc, err := grpcclient.NewMainChainClient(cfg, cfg.GetModuleConfig().RPC.GrpcBindAddr)
	require.Nil(t, err)
	cfg.SetTitleOnlyForTest("user.p.para")
	h.ru.chainCfg = cfg
	h.ru.mainChainGrpc = grpc
	h.ru.ctx = context.Background()
	txs := util.GenNoneTxs(cfg, node.GetGenesisKey(), 10)
	for i := 0; i < len(txs); i++ {
		_, err = node.GetAPI().SendTx(txs[i])
		require.Nil(t, err)
		require.Nil(t, node.WaitHeightTimeout(int64(i+1), 5))
	}

	go h.pullCrossTx()
	start := types.Now().Unix()
	for {
		h.lock.Lock()
		pulled := h.pulledHeight
		h.lock.Unlock()
		if pulled == 10-h.ru.cfg.ReservedMainHeight {
			return
		}
		if types.Now().Unix()-start >= 5 {
			t.Errorf("test timeout, pullHeight= %d", pulled)
			return
		}
		time.Sleep(time.Millisecond)
	}
}

func Test_send2Mempool(t *testing.T) {

	h := newTestHandler()

	q := queue.New("test")
	defer q.Close()
	api, _ := client.New(q.Client(), nil)
	h.ru.base.SetAPI(api)
	var expectTxs []*types.Transaction
	go func() {
		cli := q.Client()
		cli.Sub("mempool")
		count := 0
		for msg := range cli.Recv() {
			tx, ok := msg.GetData().(*types.Transaction)
			require.True(t, ok)
			require.Equal(t, expectTxs[count].Header, tx.Header)
			require.Equal(t, expectTxs[count].Hash(), tx.Hash())
			count++
			msg.Reply(&queue.Message{})
		}
	}()

	tx1 := &types.Transaction{Execer: []byte("user.p.test.coins"), Payload: []byte("test-tx1")}
	tx2 := &types.Transaction{Execer: []byte("user.p.test.paracross"), Payload: []byte("test-tx2")}
	tx3 := &types.Transaction{Execer: []byte("user.p.test.paracross")}

	txs, err := types.CreateTxGroup([]*types.Transaction{tx1, tx2}, 100)
	require.Nil(t, err)

	expectTxs = []*types.Transaction{txs.Tx(), tx3}
	h.send2Mempool(0, []*types.Transaction{tx1, tx2, tx3})
}
