/*Package commands implement dapp client commands*/
package commands

import (
	vty "github.com/assetcloud/plugin/plugin/dapp/vote/types"
	jsonrpc "github.com/assetcloud/chain/rpc/jsonclient"
	rpctypes "github.com/assetcloud/chain/rpc/types"
	"github.com/assetcloud/chain/types"
	"github.com/spf13/cobra"
)

/*
 * 实现合约对应客户端
 */

// Cmd vote client command
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote",
		Short: "vote command",
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(
		//create tx
		createGroupCMD(),
		updateGroupCMD(),
		createVoteCMD(),
		commitVoteCMD(),
		closeVoteCMD(),
		updateMemberCMD(),
		//query rpc
		groupInfoCMD(),
		voteInfoCMD(),
		memberInfoCMD(),
		listGroupCMD(),
		listVoteCMD(),
		listMemberCMD(),
	)
	return cmd
}

func markRequired(cmd *cobra.Command, params ...string) {
	for _, param := range params {
		_ = cmd.MarkFlagRequired(param)
	}
}

func sendCreateTxRPC(cmd *cobra.Command, actionName string, req types.Message) {
	rpcAddr, _ := cmd.Flags().GetString("rpc_laddr")
	paraName, _ := cmd.Flags().GetString("paraName")
	payLoad := types.MustPBToJSON(req)
	pm := &rpctypes.CreateTxIn{
		Execer:     types.GetExecName(vty.VoteX, paraName),
		ActionName: actionName,
		Payload:    payLoad,
	}

	var res string
	ctx := jsonrpc.NewRPCCtx(rpcAddr, "Chain.CreateTransaction", pm, &res)
	ctx.RunWithoutMarshal()
}

func sendQueryRPC(cmd *cobra.Command, funcName string, req, reply types.Message) {
	rpcAddr, _ := cmd.Flags().GetString("rpc_laddr")
	paraName, _ := cmd.Flags().GetString("paraName")
	payLoad := types.MustPBToJSON(req)
	query := &rpctypes.Query4Jrpc{
		Execer:   types.GetExecName(vty.VoteX, paraName),
		FuncName: funcName,
		Payload:  payLoad,
	}

	ctx := jsonrpc.NewRPCCtx(rpcAddr, "Chain.Query", query, reply)
	ctx.Run()
}
