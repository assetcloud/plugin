package l2txs

import (
	"fmt"
	"strconv"
	"strings"

	zksyncTypes "github.com/assetcloud/plugin/plugin/dapp/zksync/types"
	"github.com/spf13/cobra"
)

func sendWithdrawTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "send withdraw tx to chain",
		Run:   sendWithdraw,
	}
	sendWithdrawFlags(cmd)
	return cmd
}

func sendWithdrawFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "k", "", "private key")
	_ = cmd.MarkFlagRequired("key")

	cmd.Flags().Uint64P("tokenId", "t", 0, "eth token id")
	_ = cmd.MarkFlagRequired("tokenId")
	cmd.Flags().Uint64P("accountID", "a", 0, "L2 account id on chain")
	_ = cmd.MarkFlagRequired("accountID")
	cmd.Flags().StringP("amount", "m", "0", "deposit amount")
	_ = cmd.MarkFlagRequired("amount")
}

func sendWithdraw(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	privateKey, _ := cmd.Flags().GetString("key")

	tokenId, _ := cmd.Flags().GetUint64("tokenId")
	accountID, _ := cmd.Flags().GetUint64("accountID")
	amount, _ := cmd.Flags().GetString("amount")
	paraName, _ := cmd.Flags().GetString("paraName")

	withdraw := &zksyncTypes.ZkWithdraw{
		TokenId:   tokenId,
		Amount:    amount,
		AccountId: accountID,
	}

	action := &zksyncTypes.ZksyncAction{
		Ty: zksyncTypes.TyWithdrawAction,
		Value: &zksyncTypes.ZksyncAction_ZkWithdraw{
			ZkWithdraw: withdraw,
		},
	}

	tx, err := createChainTx(privateKey, getRealExecName(paraName, zksyncTypes.Zksync), action)
	if nil != err {
		fmt.Println("sendDeposit failed to createChainTx due to err:", err.Error())
		return
	}
	sendTx(rpcLaddr, tx)
}

func sendManyWithdrawTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw_many",
		Short: "send withdraw tx to chain",
		Run:   sendManyWithdraw,
	}
	sendManyWithdrawFlags(cmd)
	return cmd
}

func sendManyWithdrawFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("keys", "k", "", "private keys, use ',' separate")
	_ = cmd.MarkFlagRequired("keys")
	cmd.Flags().Uint64P("tokenId", "t", 0, "eth token id")
	_ = cmd.MarkFlagRequired("tokenId")
	cmd.Flags().StringP("accountIDs", "a", "0", "L2 account ids on chain, use ',' separate")
	_ = cmd.MarkFlagRequired("accountIDs")
	cmd.Flags().StringP("amount", "m", "0", "deposit amount")
	_ = cmd.MarkFlagRequired("amount")
}

func sendManyWithdraw(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	privateKeys, _ := cmd.Flags().GetString("keys")
	tokenId, _ := cmd.Flags().GetUint64("tokenId")
	accountIDs, _ := cmd.Flags().GetString("accountIDs")
	amount, _ := cmd.Flags().GetString("amount")
	paraName, _ := cmd.Flags().GetString("paraName")

	ids := strings.Split(accountIDs, ",")
	keys := strings.Split(privateKeys, ",")

	if len(ids) != len(keys) {
		fmt.Println("err len(ids) != len(keys)", len(ids), "!=", len(keys))
		return
	}

	for i := 0; i < len(ids); i++ {
		id, _ := strconv.ParseInt(ids[i], 10, 64)
		withdraw := &zksyncTypes.ZkWithdraw{
			TokenId:   tokenId,
			Amount:    amount,
			AccountId: uint64(id),
		}

		action := &zksyncTypes.ZksyncAction{
			Ty: zksyncTypes.TyWithdrawAction,
			Value: &zksyncTypes.ZksyncAction_ZkWithdraw{
				ZkWithdraw: withdraw,
			},
		}

		tx, err := createChainTx(keys[i], getRealExecName(paraName, zksyncTypes.Zksync), action)
		if nil != err {
			fmt.Println("sendDeposit failed to createChainTx due to err:", err.Error())
			return
		}
		sendTx(rpcLaddr, tx)
	}
}
