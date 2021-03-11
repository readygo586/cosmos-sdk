package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/token/types"
	"github.com/gorilla/mux"
	"net/http"
)

// NewTokenReq defines the properties of a send request's body.
type NewTokenReq struct {
	BaseReq  rest.BaseReq `json:"base_req" yaml:"base_req"`
	Decimals uint64       `json:"decimals" yaml:"decimals"`
	Amount   sdk.Coin     `json:"amount" yaml:"amount"`
}

// NewTokenRequestHandlerFn returns an HTTP REST handler for creating a MsgNewToken
// transaction.
func NewTokenRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		toAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		var req NewTokenReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		msg := types.NewMsgNewToken(fromAddr, toAddr, req.Decimals, req.Amount)
		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// InflateTokenReq defines the properties of a send request's body.
type InflateTokenReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Coin     `json:"amount" yaml:"amount"`
}

// InfalteTokenRequestHandlerFn returns an HTTP REST handler for creating a MsgInflateToken
// transaction.
func InflateTokenRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		toAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		var req InflateTokenReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		msg := types.NewMsgInflateToken(fromAddr, toAddr, req.Amount)
		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// BurnTokenReq defines the properties of a send request's body.
type BurnTokenReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Coin     `json:"amount" yaml:"amount"`
}

// InfalteTokenRequestHandlerFn returns an HTTP REST handler for creating a MsgInflateToken
// transaction.
func BurnTokenRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BurnTokenReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		msg := types.NewMsgBurnToken(fromAddr, req.Amount)
		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
