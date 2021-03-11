// +build norace

package keeper_test

import (
	gocontext "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/token/types"
)

func (suite *IntegrationTestSuite) TestQuerySybmol() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	res, err := queryClient.Symbol(gocontext.Background(), &types.QuerySymbolRequest{Denom: sdk.DefaultDenom})
	suite.Require().NoError(err)
	meta := res.Info
	suite.Equal(banktypes.DefaultMetadatas()[0], meta)

	meta.Issuer = "me"
	meta.Description = "test"
	app.TokenKeeper.SetTokenInfo(ctx, meta)

	res, err = queryClient.Symbol(gocontext.Background(), &types.QuerySymbolRequest{Denom: sdk.DefaultDenom})
	suite.Require().NoError(err)
	suite.Equal(meta, res.Info)

	res1, err := queryClient.Symbols(gocontext.Background(), &types.QuerySymbolsRequest{})
	suite.Require().NoError(err)
	suite.EqualValues(1, len(res1.Symbols))

	app.TokenKeeper.SetTokenInfo(ctx, btcmeta)

	res, err = queryClient.Symbol(gocontext.Background(), &types.QuerySymbolRequest{Denom: "btc"})
	suite.Require().NoError(err)
	suite.Equal(btcmeta, res.Info)

	res1, err = queryClient.Symbols(gocontext.Background(), &types.QuerySymbolsRequest{})
	suite.Require().NoError(err)
	suite.EqualValues(2, len(res1.Symbols))
	suite.Contains(res1.Symbols, "btc")

}

func (suite *IntegrationTestSuite) TestQueryParams() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	res, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	params := res.Params
	suite.Require().Equal(types.DefaultParams(), params)

	params.NewTokenFee = sdk.NewInt(123)
	params.TokenCacheSize = 10
	app.TokenKeeper.SetParams(ctx, params)

	res, err = queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(params, res.Params)
}
