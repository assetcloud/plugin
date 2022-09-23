package chain

import (
	"errors"
	"fmt"

	chainCommon "github.com/assetcloud/chain/common"
	"github.com/assetcloud/chain/common/address"
	"github.com/assetcloud/chain/system/crypto/secp256k1"
	chainTypes "github.com/assetcloud/chain/types"
	wcom "github.com/assetcloud/chain/wallet/common"
	x2ethTypes "github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/types"
	btcec_secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	chainAccountKey = []byte("ChainAccount4Relayer")
	start           = int(1)
)

//GetAccount ...
func (chainRelayer *Relayer4Chain) GetAccount(passphrase string) (privateKey, addr string, err error) {
	accountInfo, err := chainRelayer.db.Get(chainAccountKey)
	if nil != err {
		return "", "", err
	}
	ethAccount := &x2ethTypes.Account4Relayer{}
	if err := chainTypes.Decode(accountInfo, ethAccount); nil != err {
		return "", "", err
	}
	decryptered := wcom.CBCDecrypterPrivkey([]byte(passphrase), ethAccount.Privkey)
	privateKey = chainCommon.ToHex(decryptered)
	addr = ethAccount.Addr
	return
}

//GetAccountAddr ...
func (chainRelayer *Relayer4Chain) GetAccountAddr() (addr string, err error) {
	accountInfo, err := chainRelayer.db.Get(chainAccountKey)
	if nil != err {
		relayerLog.Info("GetValidatorAddr", "Failed to get account from db due to:", err.Error())
		return "", err
	}
	ethAccount := &x2ethTypes.Account4Relayer{}
	if err := chainTypes.Decode(accountInfo, ethAccount); nil != err {
		relayerLog.Info("GetValidatorAddr", "Failed to decode due to:", err.Error())
		return "", err
	}
	addr = ethAccount.Addr
	return
}

func (chainRelayer *Relayer4Chain) ImportPrivateKey(passphrase, privateKeyStr string) error {
	var driver secp256k1.Driver
	privateKeySli, err := chainCommon.FromHex(privateKeyStr)
	if nil != err {
		return err
	}
	priKey, err := driver.PrivKeyFromBytes(privateKeySli)
	if nil != err {
		return err
	}

	chainRelayer.rwLock.Lock()
	chainRelayer.privateKey4Chain = priKey
	temp, _ := btcec_secp256k1.PrivKeyFromBytes(btcec_secp256k1.S256(), priKey.Bytes())
	chainRelayer.privateKey4Chain_ecdsa = temp.ToECDSA()
	chainRelayer.rwLock.Unlock()
	chainRelayer.unlockChan <- start
	addr := address.PubKeyToAddr(address.DefaultID, priKey.PubKey().Bytes())

	encryptered := wcom.CBCEncrypterPrivkey([]byte(passphrase), privateKeySli)
	account := &x2ethTypes.Account4Relayer{
		Privkey: encryptered,
		Addr:    addr,
	}
	encodedInfo := chainTypes.Encode(account)
	return chainRelayer.db.SetSync(chainAccountKey, encodedInfo)
}

//StoreAccountWithNewPassphase ...
func (chainRelayer *Relayer4Chain) StoreAccountWithNewPassphase(newPassphrase, oldPassphrase string) error {
	accountInfo, err := chainRelayer.db.Get(chainAccountKey)
	if nil != err {
		relayerLog.Info("StoreAccountWithNewPassphase", "pls check account is created already, err", err)
		return err
	}
	ethAccount := &x2ethTypes.Account4Relayer{}
	if err := chainTypes.Decode(accountInfo, ethAccount); nil != err {
		return err
	}
	decryptered := wcom.CBCDecrypterPrivkey([]byte(oldPassphrase), ethAccount.Privkey)
	encryptered := wcom.CBCEncrypterPrivkey([]byte(newPassphrase), decryptered)
	ethAccount.Privkey = encryptered
	encodedInfo := chainTypes.Encode(ethAccount)
	return chainRelayer.db.SetSync(chainAccountKey, encodedInfo)
}

//RestorePrivateKeys ...
func (chainRelayer *Relayer4Chain) RestorePrivateKeys(passPhase string) (err error) {
	accountInfo, err := chainRelayer.db.Get(chainAccountKey)
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
			chainRelayer.rwLock.Lock()
			chainRelayer.privateKey4Chain = priKey
			chainRelayer.privateKey4Chain_ecdsa, err = crypto.ToECDSA(priKey.Bytes())
			if nil != err {
				return err
			}
			chainRelayer.rwLock.Unlock()
		}
	}

	chainRelayer.rwLock.RLock()
	if nil != chainRelayer.privateKey4Chain {
		chainRelayer.unlockChan <- start
	}
	chainRelayer.rwLock.RUnlock()

	return nil
}
