package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

// x/token module sentinel errors
var (
	ErrDecimalsOverFlow  = sdkerrors.Register(ModuleName, 1, "decimal overflows")
	ErrDenomAlreadyExist = sdkerrors.Register(ModuleName, 2, "denom already exist")
	ErrDenomNotExist     = sdkerrors.Register(ModuleName, 3, "denom not exist")
	ErrInvalidParams     = sdkerrors.Register(ModuleName, 4, "invalid parameter")
)
