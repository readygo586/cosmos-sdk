package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/token/types"
	"github.com/spf13/cobra"
)

// NewCmdTokenParamsChangeProposal implements a command handler for submitting a token param change proposal
func NewCmdTokenParamsChangeProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-params-change [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a token params change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a token params change proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ tx gov submit-proposal token-params-change <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Token Param Change",
  "description": "token param change proposal",
  "changes": [
    {
      "key": "is_send_enabled",
      "value": true
    },
  ],
  "deposit": [
    {
      "denom": "hbc",
      "amount": "10000"
    }
  ]
}
`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal, err := ParseTokenParamsChangeProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			changes := proposal.Changes.ToParamChanges()
			content := types.NewTokenParamsChangeProposal(proposal.Title, proposal.Description, proposal.Denom, changes)

			msg, err := govtypes.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}

// NewCmdDisableTokenProposal implements the command to submit a DisableToken proposal
func NewCmdDisableTokenProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable-token [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a disable token proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a disable token proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ tx gov submit-proposal disable-token <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Disable Token",
  "description": "disable token proposal",
  "symbol": "testtoken",
  "deposit": [
    {
      "denom": "hbc",
      "amount": "100000"
    }
  ]
}
`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal, err := ParseDisableTokenProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			content := types.NewDisableTokenProposal(proposal.Title, proposal.Description, proposal.Denom)
			msg, err := govtypes.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}
