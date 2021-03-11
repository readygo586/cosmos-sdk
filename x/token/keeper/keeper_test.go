package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/token/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

var btcmeta = banktypes.Metadata{
	Description: "btc",
	Base:        "btc",
	Decimals:    8,
	Issuer:      "btc",
	SendEnabled: true,
}

type IntegrationTestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	queryClient types.QueryClient
}

func (suite *IntegrationTestSuite) SetupTest() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	app.BankKeeper.SetParams(ctx, banktypes.DefaultParams())

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.TokenKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.ctx = ctx
	suite.queryClient = queryClient
}

func (suite *IntegrationTestSuite) TestGetSetParams() {
	app, ctx := suite.app, suite.ctx
	params := app.TokenKeeper.GetParams(ctx)
	suite.Require().Equal(types.DefaultParams(), params)

	params.NewTokenFee = sdk.NewInt(100)
	params.TokenCacheSize = 10
	app.TokenKeeper.SetParams(ctx, params)

	params1 := app.TokenKeeper.GetParams(ctx)
	suite.Require().Equal(params, params1)
}

func (suite *IntegrationTestSuite) TestGetSetSybmol() {
	app, ctx := suite.app, suite.ctx

	symbols := app.TokenKeeper.GetSymbols(ctx)
	suite.Require().Equal(1, len(symbols))
	suite.Require().Equal(sdk.DefaultDenom, symbols[0])

	meta := app.TokenKeeper.GetTokenInfo(ctx, sdk.DefaultDenom)
	suite.Require().Equal(banktypes.DefaultMetadata()[0], meta)
	suite.Require().Equal("", app.TokenKeeper.GetIssuer(ctx, sdk.DefaultDenom))
	suite.Require().EqualValues(18, app.TokenKeeper.GetDecimals(ctx, sdk.DefaultDenom))
	suite.Require().Equal(true, app.TokenKeeper.SendEnabled(ctx, sdk.DefaultDenom))
	suite.Require().Equal(true, app.TokenKeeper.IsSupported(ctx, sdk.DefaultDenom))
	suite.Require().Equal(sdk.ZeroInt(), app.TokenKeeper.GetTotalSupply(ctx, sdk.DefaultDenom))

	meta.SendEnabled = false
	meta.Issuer = "me"
	meta.Description = "test"
	app.TokenKeeper.SetTokenInfo(ctx, meta)

	got := app.TokenKeeper.GetTokenInfo(ctx, sdk.DefaultDenom)
	suite.Require().Equal(meta, got)
	suite.Require().Equal("test", got.Description)
	suite.Require().Equal("me", app.TokenKeeper.GetIssuer(ctx, sdk.DefaultDenom))
	suite.Require().EqualValues(18, app.TokenKeeper.GetDecimals(ctx, sdk.DefaultDenom))
	suite.Require().Equal(false, app.TokenKeeper.SendEnabled(ctx, sdk.DefaultDenom))
	suite.Require().Equal(true, app.TokenKeeper.IsSupported(ctx, sdk.DefaultDenom))
	suite.Require().Equal(sdk.ZeroInt(), app.TokenKeeper.GetTotalSupply(ctx, sdk.DefaultDenom))

	app.TokenKeeper.EnableSend(ctx, sdk.DefaultDenom)
	suite.Require().Equal(true, app.TokenKeeper.SendEnabled(ctx, sdk.DefaultDenom))

	app.TokenKeeper.DisableSend(ctx, sdk.DefaultDenom)
	suite.Require().Equal(false, app.TokenKeeper.SendEnabled(ctx, sdk.DefaultDenom))
}

func (suite *IntegrationTestSuite) TestGetAllTokenInfos() {
	app, ctx := suite.app, suite.ctx

	symbols := app.TokenKeeper.GetSymbols(ctx)
	suite.Require().Equal(1, len(symbols))
	suite.Require().Equal(sdk.DefaultDenom, symbols[0])

	meta := app.TokenKeeper.GetTokenInfo(ctx, sdk.DefaultDenom)
	suite.Require().Equal(banktypes.DefaultMetadata()[0], meta)

	metas := app.TokenKeeper.GetAllTokenInfo(ctx)
	suite.Require().EqualValues(1, len(metas))
	suite.Require().Equal(banktypes.DefaultMetadata()[0], metas[0])

	app.TokenKeeper.SetTokenInfo(ctx, btcmeta)

	meta = app.TokenKeeper.GetTokenInfo(ctx, "btc")
	suite.Require().Equal(btcmeta, meta)

	metas = app.TokenKeeper.GetAllTokenInfo(ctx)
	suite.Require().EqualValues(2, len(metas))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
