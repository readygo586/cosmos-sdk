package rest_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/testutil"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/token/types"
	"github.com/gogo/protobuf/proto"
)

func (s *IntegrationTestSuite) TestSymbolsGRPCHandler() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name     string
		url      string
		headers  map[string]string
		respType proto.Message
		expected proto.Message
	}{
		{
			"test GRPC symbols",
			fmt.Sprintf("%s/cosmos/token/v1beta1/symbols", baseURL),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			&types.QuerySymbolsResponse{},
			&types.QuerySymbolsResponse{
				Symbols: []string{fmt.Sprintf("%stoken", val.Moniker), s.cfg.BondDenom},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			s.Require().NoError(err)

			s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(resp, tc.respType))
			s.Require().Equal(tc.expected.String(), tc.respType.String())
		})
	}
}

func (s *IntegrationTestSuite) TestSymbolGRPCHandler() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name     string
		url      string
		headers  map[string]string
		respType proto.Message
		expected proto.Message
	}{
		{
			"test GRPC symbol,default denom",
			fmt.Sprintf("%s/cosmos/token/v1beta1/symbol/%s", baseURL, s.cfg.BondDenom),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			&types.QuerySymbolResponse{},
			&types.QuerySymbolResponse{
				Info: banktypes.Metadata{
					Description: s.cfg.BondDenom,
					Base:        s.cfg.BondDenom,
					Display:     s.cfg.BondDenom,
					Issuer:      "",
					Decimals:    18,
					SendEnabled: true,
				},
			},
		},
		{
			"test GRPC symbol, nodextoken",
			fmt.Sprintf("%s/cosmos/token/v1beta1/symbol/%s", baseURL, fmt.Sprintf("%stoken", val.Moniker)),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			&types.QuerySymbolResponse{},
			&types.QuerySymbolResponse{
				Info: banktypes.Metadata{
					Description: fmt.Sprintf("%stoken", val.Moniker),
					Base:        fmt.Sprintf("%stoken", val.Moniker),
					Display:     fmt.Sprintf("%stoken", val.Moniker),
					Issuer:      "",
					Decimals:    18,
					SendEnabled: true,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			s.Require().NoError(err)

			s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(resp, tc.respType))
			s.Require().Equal(tc.expected.String(), tc.respType.String())
		})
	}
}

func (s *IntegrationTestSuite) TestParamsGRPCHandler() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name     string
		url      string
		headers  map[string]string
		respType proto.Message
		expected proto.Message
	}{
		{
			"test GRPC symbols",
			fmt.Sprintf("%s/cosmos/token/v1beta1/params", baseURL),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			&types.QueryParamsResponse{},
			&types.QueryParamsResponse{
				Params: types.Params{
					TokenCacheSize: types.DefaultTokenCacheSize,
					NewTokenFee:    types.DefaultNewTokenFee,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			s.Require().NoError(err)

			s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(resp, tc.respType))
			s.Require().Equal(tc.expected.String(), tc.respType.String())
		})
	}
}
