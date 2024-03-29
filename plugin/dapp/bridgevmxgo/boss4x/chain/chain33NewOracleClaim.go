package chain

import (
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/assetcloud/chain/common"
	chainCommon "github.com/assetcloud/chain/common"
	"github.com/assetcloud/chain/common/address"
	chainCrypto "github.com/assetcloud/chain/common/crypto"
	"github.com/assetcloud/chain/rpc/jsonclient"
	rpctypes "github.com/assetcloud/chain/rpc/types"
	"github.com/assetcloud/chain/system/crypto/secp256k1"
	"github.com/assetcloud/chain/types"
	"github.com/assetcloud/plugin/plugin/dapp/bridgevmxgo/contracts/generated"
	"github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/relayer/events"
	"github.com/assetcloud/plugin/plugin/dapp/cross2eth/ebrelayer/utils"
	evmAbi "github.com/assetcloud/plugin/plugin/dapp/evm/executor/abi"
	evmtypes "github.com/assetcloud/plugin/plugin/dapp/evm/types"
	btcec_secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
)

func NewOracleClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn_xgo",
		Short: "burn xgo to chain evm",
		Run:   NewOracleClaim,
	}
	addNewOracleClaimFlags(cmd)
	return cmd
}

func addNewOracleClaimFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("relayerkey", "k", "", "relayer key")
	_ = cmd.MarkFlagRequired("relayerkey")
	cmd.Flags().StringP("token", "t", "", "bridge token address")
	_ = cmd.MarkFlagRequired("token")
	cmd.Flags().StringP("receiver", "r", "", "receiver address")
	_ = cmd.MarkFlagRequired("receiver")
	cmd.Flags().StringP("amount", "m", "", "amount")
	_ = cmd.MarkFlagRequired("amount")
	cmd.Flags().StringP("symbol", "s", "", "symbol")
	_ = cmd.MarkFlagRequired("symbol")
	cmd.Flags().StringP("fromAddr", "f", "", "fromAddr")
	_ = cmd.MarkFlagRequired("fromAddr")
	cmd.Flags().StringP("oracleAddr", "o", "", "oracleAddr")
	_ = cmd.MarkFlagRequired("oracleAddr")
	cmd.Flags().Int64P("nonce", "n", 0, "nonce")
	_ = cmd.MarkFlagRequired("nonce")
}

func NewOracleClaim(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	paraName, _ := cmd.Flags().GetString("paraName")
	chainID, _ := cmd.Flags().GetInt32("chainID")
	symbol, _ := cmd.Flags().GetString("symbol")
	tokenAddr, _ := cmd.Flags().GetString("token")
	amount, _ := cmd.Flags().GetString("amount")
	receiver, _ := cmd.Flags().GetString("receiver")
	nonce, _ := cmd.Flags().GetInt64("nonce")
	fromAddr, _ := cmd.Flags().GetString("fromAddr")
	oracleAddr, _ := cmd.Flags().GetString("oracleAddr")
	privateKeyStr, _ := cmd.Flags().GetString("relayerkey")

	var driver secp256k1.Driver
	privateKeySli, err := chainCommon.FromHex(privateKeyStr)
	if nil != err {
		fmt.Println("Failed to do chainCommon.FromHex")
		return
	}
	privateKey, err := driver.PrivKeyFromBytes(privateKeySli)
	if nil != err {
		fmt.Println("Failed to do PrivKeyFromBytes")
		return
	}

	temp, _ := btcec_secp256k1.PrivKeyFromBytes(btcec_secp256k1.S256(), privateKey.Bytes())
	privatekey4chainEcdsa := temp.ToECDSA()

	nonceBytes := big.NewInt(nonce).Bytes()
	bigAmount := big.NewInt(0)
	bigAmount.SetString(amount, 10)
	amountBytes := bigAmount.Bytes()
	claimID := crypto.Keccak256Hash(nonceBytes, []byte(fromAddr), []byte(receiver), []byte(symbol), amountBytes)

	signature, err := utils.SignClaim4Evm(claimID, privatekey4chainEcdsa)
	if nil != err {
		fmt.Println("SignClaim4Evm due to" + err.Error())
		return
	}

	parameter := fmt.Sprintf("newOracleClaim(%d, %s, %s, %s, %s, %s, %s, %s)",
		1,
		fromAddr,
		receiver,
		tokenAddr,
		symbol,
		amount,
		claimID.String(),
		common.ToHex(signature))

	note := fmt.Sprintf("relay with type:%s, chain-receiver:%s, ethereum-sender:%s, symbol:%s, amout:%s",
		events.ClaimType(1).String(), receiver, fromAddr, symbol, amount)
	_, packData, err := evmAbi.Pack(parameter, generated.OracleABI, false)
	if nil != err {
		fmt.Println("relayEvmTx2Chain", "Failed to do abi.Pack due to:", err.Error())
		return
	}

	action := evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: note, Para: packData, ContractAddr: oracleAddr}

	feeInt64 := int64(1 * 1e6)
	wholeEvm := getExecerName(paraName)
	toAddr := address.ExecAddress(wholeEvm)
	//name表示发给哪个执行器
	data := createEvmTx(privateKey, &action, wholeEvm, toAddr, feeInt64, chainID)
	params := rpctypes.RawParm{
		Token: "BTY",
		Data:  data,
	}
	var txhash string
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain.SendTransaction", params, &txhash)
	_, err = ctx.RunResult()
	fmt.Println(txhash)
}

func getExecerName(name string) string {
	var ret string
	names := strings.Split(name, ".")
	for _, v := range names {
		if v != "" {
			ret = ret + v + "."
		}
	}
	ret += "evm"
	return ret
}

func createEvmTx(privateKey chainCrypto.PrivKey, action proto.Message, execer, to string, fee int64, chainID int32) string {
	tx := &types.Transaction{Execer: []byte(execer), Payload: types.Encode(action), Fee: fee, To: to, ChainID: chainID}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	tx.Nonce = random.Int63()

	tx.Sign(types.SECP256K1, privateKey)
	txData := types.Encode(tx)
	dataStr := common.ToHex(txData)
	return dataStr
}
