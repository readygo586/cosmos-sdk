

package keeper_test

import (
	gocontext "context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (suite *IntegrationTestSuite) TestQueryBalance() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	_, _, addr := testdata.KeyTestPubAddr()

	_, err := queryClient.Balance(gocontext.Background(), &types.QueryBalanceRequest{})
	suite.Require().Error(err)

	_, err = queryClient.Balance(gocontext.Background(), &types.QueryBalanceRequest{Address: addr.String()})
	suite.Require().Error(err)

	req := types.NewQueryBalanceRequest(addr, fooDenom)
	res, err := queryClient.Balance(gocontext.Background(), req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.True(res.Balance.IsZero())

	origCoins := sdk.NewCoins(newFooCoin(50), newBarCoin(30))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)

	app.AccountKeeper.SetAccount(ctx, acc)
	suite.Require().NoError(app.BankKeeper.SetBalances(ctx, acc.GetAddress(), origCoins))

	res, err = queryClient.Balance(gocontext.Background(), req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.True(res.Balance.IsEqual(newFooCoin(50)))
}

func (suite *IntegrationTestSuite) TestQueryAllBalances() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	_, _, addr := testdata.KeyTestPubAddr()
	_, err := queryClient.AllBalances(gocontext.Background(), &types.QueryAllBalancesRequest{})
	suite.Require().Error(err)

	pageReq := &query.PageRequest{
		Key:        nil,
		Limit:      1,
		CountTotal: false,
	}
	req := types.NewQueryAllBalancesRequest(addr, pageReq)
	res, err := queryClient.AllBalances(gocontext.Background(), req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.True(res.Balances.IsZero())

	fooCoins := newFooCoin(50)
	barCoins := newBarCoin(30)

	origCoins := sdk.NewCoins(fooCoins, barCoins)
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)

	app.AccountKeeper.SetAccount(ctx, acc)
	suite.Require().NoError(app.BankKeeper.SetBalances(ctx, acc.GetAddress(), origCoins))

	res, err = queryClient.AllBalances(gocontext.Background(), req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Equal(res.Balances.Len(), 1)
	suite.NotNil(res.Pagination.NextKey)

	suite.T().Log("query second page with nextkey")
	pageReq = &query.PageRequest{
		Key:        res.Pagination.NextKey,
		Limit:      1,
		CountTotal: true,
	}
	req = types.NewQueryAllBalancesRequest(addr, pageReq)
	res, err = queryClient.AllBalances(gocontext.Background(), req)
	suite.Equal(res.Balances.Len(), 1)
	suite.Nil(res.Pagination.NextKey)
}

func (suite *IntegrationTestSuite) TestQueryTotalSupply() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	expectedTotalSupply := types.NewSupplys(sdk.NewCoins(sdk.NewInt64Coin("test", 400000000)))
	app.BankKeeper.SetSupplys(ctx, expectedTotalSupply)

	res, err := queryClient.TotalSupply(gocontext.Background(), &types.QueryTotalSupplyRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	suite.Require().Equal(expectedTotalSupply.Total, res.Supply)
}

func (suite *IntegrationTestSuite) TestQueryTotalSupply2() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	test1Supply := sdk.NewInt64Coin("test1", 4000000)
	test2Supply := sdk.NewInt64Coin("test2", 700000000)
	expectedTotalSupply := types.NewSupplys(sdk.NewCoins(test1Supply, test2Supply))
	app.BankKeeper.SetSupplys(ctx, expectedTotalSupply)

	res, err := queryClient.TotalSupply(gocontext.Background(), &types.QueryTotalSupplyRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(expectedTotalSupply.GetTotal(), res.Supply)
}

func (suite *IntegrationTestSuite) TestQueryTotalSupply3() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	test1Supply := sdk.NewInt64Coin("test1", 4000000)
	test2Supply := sdk.NewInt64Coin("test2", 700000000)
	expectedTotalSupply := types.NewSupplys(sdk.NewCoins(test1Supply, test2Supply))
	app.BankKeeper.SetSupplys(ctx, expectedTotalSupply)

	res, err := queryClient.TotalSupply(gocontext.Background(), &types.QueryTotalSupplyRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(expectedTotalSupply.GetTotal(), res.Supply)

	test3Supply := sdk.NewInt64Coin("test3", 8000000)
	expectedTotalSupply = types.NewSupplys(sdk.NewCoins(test1Supply, test2Supply, test3Supply))
	app.BankKeeper.SetSupply(ctx, types.NewSupply(test3Supply))

	res, err = queryClient.TotalSupply(gocontext.Background(), &types.QueryTotalSupplyRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(expectedTotalSupply.GetTotal(), res.Supply)

	test1Supply = sdk.NewInt64Coin("test1", 2000000)
	expectedTotalSupply = types.NewSupplys(sdk.NewCoins(test1Supply, test2Supply, test3Supply))
	app.BankKeeper.SetSupply(ctx, types.NewSupply(test1Supply))

	res, err = queryClient.TotalSupply(gocontext.Background(), &types.QueryTotalSupplyRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(expectedTotalSupply.GetTotal(), res.Supply)
}


func (suite *IntegrationTestSuite) TestQueryTotalSupplyOf() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	test1Supply := sdk.NewInt64Coin("test1", 4000000)
	test2Supply := sdk.NewInt64Coin("test2", 700000000)
	expectedTotalSupply := types.NewSupplys(sdk.NewCoins(test1Supply, test2Supply))
	app.BankKeeper.SetSupplys(ctx, expectedTotalSupply)

        total, err := queryClient.TotalSupply(gocontext.Background(), &types.QueryTotalSupplyRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(expectedTotalSupply.GetTotal(), total.Supply)
	res, err := queryClient.SupplyOf(gocontext.Background(), &types.QuerySupplyOfRequest{})
	suite.Require().Error(err)
	suite.Require().Nil(res)

	res, err = queryClient.SupplyOf(gocontext.Background(), &types.QuerySupplyOfRequest{Denom: test1Supply.Denom})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	suite.Require().Equal(test1Supply, res.Amount)
	
	res, err = queryClient.SupplyOf(gocontext.Background(), &types.QuerySupplyOfRequest{Denom: test2Supply.Denom})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(test2Supply, res.Amount)
	
	res, err = queryClient.SupplyOf(gocontext.Background(), &types.QuerySupplyOfRequest{Denom: "test3"})
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewCoin("test3", sdk.ZeroInt()), res.Amount)
}

func (suite *IntegrationTestSuite) TestQueryParams() {
	res, err := suite.queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(suite.app.BankKeeper.GetParams(suite.ctx), res.GetParams())
}

func (suite *IntegrationTestSuite) TestQueryDenomsMetadataRequest() {
	var (
		req         *types.QueryDenomsMetadataRequest
		expMetadata = []types.Metadata{}
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty pagination",
			func() {
				expMetadata = []types.Metadata{}
				expMetadata = append(expMetadata,types.DefaultMetadatas()...)
				req = &types.QueryDenomsMetadataRequest{}
			},
			true,
		},
		{
			"success, no results",
			func() {
				expMetadata = []types.Metadata{}
				expMetadata = append(expMetadata,types.DefaultMetadatas()...)
				req = &types.QueryDenomsMetadataRequest{
					Pagination: &query.PageRequest{
						Limit:      3,
						CountTotal: true,
					},
				}
			},
			true,
		},
		{
			"success",
			func() {
				expMetadata = []types.Metadata{}
				metadataAtom := types.Metadata{
					Description: "The native staking token of the Cosmos Hub.",
					DenomUnits: []*types.DenomUnit{
						{
							Denom:    "uatom",
							Exponent: 0,
							Aliases:  []string{"microatom"},
						},
						{
							Denom:    "atom",
							Exponent: 6,
							Aliases:  []string{"ATOM"},
						},
					},
					Base:    "uatom",
					Display: "atom",
					Decimals: 18,
					SendEnabled: false,
				}

				metadataEth := types.Metadata{
					Description: "Ethereum native token",
					DenomUnits: []*types.DenomUnit{
						{
							Denom:    "wei",
							Exponent: 0,
						},
						{
							Denom:    "eth",
							Exponent: 18,
							Aliases:  []string{"ETH", "ether"},
						},
					},
					Base:    "eth",
					Display: "eth",
					Decimals: 18,
					SendEnabled: true,
				}

				metadataTrx := types.Metadata{
					Description: "trx native token",
					Base:    "trx",
					Display: "trx",
					Issuer: "Sun",
					Decimals: 18,
					SendEnabled: true,

				}

				metadataBhc := types.Metadata{
					Description: "bhc native token",
					Base:    "bhc",
					Display: "bhc",
					Issuer: "bhc",
					Decimals: 18,
					SendEnabled: true,

				}

				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, metadataAtom)
				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, metadataEth)
				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, metadataTrx)
				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, metadataBhc)
				expMetadata = append(expMetadata, types.DefaultMetadatas()...)
				expMetadata = append(expMetadata,metadataAtom, metadataEth,metadataTrx,metadataBhc)
				expMetadata = types.Metadatas(expMetadata).Sort()
				req = &types.QueryDenomsMetadataRequest{
					Pagination: &query.PageRequest{
						Limit:      7,
						CountTotal: true,
					},
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.ctx)

			res, err := suite.queryClient.DenomsMetadata(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expMetadata, res.Metadatas)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) QueryDenomMetadataRequest() {
	var (
		req         *types.QueryDenomMetadataRequest
		expMetadata = types.Metadata{}
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty denom",
			func() {
				req = &types.QueryDenomMetadataRequest{}
			},
			false,
		},
		{
			"not found denom",
			func() {
				req = &types.QueryDenomMetadataRequest{
					Denom: "foo",
				}
			},
			false,
		},
		{
			"success",
			func() {
				expMetadata := types.Metadata{
					Description: "The native staking token of the Cosmos Hub.",
					DenomUnits: []*types.DenomUnit{
						{
							Denom:    "uatom",
							Exponent: 0,
							Aliases:  []string{"microatom"},
						},
						{
							Denom:    "atom",
							Exponent: 6,
							Aliases:  []string{"ATOM"},
						},
					},
					Base:    "uatom",
					Display: "atom",
				}

				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, expMetadata)
				req = &types.QueryDenomMetadataRequest{
					Denom: expMetadata.Base,
				}
			},
			true,
		},
		{
			"success without denomuint",
			func() {
				expMetadata := types.Metadata{
					Description: "eth",
					Base:    "eth",
					Display: "eth",
					Decimals: 18,
					SendEnabled: true,
					Issuer: "Viltak",
				}

				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, expMetadata)
				req = &types.QueryDenomMetadataRequest{
					Denom: expMetadata.Base,
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.ctx)

			res, err := suite.queryClient.DenomMetadata(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expMetadata, res.Metadata)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
