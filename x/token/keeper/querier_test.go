package keeper_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/token/keeper"
	"github.com/cosmos/cosmos-sdk/x/token/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (suite *IntegrationTestSuite) TestQuerier_QuerySymbol() {
	app, ctx := suite.app, suite.ctx
	legacyAmino := app.LegacyAmino()

	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QuerySymbol),
		Data: []byte{},
	}

	querier := keeper.NewQuerier(app.TokenKeeper, legacyAmino)
	res, err := querier(ctx, []string{types.QuerySymbol}, req)
	suite.Require().NotNil(err)
	suite.Require().Nil(res)

	req.Data = legacyAmino.MustMarshalJSON(types.NewQuerySymbolRequest(sdk.DefaultDenom))
	res, err = querier(ctx, []string{types.QuerySymbol}, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	var meta banktypes.Metadata
	suite.Require().NoError(legacyAmino.UnmarshalJSON(res, &meta))
	suite.Require().Equal(banktypes.DefaultMetadata()[0], meta)

	app.TokenKeeper.SetTokenInfo(ctx, btcmeta)

	req.Data = legacyAmino.MustMarshalJSON(types.NewQuerySymbolRequest("btc"))
	res, err = querier(ctx, []string{types.QuerySymbol}, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	var meta1 banktypes.Metadata
	suite.Require().NoError(legacyAmino.UnmarshalJSON(res, &meta1))
	suite.Require().Equal(btcmeta, meta1)
}

func (suite *IntegrationTestSuite) TestQuerier_QuerySymbols() {
	app, ctx := suite.app, suite.ctx
	legacyAmino := app.LegacyAmino()

	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QuerySymbols),
		Data: []byte{},
	}

	querier := keeper.NewQuerier(app.TokenKeeper, legacyAmino)
	res, err := querier(ctx, []string{types.QuerySymbols}, req)
	suite.Require().NoError(err)

	var symbols []string
	suite.Require().NoError(legacyAmino.UnmarshalJSON(res, &symbols))
	suite.Require().EqualValues(1, len(symbols))
	suite.Require().EqualValues(sdk.DefaultDenom, symbols[0])

	app.TokenKeeper.SetTokenInfo(ctx, btcmeta)
	res, err = querier(ctx, []string{types.QuerySymbols}, req)
	suite.Require().NoError(err)

	suite.Require().NoError(legacyAmino.UnmarshalJSON(res, &symbols))
	suite.Require().EqualValues(2, len(symbols))
	suite.Require().Contains(symbols, "btc")
}

func (suite *IntegrationTestSuite) TestQuerier_QueryParams() {
	app, ctx := suite.app, suite.ctx
	legacyAmino := app.LegacyAmino()

	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryParameters),
		Data: []byte{},
	}

	querier := keeper.NewQuerier(app.TokenKeeper, legacyAmino)
	res, err := querier(ctx, []string{types.QueryParameters}, req)
	suite.Require().NoError(err)

	var params types.Params
	suite.Require().NoError(legacyAmino.UnmarshalJSON(res, &params))
	suite.Require().Equal(types.DefaultParams(), params)

	params.TokenCacheSize = 10
	params.NewTokenFee = sdk.NewInt(123)
	app.TokenKeeper.SetParams(ctx, params)
	res, err = querier(ctx, []string{types.QueryParameters}, req)
	suite.Require().NoError(err)

	var got types.Params
	suite.Require().NoError(legacyAmino.UnmarshalJSON(res, &got))
	suite.Require().Equal(params, got)
}
