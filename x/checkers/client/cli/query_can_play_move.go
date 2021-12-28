package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/zjing20/checkers/x/checkers/types"
)

var _ = strconv.Itoa(0)

func CmdCanPlayMove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "can-play-move [id-value] [player] [to-x] [to-y]",
		Short: "Query canPlayMove",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqIdValue := args[0]
			reqPlayer := args[1]
			reqToX, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			reqToY, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryCanPlayMoveRequest{

				IdValue: reqIdValue,
				Player:  reqPlayer,
				ToX:     reqToX,
				ToY:     reqToY,
			}

			res, err := queryClient.CanPlayMove(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
