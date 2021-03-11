package types

// NewQueryBalanceRequest creates a new instance of QueryBalanceRequest.
//nolint:interfacer
func NewQuerySymbolRequest(denom string) *QuerySymbolRequest {
	return &QuerySymbolRequest{Denom: denom}
}

// NewQueryAllBalancesRequest creates a new instance of QueryAllBalancesRequest.
//nolint:interfacer
func NewQuerySymbolsRequest() *QuerySymbolsRequest {
	return &QuerySymbolsRequest{}
}

// NewQueryParams creates a new instance to query the params
func NewQueryParams() *QueryParamsRequest {
	return &QueryParamsRequest{}
}
