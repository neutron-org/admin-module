package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/admin-module/v2/x/adminmodule/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
)

var _ = strconv.Itoa(0)

func CmdArchivedProposals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archivedproposals",
		Short: "Query archived proposals",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryArchivedProposalsRequest{}

			res, err := queryClient.ArchivedProposals(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdArchivedProposalsLegacy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archivedproposalslegacy",
		Short: "Query archived proposals legacy",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryArchivedProposalsLegacyRequest{}

			res, err := queryClient.ArchivedProposalsLegacy(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
