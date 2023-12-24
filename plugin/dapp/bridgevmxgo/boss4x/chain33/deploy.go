package chain

import (
	"github.com/assetcloud/plugin/plugin/dapp/bridgevmxgo/boss4x/chain/offline"
	"github.com/spf13/cobra"
)

func ChainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain",
		Short: "deploy to chain",
	}
	cmd.AddCommand(
		offline.Boss4xOfflineCmd(),
		NewOracleClaimCmd(),
	)
	return cmd

}
