package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/token/types"
)

func (suite *IntegrationTestSuite) TestExportAndInitGeneis() {
	app, ctx := suite.app, suite.ctx
	orig := types.DefaultGenesisState()
	got := app.TokenKeeper.ExportGenesis(ctx)
	suite.Require().Equal(orig, got)

	orig.Params.TokenCacheSize = 10
	orig.Params.NewTokenFee = sdk.NewInt(123)
	app.TokenKeeper.InitGenesis(ctx, *orig)
	got = app.TokenKeeper.ExportGenesis(ctx)
	suite.Require().Equal(orig, got)
}
