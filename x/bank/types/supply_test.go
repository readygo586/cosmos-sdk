package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestSupplysMarshalYAML(t *testing.T) {
	supply := DefaultSupplys()
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt()))
	supply.Inflate(coins)

	bz, err := yaml.Marshal(supply)
	require.NoError(t, err)
	bzCoins, err := yaml.Marshal(coins)
	require.NoError(t, err)

	want := fmt.Sprintf(`total:
%s`, string(bzCoins))

	require.Equal(t, want, string(bz))
	require.Equal(t, want, supply.String())
}
func TestSupplyMarshalYAML(t *testing.T) {
	supply := DefaultSupply()
	coin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())
	supply.Inflate(coin)
	bz, err := yaml.Marshal(supply)
	require.NoError(t, err)
	_, err = yaml.Marshal(coin)
	require.NoError(t, err)
	want := "total:\n  denom: stake\n  amount: \"1\"\n"
	require.Equal(t, want, string(bz))
	require.Equal(t, want, supply.String())
}
func TestSupplys_SetGet(t *testing.T) {
	origCoins := sdk.NewCoins(sdk.NewCoin("test1", sdk.NewInt(100)), sdk.NewCoin("test2", sdk.NewInt(200)))
	supply := NewSupplys(origCoins)
	require.Equal(t, origCoins, supply.GetTotal())
	coins1 := sdk.NewCoins(sdk.NewCoin("test1", sdk.NewInt(1000)), sdk.NewCoin("test2", sdk.NewInt(2000)))
	supply.SetTotal(coins1)
	require.Equal(t, coins1, supply.GetTotal())
	coins2 := sdk.NewCoins(sdk.NewCoin("test3", sdk.NewInt(1000)), sdk.NewCoin("test4", sdk.NewInt(2000)))
	supply.SetTotal(coins2)
	require.Equal(t, coins2, supply.GetTotal())
	coins3 := sdk.NewCoins()
	supply.SetTotal(coins3)
	require.Equal(t, coins3, supply.GetTotal())
}
func TestSupplys_InfalteDeflate(t *testing.T) {
	origCoins := sdk.NewCoins(sdk.NewCoin("test1", sdk.NewInt(100)), sdk.NewCoin("test2", sdk.NewInt(200)))
	supply := NewSupplys(origCoins)
	require.Equal(t, origCoins, supply.GetTotal())
	coins1 := origCoins.Add(sdk.NewCoin("test1", sdk.NewInt(300)))
	supply.Inflate(sdk.NewCoins(sdk.NewCoin("test1", sdk.NewInt(300))))
	require.Equal(t, coins1, supply.GetTotal())
	coins2 := coins1.Add(sdk.NewCoin("test2", sdk.NewInt(400)))
	supply.Inflate(sdk.NewCoins(sdk.NewCoin("test2", sdk.NewInt(400))))
	require.Equal(t, coins2, supply.GetTotal())
	coins3 := coins2.Add(sdk.NewCoin("test1", sdk.NewInt(500)), sdk.NewCoin("test2", sdk.NewInt(600)))
	supply.Inflate(sdk.NewCoins(sdk.NewCoin("test1", sdk.NewInt(500)), sdk.NewCoin("test2", sdk.NewInt(600))))
	require.Equal(t, coins3, supply.GetTotal())
	supply.Deflate(sdk.NewCoins(sdk.NewCoin("test1", sdk.NewInt(500)), sdk.NewCoin("test2", sdk.NewInt(600))))
	require.Equal(t, coins2, supply.GetTotal())
	supply.Deflate(sdk.NewCoins(sdk.NewCoin("test2", sdk.NewInt(400))))
	require.Equal(t, coins1, supply.GetTotal())
	supply.Deflate(sdk.NewCoins(sdk.NewCoin("test1", sdk.NewInt(300))))
	require.Equal(t, origCoins, supply.GetTotal())
	supply.Deflate(sdk.NewCoins(sdk.NewCoin("test2", sdk.NewInt(200))))
	require.Equal(t, sdk.NewCoins(sdk.NewCoin("test1", sdk.NewInt(100))), supply.GetTotal())
	supply.Deflate(sdk.NewCoins(sdk.NewCoin("test1", sdk.NewInt(100))))
	require.EqualValues(t, sdk.Coins(nil), supply.GetTotal())
}
func TestSupplys_InfalteDeflate2(t *testing.T) {
	origCoins := sdk.NewCoins(sdk.NewCoin("test1", sdk.NewInt(100)), sdk.NewCoin("test2", sdk.NewInt(200)))
	supply := NewSupplys(origCoins)
	coins1 := origCoins.Add(sdk.NewCoin("test3", sdk.NewInt(300)))
	supply.Inflate(sdk.NewCoins(sdk.NewCoin("test3", sdk.NewInt(300))))
	require.Equal(t, coins1, supply.GetTotal())
	supply.Deflate(sdk.NewCoins(sdk.NewCoin("test3", sdk.NewInt(300))))
	require.Equal(t, origCoins, supply.GetTotal())
}
func TestSupplys_InfalteDeflate3(t *testing.T) {
	origCoins := sdk.NewCoins(sdk.NewCoin("test1", sdk.NewInt(100)), sdk.NewCoin("test2", sdk.NewInt(200)))
	supply := NewSupplys(origCoins)
	require.Panics(t, func() { supply.Deflate(sdk.NewCoins(sdk.NewCoin("test3", sdk.NewInt(300)))) })
}
func TestSupply_SetGet(t *testing.T) {
	origCoin := sdk.NewCoin("test1", sdk.NewInt(100))
	supply := NewSupply(origCoin)
	coin1 := sdk.NewCoin("test1", sdk.NewInt(1000))
	supply.SetTotal(coin1)
	require.Equal(t, coin1, supply.GetTotal())
	coin2 := sdk.NewCoin("test2", sdk.NewInt(100))
	supply.SetTotal(coin2)
	require.Equal(t, coin2, supply.GetTotal())
	coin3 := sdk.Coin{}
	supply.SetTotal(coin3)
	require.Equal(t, coin3, supply.GetTotal())
}
func TestSupply_InfalteDeflate(t *testing.T) {
	origCoin := sdk.NewCoin("test1", sdk.NewInt(100))
	supply := NewSupply(origCoin)
	coin1 := origCoin.Add(sdk.NewCoin("test1", sdk.NewInt(200)))
	supply.Inflate(sdk.NewCoin("test1", sdk.NewInt(200)))
	require.Equal(t, coin1, supply.GetTotal())
	coin2 := coin1.Add(sdk.NewCoin("test1", sdk.NewInt(300)))
	supply.Inflate(sdk.NewCoin("test1", sdk.NewInt(300)))
	require.Equal(t, coin2, supply.GetTotal())
	supply.Deflate(sdk.NewCoin("test1", sdk.NewInt(300)))
	require.Equal(t, coin1, supply.GetTotal())
	supply.Deflate(sdk.NewCoin("test1", sdk.NewInt(200)))
	require.Equal(t, origCoin, supply.GetTotal())
	supply.Deflate(sdk.NewCoin("test1", sdk.NewInt(100)))
	require.Equal(t, "test1", supply.GetTotal().Denom)
	require.EqualValues(t, int64(0), supply.GetTotal().Amount.Int64())
}
func TestSupply_InfalteDeflate2(t *testing.T) {
	origCoin := sdk.NewCoin("test1", sdk.NewInt(100))
	supply := NewSupply(origCoin)
	require.Panics(t, func() { supply.Deflate(sdk.NewCoin("test1", sdk.NewInt(101))) })
	require.Panics(t, func() { supply.Inflate(sdk.NewCoin("test2", sdk.NewInt(200))) })
	require.Panics(t, func() { supply.Deflate(sdk.NewCoin("test2", sdk.NewInt(200))) })
}
