package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParamsEqual(t *testing.T) {
	p1 := DefaultParams()
	p2 := DefaultParams()

	require.True(t, p1.Equal(p2))

	p2.TokenCacheSize--
	require.False(t, p1.Equal(p2))

	p2.TokenCacheSize++
	require.True(t, p1.Equal(p2))
}

func TestParamsString(t *testing.T) {
	expectedStr := "Params:TokenCacheSize:32\tNewTokenFee:100000000\t"
	p := DefaultParams()

	require.Equal(t, expectedStr, p.String())
}

func TestParamsMarshalJSON(t *testing.T) {
	p := DefaultParams()
	bz, err := ModuleCdc.Amino.MarshalJSON(p)
	require.Nil(t, err)

	p1 := Params{}
	ModuleCdc.Amino.UnmarshalJSON(bz, &p1)
	require.True(t, p.Equal(p1))
}

func TestParamValidate(t *testing.T) {
	p := DefaultParams()
	require.Nil(t, p.Validate())

	p.NewTokenFee = sdk.NewInt(1)
	require.Nil(t, p.Validate())

	p.NewTokenFee = sdk.NewInt(0)
	require.NotNil(t, p.Validate())

	p.NewTokenFee = sdk.NewInt(-1)
	require.NotNil(t, p.Validate())

}
