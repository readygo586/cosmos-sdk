
package cli_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/token/client/cli"
	"github.com/gogo/protobuf/proto"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"io/ioutil"
)

func (s *IntegrationTestSuite) TestParseDisableTokenProposal() {
	okJSON := testutil.WriteToNewTempFile(s.T(), `
{
  "title": "Disable a token",
  "description": "disable a token!",
  "denom": "eos",
  "deposit": [
	{
      "denom": "stake",
      "amount": "1000"
    }
   ]
}
`)

	proposal, err := cli.ParseDisableTokenProposalJSON(s.cfg.LegacyAmino, okJSON.Name())
	s.Require().NoError(err)
	s.Require().Equal("Disable a token", proposal.Title)
	s.Require().Equal("disable a token!", proposal.Description)
	s.Require().Equal("eos", proposal.Denom)
	s.Require().Equal(sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1000))), proposal.Deposit)
}

func (s *IntegrationTestSuite) TestParseTokenParamsChangeProposal() {
	okJSON := testutil.WriteToNewTempFile(s.T(), `
{
  "title": "Token Parameters Change",
  "description": "token parameter change proposal",
  "denom": "eos",
  "changes": [
    {
      "key": "send_enabled",
      "value": true
    }
  ],
  "deposit": [
    {
      "denom": "stake",
      "amount": "1000"
    }
  ]
}
`)

	proposal, err := cli.ParseTokenParamsChangeProposalJSON(s.cfg.LegacyAmino, okJSON.Name())
	s.Require().NoError(err)
	s.Require().Equal("Token Parameters Change", proposal.Title)
	s.Require().Equal("token parameter change proposal", proposal.Description)
	s.Require().Equal("eos", proposal.Denom)
	s.Require().Equal("send_enabled", proposal.Changes[0].Key)
	s.Require().Equal(sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1000))), proposal.Deposit)
	val := false
	s.cfg.LegacyAmino.UnmarshalJSON(proposal.Changes[0].Value, &val)
	s.Require().Equal(true, val)
}

func (s *IntegrationTestSuite) TestDisableTokenProposalSuccess() {
	val := s.network.Validators[0]
	invalidPropFile, err := ioutil.TempFile(s.T().TempDir(), "invalid_text_proposal.*.json")
	s.Require().NoError(err)
	invalidProp := `
{
  "title": "Disable a token",
  "description": "disable a token!",
  "denom": "node0token",
  "deposit": [
	{
      "denom": "stake",
      "amount": "1000"
    }
   ]
}
`
	_, err = invalidPropFile.WriteString(invalidProp)
	s.Require().NoError(err)

	cmd := cli.NewCmdDisableTokenProposal()
	clientCtx := val.ClientCtx
	flags.AddTxFlagsToCmd(cmd)
	args := []string{
		invalidPropFile.Name(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation), //Skip confirmation
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
	s.Require().NoError(err)

	cmd = cli.GetCmdSymbol()
	args = []string{
		fmt.Sprintf("%stoken", val.Moniker),
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		fmt.Sprintf("--%s=1", flags.FlagHeight),
	}
	out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
	meta := banktypes.Metadata{}
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &meta))
	s.Require().Equal(true, meta.SendEnabled) //proposal will effect only after passing vote
	s.Require().Equal("node0token", meta.Base)

}

func (s *IntegrationTestSuite) TestDisableTokenProposal() {
	val := s.network.Validators[0]

	denomNotExistPropFile, err := ioutil.TempFile(s.T().TempDir(), "denom_not_exist_proposal.*.json")
	s.Require().NoError(err)
	denomNotExistProp := `
{
  "title": "Disable a token",
  "description": "disable a token!",
  "denom": "eos",
  "deposit": [
	{
      "denom": "stake",
      "amount": "1000"
    }
   ]
}
`
	_, err = denomNotExistPropFile.WriteString(denomNotExistProp)
	s.Require().NoError(err)

	disableDefaultDenomPropFile, err := ioutil.TempFile(s.T().TempDir(), "disable_default_denom_proposal.*.json")
	s.Require().NoError(err)
	disableDefaultProp := `
{
  "title": "Disable a token",
  "description": "disable a token!",
  "denom": "stake",
  "deposit": [
	{
      "denom": "stake",
      "amount": "1000"
    }
   ]
}
`

	_, err = disableDefaultDenomPropFile.WriteString(disableDefaultProp)
	s.Require().NoError(err)

	disableNodeDenomPropFile, err := ioutil.TempFile(s.T().TempDir(), "disable_node_denom_proposal.*.json")
	s.Require().NoError(err)
	disableNodeDenomProp := `
{
  "title": "Disable a token",
  "description": "disable a token!",
  "denom": "node0token",
  "deposit": [
	{
      "denom": "stake",
      "amount": "1000"
    }
   ]
}
`
	_, err = disableNodeDenomPropFile.WriteString(disableNodeDenomProp)
	s.Require().NoError(err)

	testCases := []struct {
		name              string
		args              []string
		expectErr         bool
		respType          proto.Message
		expectedCodeSpace string
		expectedCode      uint32
	}{
		{
			"disable not exist denom proposal",
			[]string{
				denomNotExistPropFile.Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation), //Skip confirmation
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, "gov", 5,
		},
		{
			"disable default denom proposal",
			[]string{
				disableDefaultDenomPropFile.Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation), //Skip confirmation
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, &sdk.TxResponse{}, "gov", 5, // return error in DisableTokenProposal's // ValidateBasic
		},
		{
			"insufficient gas fee",
			[]string{
				disableNodeDenomPropFile.Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation), //Skip confirmation
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1))).String()),
			},
			false, &sdk.TxResponse{}, "sdk", 13,
		},
		{
			"disable_node_denom_proposal",
			[]string{
				disableNodeDenomPropFile.Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation), //Skip confirmation
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, "", 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewCmdDisableTokenProposal()
			clientCtx := val.ClientCtx
			flags.AddTxFlagsToCmd(cmd)

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
				s.Require().Equal(tc.expectedCodeSpace, txResp.Codespace)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestTokenParamsChangeProposal() {
	val := s.network.Validators[0]

	denomNotExistPropFile, err := ioutil.TempFile(s.T().TempDir(), "denom_not_exist_proposal.*.json")
	s.Require().NoError(err)
	denomNotExistProp := `
{
  "title": "Token Parameters Change",
  "description": "token parameter change proposal",
  "denom": "eos",
  "changes": [
    {
      "key": "send_enabled",
      "value": true
    }
  ],
  "deposit": [
    {
      "denom": "stake",
      "amount": "1000"
    }
  ]
}
`
	_, err = denomNotExistPropFile.WriteString(denomNotExistProp)
	s.Require().NoError(err)

	disableDefaultDenomPropFile, err := ioutil.TempFile(s.T().TempDir(), "disable_default_denom_proposal.*.json")
	s.Require().NoError(err)
	disableDefaultProp := `
{
  "title": "Token Parameters Change",
  "description": "token parameter change proposal",
  "denom": "stake",
  "changes": [
    {
      "key": "send_enabled",
      "value": true
    }
  ],
  "deposit": [
    {
      "denom": "stake",
      "amount": "1000"
    }
  ]
}
`

	_, err = disableDefaultDenomPropFile.WriteString(disableDefaultProp)
	s.Require().NoError(err)

	disableNodeDenomPropFile, err := ioutil.TempFile(s.T().TempDir(), "disable_node_denom_proposal.*.json")
	s.Require().NoError(err)
	disableNodeDenomProp := `
{
  "title": "Token Parameters Change",
  "description": "token parameter change proposal",
  "denom": "node0token",
  "changes": [
    {
      "key": "send_enabled",
      "value": false
    }
  ],
  "deposit": [
    {
      "denom": "stake",
      "amount": "1000"
    }
  ]
}
`
	_, err = disableNodeDenomPropFile.WriteString(disableNodeDenomProp)
	s.Require().NoError(err)

	enableNodeDenomPropFile, err := ioutil.TempFile(s.T().TempDir(), "enable_node_denom_proposal.*.json")
	s.Require().NoError(err)
	enableNodeDenomProp := `
{
  "title": "Token Parameters Change",
  "description": "token parameter change proposal",
  "denom": "node0token",
  "changes": [
    {
      "key": "send_enabled",
      "value": false
    }
  ],
  "deposit": [
    {
      "denom": "stake",
      "amount": "1000"
    }
  ]
}
`
	_, err = enableNodeDenomPropFile.WriteString(enableNodeDenomProp)
	s.Require().NoError(err)

	testCases := []struct {
		name              string
		args              []string
		expectErr         bool
		respType          proto.Message
		expectedCodeSpace string
		expectedCode      uint32
	}{
		{
			"disable not exist denom proposal",
			[]string{
				denomNotExistPropFile.Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation), //Skip confirmation
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, "gov", 5,
		},
		{
			"disable default denom proposal",
			[]string{
				disableDefaultDenomPropFile.Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation), //Skip confirmation
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, &sdk.TxResponse{}, "gov", 5, // return error in DisableTokenProposal's // ValidateBasic
		},
		{
			"insufficient gas fee",
			[]string{
				disableNodeDenomPropFile.Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation), //Skip confirmation
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1))).String()),
			},
			false, &sdk.TxResponse{}, "sdk", 13,
		},
		{
			"disable_node_denom_proposal",
			[]string{
				disableNodeDenomPropFile.Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation), //Skip confirmation
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, "", 0,
		},
		{
			"enable_node_denom_proposal",
			[]string{
				enableNodeDenomPropFile.Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation), //Skip confirmation
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, "", 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewCmdTokenParamsChangeProposal()
			clientCtx := val.ClientCtx
			flags.AddTxFlagsToCmd(cmd)

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
				s.Require().Equal(tc.expectedCodeSpace, txResp.Codespace)
			}
		})
	}
}
