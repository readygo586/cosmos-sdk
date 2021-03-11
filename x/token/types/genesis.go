package types

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
)

// ValidateGenesis performs basic validation of supply genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	return data.Params.Validate()
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{
		Params: params,
	}
}

// DefaultGenesisState returns a default bank module genesis state.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams())
}

// GetGenesisStateFromAppState returns x/token GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.JSONMarshaler, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}
