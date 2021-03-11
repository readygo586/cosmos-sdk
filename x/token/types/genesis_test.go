package types_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/token/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefaultGenesisState(t *testing.T) {
	data := types.DefaultGenesisState()
	require.NoError(t, types.ValidateGenesis(*data))
	require.Equal(t, types.DefaultNewTokenFee, data.Params.NewTokenFee)
	require.Equal(t, types.DefaultTokenCacheSize, data.Params.TokenCacheSize)

	data.Params.TokenCacheSize++
	require.NotEqual(t, types.DefaultGenesisState(), data)

	data.Params.TokenCacheSize--
	require.Equal(t, types.DefaultGenesisState(), data)

	data.Params.NewTokenFee = sdk.NewInt(100)
	require.NotEqual(t, types.DefaultGenesisState(), data)
}

func TestMarshalJSONGeneisState(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	orig := types.NewGenesisState(
		types.Params{
			NewTokenFee:    sdk.NewInt(10),
			TokenCacheSize: 10,
		})
	bz, err := cdc.MarshalJSON(orig)
	require.NoError(t, err)

	got := types.GenesisState{}
	cdc.UnmarshalJSON(bz, &got)
	require.Equal(t, *orig, got)

}
