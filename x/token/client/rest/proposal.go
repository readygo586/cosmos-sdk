package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/token/client/cli"
	"github.com/cosmos/cosmos-sdk/x/token/types"
	"net/http"
)

type (

	// TokenParamsChangeProposalReq defines a token params change request body.
	TokenParamsChangeProposalReq struct {
		BaseReq     rest.BaseReq         `json:"base_req" yaml:"base_req"`
		Title       string               `json:"title" yaml:"title"`
		Description string               `json:"description" yaml:"description"`
		Symbol      string               `json:"symbol" yaml:"symbol"`
		Changes     cli.ParamChangesJSON `json:"changes" yaml:"changes"`
		Deposit     sdk.Coins            `json:"deposit" yaml:"deposit"`
	}

	// DisableTokenProposalReq defines a disable token request body.
	DisableTokenProposalReq struct {
		BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
		Title       string       `json:"title" yaml:"title"`
		Description string       `json:"description" yaml:"description"`
		Symbol      string       `json:"symbol" yaml:"symbol"`
		Deposit     sdk.Coins    `json:"deposit" yaml:"deposit`
	}
)

// ProposalRESTHandler returns a ProposalRESTHandler that exposes the community pool spend REST handler with a given sub-route.
func DisableTokenProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "disable_token",
		Handler:  disableTokenProposalHandlerFn(clientCtx),
	}
}

func disableTokenProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DisableTokenProposalReq
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

		content := types.NewDisableTokenProposal(req.Title, req.Description, req.Symbol)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func TokenParamsChangeProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "token_params_change",
		Handler:  tokenParamsChangeProposalHandlerFn(clientCtx),
	}
}

func tokenParamsChangeProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TokenParamsChangeProposalReq
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

		changes := req.Changes.ToParamChanges()
		content := types.NewTokenParamsChangeProposal(req.Title, req.Description, req.Symbol, changes)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
