/*Package commands implement dapp client commands*/
package commands

import (
	"bytes"
	"fmt"
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/types"
	"os"

	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	zt "github.com/33cn/plugin/plugin/dapp/zksync/types"
	"github.com/33cn/plugin/plugin/dapp/zksync/wallet"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

/*
 * 实现合约对应客户端
 */

// ZksyncCmd zksync client command
func ZksyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zksync",
		Short: "zksync command",
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(
		depositCmd(),
		withdrawCmd(),
		contractToLeafCmd(),
		leafToContractCmd(),
		transferCmd(),
		transferToNewCmd(),
		forceExitCmd(),
		setPubKeyCmd(),
		getChain33AddrCmd(),
		getAccountTreeCmd(),
		getTxProofCmd(),
		getTxProofByHeightCmd(),
		getAccountByIdCmd(),
		getAccountByEthCmd(),
		getAccountByChain33Cmd(),
		getContractAccountCmd(),
		getTokenBalanceCmd(),
	)
	return cmd
}

func depositCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit",
		Short: "get deposit tx",
		Run:   deposit,
	}
	depositFlag(cmd)
	return cmd
}

func depositFlag(cmd *cobra.Command) {
	cmd.Flags().Uint64P("tokenId", "t", 1, "deposit tokenId")
	cmd.MarkFlagRequired("tokenId")
	cmd.Flags().StringP("amount", "a", "0", "deposit amount")
	cmd.MarkFlagRequired("amount")
	cmd.Flags().StringP("ethAddress", "e", "", "deposit ethaddress")
	cmd.MarkFlagRequired("ethAddress")
	cmd.Flags().StringP("chain33Addr", "c", "", "deposit chain33Addr")
	cmd.MarkFlagRequired("chain33Addr")

}

func deposit(cmd *cobra.Command, args []string) {
	tokenId, _ := cmd.Flags().GetUint64("tokenId")
	amount, _ := cmd.Flags().GetString("amount")
	ethAddress, _ := cmd.Flags().GetString("ethAddress")
	chain33Addr, _ := cmd.Flags().GetString("chain33Addr")

	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	payload, err := wallet.CreateRawTx(zt.TyDepositAction, tokenId, amount, ethAddress, "", chain33Addr, 0, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "createRawTx"))
		return
	}
	params := &rpctypes.CreateTxIn{
		Execer:     zt.Zksync,
		ActionName: "Deposit",
		Payload:    payload,
	}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.CreateTransaction", params, nil)
	ctx.RunWithoutMarshal()
}

func withdrawCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "get withdraw tx",
		Run:   withdraw,
	}
	withdrawFlag(cmd)
	return cmd
}

func withdrawFlag(cmd *cobra.Command) {
	cmd.Flags().Uint64P("tokenId", "t", 1, "withdraw tokenId")
	cmd.MarkFlagRequired("tokenId")
	cmd.Flags().StringP("amount", "a", "0", "withdraw amount")
	cmd.MarkFlagRequired("amount")
	cmd.Flags().Uint64P("accountId", "", 0, "withdraw accountId")
	cmd.MarkFlagRequired("accountId")

}

func withdraw(cmd *cobra.Command, args []string) {
	tokenId, _ := cmd.Flags().GetUint64("tokenId")
	amount, _ := cmd.Flags().GetString("amount")
	accountId, _ := cmd.Flags().GetUint64("accountId")

	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	payload, err := wallet.CreateRawTx(zt.TyWithdrawAction, tokenId, amount, "", "", "", accountId, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "createRawTx"))
		return
	}
	params := &rpctypes.CreateTxIn{
		Execer:     zt.Zksync,
		ActionName: "Withdraw",
		Payload:    payload,
	}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.CreateTransaction", params, nil)
	ctx.RunWithoutMarshal()
}

func leafToContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "leafToContract",
		Short: "get leafToContract tx",
		Run:   leafToContract,
	}
	leafToContractFlag(cmd)
	return cmd
}

func leafToContractFlag(cmd *cobra.Command) {
	cmd.Flags().Uint64P("tokenId", "t", 1, "leafToContract tokenId")
	cmd.MarkFlagRequired("tokenId")
	cmd.Flags().StringP("amount", "a", "0", "leafToContract amount")
	cmd.MarkFlagRequired("amount")
	cmd.Flags().Uint64P("accountId", "", 0, "leafToContract accountId")
	cmd.MarkFlagRequired("accountId")

}

func leafToContract(cmd *cobra.Command, args []string) {
	tokenId, _ := cmd.Flags().GetUint64("tokenId")
	amount, _ := cmd.Flags().GetString("amount")
	accountId, _ := cmd.Flags().GetUint64("accountId")

	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	payload, err := wallet.CreateRawTx(zt.TyTreeToContractAction, tokenId, amount, "", "", "", accountId, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "createRawTx"))
		return
	}
	params := &rpctypes.CreateTxIn{
		Execer:     zt.Zksync,
		ActionName: "LeafToContract",
		Payload:    payload,
	}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.CreateTransaction", params, nil)
	ctx.RunWithoutMarshal()
}

func contractToLeafCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contractToLeaf",
		Short: "get contractToLeaf tx",
		Run:   contractToLeaf,
	}
	contractToLeafFlag(cmd)
	return cmd
}

func contractToLeafFlag(cmd *cobra.Command) {
	cmd.Flags().Uint64P("tokenId", "t", 1, "contractToLeaf tokenId")
	cmd.MarkFlagRequired("tokenId")
	cmd.Flags().StringP("amount", "a", "0", "contractToLeaf amount")
	cmd.MarkFlagRequired("amount")
	cmd.Flags().Uint64P("accountId", "", 0, "contractToLeaf accountId")
	cmd.MarkFlagRequired("accountId")
	cmd.Flags().StringP("ethAddress", "e", "", "deposit ethaddress")
	cmd.MarkFlagRequired("ethAddress")
	cmd.Flags().StringP("chain33Addr", "c", "", "deposit chain33Addr")
	cmd.MarkFlagRequired("chain33Addr")

}

func contractToLeaf(cmd *cobra.Command, args []string) {
	tokenId, _ := cmd.Flags().GetUint64("tokenId")
	amount, _ := cmd.Flags().GetString("amount")
	accountId, _ := cmd.Flags().GetUint64("accountId")
	ethAddress, _ := cmd.Flags().GetString("ethAddress")
	chain33Addr, _ := cmd.Flags().GetString("chain33Addr")

	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	payload, err := wallet.CreateRawTx(zt.TyContractToTreeAction, tokenId, amount, ethAddress, "", chain33Addr, accountId, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "createRawTx"))
		return
	}
	params := &rpctypes.CreateTxIn{
		Execer:     zt.Zksync,
		ActionName: "ContractToLeaf",
		Payload:    payload,
	}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.CreateTransaction", params, nil)
	ctx.RunWithoutMarshal()
}

func transferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "get transfer tx",
		Run:   transfer,
	}
	transferFlag(cmd)
	return cmd
}

func transferFlag(cmd *cobra.Command) {
	cmd.Flags().Uint64P("tokenId", "t", 1, "transfer tokenId")
	cmd.MarkFlagRequired("tokenId")
	cmd.Flags().StringP("amount", "a", "0", "transfer amount")
	cmd.MarkFlagRequired("amount")
	cmd.Flags().Uint64P("accountId", "", 0, "transfer fromAccountId")
	cmd.MarkFlagRequired("accountId")
	cmd.Flags().Uint64P("toAccountId", "", 0, "transfer toAccountId")
	cmd.MarkFlagRequired("toAccountId")

}

func transfer(cmd *cobra.Command, args []string) {
	tokenId, _ := cmd.Flags().GetUint64("tokenId")
	amount, _ := cmd.Flags().GetString("amount")
	accountId, _ := cmd.Flags().GetUint64("accountId")
	toAccountId, _ := cmd.Flags().GetUint64("toAccountId")

	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	payload, err := wallet.CreateRawTx(zt.TyTransferAction, tokenId, amount, "", "", "", accountId, toAccountId)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "createRawTx"))
		return
	}
	params := &rpctypes.CreateTxIn{
		Execer:     zt.Zksync,
		ActionName: "Transfer",
		Payload:    payload,
	}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.CreateTransaction", params, nil)
	ctx.RunWithoutMarshal()
}

func transferToNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transferToNew",
		Short: "get transferToNew tx",
		Run:   transferToNew,
	}
	transferToNewFlag(cmd)
	return cmd
}

func transferToNewFlag(cmd *cobra.Command) {
	cmd.Flags().Uint64P("tokenId", "t", 1, "transferToNew tokenId")
	cmd.MarkFlagRequired("tokenId")
	cmd.Flags().StringP("amount", "a", "0", "transferToNew amount")
	cmd.MarkFlagRequired("amount")
	cmd.Flags().Uint64P("accountId", "", 0, "transferToNew fromAccountId")
	cmd.MarkFlagRequired("accountId")
	cmd.Flags().StringP("toEthAddress", "", "", "transferToNew toEthAddress")
	cmd.MarkFlagRequired("toEthAddress")
	cmd.Flags().StringP("chain33Addr", "c", "", "transferToNew toChain33Addr")
	cmd.MarkFlagRequired("chain33Addr")
}

func transferToNew(cmd *cobra.Command, args []string) {
	tokenId, _ := cmd.Flags().GetUint64("tokenId")
	amount, _ := cmd.Flags().GetString("amount")
	accountId, _ := cmd.Flags().GetUint64("accountId")
	toEthAddress, _ := cmd.Flags().GetString("toEthAddress")
	chain33Addr, _ := cmd.Flags().GetString("chain33Addr")

	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	payload, err := wallet.CreateRawTx(zt.TyTransferToNewAction, tokenId, amount, "", toEthAddress, chain33Addr, accountId, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "createRawTx"))
		return
	}
	params := &rpctypes.CreateTxIn{
		Execer:     zt.Zksync,
		ActionName: "TransferToNew",
		Payload:    payload,
	}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.CreateTransaction", params, nil)
	ctx.RunWithoutMarshal()
}

func forceExitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "forceExit",
		Short: "get forceExit tx",
		Run:   forceExit,
	}
	forceExitFlag(cmd)
	return cmd
}

func forceExitFlag(cmd *cobra.Command) {
	cmd.Flags().Uint64P("tokenId", "t", 1, "forceExit tokenId")
	cmd.MarkFlagRequired("tokenId")
	cmd.Flags().Uint64P("accountId", "a", 0, "forceExit accountId")
	cmd.MarkFlagRequired("accountId")

}

func forceExit(cmd *cobra.Command, args []string) {
	tokenId, _ := cmd.Flags().GetUint64("tokenId")
	accountId, _ := cmd.Flags().GetUint64("accountId")

	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	payload, err := wallet.CreateRawTx(zt.TyForceExitAction, tokenId, "0", "", "", "", accountId, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "createRawTx"))
		return
	}
	params := &rpctypes.CreateTxIn{
		Execer:     zt.Zksync,
		ActionName: "ForceExit",
		Payload:    payload,
	}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.CreateTransaction", params, nil)
	ctx.RunWithoutMarshal()
}

func setPubKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setPubKey",
		Short: "get setPubKey tx",
		Run:   setPubKey,
	}
	setPubKeyFlag(cmd)
	return cmd
}

func setPubKeyFlag(cmd *cobra.Command) {
	cmd.Flags().Uint64P("accountId", "a", 0, "setPubKeyFlag accountId")
	cmd.MarkFlagRequired("accountId")

}

func setPubKey(cmd *cobra.Command, args []string) {
	accountId, _ := cmd.Flags().GetUint64("accountId")

	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	payload, err := wallet.CreateRawTx(zt.TySetPubKeyAction, 0, "0", "", "", "", accountId, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "createRawTx"))
		return
	}
	params := &rpctypes.CreateTxIn{
		Execer:     zt.Zksync,
		ActionName: "SetPubKey",
		Payload:    payload,
	}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.CreateTransaction", params, nil)
	ctx.RunWithoutMarshal()
}

func getChain33AddrCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getChain33Addr",
		Short: "get chain33 address by privateKey",
		Run:   getChain33Addr,
	}
	getChain33AddrFlag(cmd)
	return cmd
}

func getChain33AddrFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("private", "k", "", "private key")
}

func getChain33Addr(cmd *cobra.Command, args []string) {
	privateKeyString, _ := cmd.Flags().GetString("private")
	privateKeyBytes, err := common.FromHex(privateKeyString)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "hex.DecodeString"))
		return
	}
	privateKey, err := eddsa.GenerateKey(bytes.NewReader(privateKeyBytes))

	hash := mimc.NewMiMC(zt.ZkMimcHashSeed)
	hash.Write(privateKey.PublicKey.Bytes())
	fmt.Println(common.ToHex(hash.Sum(nil)))
}

func getAccountTreeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tree",
		Short: "get current accountTree",
		Run:   getAccountTree,
	}
	getAccountTreeFlag(cmd)
	return cmd
}

func getAccountTreeFlag(cmd *cobra.Command) {
}

func getAccountTree(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	var params rpctypes.Query4Jrpc

	params.Execer = zt.Zksync
	req := new(types.ReqNil)

	params.FuncName = "GetAccountTree"
	params.Payload = types.MustPBToJSON(req)

	var resp zt.AccountTree
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, &resp)
	ctx.Run()
}

func getTxProofCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proof",
		Short: "get tx proof",
		Run:   getTxProof,
	}
	getTxProofFlag(cmd)
	return cmd
}

func getTxProofFlag(cmd *cobra.Command) {
	cmd.Flags().Uint64P("height", "", 0, "zksync proof height")
	cmd.MarkFlagRequired("height")
	cmd.Flags().Uint32P("index", "i", 0, "tx index")
	cmd.MarkFlagRequired("index")
}

func getTxProof(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	height, _ := cmd.Flags().GetUint64("height")
	index, _ := cmd.Flags().GetUint32("index")

	var params rpctypes.Query4Jrpc

	params.Execer = zt.Zksync
	req := &zt.ZkQueryReq {
		BlockHeight: height,
		TxIndex: index,
	}

	params.FuncName = "GetTxProof"
	params.Payload = types.MustPBToJSON(req)

	var resp zt.OperationInfo
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, &resp)
	ctx.Run()
}

func getTxProofByHeightCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proofs",
		Short: "get block proofs",
		Run:   getTxProofByHeight,
	}
	getTxProofByHeightFlag(cmd)
	return cmd
}

func getTxProofByHeightFlag(cmd *cobra.Command) {
	cmd.Flags().Uint64P("height", "", 0, "zksync proof height")
	cmd.MarkFlagRequired("height")
}

func getTxProofByHeight(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	height, _ := cmd.Flags().GetUint64("height")

	var params rpctypes.Query4Jrpc

	params.Execer = zt.Zksync
	req := &zt.ZkQueryReq {
		BlockHeight: height,
	}

	params.FuncName = "GetTxProofByHeight"
	params.Payload = types.MustPBToJSON(req)

	var resp zt.ZkQueryResp
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, &resp)
	ctx.Run()
}

func getAccountByIdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "get zksync account by id",
		Run:   getAccountById,
	}
	getAccountByIdFlag(cmd)
	return cmd
}

func getAccountByIdFlag(cmd *cobra.Command) {
	cmd.Flags().Uint64P("accountId", "a", 0, "zksync accountId")
	cmd.MarkFlagRequired("accountId")
}

func getAccountById(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	accountId, _ := cmd.Flags().GetUint64("accountId")

	var params rpctypes.Query4Jrpc

	params.Execer = zt.Zksync
	req := &zt.ZkQueryReq {
		AccountId: accountId,
	}

	params.FuncName = "GetAccountById"
	params.Payload = types.MustPBToJSON(req)

	var resp zt.Leaf
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, &resp)
	ctx.Run()
}

func getAccountByEthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accountE",
		Short: "get zksync account by ethAddress",
		Run:   getAccountByEth,
	}
	getAccountByEthFlag(cmd)
	return cmd
}

func getAccountByEthFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("ethAddress", "e", " ", "zksync account ethAddress")
	cmd.MarkFlagRequired("ethAddress")
}

func getAccountByEth(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	ethAddress, _ := cmd.Flags().GetString("ethAddress")

	var params rpctypes.Query4Jrpc

	params.Execer = zt.Zksync
	req := &zt.ZkQueryReq {
		EthAddress: ethAddress,
	}

	params.FuncName = "GetAccountByEth"
	params.Payload = types.MustPBToJSON(req)

	var resp zt.ZkQueryResp
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, &resp)
	ctx.Run()
}

func getAccountByChain33Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accountC",
		Short: "get zksync account by chain33Addr",
		Run:   getAccountByChain33,
	}
	getAccountByChain33Flag(cmd)
	return cmd
}

func getAccountByChain33Flag(cmd *cobra.Command) {
	cmd.Flags().StringP("chain33Addr", "c", "", "zksync account chain33Addr")
	cmd.MarkFlagRequired("chain33Addr")
}

func getAccountByChain33(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	chain33Addr, _ := cmd.Flags().GetString("chain33Addr")

	var params rpctypes.Query4Jrpc

	params.Execer = zt.Zksync
	req := &zt.ZkQueryReq {
		Chain33Addr: chain33Addr,
	}

	params.FuncName = "GetAccountByChain33"
	params.Payload = types.MustPBToJSON(req)

	var resp zt.ZkQueryResp
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, &resp)
	ctx.Run()
}


func getContractAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contractAccount",
		Short: "get zksync contractAccount by chain33WalletAddr and token symbol",
		Run:   getContractAccount,
	}
	getContractAccountFlag(cmd)
	return cmd
}

func getContractAccountFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("chain33Addr", "c", "", "chain33 wallet address")
	cmd.MarkFlagRequired("chain33Addr")
	cmd.Flags().StringP("token", "t", "", "token symbol")
	cmd.MarkFlagRequired("token")
}

func getContractAccount(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	chain33Addr, _ := cmd.Flags().GetString("chain33Addr")
	token, _ := cmd.Flags().GetString("token")

	var params rpctypes.Query4Jrpc

	params.Execer = zt.Zksync
	req := &zt.ZkQueryReq {
		TokenSymbol: token,
		Chain33WalletAddr: chain33Addr,
	}

	params.FuncName = "GetZkContractAccount"
	params.Payload = types.MustPBToJSON(req)

	var resp types.Account
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, &resp)
	ctx.Run()
}



func getTokenBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "get zksync tokenBalance by accountId and tokenId",
		Run:   getTokenBalance,
	}
	getTokenBalanceFlag(cmd)
	return cmd
}

func getTokenBalanceFlag(cmd *cobra.Command) {
	cmd.Flags().Uint64P("accountId", "a", 0, "zksync account id")
	cmd.MarkFlagRequired("accountId")
	cmd.Flags().Uint64P("token", "t", 1, "zksync token id")
	cmd.MarkFlagRequired("token")
}

func getTokenBalance(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	accountId, _ := cmd.Flags().GetUint64("accountId")
	token, _ := cmd.Flags().GetUint64("token")

	var params rpctypes.Query4Jrpc

	params.Execer = zt.Zksync
	req := &zt.ZkQueryReq {
		TokenId: token,
		AccountId: accountId,
	}

	params.FuncName = "GetTokenBalance"
	params.Payload = types.MustPBToJSON(req)

	var resp zt.ZkQueryResp
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, &resp)
	ctx.Run()
}