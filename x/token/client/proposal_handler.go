package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/cosmos/cosmos-sdk/x/token/client/cli"
	"github.com/cosmos/cosmos-sdk/x/token/client/rest"
)

var (
	DisableTokenProposalHandler      = govclient.NewProposalHandler(cli.NewCmdDisableTokenProposal, rest.DisableTokenProposalRESTHandler)
	TokenParamsChangeProposalHandler = govclient.NewProposalHandler(cli.NewCmdTokenParamsChangeProposal, rest.TokenParamsChangeProposalRESTHandler)
)
