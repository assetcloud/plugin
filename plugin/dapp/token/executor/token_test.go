package executor

import (
	"testing"

	"github.com/assetcloud/chain/account"
	"github.com/assetcloud/chain/common/address"
	"github.com/assetcloud/chain/types"
	"github.com/assetcloud/chain/util"

	"github.com/assetcloud/chain/common"
	"github.com/assetcloud/chain/common/crypto"
	dbm "github.com/assetcloud/chain/common/db"
	pty "github.com/assetcloud/plugin/plugin/dapp/token/types"
	"github.com/stretchr/testify/assert"

	//"github.com/assetcloud/chain/types/jsonpb"
	"strings"

	apimock "github.com/assetcloud/chain/client/mocks"
	"github.com/stretchr/testify/mock"
)

type execEnv struct {
	blockTime   int64
	blockHeight int64
	difficulty  uint64
}

var (
	Symbol         = "TEST"
	AssetExecToken = "token"
	AssetExecPara  = "paracross"

	PrivKeyA = "0x6da92a632ab7deb67d38c0f6560bcfed28167998f6496db64c258d5e8393a81b" // 1KSBd17H7ZK8iT37aJztFB22XGwsPTdwE4
	PrivKeyB = "0x19c069234f9d3e61135fefbeb7791b149cdf6af536f26bebb310d4cd22c3fee4" // 1JRNjdEqp4LJ5fqycUBm9ayCKSeeskgMKR
	PrivKeyC = "0x7a80a1f75d7360c6123c32a78ecf978c1ac55636f87892df38d8b85a9aeff115" // 1NLHPEcbTWWxxU3dGUZBhayjrCHD3psX7k
	PrivKeyD = "0xcacb1f5d51700aea07fca2246ab43b0917d70405c65edea9b5063d72eb5c6b71" // 1MCftFynyvG2F4ED5mdHYgziDxx6vDrScs
	Nodes    = [][]byte{
		[]byte("1KSBd17H7ZK8iT37aJztFB22XGwsPTdwE4"),
		[]byte("1JRNjdEqp4LJ5fqycUBm9ayCKSeeskgMKR"),
		[]byte("1NLHPEcbTWWxxU3dGUZBhayjrCHD3psX7k"),
		[]byte("1MCftFynyvG2F4ED5mdHYgziDxx6vDrScs"),
	}
)

func TestToken(t *testing.T) {
	cfg := types.NewChainConfig(strings.Replace(types.GetDefaultCfgstring(), "Title=\"local\"", "Title=\"chain\"", 1))
	cfg.SetDappFork(pty.TokenX, pty.ForkTokenCheckX, 1600000)
	Init(pty.TokenX, cfg, nil)
	tokenTotal := int64(10000 * 1e8)
	tokenBurn := int64(10 * 1e8)
	tokenMint := int64(20 * 1e8)
	total := int64(100000)
	accountA := types.Account{
		Balance: total,
		Frozen:  0,
		Addr:    string(Nodes[0]),
	}
	accountB := types.Account{
		Balance: total,
		Frozen:  0,
		Addr:    string(Nodes[1]),
	}

	execAddr := address.ExecAddress(pty.TokenX)
	stateDB, _ := dbm.NewGoMemDB("1", "2", 100)
	_, _, kvdb := util.CreateTestDB()

	accA, _ := account.NewAccountDB(cfg, AssetExecPara, Symbol, stateDB)
	accA.SaveExecAccount(execAddr, &accountA)

	accB, _ := account.NewAccountDB(cfg, AssetExecPara, Symbol, stateDB)
	accB.SaveExecAccount(execAddr, &accountB)

	env := execEnv{
		10,
		cfg.GetDappFork(pty.TokenX, pty.ForkTokenCheckX),
		1539918074,
	}

	// set config key
	item := &types.ConfigItem{
		Key: "mavl-manage-token-blacklist",
		Value: &types.ConfigItem_Arr{
			Arr: &types.ArrayConfig{Value: []string{"bty"}},
		},
	}
	stateDB.Set([]byte(item.Key), types.Encode(item))

	item2 := &types.ConfigItem{
		Key: "mavl-manage-token-finisher",
		Value: &types.ConfigItem_Arr{
			Arr: &types.ArrayConfig{Value: []string{string(Nodes[0])}},
		},
	}
	stateDB.Set([]byte(item2.Key), types.Encode(item2))

	// create token
	// 创建
	//ty := pty.TokenType{}
	p1 := &pty.TokenPreCreate{
		Name:         Symbol,
		Symbol:       Symbol,
		Introduction: Symbol,
		Total:        tokenTotal,
		Price:        0,
		Owner:        string(Nodes[0][1:]),
		Category:     pty.CategoryMintBurnSupport,
	}
	//v, _ := types.PBToJSON(p1)
	createTx, err := types.CallCreateTransaction(pty.TokenX, "TokenPreCreate", p1)
	if err != nil {
		t.Error("RPC_Default_Process", "err", err)
	}
	createTx, err = signTx(createTx, PrivKeyA)
	if err != nil {
		t.Error("RPC_Default_Process sign", "err", err)
	}
	exec := newToken()
	api := new(apimock.QueueProtocolAPI)
	api.On("GetConfig", mock.Anything).Return(cfg, nil)
	exec.SetAPI(api)
	exec.SetStateDB(stateDB)
	exec.SetLocalDB(kvdb)
	exec.SetEnv(env.blockHeight, env.blockTime, env.difficulty)
	receipt, err := exec.Exec(createTx, int(1))
	assert.NotNil(t, err)
	assert.Nil(t, receipt)

	p1 = &pty.TokenPreCreate{
		Name:         Symbol,
		Symbol:       Symbol,
		Introduction: Symbol,
		Total:        tokenTotal,
		Price:        0,
		Owner:        string(Nodes[0]),
		Category:     pty.CategoryMintBurnSupport,
	}
	//v, _ := types.PBToJSON(p1)
	createTx, err = types.CallCreateTransaction(pty.TokenX, "TokenPreCreate", p1)
	if err != nil {
		t.Error("RPC_Default_Process", "err", err)
	}
	createTx, err = signTx(createTx, PrivKeyA)
	if err != nil {
		t.Error("RPC_Default_Process sign", "err", err)
	}
	exec = newToken()
	exec.SetAPI(api)
	exec.SetStateDB(stateDB)
	exec.SetLocalDB(kvdb)
	exec.SetEnv(env.blockHeight, env.blockTime, env.difficulty)
	receipt, err = exec.Exec(createTx, int(1))
	assert.Nil(t, err)
	assert.NotNil(t, receipt)
	t.Log(receipt)
	for _, kv := range receipt.KV {
		stateDB.Set(kv.Key, kv.Value)
	}

	receiptDate := &types.ReceiptData{Ty: receipt.Ty, Logs: receipt.Logs}
	set, err := exec.ExecLocal(createTx, receiptDate, int(1))
	assert.Nil(t, err)
	assert.NotNil(t, set)
	for _, kv := range set.KV {
		kvdb.Set(kv.Key, kv.Value)
	}

	p2 := &pty.TokenFinishCreate{
		Symbol: Symbol,
		Owner:  string(Nodes[0]),
	}
	//v, _ := types.PBToJSON(p1)
	createTx2, err := types.CallCreateTransaction(pty.TokenX, "TokenFinishCreate", p2)
	if err != nil {
		t.Error("RPC_Default_Process", "err", err)
	}
	createTx2, err = signTx(createTx2, PrivKeyA)
	if err != nil {
		t.Error("RPC_Default_Process sign", "err", err)
	}

	exec.SetEnv(env.blockHeight+1, env.blockTime+1, env.difficulty)
	receipt, err = exec.Exec(createTx2, int(1))
	assert.Nil(t, err)
	assert.NotNil(t, receipt)
	//t.Log(receipt)
	for _, kv := range receipt.KV {
		stateDB.Set(kv.Key, kv.Value)
	}
	accDB, _ := account.NewAccountDB(cfg, pty.TokenX, Symbol, stateDB)
	accCheck := accDB.LoadAccount(string(Nodes[0]))
	assert.Equal(t, tokenTotal, accCheck.Balance)

	receiptDate = &types.ReceiptData{Ty: receipt.Ty, Logs: receipt.Logs}
	set, err = exec.ExecLocal(createTx2, receiptDate, int(1))
	assert.Nil(t, err)
	assert.NotNil(t, set)
	for _, kv := range set.KV {
		kvdb.Set(kv.Key, kv.Value)
	}

	// mint burn
	p3 := &pty.TokenMint{
		Symbol: Symbol,
		Amount: tokenMint,
	}
	//v, _ := types.PBToJSON(p1)
	createTx3, err := types.CallCreateTransaction(pty.TokenX, "TokenMint", p3)
	if err != nil {
		t.Error("RPC_Default_Process", "err", err)
	}
	createTx3, err = signTx(createTx3, PrivKeyA)
	if err != nil {
		t.Error("RPC_Default_Process sign", "err", err)
	}

	exec.SetEnv(env.blockHeight+2, env.blockTime+2, env.difficulty)
	receipt, err = exec.Exec(createTx3, int(1))
	assert.Nil(t, err)
	assert.NotNil(t, receipt)
	//t.Log(receipt)
	for _, kv := range receipt.KV {
		stateDB.Set(kv.Key, kv.Value)
	}

	accCheck = accDB.LoadAccount(string(Nodes[0]))
	assert.Equal(t, tokenTotal+tokenMint, accCheck.Balance)

	receiptDate = &types.ReceiptData{Ty: receipt.Ty, Logs: receipt.Logs}
	set, err = exec.ExecLocal(createTx3, receiptDate, int(1))
	assert.Nil(t, err)
	assert.NotNil(t, set)
	for _, kv := range set.KV {
		kvdb.Set(kv.Key, kv.Value)
	}

	p4 := &pty.TokenBurn{
		Symbol: Symbol,
		Amount: tokenBurn,
	}
	//v, _ := types.PBToJSON(p1)
	createTx4, err := types.CallCreateTransaction(pty.TokenX, "TokenBurn", p4)
	if err != nil {
		t.Error("RPC_Default_Process", "err", err)
	}
	createTx4, err = signTx(createTx4, PrivKeyA)
	if err != nil {
		t.Error("RPC_Default_Process sign", "err", err)
	}

	exec.SetEnv(env.blockHeight+1, env.blockTime+1, env.difficulty)
	receipt, err = exec.Exec(createTx4, int(1))
	assert.Nil(t, err)
	assert.NotNil(t, receipt)
	//t.Log(receipt)
	for _, kv := range receipt.KV {
		stateDB.Set(kv.Key, kv.Value)
	}
	accCheck = accDB.LoadAccount(string(Nodes[0]))
	assert.Equal(t, tokenTotal+tokenMint-tokenBurn, accCheck.Balance)

	receiptDate = &types.ReceiptData{Ty: receipt.Ty, Logs: receipt.Logs}
	set, err = exec.ExecLocal(createTx4, receiptDate, int(1))
	assert.Nil(t, err)
	assert.NotNil(t, set)
	for _, kv := range set.KV {
		kvdb.Set(kv.Key, kv.Value)
	}

	tokenExec, ok := exec.(*token)
	assert.True(t, ok)

	in := pty.ReqAccountTokenAssets{
		Address: string(Nodes[0]),
		Execer:  pty.TokenX,
	}
	out, err := tokenExec.Query_GetAccountTokenAssets(&in)
	assert.Nil(t, err)
	reply := out.(*pty.ReplyAccountTokenAssets)
	assert.Equal(t, 1, len(reply.TokenAssets))
	assert.NotEqual(t, 0, reply.TokenAssets[0].Account.Balance)
	assert.Equal(t, string(Nodes[0]), reply.TokenAssets[0].Account.Addr)
	t.Log(reply.TokenAssets)
}

func signTx(tx *types.Transaction, hexPrivKey string) (*types.Transaction, error) {
	signType := types.SECP256K1
	c, err := crypto.Load(types.GetSignName(pty.TokenX, signType), -1)
	if err != nil {
		return tx, err
	}

	bytes, err := common.FromHex(hexPrivKey[:])
	if err != nil {
		return tx, err
	}

	privKey, err := c.PrivKeyFromBytes(bytes)
	if err != nil {
		return tx, err
	}

	tx.Sign(int32(signType), privKey)
	return tx, nil
}
