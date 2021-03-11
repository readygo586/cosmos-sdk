package token_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/token/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

var (
	priv1 = secp256k1.GenPrivKey()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())
	priv2 = secp256k1.GenPrivKey()
	addr2 = sdk.AccAddress(priv2.PubKey().Address())
	addr3 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	priv4 = secp256k1.GenPrivKey()
	addr4 = sdk.AccAddress(priv4.PubKey().Address())
)

func TestNewTokenSuccess1(t *testing.T) {
	acc := &authtypes.BaseAccount{
		Address: addr1.String(),
	}

	genAccs := []authtypes.GenesisAccount{acc}
	initBalance := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(100)))
	app := simapp.SetupWithGenesisAccounts(genAccs, banktypes.Balance{addr1.String(), initBalance})
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	app.TokenKeeper.SetParams(ctx, types.Params{
		TokenCacheSize: 10,
		NewTokenFee:    sdk.NewInt(10),
	})

	app.Commit()

	res1 := app.AccountKeeper.GetAccount(ctx, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*authtypes.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	btcAmount := sdk.NewCoin("btc", sdk.NewInt(100))
	newMsg := types.NewMsgNewToken(addr1, addr1, 18, btcAmount)
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeTestEncodingConfig().TxConfig
	_, _, err := simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{newMsg}, "", []uint64{origAccNum}, []uint64{origSeq}, true, true, priv1)
	require.NoError(t, err)

	newTokenFee := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, app.TokenKeeper.GetParams(ctx).NewTokenFee))
	balances := initBalance.Sub(newTokenFee).Add(btcAmount)
	simapp.CheckBalance(t, app, addr1, balances)

	//feeCollect can be omit because initBalance is too small
	simapp.CheckBalance(t, app, app.AccountKeeper.GetModuleAddress(distrtypes.ModuleName), newTokenFee)
	feePoolAmount, _ := app.DistrKeeper.GetFeePoolCommunityCoins(ctx).TruncateDecimal()
	require.Equal(t, feePoolAmount, newTokenFee)

	res2 := app.AccountKeeper.GetAccount(app.NewContext(true, tmproto.Header{}), addr1)
	require.NotNil(t, res2)

	require.Equal(t, res2.GetAccountNumber(), origAccNum)
	require.Equal(t, res2.GetSequence(), origSeq+1)
}

func TestNewTokenSuccess2(t *testing.T) {
	acc := &authtypes.BaseAccount{
		Address: addr1.String(),
	}

	genAccs := []authtypes.GenesisAccount{acc}
	initBalance := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(100000)))
	app := simapp.SetupWithGenesisAccounts(genAccs, banktypes.Balance{addr1.String(), initBalance})
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.Commit()

	res1 := app.AccountKeeper.GetAccount(ctx, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*authtypes.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	btcAmount := sdk.NewCoin("btc", sdk.NewInt(100))
	newMsg := types.NewMsgNewToken(addr1, addr1, 18, btcAmount)
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeTestEncodingConfig().TxConfig
	_, _, err := simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{newMsg}, "", []uint64{origAccNum}, []uint64{origSeq}, true, true, priv1)
	require.NoError(t, err)

	newTokenFee := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, app.TokenKeeper.GetParams(ctx).NewTokenFee))
	balances := initBalance.Sub(newTokenFee).Add(btcAmount)
	simapp.CheckBalance(t, app, addr1, balances)

	//feeCollected can not be omitted
	require.True(t, app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(distrtypes.ModuleName)).IsAllGT((newTokenFee)))
	//simapp.CheckBalance(t, app, app.AccountKeeper.GetModuleAddress(distrtypes.ModuleName), newTokenFee)
	feePoolAmount, _ := app.DistrKeeper.GetFeePoolCommunityCoins(ctx).TruncateDecimal()
	require.True(t, feePoolAmount.IsAllGT(newTokenFee))
	require.Equal(t, app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(distrtypes.ModuleName)), feePoolAmount)

	res2 := app.AccountKeeper.GetAccount(app.NewContext(true, tmproto.Header{}), addr1)
	require.NotNil(t, res2)

	require.Equal(t, res2.GetAccountNumber(), origAccNum)
	require.Equal(t, res2.GetSequence(), origSeq+1)
}

func TestNewTokenNotEnoughFee(t *testing.T) {
	acc := &authtypes.BaseAccount{
		Address: addr1.String(),
	}

	genAccs := []authtypes.GenesisAccount{acc}
	initBalance := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(10))) //defaultFee is sdk.TokensFromConsensusPower(100))
	app := simapp.SetupWithGenesisAccounts(genAccs, banktypes.Balance{addr1.String(), initBalance})
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.Commit()

	res1 := app.AccountKeeper.GetAccount(ctx, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*authtypes.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	btcAmount := sdk.NewCoin("btc", sdk.NewInt(100))
	newMsg := types.NewMsgNewToken(addr1, addr1, 18, btcAmount)
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeTestEncodingConfig().TxConfig
	_, _, err := simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{newMsg}, "", []uint64{origAccNum}, []uint64{origSeq}, false, false, priv1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "insufficient funds")

	//newtoken Fee is not deducted
	simapp.CheckBalance(t, app, addr1, initBalance)
	res2 := app.AccountKeeper.GetAccount(app.NewContext(true, tmproto.Header{}), addr1)
	require.NotNil(t, res2)
	require.Equal(t, res2.GetAccountNumber(), origAccNum)
	require.Equal(t, res2.GetSequence(), origSeq+1)
}

func TestNewTokenDenomAlreadyExist(t *testing.T) {
	acc := &authtypes.BaseAccount{
		Address: addr1.String(),
	}

	genAccs := []authtypes.GenesisAccount{acc}
	initBalance := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(10000)))
	app := simapp.SetupWithGenesisAccounts(genAccs, banktypes.Balance{addr1.String(), initBalance})
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.Commit()

	res1 := app.AccountKeeper.GetAccount(ctx, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*authtypes.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	btcAmount := sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(100))
	newMsg := types.NewMsgNewToken(addr1, addr1, 18, btcAmount)
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeTestEncodingConfig().TxConfig
	_, _, err := simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{newMsg}, "", []uint64{origAccNum}, []uint64{origSeq}, false, false, priv1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "denom already exist")

	//newtoken Fee is nit deducted
	simapp.CheckBalance(t, app, addr1, initBalance)
	res2 := app.AccountKeeper.GetAccount(app.NewContext(true, tmproto.Header{}), addr1)
	require.NotNil(t, res2)
	require.Equal(t, res2.GetAccountNumber(), origAccNum)
	require.Equal(t, res2.GetSequence(), origSeq+1)
}

func TestNewTokenDenomAlreadyExist1(t *testing.T) {
	acc := &authtypes.BaseAccount{
		Address: addr1.String(),
	}

	genAccs := []authtypes.GenesisAccount{acc}
	initBalance := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(10000)))
	app := simapp.SetupWithGenesisAccounts(genAccs, banktypes.Balance{addr1.String(), initBalance})
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.Commit()

	res1 := app.AccountKeeper.GetAccount(ctx, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*authtypes.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	//first, new btc token, will success
	btcAmount := sdk.NewCoin("btc", sdk.NewInt(100))
	newMsg := types.NewMsgNewToken(addr1, addr1, 18, btcAmount)
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeTestEncodingConfig().TxConfig
	_, _, err := simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{newMsg}, "", []uint64{origAccNum}, []uint64{origSeq}, true, true, priv1)
	require.NoError(t, err)

	newTokenFee := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, app.TokenKeeper.GetParams(ctx).NewTokenFee))
	balances := initBalance.Sub(newTokenFee).Add(btcAmount)
	simapp.CheckBalance(t, app, addr1, balances)

	//new btc again, will fail
	header = tmproto.Header{Height: app.LastBlockHeight() + 1}
	origSeq++
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{newMsg}, "", []uint64{origAccNum}, []uint64{origSeq}, false, false, priv1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "denom already exist")
}

func TestInflateTokenSuccess(t *testing.T) {
	acc := &authtypes.BaseAccount{
		Address: addr1.String(),
	}

	genAccs := []authtypes.GenesisAccount{acc}
	initBalance := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(10000)))
	app := simapp.SetupWithGenesisAccounts(genAccs, banktypes.Balance{addr1.String(), initBalance})
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.Commit()

	res1 := app.AccountKeeper.GetAccount(ctx, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*authtypes.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	btcAmount := sdk.NewCoin("btc", sdk.NewInt(100))
	newMsg := types.NewMsgNewToken(addr1, addr1, 18, btcAmount)
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeTestEncodingConfig().TxConfig
	_, _, err := simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{newMsg}, "", []uint64{origAccNum}, []uint64{origSeq}, true, true, priv1)
	require.NoError(t, err)

	newTokenFee := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, app.TokenKeeper.GetParams(ctx).NewTokenFee))
	balances := initBalance.Sub(newTokenFee).Add(btcAmount)
	simapp.CheckBalance(t, app, addr1, balances)

	header = tmproto.Header{Height: app.LastBlockHeight() + 1}
	origSeq++
	inflateMsg := types.NewMsgInflateToken(addr1, addr1, sdk.NewCoin("btc", sdk.NewInt(10)))
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{inflateMsg}, "", []uint64{origAccNum}, []uint64{origSeq}, true, true, priv1)
	require.NoError(t, err)
	inflateBalances := initBalance.Sub(newTokenFee).Add(sdk.NewCoin("btc", sdk.NewInt(110)))
	simapp.CheckBalance(t, app, addr1, inflateBalances)
}

func TestInflateTokenFailDenomNotExists(t *testing.T) {
	acc := &authtypes.BaseAccount{
		Address: addr1.String(),
	}

	genAccs := []authtypes.GenesisAccount{acc}
	initBalance := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(10000)))
	app := simapp.SetupWithGenesisAccounts(genAccs, banktypes.Balance{addr1.String(), initBalance})
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.Commit()

	res1 := app.AccountKeeper.GetAccount(ctx, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*authtypes.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	btcAmount := sdk.NewCoin("btc", sdk.NewInt(100))
	inflateMsg := types.NewMsgInflateToken(addr1, addr1, btcAmount)
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeTestEncodingConfig().TxConfig
	_, _, err := simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{inflateMsg}, "", []uint64{origAccNum}, []uint64{origSeq}, false, false, priv1)
	require.Error(t, err)

	simapp.CheckBalance(t, app, addr1, initBalance)
}

func TestInflateFailNotIssuer(t *testing.T) {
	acc1 := &authtypes.BaseAccount{
		Address: addr1.String(),
	}

	acc2 := &authtypes.BaseAccount{
		Address: addr2.String(),
	}

	genAccs := []authtypes.GenesisAccount{acc1, acc2}
	initBalance := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(10000)))
	app := simapp.SetupWithGenesisAccounts(genAccs, banktypes.Balance{addr1.String(), initBalance})
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.Commit()

	res1 := app.AccountKeeper.GetAccount(ctx, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc1, res1.(*authtypes.BaseAccount))

	origAccNum1 := res1.GetAccountNumber()
	origSeq1 := res1.GetSequence()

	btcAmount := sdk.NewCoin("btc", sdk.NewInt(100))
	newMsg := types.NewMsgNewToken(addr1, addr1, 18, btcAmount)
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeTestEncodingConfig().TxConfig
	_, _, err := simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{newMsg}, "", []uint64{origAccNum1}, []uint64{origSeq1}, true, true, priv1)
	require.NoError(t, err)

	newTokenFee := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, app.TokenKeeper.GetParams(ctx).NewTokenFee))
	balances := initBalance.Sub(newTokenFee).Add(btcAmount)
	simapp.CheckBalance(t, app, addr1, balances)

	header = tmproto.Header{Height: app.LastBlockHeight() + 1}
	res2 := app.AccountKeeper.GetAccount(ctx, addr2)
	require.NotNil(t, res1)
	require.Equal(t, acc1, res1.(*authtypes.BaseAccount))

	origAccNum2 := res2.GetAccountNumber()
	origSeq2 := res2.GetSequence()
	inflateMsg := types.NewMsgInflateToken(addr2, addr2, sdk.NewCoin("btc", sdk.NewInt(10)))
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{inflateMsg}, "", []uint64{origAccNum2}, []uint64{origSeq2}, false, false, priv2)
	require.Error(t, err)
	require.Contains(t, err.Error(), "is not authorized to inflate")
}

func TestBurnTokenSuccess(t *testing.T) {
	acc := &authtypes.BaseAccount{
		Address: addr1.String(),
	}

	genAccs := []authtypes.GenesisAccount{acc}
	initBalance := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(10000)))
	app := simapp.SetupWithGenesisAccounts(genAccs, banktypes.Balance{addr1.String(), initBalance})
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.Commit()

	res1 := app.AccountKeeper.GetAccount(ctx, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*authtypes.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	burnMsg := types.NewMsgBurnToken(addr1, sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(100)))
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeTestEncodingConfig().TxConfig
	_, _, err := simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{burnMsg}, "", []uint64{origAccNum}, []uint64{origSeq}, true, true, priv1)
	require.NoError(t, err)
	simapp.CheckBalance(t, app, addr1, sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(9900))))
}

func TestBurnTokenSuccess1(t *testing.T) {
	acc := &authtypes.BaseAccount{
		Address: addr1.String(),
	}

	genAccs := []authtypes.GenesisAccount{acc}
	initBalance := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(10000)))
	app := simapp.SetupWithGenesisAccounts(genAccs, banktypes.Balance{addr1.String(), initBalance})
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.Commit()

	res1 := app.AccountKeeper.GetAccount(ctx, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*authtypes.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	burnMsg := types.NewMsgBurnToken(addr1, sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(10000)))
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeTestEncodingConfig().TxConfig
	_, _, err := simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{burnMsg}, "", []uint64{origAccNum}, []uint64{origSeq}, true, true, priv1)
	require.NoError(t, err)
	simapp.CheckBalance(t, app, addr1, sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(0))))
}

func TestBurnTokenSuccess2(t *testing.T) {
	acc := &authtypes.BaseAccount{
		Address: addr1.String(),
	}

	genAccs := []authtypes.GenesisAccount{acc}
	initBalance := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(100000)))
	app := simapp.SetupWithGenesisAccounts(genAccs, banktypes.Balance{addr1.String(), initBalance})
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.Commit()

	res1 := app.AccountKeeper.GetAccount(ctx, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*authtypes.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	btcAmount := sdk.NewCoin("btc", sdk.NewInt(100))
	newMsg := types.NewMsgNewToken(addr1, addr1, 18, btcAmount)
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeTestEncodingConfig().TxConfig
	_, _, err := simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{newMsg}, "", []uint64{origAccNum}, []uint64{origSeq}, true, true, priv1)
	require.NoError(t, err)

	newTokenFee := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, app.TokenKeeper.GetParams(ctx).NewTokenFee))
	balances := initBalance.Sub(newTokenFee).Add(btcAmount)
	simapp.CheckBalance(t, app, addr1, balances)

	//feeCollected can not be omitted
	require.True(t, app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(distrtypes.ModuleName)).IsAllGT((newTokenFee)))
	//simapp.CheckBalance(t, app, app.AccountKeeper.GetModuleAddress(distrtypes.ModuleName), newTokenFee)
	feePoolAmount, _ := app.DistrKeeper.GetFeePoolCommunityCoins(ctx).TruncateDecimal()
	require.True(t, feePoolAmount.IsAllGT(newTokenFee))
	require.Equal(t, app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(distrtypes.ModuleName)), feePoolAmount)

	res2 := app.AccountKeeper.GetAccount(app.NewContext(true, tmproto.Header{}), addr1)
	require.NotNil(t, res2)

	require.Equal(t, res2.GetAccountNumber(), origAccNum)
	require.Equal(t, res2.GetSequence(), origSeq+1)

	origSeq++
	burnMsg := types.NewMsgBurnToken(addr1, sdk.NewCoin("btc", sdk.NewInt(10)))
	header = tmproto.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{burnMsg}, "", []uint64{origAccNum}, []uint64{origSeq}, true, true, priv1)
	require.NoError(t, err)
	balances = balances.Sub(sdk.NewCoins(sdk.NewCoin("btc", sdk.NewInt(10))))
	simapp.CheckBalance(t, app, addr1, balances)
}

func TestBurnTokenFailNotEnough(t *testing.T) {
	acc := &authtypes.BaseAccount{
		Address: addr1.String(),
	}

	genAccs := []authtypes.GenesisAccount{acc}
	initBalance := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(10000)))
	app := simapp.SetupWithGenesisAccounts(genAccs, banktypes.Balance{addr1.String(), initBalance})
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.Commit()

	res1 := app.AccountKeeper.GetAccount(ctx, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*authtypes.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	burnMsg := types.NewMsgBurnToken(addr1, sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(10001)))
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeTestEncodingConfig().TxConfig
	_, _, err := simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{burnMsg}, "", []uint64{origAccNum}, []uint64{origSeq}, false, false, priv1)
	require.Error(t, err)
	simapp.CheckBalance(t, app, addr1, sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, sdk.TokensFromConsensusPower(10000))))
}
