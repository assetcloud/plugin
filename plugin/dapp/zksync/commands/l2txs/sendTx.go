package l2txs

import (
	"github.com/assetcloud/chain/common"
	"github.com/assetcloud/chain/rpc/jsonclient"
	rpctypes "github.com/assetcloud/chain/rpc/types"
	"github.com/assetcloud/chain/types"
	"github.com/spf13/cobra"
)

func SendChainL2TxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sendl2",
		Short: "send l2 tx to chain ",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(
		sendDepositTxCmd(),
		batchSendDepositTxCmd(),
		sendWithdrawTxCmd(),
		BatchSendTransferTxCmd(),
		SendTransferTxCmd(),
		sendManyDepositTxCmd(),
		sendManyWithdrawTxCmd(),
		treeManyToContractCmd(),
		contractManyToTreeCmd(),
		SendManyTransferTxCmd(),
		SendManyTransferTxFromOneCmd(),
		transferManyToNewCmd(),
		transferToNewManyCmd(),
		proxyManyExitCmd(),
		nftManyCmd(),
		setManyPubKeyCmd(),
		fetchL2BlockCmd(),
	)

	return cmd
}

func sendTx(rpcLaddr string, tx *types.Transaction) {
	txData := types.Encode(tx)
	dataStr := common.ToHex(txData)

	//fmt.Println("sendTx", "dataStr", dataStr)
	params := rpctypes.RawParm{
		Token: "BTY",
		Data:  dataStr,
	}

	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain.SendTransaction", params, nil)
	ctx.RunWithoutMarshal()
}
