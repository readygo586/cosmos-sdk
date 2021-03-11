package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/token/types"
	"github.com/spf13/cobra"
)

// NewTxCmd returns a root CLI command handler for all x/bank transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "token transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewNewTokenTxCmd(),
		NewInflateTokenTxCmd(),
		NewBurnTokenTxCmd(),
	)

	return txCmd
}

// NewNewTokenTxCmd returns a CLI command handler for creating a new token transaction.
func NewNewTokenTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [to_address][decimals][amount]",
		Short: "new a token",
		Long:  ` Example: new  cosmos1pxxxx 18 1000000000000000000000000000btc`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			toAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			decimals, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid decimals:%v", args[1])
			}

			coin, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgNewToken(clientCtx.GetFromAddress(), toAddr, uint64(decimals.Int64()), coin)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewInflateTokenTxCmd returns a CLI command handler for creating a inflate token transaction.
func NewInflateTokenTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inflate [to_address] [amount]",
		Short: `inflate a token and send to to_address.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			toAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			coin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgInflateToken(clientCtx.GetFromAddress(), toAddr, coin)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewBurnTokenTxCmd returns a CLI command handler for creating a burn token transaction.
func NewBurnTokenTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn coin",
		Short: "burn some token",
		Long:  ` Example: burn 10000000000btc --from alice`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgBurnToken(clientCtx.GetFromAddress(), coin)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
