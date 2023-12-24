package testnode

import (
	"testing"

	"github.com/assetcloud/chain/util"

	_ "github.com/assetcloud/chain/system"
	"github.com/assetcloud/chain/types"
	_ "github.com/assetcloud/plugin/plugin"
	pt "github.com/assetcloud/plugin/plugin/dapp/paracross/types"
	"github.com/stretchr/testify/assert"
)

func TestParaNode(t *testing.T) {
	para := NewParaNode(nil, nil)
	paraCfg := para.Para.GetAPI().GetConfig()
	defer para.Close()
	//通过rpc 发生信息
	tx := util.CreateTxWithExecer(paraCfg, para.Para.GetGenesisKey(), "user.p.test.none")
	assert.NotNil(t, tx)
	para.Para.SendTxRPC(tx)
	para.Para.WaitHeight(1)
	tx = util.CreateTxWithExecer(paraCfg, para.Para.GetGenesisKey(), "user.p.test.none")
	assert.NotNil(t, tx)
	para.Para.SendTxRPC(tx)
	para.Para.WaitHeight(2)

	res, err := para.Para.GetAPI().Query(pt.ParaX, "GetTitle", &types.ReqString{Data: "user.p.test."})
	assert.Nil(t, err)
	assert.Equal(t, int64(-1), res.(*pt.ParacrossStatus).Height)
}
