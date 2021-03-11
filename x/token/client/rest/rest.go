package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/gorilla/mux"
)

// RegisterHandlers registers all x/bank transaction and query HTTP REST handlers
// on the provided mux router.
func RegisterHandlers(clientCtx client.Context, rtr *mux.Router) {
	r := rest.WithHTTPDeprecationHeaders(rtr)
	r.HandleFunc("/token/newtoken/{address}", NewTokenRequestHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc("/token/infaltetoken/{address}", InflateTokenRequestHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc("/token/burntoken", BurnTokenRequestHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc("/token/symbols", querySymbolsHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/token/params", queryParamsHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/token/symbol/{denom}", querySymbolHandlerFn(clientCtx)).Methods("GET")
}
