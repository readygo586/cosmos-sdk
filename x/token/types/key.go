package types

const (
	// module name
	ModuleName = "token"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey is the message route for token
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// Parameter store default paramter store
	DefaultParamspace = ModuleName

	// query endpoints supported by the nameservice Querier
	QuerySymbol     = "symbol"
	QuerySymbols    = "symbols"
	QueryParameters = "parameters"

	// MsgNewToken
	TypeMsgNewToken     = "new"
	TypeMsgInflateToken = "inflate"
	TypeMsgBurnToken    = "burn"
)
