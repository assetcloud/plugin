// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authority_test

import (
	"fmt"
	"testing"

	"github.com/assetcloud/chain/common"
	"github.com/assetcloud/chain/common/address"
	"github.com/assetcloud/chain/common/crypto"
	drivers "github.com/assetcloud/chain/system/dapp"
	cty "github.com/assetcloud/chain/system/dapp/coins/types"
	"github.com/assetcloud/chain/types"
	"github.com/assetcloud/plugin/plugin/dapp/cert/authority"
	"github.com/assetcloud/plugin/plugin/dapp/cert/authority/utils"
	ct "github.com/assetcloud/plugin/plugin/dapp/cert/types"
	"github.com/stretchr/testify/assert"

	_ "github.com/assetcloud/chain/system"
	_ "github.com/assetcloud/plugin/plugin"
)

var (
	transfer = &ct.CertAction{Value: nil, Ty: ct.CertActionNormal}
	to       = drivers.ExecAddress("cert")
	tx1      = &types.Transaction{Execer: []byte("cert"), Payload: types.Encode(transfer), Fee: 1000000, Expire: 2, To: to}
	tx2      = &types.Transaction{Execer: []byte("cert"), Payload: types.Encode(transfer), Fee: 100000000, Expire: 0, To: to}
	tx3      = &types.Transaction{Execer: []byte("cert"), Payload: types.Encode(transfer), Fee: 200000000, Expire: 0, To: to}
	tx4      = &types.Transaction{Execer: []byte("cert"), Payload: types.Encode(transfer), Fee: 300000000, Expire: 0, To: to}
	tx5      = &types.Transaction{Execer: []byte("cert"), Payload: types.Encode(transfer), Fee: 400000000, Expire: 0, To: to}
	tx6      = &types.Transaction{Execer: []byte("cert"), Payload: types.Encode(transfer), Fee: 500000000, Expire: 0, To: to}
	tx7      = &types.Transaction{Execer: []byte("cert"), Payload: types.Encode(transfer), Fee: 600000000, Expire: 0, To: to}
	tx8      = &types.Transaction{Execer: []byte("cert"), Payload: types.Encode(transfer), Fee: 700000000, Expire: 0, To: to}
	tx9      = &types.Transaction{Execer: []byte("cert"), Payload: types.Encode(transfer), Fee: 800000000, Expire: 0, To: to}
	tx10     = &types.Transaction{Execer: []byte("cert"), Payload: types.Encode(transfer), Fee: 900000000, Expire: 0, To: to}
	tx11     = &types.Transaction{Execer: []byte("cert"), Payload: types.Encode(transfer), Fee: 450000000, Expire: 0, To: to}
	tx12     = &types.Transaction{Execer: []byte("cert"), Payload: types.Encode(transfer), Fee: 460000000, Expire: 0, To: to}
	tx13     = &types.Transaction{Execer: []byte("cert"), Payload: types.Encode(transfer), Fee: 100, Expire: 0, To: to}
	txs      = []*types.Transaction{tx1, tx2, tx3, tx4, tx5, tx6, tx7, tx8, tx9, tx10, tx11, tx12}

	privRaw, _  = common.FromHex("CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944")
	tr          = &cty.CoinsAction_Transfer{Transfer: &types.AssetsTransfer{Amount: int64(1e8)}}
	secpp256, _ = crypto.Load(types.GetSignName("", types.SECP256K1), -1)
	privKey, _  = secpp256.PrivKeyFromBytes(privRaw)
	tx14        = &types.Transaction{
		Execer:  []byte("coins"),
		Payload: types.Encode(&cty.CoinsAction{Value: tr, Ty: cty.CoinsActionTransfer}),
		Fee:     1000000,
		Expire:  2,
		To:      address.PubKeyToAddr(address.DefaultID, privKey.PubKey().Bytes()),
	}
)

var USERNAME = "user1"
var ORGNAME = "org1"
var SIGNTYPE = ct.AuthSM2

func signtx(tx *types.Transaction, priv crypto.PrivKey, cert []byte) {
	tx.Sign(int32(SIGNTYPE), priv)
	tx.Signature.Signature = utils.EncodeCertToSignature(tx.Signature.Signature, cert, nil)
}

func signtxs(priv crypto.PrivKey, cert []byte) {
	signtx(tx1, priv, cert)
	signtx(tx2, priv, cert)
	signtx(tx3, priv, cert)
	signtx(tx4, priv, cert)
	signtx(tx5, priv, cert)
	signtx(tx6, priv, cert)
	signtx(tx7, priv, cert)
	signtx(tx8, priv, cert)
	signtx(tx9, priv, cert)
	signtx(tx10, priv, cert)
	signtx(tx11, priv, cert)
	signtx(tx12, priv, cert)
	signtx(tx13, priv, cert)
}

/**
初始化Author实例和userloader
*/
func initEnv() (*types.ChainConfig, error) {
	cfg := types.NewChainConfig(types.ReadFile("./test/chain.auth.test.toml"))
	sub := cfg.GetSubConfig()
	var subcfg ct.Authority
	if sub.Exec["cert"] != nil {
		types.MustDecode(sub.Exec["cert"], &subcfg)
	}
	authority.Author.Init(&subcfg)
	SIGNTYPE = types.GetSignType("cert", subcfg.SignType)

	userLoader := &authority.UserLoader{}
	err := userLoader.Init(subcfg.CryptoPath, subcfg.SignType)
	if err != nil {
		fmt.Printf("Init user loader falied -> %v", err)
		return nil, err
	}

	user, err := userLoader.Get(USERNAME, ORGNAME)
	if err != nil {
		fmt.Printf("Get user failed")
		return nil, err
	}

	signtxs(user.Key, user.Cert)
	if err != nil {
		fmt.Printf("Init authority failed")
		return nil, err
	}

	return cfg, nil
}

/**
TestCase01 带证书的交易验签
*/
func TestChckSign(t *testing.T) {
	cfg, err := initEnv()
	if err != nil {
		t.Errorf("init env failed, error:%s", err)
		return
	}
	cfg.SetMinFee(0)

	assert.Equal(t, true, tx1.CheckSign(0))
}

/**
TestCase10 带证书的多交易验签
*/
func TestChckSigns(t *testing.T) {
	cfg, err := initEnv()
	if err != nil {
		t.Errorf("init env failed, error:%s", err)
		return
	}
	cfg.SetMinFee(0)

	for i, tx := range txs {
		if !tx.CheckSign(0) {
			t.Error(fmt.Sprintf("error check tx[%d]", i+1))
			return
		}
	}
}

/**
TestCase02 带证书的交易并行验签
*/
func TestChckSignsPara(t *testing.T) {
	cfg, err := initEnv()
	if err != nil {
		t.Errorf("init env failed, error:%s", err)
		return
	}
	cfg.SetMinFee(0)

	block := types.Block{}
	block.Txs = txs
	if !block.CheckSign(cfg) {
		t.Error("error check txs")
		return
	}
}

/**
TestCase03 不带证书，公链签名算法验证
*/
func TestChckSignWithNoneAuth(t *testing.T) {
	cfg, err := initEnv()
	if err != nil {
		t.Errorf("init env failed, error:%s", err)
		return
	}
	cfg.SetMinFee(0)

	tx14.Sign(types.SECP256K1, privKey)
	if !tx14.CheckSign(0) {
		t.Error("check signature failed")
		return
	}
}

/**
TestCase04 不带证书，SM2签名验证
*/
func TestChckSignWithSm2(t *testing.T) {
	sm2, err := crypto.Load(types.GetSignName("cert", ct.AuthSM2), -1)
	assert.Nil(t, err)
	privKeysm2, _ := sm2.PrivKeyFromBytes(privRaw)
	tx15 := &types.Transaction{Execer: []byte("coins"),
		Payload: types.Encode(&cty.CoinsAction{Value: tr, Ty: cty.CoinsActionTransfer}),
		Fee:     1000000, Expire: 2, To: address.PubKeyToAddr(address.DefaultID, privKeysm2.PubKey().Bytes())}

	cfg, err := initEnv()
	if err != nil {
		t.Errorf("init env failed, error:%s", err)
		return
	}
	cfg.SetMinFee(0)

	tx15.Sign(ct.AuthSM2, privKeysm2)
	if !tx15.CheckSign(0) {
		t.Error("check signature failed")
		return
	}
}

/**
TestCase05 不带证书，secp256r1签名验证
*/
func TestChckSignWithEcdsa(t *testing.T) {
	ecdsacrypto, _ := crypto.Load(types.GetSignName("cert", ct.AuthECDSA), -1)
	privKeyecdsa, _ := ecdsacrypto.PrivKeyFromBytes(privRaw)
	tx16 := &types.Transaction{Execer: []byte("coins"),
		Payload: types.Encode(&cty.CoinsAction{Value: tr, Ty: cty.CoinsActionTransfer}),
		Fee:     1000000, Expire: 2, To: address.PubKeyToAddr(address.DefaultID, privKeyecdsa.PubKey().Bytes())}

	cfg, err := initEnv()
	if err != nil {
		t.Errorf("init env failed, error:%s", err)
		return
	}
	cfg.SetMinFee(0)

	tx16.Sign(ct.AuthECDSA, privKeyecdsa)
	if !tx16.CheckSign(0) {
		t.Error("check signature failed")
		return
	}
}

/**
TestCase 06 证书检验
*/
func TestValidateCert(t *testing.T) {
	cfg, err := initEnv()
	if err != nil {
		t.Errorf("init env failed, error:%s", err)
		return
	}

	cfg.SetMinFee(0)

	for _, tx := range txs {
		err = authority.Author.Validate(tx.Signature)
		if err != nil {
			t.Error("error cert validate", err.Error())
			return
		}
	}
}

/**
Testcase07 noneimpl校验器验证（回滚到未开启证书验证的区块使用）
*/
func TestValidateTxWithNoneAuth(t *testing.T) {
	cfg, err := initEnv()
	if err != nil {
		t.Errorf("init env failed, error:%s", err)
		return
	}
	noneCertdata := &types.HistoryCertStore{}
	noneCertdata.CurHeigth = 0
	authority.Author.ReloadCert(noneCertdata)

	cfg.SetMinFee(0)

	err = authority.Author.Validate(tx14.Signature)
	if err != nil {
		t.Error("error cert validate", err.Error())
		return
	}
}

/**
Testcase08 重载历史证书
*/
func TestReloadCert(t *testing.T) {
	cfg, err := initEnv()
	if err != nil {
		t.Errorf("init env failed, error:%s", err)
		return
	}

	cfg.SetMinFee(0)

	store := &types.HistoryCertStore{}

	authority.Author.ReloadCert(store)

	err = authority.Author.Validate(tx1.Signature)
	if err != nil {
		t.Error(err.Error())
	}
}

/**
Testcase09 根据高度重载历史证书
*/
func TestReloadByHeight(t *testing.T) {
	cfg, err := initEnv()
	if err != nil {
		t.Errorf("init env failed, error:%s", err)
		return
	}
	cfg.SetMinFee(0)

	authority.Author.ReloadCertByHeght(30)
	if authority.Author.HistoryCertCache.CurHeight != 30 {
		t.Error("reload by height failed")
	}
}

//FIXME 有并发校验的场景需要考虑竞争，暂时没有并发校验的场景
/*
func TestValidateCerts(t *testing.T) {
	err := initEnv()
	if err != nil {
		t.Errorf("init env failed, error:%s", err)
	}

	prev := types.GetMinTxFeeRate()
	types.SetMinFee(0)
	defer types.SetMinFee(prev)

	signatures := []*types.Signature{tx1.Signature, tx2.Signature, tx3.Signature, tx4.Signature, tx5.Signature,
		tx6.Signature, tx7.Signature, tx8.Signature, tx9.Signature, tx10.Signature, tx11.Signature,
		tx12.Signature, tx13.Signature}

	result := Author.ValidateCerts(signatures)
	if !result {
		t.Error("error process txs signature validate")
	}
}
*/
