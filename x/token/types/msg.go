package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgNewToken{}

// NewMsgNewToken - construct a msg to new coin
//nolint:interfacer
func NewMsgNewToken(fromAddr, toAddr sdk.AccAddress, decimals uint64, amount sdk.Coin) *MsgNewToken {
	return &MsgNewToken{FromAddress: fromAddr.String(), ToAddress: toAddr.String(), Decimals: decimals, Amount: amount}
}

// Route Implements Msg.
func (msg MsgNewToken) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgNewToken) Type() string { return TypeMsgNewToken }

// ValidateBasic Implements Msg.
func (msg MsgNewToken) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid recipient address (%s)", err)
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if msg.Decimals > sdk.Precision {
		return sdkerrors.Wrapf(ErrDecimalsOverFlow, "overflow decimals (%v)", msg.Decimals)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgNewToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgNewToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

var _ sdk.Msg = &MsgInflateToken{}

// NewMsgInflateToken - construct a msg to inflate coin
//nolint:interfacer
func NewMsgInflateToken(fromAddr, toAddr sdk.AccAddress, amount sdk.Coin) *MsgInflateToken {
	return &MsgInflateToken{FromAddress: fromAddr.String(), ToAddress: toAddr.String(), Amount: amount}
}

// Route Implements Msg.
func (msg MsgInflateToken) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgInflateToken) Type() string { return TypeMsgInflateToken }

// ValidateBasic Implements Msg.
func (msg MsgInflateToken) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid recipient address (%s)", err)
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgInflateToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgInflateToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

var _ sdk.Msg = &MsgBurnToken{}

// NewMsgBurnToken - construct a msg to burn coin
//nolint:interfacer
func NewMsgBurnToken(fromAddr sdk.AccAddress, amount sdk.Coin) *MsgBurnToken {
	return &MsgBurnToken{FromAddress: fromAddr.String(), Amount: amount}
}

// Route Implements Msg.
func (msg MsgBurnToken) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgBurnToken) Type() string { return TypeMsgBurnToken }

// ValidateBasic Implements Msg.
func (msg MsgBurnToken) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBurnToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgBurnToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
