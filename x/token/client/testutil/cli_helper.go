package testutil

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	tokencli "github.com/cosmos/cosmos-sdk/x/token/client/cli"
)

func QuerySymbolExec(clientCtx client.Context, denom string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{denom}
	args = append(args, extraArgs...)
	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdSymbol(), args)
}
