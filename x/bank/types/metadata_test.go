package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

func TestMetadataValidate(t *testing.T) {
	testCases := []struct {
		name     string
		metadata types.Metadata
		expErr   bool
	}{
		{
			"non-empty coins",
			types.Metadata{
				Description: "The native staking token of the Cosmos Hub.",
				DenomUnits: []*types.DenomUnit{
					{"uatom", uint32(0), []string{"microatom"}},
					{"matom", uint32(3), []string{"milliatom"}},
					{"atom", uint32(6), nil},
				},
				Base:    "uatom",
				Display: "atom",
			},
			false,
		},
		{"empty metadata", types.Metadata{}, true},
		{
			"invalid base denom",
			types.Metadata{
				Base: "",
			},
			true,
		},
		{
			"invalid display denom",
			types.Metadata{
				Base:    "uatom",
				Display: "",
			},
			true,
		},
		{
			"duplicate denom unit",
			types.Metadata{
				Description: "The native staking token of the Cosmos Hub.",
				DenomUnits: []*types.DenomUnit{
					{"uatom", uint32(0), []string{"microatom"}},
					{"uatom", uint32(1), []string{"microatom"}},
				},
				Base:    "uatom",
				Display: "atom",
			},
			true,
		},
		{
			"invalid denom unit",
			types.Metadata{
				Description: "The native staking token of the Cosmos Hub.",
				DenomUnits: []*types.DenomUnit{
					{"", uint32(0), []string{"microatom"}},
				},
				Base:    "uatom",
				Display: "atom",
			},
			true,
		},
		{
			"invalid denom unit alias",
			types.Metadata{
				Description: "The native staking token of the Cosmos Hub.",
				DenomUnits: []*types.DenomUnit{
					{"uatom", uint32(0), []string{""}},
				},
				Base:    "uatom",
				Display: "atom",
			},
			true,
		},
		{
			"duplicate denom unit alias",
			types.Metadata{
				Description: "The native staking token of the Cosmos Hub.",
				DenomUnits: []*types.DenomUnit{
					{"uatom", uint32(0), []string{"microatom", "microatom"}},
				},
				Base:    "uatom",
				Display: "atom",
			},
			true,
		},
		{
			"no base denom unit",
			types.Metadata{
				Description: "The native staking token of the Cosmos Hub.",
				DenomUnits: []*types.DenomUnit{
					{"matom", uint32(3), []string{"milliatom"}},
					{"atom", uint32(6), nil},
				},
				Base:    "uatom",
				Display: "atom",
			},
			true,
		},
		{
			"base denom exponent not zero",
			types.Metadata{
				Description: "The native staking token of the Cosmos Hub.",
				DenomUnits: []*types.DenomUnit{
					{"uatom", uint32(1), []string{"microatom"}},
					{"matom", uint32(3), []string{"milliatom"}},
					{"atom", uint32(6), nil},
				},
				Base:    "uatom",
				Display: "atom",
			},
			true,
		},
		{
			"no display denom unit",
			types.Metadata{
				Description: "The native staking token of the Cosmos Hub.",
				DenomUnits: []*types.DenomUnit{
					{"uatom", uint32(0), []string{"microatom"}},
				},
				Base:    "uatom",
				Display: "atom",
			},
			true,
		},
		{
			"denom units not sorted",
			types.Metadata{
				Description: "The native staking token of the Cosmos Hub.",
				DenomUnits: []*types.DenomUnit{
					{"uatom", uint32(0), []string{"microatom"}},
					{"atom", uint32(6), nil},
					{"matom", uint32(3), []string{"milliatom"}},
				},
				Base:    "uatom",
				Display: "atom",
			},
			true,
		},
		{
			"without denomuints, sendenable =true",
			types.Metadata{
				Description: "The native staking token of the Cosmos Hub.",
				Base:    sdk.DefaultBondDenom,
				Display: sdk.DefaultBondDenom,
				Decimals: 18,
				SendEnabled: true,
			},
			false,
		},
		{
			"without denomuints sendenable = false",
			types.Metadata{
				Description: "The native staking token of the Cosmos Hub.",
				Base:    sdk.DefaultBondDenom,
				Display: sdk.DefaultBondDenom,
				Decimals: 18,
				SendEnabled: false,
			},
			false,
		},
		{
			"decimal >  sdk.Precision",
			types.Metadata{
				Description: "The native staking token of the Cosmos Hub.",
				Base:    sdk.DefaultBondDenom,
				Display: sdk.DefaultBondDenom,
				Decimals: 19,
			},
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			err := tc.metadata.Validate()

			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMarshalJSONMetaData(t *testing.T) {
	cdc := codec.NewLegacyAmino()

	testCases := []struct {
		name      string
		input     []types.Metadata
		strOutput string
	}{
		{"nil metadata", nil, `null`},
		{"empty metadata", []types.Metadata{}, `[]`},
		{"non-empty coins", []types.Metadata{{
			Description: "The native staking token of the Cosmos Hub.",
			DenomUnits: []*types.DenomUnit{
				{"uatom", uint32(0), []string{"microatom"}}, // The default exponent value 0 is omitted in the json
				{"matom", uint32(3), []string{"milliatom"}},
				{"atom", uint32(6), nil},
			},
			Base:    "uatom",
			Display: "atom",
		},
		},
			`[{"description":"The native staking token of the Cosmos Hub.","denom_units":[{"denom":"uatom","aliases":["microatom"]},{"denom":"matom","exponent":3,"aliases":["milliatom"]},{"denom":"atom","exponent":6}],"base":"uatom","display":"atom"}]`},

		{"without denomuints", []types.Metadata{{
			Description: sdk.DefaultDenom,
			Base:    sdk.DefaultDenom,
			Display: sdk.DefaultDenom,
			Decimals: 18,
		},
		},
			`[{"description":"stake","base":"stake","display":"stake","decimals":"18"}]`},

		{"with issuer", []types.Metadata{{
			Description: sdk.DefaultDenom,
			Base:    sdk.DefaultDenom,
			Display: sdk.DefaultDenom,
			Issuer: "AAAA",
			Decimals: 18,
		},
		},
			`[{"description":"stake","base":"stake","display":"stake","issuer":"AAAA","decimals":"18"}]`},
			{"with sendenable = true", []types.Metadata{{
			Description: sdk.DefaultDenom,
			Base:    sdk.DefaultDenom,
			Display: sdk.DefaultDenom,
			Issuer: "AAAA",
			Decimals: 18,
			SendEnabled: true,
		},
		},
			`[{"description":"stake","base":"stake","display":"stake","issuer":"AAAA","decimals":"18","send_enabled":true}]`},

		{"with sendenable = false", []types.Metadata{{
			Description: sdk.DefaultDenom,
			Base:    sdk.DefaultDenom,
			Display: sdk.DefaultDenom,
			Issuer: "AAAA",
			Decimals: 18,
			SendEnabled: false,
		},
		},
			`[{"description":"stake","base":"stake","display":"stake","issuer":"AAAA","decimals":"18"}]`},

	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			bz, err := cdc.MarshalJSON(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.strOutput, string(bz))

			var newMetadata []types.Metadata
			require.NoError(t, cdc.UnmarshalJSON(bz, &newMetadata))

			if len(tc.input) == 0 {
				require.Nil(t, newMetadata)
			} else {
				require.Equal(t, tc.input, newMetadata)
			}
		})
	}
}

func TestDefaultMetadata(t *testing.T) {
	data := types.DefaultMetadatas()
	require.Equal(t, 1, len(data))
	require.Equal(t, sdk.DefaultDenom, data[0].Description)
	require.Equal(t, sdk.DefaultDenom, data[0].Base)
	require.Equal(t, sdk.DefaultDenom, data[0].Display)
	require.EqualValues(t, 18, data[0].Decimals)
	require.Equal(t, true, data[0].SendEnabled)
	require.NoError(t, data[0].Validate())
	require.NoError(t, types.Metadatas(data).Validate())
}


func TestMetadataSorting(t *testing.T) {
	metaDatas := []types.Metadata{}

	metaEth := types.Metadata{
		Description: "eth",
		Base: "eth",
		Display: "eth",
		Decimals: 18,
		SendEnabled: true,
	}

	metaBtc := types.Metadata{
		Description: "btc",
		Base: "btc",
		Display: "btc",
		Decimals: 18,
		SendEnabled: true,
	}

	metaBhc := types.Metadata{
		Description: "bhc",
		Base: "bhc",
		Display: "bhc",
		Decimals: 18,
		SendEnabled: true,
	}

	metaAtom := types.Metadata{
		Description: "atom",
		Base: "atom",
		Display: "atom",
		Decimals: 18,
		SendEnabled: true,
	}

	metaTrx := types.Metadata{
		Description: "trx",
		Base: "trx",
		Display: "trx",
		Decimals: 18,
		SendEnabled: true,
	}

	metaDatas = append(metaDatas, metaTrx, metaBhc, metaBtc,metaEth, metaAtom)
	res := types.Metadatas(metaDatas).Sort()
	require.Equal(t, metaAtom, res[0])
	require.Equal(t, metaBhc, res[1])
	require.Equal(t, metaBtc, res[2])
	require.Equal(t, metaEth, res[3])
	require.Equal(t, metaTrx, res[4])
}