package rest_test

import (
	"encoding/hex"
	"fmt"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/token/types"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := network.DefaultConfig()
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestQuerySymbolsHandlerFn() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name     string
		url      string
		respType []string
		expected []string
	}{
		{
			"Get all symbols",
			fmt.Sprintf("%s/token/symbols?height=1", baseURL),
			[]string{},
			[]string{fmt.Sprintf("%stoken", val.Moniker), s.cfg.BondDenom},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)

			bz, err := rest.ParseResponseWithHeight(val.ClientCtx.LegacyAmino, resp)
			s.Require().NoError(err)
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(bz, &tc.respType))
			s.Require().Equal(tc.expected, tc.respType)
		})
	}
}

func (s *IntegrationTestSuite) TestQuerySymbolHandlerFn() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name     string
		url      string
		respType fmt.Stringer
		expected fmt.Stringer
	}{
		{
			"Get default denom info",
			fmt.Sprintf("%s/token/symbol/%s?height=1", baseURL, s.cfg.BondDenom),
			&banktypes.Metadata{},
			&banktypes.Metadata{
				Description: s.cfg.BondDenom,
				Base:        s.cfg.BondDenom,
				Display:     s.cfg.BondDenom,
				Issuer:      "",
				Decimals:    18,
				SendEnabled: true,
			},
		},
		{
			"Get node's token denom info",
			fmt.Sprintf("%s/token/symbol/%s?height=1", baseURL, fmt.Sprintf("%stoken", val.Moniker)),
			&banktypes.Metadata{},
			&banktypes.Metadata{
				Description: fmt.Sprintf("%stoken", val.Moniker),
				Base:        fmt.Sprintf("%stoken", val.Moniker),
				Display:     fmt.Sprintf("%stoken", val.Moniker),
				Issuer:      "",
				Decimals:    18,
				SendEnabled: true,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)

			bz, err := rest.ParseResponseWithHeight(val.ClientCtx.LegacyAmino, resp)
			s.Require().NoError(err)
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(bz, tc.respType))
			s.Require().Equal(tc.expected, tc.respType)
		})
	}
}

func (s *IntegrationTestSuite) TestQueryParamsHandlerFn() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name     string
		url      string
		respType fmt.Stringer
		expected fmt.Stringer
	}{
		{
			"Get params",
			fmt.Sprintf("%s/token/params?height=1", baseURL),
			&types.Params{},
			&types.Params{
				TokenCacheSize: types.DefaultTokenCacheSize,
				NewTokenFee:    types.DefaultNewTokenFee,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)

			bz, err := rest.ParseResponseWithHeight(val.ClientCtx.LegacyAmino, resp)
			fmt.Printf("resp:%+v, bz:%v", hex.EncodeToString(resp), hex.EncodeToString(bz))
			s.Require().NoError(err)
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(bz, tc.respType))
			s.Require().Equal(tc.expected, tc.respType)
		})
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
