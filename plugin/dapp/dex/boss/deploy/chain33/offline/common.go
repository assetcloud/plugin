package offline

import (
	"github.com/assetcloud/plugin/plugin/dapp/dex/utils"
	"github.com/spf13/cobra"
)

func ChainOfflineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain",
		Short: "create and sign offline tx to deploy and set dex contracts to chain",
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(
		createERC20ContractCmd(),
		createRouterCmd(),
		farmofflineCmd(),
		sendSignTxs2ChainCmd(),
	)
	return cmd
}

func createRouterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "router",
		Short: "create and sign offline weth9, factory and router contracts",
		Run:   createRouterContract,
	}
	addCreateRouterFlags(cmd)
	return cmd
}

func addCreateRouterFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "k", "", "the private key")
	cmd.MarkFlagRequired("key")

	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")

	cmd.Flags().StringP("feeToSetter", "a", "", "address for fee to Setter")
	cmd.MarkFlagRequired("feeToSetter")

}

func sendSignTxs2ChainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "send one or serval dex txs to chain in serial",
		Run:   sendSignTxs2Chain,
	}
	addSendSignTxs2ChainFlags(cmd)
	return cmd
}

func addSendSignTxs2ChainFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("path", "p", "./", "(optional)path of txs file,default to current directroy")
	cmd.Flags().StringP("file", "f", "", "file name which contains the txs to be sent to chain")
	_ = cmd.MarkFlagRequired("file")
}

func sendSignTxs2Chain(cmd *cobra.Command, args []string) {
	filePath, _ := cmd.Flags().GetString("path")
	file, _ := cmd.Flags().GetString("file")
	url, _ := cmd.Flags().GetString("rpc_laddr")
	filePath += file
	utils.SendSignTxs2Chain(filePath, url)
}
