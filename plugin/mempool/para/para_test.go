package para_test

import (
	"testing"

	"github.com/assetcloud/chain/types"
	"github.com/assetcloud/chain/util"
	"github.com/assetcloud/chain/util/testnode"
	paratest "github.com/assetcloud/plugin/plugin/dapp/paracross/testnode"
	"github.com/assetcloud/plugin/plugin/mempool/para"
	"github.com/stretchr/testify/assert"

	_ "github.com/assetcloud/chain/system"
	_ "github.com/assetcloud/plugin/plugin"
)

func TestClose(t *testing.T) {
	mem := para.NewMempool(nil)
	n := 1000
	done := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		go func() {
			mem.Close()
			done <- struct{}{}
		}()
	}
	for i := 0; i < n; i++ {
		<-done
	}
}

func TestParaNodeMempool(t *testing.T) {
	main := testnode.New("", nil)
	main.Listen()

	chainCfg := types.NewChainConfigNoInit(paratest.DefaultConfig)
	chainCfg.GetModuleConfig().RPC.ParaChain.MainChainGrpcAddr = main.GetCfg().RPC.GrpcBindAddr
	cfg := chainCfg.GetModuleConfig()
	cfg.Mempool.Name = "para"
	para := testnode.NewWithConfig(chainCfg, nil)
	para.Listen()
	mockpara := paratest.NewParaNode(main, para)
	tx := util.CreateTxWithExecer(chainCfg, mockpara.Para.GetGenesisKey(), "user.p.guodun.none")
	hash := mockpara.Para.SendTx(tx)
	assert.Equal(t, tx.Hash(), hash)
	msg := para.GetClient().NewMessage("mempool", types.EventGetMempoolSize, nil)
	para.GetClient().Send(msg, true)
	reply, err := para.GetClient().Wait(msg)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("TestGetMempoolSize ", reply.GetData().(*types.MempoolSize).Size)

}
