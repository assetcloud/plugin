package ethereum

import (
	"crypto/ecdsa"
	crand "crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	chainCommon "github.com/assetcloud/chain/common"
	dbm "github.com/assetcloud/chain/common/db"
	"github.com/assetcloud/chain/system/crypto/secp256k1"
	chainTypes "github.com/assetcloud/chain/types"
	wcom "github.com/assetcloud/chain/wallet/common"
	x2ethTypes "github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58/base58"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/ripemd160"
)

var (
	chainAccountKey = []byte("ChainAccount4EthRelayer")
	start             = int(1)
)

//Key ...
type Key struct {
	ID uuid.UUID // Version 4 "random" for unique id not derived from key data
	// to simplify lookups we also store the address
	Address common.Address
	// we only store privkey as pubkey/address can be derived from it
	// privkey in this struct is always in plaintext
	PrivateKey *ecdsa.PrivateKey
}

//NewAccount ...
func (ethRelayer *Relayer4Ethereum) NewAccount(passphrase string) (privateKeystr, addr string, err error) {
	_, privateKeystr, addr, err = newKeyAndStore(ethRelayer.db, crand.Reader, passphrase)
	if err != nil {
		return "", "", err
	}
	return
}

//GetAccount ...
func (ethRelayer *Relayer4Ethereum) GetAccount(passphrase string) (privateKey, addr string, err error) {
	accountInfo, err := ethRelayer.db.Get(chainAccountKey)
	if nil != err {
		return "", "", err
	}
	ChainAccount := &x2ethTypes.Account4Relayer{}
	if err := chainTypes.Decode(accountInfo, ChainAccount); nil != err {
		return "", "", err
	}
	decryptered := wcom.CBCDecrypterPrivkey([]byte(passphrase), ChainAccount.Privkey)
	privateKey = chainCommon.ToHex(decryptered)
	addr = ChainAccount.Addr
	return
}

//GetValidatorAddr ...
func (ethRelayer *Relayer4Ethereum) GetValidatorAddr() (validators x2ethTypes.ValidatorAddr4EthRelayer, err error) {
	var chainAccountAddr string
	accountInfo, err := ethRelayer.db.Get(chainAccountKey)
	if nil == err {
		ethAccount := &x2ethTypes.Account4Relayer{}
		if err := chainTypes.Decode(accountInfo, ethAccount); nil == err {
			chainAccountAddr = ethAccount.Addr
		}
	}

	if 0 == len(chainAccountAddr) {
		return x2ethTypes.ValidatorAddr4EthRelayer{}, x2ethTypes.ErrNoValidatorConfigured
	}

	validators = x2ethTypes.ValidatorAddr4EthRelayer{
		ChainValidator: chainAccountAddr,
	}
	return
}

//RestorePrivateKeys ...
func (ethRelayer *Relayer4Ethereum) RestorePrivateKeys(passPhase string) (err error) {
	accountInfo, err := ethRelayer.db.Get(chainAccountKey)
	if nil == err {
		ChainAccount := &x2ethTypes.Account4Relayer{}
		if err := chainTypes.Decode(accountInfo, ChainAccount); nil == err {
			decryptered := wcom.CBCDecrypterPrivkey([]byte(passPhase), ChainAccount.Privkey)
			var driver secp256k1.Driver
			priKey, err := driver.PrivKeyFromBytes(decryptered)
			if nil != err {
				errInfo := fmt.Sprintf("Failed to PrivKeyFromBytes due to:%s", err.Error())
				relayerLog.Info("RestorePrivateKeys", "Failed to PrivKeyFromBytes:", err.Error())
				return errors.New(errInfo)
			}
			ethRelayer.rwLock.Lock()
			ethRelayer.privateKey4Chain = priKey
			ethRelayer.rwLock.Unlock()
		}
	}

	ethRelayer.rwLock.RLock()
	if nil != ethRelayer.privateKey4Chain {
		ethRelayer.unlockchan <- start
	}
	ethRelayer.rwLock.RUnlock()

	return nil
}

//StoreAccountWithNewPassphase ...
func (ethRelayer *Relayer4Ethereum) StoreAccountWithNewPassphase(newPassphrase, oldPassphrase string) error {
	accountInfo, err := ethRelayer.db.Get(chainAccountKey)
	if nil != err {
		relayerLog.Info("StoreAccountWithNewPassphase", "pls check account is created already, err", err)
		return err
	}
	ChainAccount := &x2ethTypes.Account4Relayer{}
	if err := chainTypes.Decode(accountInfo, ChainAccount); nil != err {
		return err
	}
	decryptered := wcom.CBCDecrypterPrivkey([]byte(oldPassphrase), ChainAccount.Privkey)
	encryptered := wcom.CBCEncrypterPrivkey([]byte(newPassphrase), decryptered)
	ChainAccount.Privkey = encryptered
	encodedInfo := chainTypes.Encode(ChainAccount)
	return ethRelayer.db.SetSync(chainAccountKey, encodedInfo)
}

//ImportChainPrivateKey ...
func (ethRelayer *Relayer4Ethereum) ImportChainPrivateKey(passphrase, privateKeyStr string) error {
	var driver secp256k1.Driver
	privateKeySli, err := chainCommon.FromHex(privateKeyStr)
	if nil != err {
		return err
	}
	priKey, err := driver.PrivKeyFromBytes(privateKeySli)
	if nil != err {
		return err
	}

	ethRelayer.rwLock.Lock()
	ethRelayer.privateKey4Chain = priKey
	ethRelayer.rwLock.Unlock()
	ethRelayer.unlockchan <- start
	addr, err := pubKeyToAddress4Bty(priKey.PubKey().Bytes())
	if nil != err {
		return err
	}

	encryptered := wcom.CBCEncrypterPrivkey([]byte(passphrase), privateKeySli)
	account := &x2ethTypes.Account4Relayer{
		Privkey: encryptered,
		Addr:    addr,
	}
	encodedInfo := chainTypes.Encode(account)
	return ethRelayer.db.SetSync(chainAccountKey, encodedInfo)
}

//checksum: first four bytes of double-SHA256.
func checksum(input []byte) (cksum [4]byte) {
	h := sha256.New()
	_, err := h.Write(input)
	if err != nil {
		return
	}
	intermediateHash := h.Sum(nil)
	h.Reset()
	_, err = h.Write(intermediateHash)
	if err != nil {
		return
	}
	finalHash := h.Sum(nil)
	copy(cksum[:], finalHash[:])
	return
}

func pubKeyToAddress4Bty(pub []byte) (addr string, err error) {
	if len(pub) != 33 && len(pub) != 65 { //压缩格式 与 非压缩格式
		return "", fmt.Errorf("invalid public key byte")
	}

	sha256h := sha256.New()
	_, err = sha256h.Write(pub)
	if err != nil {
		return "", err
	}
	//160hash
	ripemd160h := ripemd160.New()
	_, err = ripemd160h.Write(sha256h.Sum([]byte("")))
	if err != nil {
		return "", err
	}
	//添加版本号
	hash160res := append([]byte{0}, ripemd160h.Sum([]byte(""))...)

	//添加校验码
	cksum := checksum(hash160res)
	address := append(hash160res, cksum[:]...)

	//地址进行base58编码
	addr = base58.Encode(address)
	return
}

func newKeyAndStore(db dbm.DB, rand io.Reader, passphrase string) (privateKey *ecdsa.PrivateKey, privateKeyStr, addr string, err error) {
	key, err := newKey(rand)
	if err != nil {
		return nil, "", "", err
	}
	privateKey = key.PrivateKey
	privateKeyBytes := math.PaddedBigBytes(key.PrivateKey.D, 32)
	Encryptered := wcom.CBCEncrypterPrivkey([]byte(passphrase), privateKeyBytes)
	ethAccount := &x2ethTypes.Account4Relayer{
		Privkey: Encryptered,
		Addr:    key.Address.Hex(),
	}
	_ = db

	privateKeyStr = chainCommon.ToHex(privateKeyBytes)
	addr = ethAccount.Addr
	return
}

func newKey(rand io.Reader) (*Key, error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand)
	if err != nil {
		return nil, err
	}
	return newKeyFromECDSA(privateKeyECDSA), nil
}

func newKeyFromECDSA(privateKeyECDSA *ecdsa.PrivateKey) *Key {
	id := uuid.NewRandom()
	key := &Key{
		ID:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	return key
}
