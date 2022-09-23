package chain

import (
	chainCommon "github.com/assetcloud/chain/common"
	"github.com/ethereum/go-ethereum/crypto"

	//dbm "github.com/assetcloud/chain/common/db"
	chainTypes "github.com/assetcloud/chain/types"
	wcom "github.com/assetcloud/chain/wallet/common"
	x2ethTypes "github.com/assetcloud/plugin/plugin/dapp/x2ethereum/ebrelayer/types"
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

//ImportPrivateKey ...
func (chainRelayer *Relayer4Chain) ImportPrivateKey(passphrase, privateKeyStr string) (addr string, err error) {
	privateKeySlice, err := chainCommon.FromHex(privateKeyStr)
	if nil != err {
		return "", err
	}
	privateKey, err := crypto.ToECDSA(privateKeySlice)
	if nil != err {
		return "", err
	}

	ethSender := crypto.PubkeyToAddress(privateKey.PublicKey)
	chainRelayer.privateKey4Ethereum = privateKey
	chainRelayer.ethSender = ethSender
	chainRelayer.unlock <- start

	addr = chainCommon.ToHex(ethSender.Bytes())
	encryptered := wcom.CBCEncrypterPrivkey([]byte(passphrase), privateKeySlice)
	ethAccount := &x2ethTypes.Account4Relayer{
		Privkey: encryptered,
		Addr:    addr,
	}
	encodedInfo := chainTypes.Encode(ethAccount)
	err = chainRelayer.db.SetSync(chainAccountKey, encodedInfo)

	return
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
func (chainRelayer *Relayer4Chain) RestorePrivateKeys(passphrase string) error {
	accountInfo, err := chainRelayer.db.Get(chainAccountKey)
	if nil != err {
		relayerLog.Info("No private key saved for Relayer4Chain")
		return nil
	}
	ethAccount := &x2ethTypes.Account4Relayer{}
	if err := chainTypes.Decode(accountInfo, ethAccount); nil != err {
		relayerLog.Info("RestorePrivateKeys", "Failed to decode due to:", err.Error())
		return err
	}
	decryptered := wcom.CBCDecrypterPrivkey([]byte(passphrase), ethAccount.Privkey)
	privateKey, err := crypto.ToECDSA(decryptered)
	if nil != err {
		relayerLog.Info("RestorePrivateKeys", "Failed to ToECDSA:", err.Error())
		return err
	}

	chainRelayer.rwLock.Lock()
	chainRelayer.privateKey4Ethereum = privateKey
	chainRelayer.ethSender = crypto.PubkeyToAddress(privateKey.PublicKey)
	chainRelayer.rwLock.Unlock()
	chainRelayer.unlock <- start
	return nil
}

//func (chainRelayer *Relayer4Chain) UpdatePrivateKey(Passphrase, privateKey string) error {
//	return nil
//}
