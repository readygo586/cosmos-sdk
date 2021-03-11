package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"strings"
)

// Default parameter values
const (
	DefaultTokenCacheSize uint64 = 32 //cache size for token
)

var (
	DefaultNewTokenFee = sdk.TokensFromConsensusPower(100) //100 DefaultDenom
)

// Parameter keys
var (
	KeyTokenCacheSize = []byte("TokenCacheSize")
	KeyNewTokenFee    = []byte("NewTokenFee")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable for auth module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		{KeyTokenCacheSize, &p.TokenCacheSize, validateTokenCacheSize},
		{KeyNewTokenFee, &p.NewTokenFee, valiateNewTokenFee},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		TokenCacheSize: DefaultTokenCacheSize,
		NewTokenFee:    DefaultNewTokenFee,
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params:")
	sb.WriteString(fmt.Sprintf("TokenCacheSize:%v\t", p.TokenCacheSize))
	sb.WriteString(fmt.Sprintf("NewTokenFee:%v\t", p.NewTokenFee))

	return sb.String()
}

func (p Params) Validate() error {
	if p.NewTokenFee.IsPositive() {
		return nil
	}
	return fmt.Errorf("NewTokenFee %v is not valid", p.NewTokenFee)
}

func valiateNewTokenFee(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("new token fee must be not nil")
	}
	if !v.IsPositive() {
		return fmt.Errorf("new token fee must be positive: %s", v)
	}

	return nil
}

func validateTokenCacheSize(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
