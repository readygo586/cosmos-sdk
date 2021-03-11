package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (suite *IntegrationTestSuite) TestExportGenesis() {
	app, ctx := suite.app, suite.ctx

	expectedMetadata := suite.getTestMetadata()
	expectedBalances := suite.getTestBalances()

	for i, _ := range expectedMetadata {
		app.BankKeeper.SetDenomMetaData(ctx, expectedMetadata[i])
	}
	//in suite setup, default meta will also setup
	expectedMetadata = append(expectedMetadata, types.DefaultMetadatas()...)
	for i := range []int{1, 2} {
		accAddr, err1 := sdk.AccAddressFromBech32(expectedBalances[i].Address)
		if err1 != nil {
			panic(err1)
		}
		err := app.BankKeeper.SetBalances(ctx, accAddr, expectedBalances[i].Coins)
		suite.Require().NoError(err)
	}

	totalSupply := types.NewSupplys(sdk.NewCoins(sdk.NewInt64Coin("test", 400000000)))
	app.BankKeeper.SetSupplys(ctx, totalSupply)
	app.BankKeeper.SetParams(ctx, types.DefaultParams())

	exportGenesis := app.BankKeeper.ExportGenesis(ctx)

	//suite.Require().Len(exportGenesis.Params.SendEnabled, 0)
	suite.Require().Equal(types.DefaultParams().DefaultSendEnabled, exportGenesis.Params.DefaultSendEnabled)
	suite.Require().Equal(totalSupply.GetTotal(), exportGenesis.Supply)
	suite.Require().Equal(expectedBalances, exportGenesis.Balances)
	suite.Require().Equal(types.Metadatas(expectedMetadata).Sort(), types.Metadatas(exportGenesis.DenomMetadata))
}

func (suite *IntegrationTestSuite) getTestBalances() []types.Balance {
	addr2, _ := sdk.AccAddressFromBech32("cosmos1f9xjhxm0plzrh9cskf4qee4pc2xwp0n0556gh0")
	addr1, _ := sdk.AccAddressFromBech32("cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh")
	return []types.Balance{
		{Address: addr2.String(), Coins: sdk.Coins{sdk.NewInt64Coin("testcoin1", 32), sdk.NewInt64Coin("testcoin2", 34)}},
		{Address: addr1.String(), Coins: sdk.Coins{sdk.NewInt64Coin("testcoin3", 10)}},
	}

}

func (suite *IntegrationTestSuite) TestInitGenesis() {
	app, ctx := suite.app, suite.ctx
	orig := types.DefaultGenesisState()
	app.BankKeeper.InitGenesis(ctx, orig)
	got := app.BankKeeper.ExportGenesis(ctx)
	suite.Require().Equal(true, got.DenomMetadata[0].SendEnabled)
	suite.Require().Equal(true, got.Params.DefaultSendEnabled)
//	suite.Require().Equal(orig, got)

	orig.Params.DefaultSendEnabled = false
	orig.DenomMetadata[0].SendEnabled = false
	app.BankKeeper.InitGenesis(ctx, orig)
	got = app.BankKeeper.ExportGenesis(ctx)
	suite.Require().Equal(false, got.DenomMetadata[0].SendEnabled)
	suite.Require().Equal(false, got.Params.DefaultSendEnabled)
	//suite.Require().Equal(orig, got)
}
