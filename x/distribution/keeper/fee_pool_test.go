package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

func TestAddCoinsFromAccountToFeePool(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	initAmout := sdk.NewInt(1000000)
	addrs := simapp.AddTestAddrs(app, ctx, 2, initAmout)

	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, initAmout))
	coins1 := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(1)))
	coins2 := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(2)))
	coins3 := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(3)))
	coins4 := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(4)))

	coins6 := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(6)))
	coins10 := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(10)))

	require.Equal(t, sdk.DecCoins(nil), app.DistrKeeper.GetFeePoolCommunityCoins(ctx))
	require.Equal(t, sdk.Coins{}, app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName)))
	err := app.DistrKeeper.AddCoinsFromAccountToFeePool(ctx, addrs[0], coins1)
	require.NoError(t, err)
	require.Equal(t, sdk.NewDecCoinsFromCoins(coins1...), app.DistrKeeper.GetFeePoolCommunityCoins(ctx))
	require.Equal(t, coins1, app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName)))
	require.Equal(t, initCoins.Sub(coins1), app.BankKeeper.GetAllBalances(ctx, addrs[0]))

	err = app.DistrKeeper.AddCoinsFromAccountToFeePool(ctx, addrs[0], coins2)
	require.NoError(t, err)
	require.Equal(t, sdk.NewDecCoinsFromCoins(coins3...), app.DistrKeeper.GetFeePoolCommunityCoins(ctx))
	require.Equal(t, coins3, app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName)))

	require.Equal(t, initCoins.Sub(coins3), app.BankKeeper.GetAllBalances(ctx, addrs[0]))

	err = app.DistrKeeper.AddCoinsFromAccountToFeePool(ctx, addrs[0], coins3)
	require.NoError(t, err)
	require.Equal(t, sdk.NewDecCoinsFromCoins(coins6...), app.DistrKeeper.GetFeePoolCommunityCoins(ctx))
	require.Equal(t, coins6, app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName)))

	require.Equal(t, initCoins.Sub(coins6), app.BankKeeper.GetAllBalances(ctx, addrs[0]))

	err = app.DistrKeeper.AddCoinsFromAccountToFeePool(ctx, addrs[0], coins4)
	require.NoError(t, err)
	require.Equal(t, sdk.NewDecCoinsFromCoins(coins10...), app.DistrKeeper.GetFeePoolCommunityCoins(ctx))
	require.Equal(t, coins10, app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName)))
	require.Equal(t, initCoins.Sub(coins10), app.BankKeeper.GetAllBalances(ctx, addrs[0]))
}

func TestDistributeFromFeePool(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	initAmout := sdk.NewInt(1000000)
	addrs := simapp.AddTestAddrs(app, ctx, 2, initAmout)

	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, initAmout))
	coins1 := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(1)))
	coins2 := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(2)))
	coins3 := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(3)))
	coins4 := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(4)))

	coins6 := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(6)))
	coins10 := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(10)))

	require.Equal(t, sdk.DecCoins(nil), app.DistrKeeper.GetFeePoolCommunityCoins(ctx))
	err := app.DistrKeeper.AddCoinsFromAccountToFeePool(ctx, addrs[0], initCoins)
	require.NoError(t, err)
	require.Equal(t, sdk.NewDecCoinsFromCoins(initCoins...), app.DistrKeeper.GetFeePoolCommunityCoins(ctx))
	require.Equal(t, initCoins.Sub(initCoins), app.BankKeeper.GetAllBalances(ctx, addrs[0]))

	err = app.DistrKeeper.DistributeFromFeePool(ctx, coins1, addrs[0])
	require.NoError(t, err)
	require.Equal(t, sdk.NewDecCoinsFromCoins(initCoins.Sub(coins1)...), app.DistrKeeper.GetFeePoolCommunityCoins(ctx))
	require.Equal(t, coins1, app.BankKeeper.GetAllBalances(ctx, addrs[0]))

	err = app.DistrKeeper.DistributeFromFeePool(ctx, coins2, addrs[0])
	require.NoError(t, err)
	require.Equal(t, sdk.NewDecCoinsFromCoins(initCoins.Sub(coins3)...), app.DistrKeeper.GetFeePoolCommunityCoins(ctx))
	require.Equal(t, coins3, app.BankKeeper.GetAllBalances(ctx, addrs[0]))

	err = app.DistrKeeper.DistributeFromFeePool(ctx, coins3, addrs[0])
	require.NoError(t, err)
	require.Equal(t, sdk.NewDecCoinsFromCoins(initCoins.Sub(coins6)...), app.DistrKeeper.GetFeePoolCommunityCoins(ctx))
	require.Equal(t, coins6, app.BankKeeper.GetAllBalances(ctx, addrs[0]))

	err = app.DistrKeeper.DistributeFromFeePool(ctx, coins4, addrs[0])
	require.NoError(t, err)
	require.Equal(t, sdk.NewDecCoinsFromCoins(initCoins.Sub(coins10)...), app.DistrKeeper.GetFeePoolCommunityCoins(ctx))
	require.Equal(t, coins10, app.BankKeeper.GetAllBalances(ctx, addrs[0]))

	coins100 := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(100)))
	err = app.DistrKeeper.DistributeFromFeePool(ctx, coins100, addrs[1])
	require.NoError(t, err)
	require.Equal(t, sdk.NewDecCoinsFromCoins(initCoins.Sub(coins100).Sub(coins10)...), app.DistrKeeper.GetFeePoolCommunityCoins(ctx))
	require.Equal(t, coins10, app.BankKeeper.GetAllBalances(ctx, addrs[0]))
	require.Equal(t, initCoins.Add(coins100...), app.BankKeeper.GetAllBalances(ctx, addrs[1]))
}
