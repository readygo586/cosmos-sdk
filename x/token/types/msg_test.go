package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgNewTokenRouteAndType(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from________________"))
	addr2 := sdk.AccAddress([]byte("to__________________"))
	var msg = NewMsgNewToken(addr1, addr2, 8, sdk.NewCoin("btc", sdk.NewInt(1000000000000)))

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), TypeMsgNewToken)
}

func TestMsgNewTokenValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from________________"))
	addr2 := sdk.AccAddress([]byte("to__________________"))
	addrEmpty := sdk.AccAddress([]byte(""))
	addrNil := sdk.AccAddress(nil)
	addrTooLong := sdk.AccAddress([]byte("Accidentally used 33 bytes pubkey"))
	validCoin := sdk.NewCoin("btc", sdk.NewInt(1000000000000))
	zeroCoin := sdk.NewCoin("btc", sdk.NewInt(0))

	cases := []struct {
		expectedErr string // empty means no error expected
		msg         *MsgNewToken
	}{
		{"", NewMsgNewToken(addr1, addr2, 0, validCoin)},
		{"", NewMsgNewToken(addr1, addr2, 18, validCoin)},
		{"overflow decimals (19): decimal overflows", NewMsgNewToken(addr1, addr2, 19, validCoin)},
		{"", NewMsgNewToken(addr1, addr2, 8, validCoin)},                   // valid send
		{"0btc: invalid coins", NewMsgNewToken(addr1, addr2, 8, zeroCoin)}, // non positive coin
		{"Invalid sender address (empty address string is not allowed): invalid address", NewMsgNewToken(addrEmpty, addr2, 8, validCoin)},
		{"Invalid sender address (empty address string is not allowed): invalid address", NewMsgNewToken(addrNil, addr2, 8, validCoin)},
		{"Invalid sender address (empty address string is not allowed): invalid address", NewMsgNewToken(nil, addr2, 8, validCoin)},
		{"Invalid sender address (incorrect address length (expected: 20, actual: 33)): invalid address", NewMsgNewToken(addrTooLong, addr2, 8, validCoin)},
		{"Invalid recipient address (empty address string is not allowed): invalid address", NewMsgNewToken(addr1, addrEmpty, 8, validCoin)},
		{"Invalid recipient address (empty address string is not allowed): invalid address", NewMsgNewToken(addr1, addrNil, 8, validCoin)},
		{"Invalid recipient address (empty address string is not allowed): invalid address", NewMsgNewToken(addr1, nil, 8, validCoin)},
		{"Invalid recipient address (incorrect address length (expected: 20, actual: 33)): invalid address", NewMsgNewToken(addr1, addrTooLong, 8, validCoin)},
	}

	for _, tc := range cases {
		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgNewTokenGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from________________"))
	addr2 := sdk.AccAddress([]byte("to__________________"))
	var msg = NewMsgNewToken(addr1, addr2, 8, sdk.NewCoin("btc", sdk.NewInt(1000000000000)))
	res := msg.GetSignBytes()

	expected := `{"type":"cosmos-sdk/MsgNewToken","value":{"amount":{"amount":"1000000000000","denom":"btc"},"decimals":"8","from_address":"cosmos1veex7m2lta047h6lta047h6lta047h6lt50pqc","to_address":"cosmos1w3h47h6lta047h6lta047h6lta047h6l620gq6"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgNewTokenGetSigners(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from________________"))
	addr2 := sdk.AccAddress([]byte("to__________________"))
	var msg = NewMsgNewToken(addr1, addr2, 8, sdk.NewCoin("btc", sdk.NewInt(1000000000000)))
	res := msg.GetSigners()
	require.Equal(t, "[66726F6D5F5F5F5F5F5F5F5F5F5F5F5F5F5F5F5F]", fmt.Sprintf("%v", res))
}

func TestMsgInflateTokenRouteAndType(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from________________"))
	addr2 := sdk.AccAddress([]byte("to__________________"))
	var msg = NewMsgInflateToken(addr1, addr2, sdk.NewCoin("btc", sdk.NewInt(1000000000000)))

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), TypeMsgInflateToken)
}

func TestMsgInflateTokenValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from________________"))
	addr2 := sdk.AccAddress([]byte("to__________________"))
	addrEmpty := sdk.AccAddress([]byte(""))
	addrNil := sdk.AccAddress(nil)
	addrTooLong := sdk.AccAddress([]byte("Accidentally used 33 bytes pubkey"))
	validCoin := sdk.NewCoin("btc", sdk.NewInt(1000000000000))
	zeroCoin := sdk.NewCoin("btc", sdk.NewInt(0))

	cases := []struct {
		expectedErr string // empty means no error expected
		msg         *MsgInflateToken
	}{
		{"", NewMsgInflateToken(addr1, addr2, validCoin)},
		{"", NewMsgInflateToken(addr1, addr2, validCoin)},
		{"0btc: invalid coins", NewMsgInflateToken(addr1, addr2, zeroCoin)}, // non positive coin
		{"Invalid sender address (empty address string is not allowed): invalid address", NewMsgInflateToken(addrEmpty, addr2, validCoin)},
		{"Invalid sender address (empty address string is not allowed): invalid address", NewMsgInflateToken(addrNil, addr2, validCoin)},
		{"Invalid sender address (empty address string is not allowed): invalid address", NewMsgInflateToken(nil, addr2, validCoin)},
		{"Invalid sender address (incorrect address length (expected: 20, actual: 33)): invalid address", NewMsgInflateToken(addrTooLong, addr2, validCoin)},
		{"Invalid recipient address (empty address string is not allowed): invalid address", NewMsgInflateToken(addr1, addrEmpty, validCoin)},
		{"Invalid recipient address (empty address string is not allowed): invalid address", NewMsgInflateToken(addr1, addrNil, validCoin)},
		{"Invalid recipient address (empty address string is not allowed): invalid address", NewMsgInflateToken(addr1, nil, validCoin)},
		{"Invalid recipient address (incorrect address length (expected: 20, actual: 33)): invalid address", NewMsgInflateToken(addr1, addrTooLong, validCoin)},
	}

	for _, tc := range cases {
		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgInflateTokenGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from________________"))
	addr2 := sdk.AccAddress([]byte("to__________________"))
	var msg = NewMsgInflateToken(addr1, addr2, sdk.NewCoin("btc", sdk.NewInt(1000000000000)))
	res := msg.GetSignBytes()

	expected := `{"type":"cosmos-sdk/MsgInflateToken","value":{"amount":{"amount":"1000000000000","denom":"btc"},"from_address":"cosmos1veex7m2lta047h6lta047h6lta047h6lt50pqc","to_address":"cosmos1w3h47h6lta047h6lta047h6lta047h6l620gq6"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgInflateTokenGetSigners(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from________________"))
	addr2 := sdk.AccAddress([]byte("to__________________"))
	var msg = NewMsgInflateToken(addr1, addr2, sdk.NewCoin("btc", sdk.NewInt(1000000000000)))
	res := msg.GetSigners()
	require.Equal(t, "[66726F6D5F5F5F5F5F5F5F5F5F5F5F5F5F5F5F5F]", fmt.Sprintf("%v", res))
}

func TestMsgBurnTokenRouteAndType(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from________________"))
	var msg = NewMsgBurnToken(addr1, sdk.NewCoin("btc", sdk.NewInt(1000000000000)))

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), TypeMsgBurnToken)
}

func TestMsgBurnTokenValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from________________"))
	addrEmpty := sdk.AccAddress([]byte(""))
	addrNil := sdk.AccAddress(nil)
	addrTooLong := sdk.AccAddress([]byte("Accidentally used 33 bytes pubkey"))
	validCoin := sdk.NewCoin("btc", sdk.NewInt(1000000000000))
	zeroCoin := sdk.NewCoin("btc", sdk.NewInt(0))

	cases := []struct {
		expectedErr string // empty means no error expected
		msg         *MsgBurnToken
	}{
		{"", NewMsgBurnToken(addr1, validCoin)},
		{"", NewMsgBurnToken(addr1, validCoin)},
		{"0btc: invalid coins", NewMsgBurnToken(addr1, zeroCoin)}, // non positive coin
		{"Invalid sender address (empty address string is not allowed): invalid address", NewMsgBurnToken(addrEmpty, validCoin)},
		{"Invalid sender address (empty address string is not allowed): invalid address", NewMsgBurnToken(addrNil, validCoin)},
		{"Invalid sender address (empty address string is not allowed): invalid address", NewMsgBurnToken(nil, validCoin)},
		{"Invalid sender address (incorrect address length (expected: 20, actual: 33)): invalid address", NewMsgBurnToken(addrTooLong, validCoin)},
	}

	for _, tc := range cases {
		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgBurnTokenGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from________________"))
	var msg = NewMsgBurnToken(addr1, sdk.NewCoin("btc", sdk.NewInt(1000000000000)))
	res := msg.GetSignBytes()

	expected := `{"type":"cosmos-sdk/MsgBurnToken","value":{"amount":{"amount":"1000000000000","denom":"btc"},"from_address":"cosmos1veex7m2lta047h6lta047h6lta047h6lt50pqc"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgBurnTokenGetSigners(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from________________"))
	var msg = NewMsgBurnToken(addr1, sdk.NewCoin("btc", sdk.NewInt(1000000000000)))
	res := msg.GetSigners()
	require.Equal(t, "[66726F6D5F5F5F5F5F5F5F5F5F5F5F5F5F5F5F5F]", fmt.Sprintf("%v", res))
}
