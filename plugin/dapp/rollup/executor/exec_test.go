package executor

import (
	"testing"

	"github.com/assetcloud/chain/client/mocks"
	"github.com/assetcloud/chain/types"
	"github.com/assetcloud/chain/util"
	rtypes "github.com/assetcloud/plugin/plugin/dapp/rollup/types"
	"github.com/stretchr/testify/require"
)

var cfg *types.ChainConfig

func init() {

	cfg = types.NewChainConfig(types.ReadFile("../../../../chain.toml"))
	Init(driverName, cfg, nil)
}

func Test_rollup(t *testing.T) {
	r := newRollup()
	require.Equal(t, driverName, r.GetDriverName())
	require.Equal(t, driverName, GetName())
}

func TestRollup_Exec_Commit(t *testing.T) {

	r := newRollup()
	dir, state, _ := util.CreateTestDB()
	defer util.CloseTestDB(dir, state)
	api := &mocks.QueueProtocolAPI{}
	r.SetAPI(api)
	cfg := types.NewChainConfig(types.GetDefaultCfgstring())
	api.On("GetConfig").Return(cfg)
	r.SetStateDB(state)

	_ = state.Set(formatRollupStatusKey(""), []byte("test"))
	tx, err := r.GetExecutorType().CreateTransaction(rtypes.NameCommitAction, &rtypes.CheckPoint{})
	require.Nil(t, err)
	_, err = r.Exec(tx, 0)
	require.Equal(t, ErrGetRollupStatus, err)

	_ = state.Set(formatRollupStatusKey(""), types.Encode(&rtypes.RollupStatus{}))
	cp := &rtypes.CheckPoint{Batch: &rtypes.BlockBatch{BlockHeaders: []*types.Header{{}}}}
	tx, err = r.GetExecutorType().CreateTransaction(rtypes.NameCommitAction, cp)
	require.Nil(t, err)
	_, err = r.Exec(tx, 0)
	require.Nil(t, err)
}

func TestRollup_ExecLocal_CommitBatch(t *testing.T) {

	r := newRollup()
	tx, err := r.GetExecutorType().CreateTransaction(rtypes.NameCommitAction, &rtypes.CheckPoint{})
	require.Nil(t, err)
	_, err = r.ExecLocal(tx, nil, 0)
	require.Nil(t, err)
}

func TestRollup_ExecDelLocal(t *testing.T) {

	r := newRollup()
	dir, state, local := util.CreateTestDB()
	defer util.CloseTestDB(dir, state)

	r.SetLocalDB(local)
	tx, err := r.GetExecutorType().CreateTransaction(rtypes.NameCommitAction, &rtypes.CheckPoint{})
	require.Nil(t, err)
	_, err = r.ExecDelLocal(tx, nil, 0)
	require.Nil(t, err)
}
