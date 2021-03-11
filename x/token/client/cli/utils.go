package cli

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/token/types"
	"io/ioutil"
)

type (
	ParamChangesJSON []ParamChangeJSON

	// ParamChangeJSON defines a parameter change used in JSON input. This
	// allows values to be specified in raw JSON instead of being string encoded.
	ParamChangeJSON struct {
		Key   string          `json:"key" yaml:"key"`
		Value json.RawMessage `json:"value" yaml:"value"`
	}

	TokenParamsChangeProposalJSON struct {
		Title       string           `json:"title" yaml:"title"`
		Description string           `json:"description" yaml:"description"`
		Denom       string           `json:"denom" yaml:"denom"`
		Changes     ParamChangesJSON `json:"changes" yaml:"changes"`
		Deposit     sdk.Coins        `json:"deposit,omitempty" yaml:"deposit"`
	}

	DisableTokenProposalJSON struct {
		Title       string    `json:"title" yaml:"title"`
		Description string    `json:"description" yaml:"description"`
		Denom       string    `json:"denom" yaml:"denom"`
		Deposit     sdk.Coins `json:"deposit,omitempty" yaml:"deposit`
	}
)

// ToParamChange converts a ParamChangeJSON object to ParamChange.
func (pcj ParamChangeJSON) ToParamChange() types.ParamChange {
	return types.NewParamChange(pcj.Key, string(pcj.Value))
}

// ToParamChanges converts a slice of paramChangesJSON objects to a slice of
// ParamChange.
func (pcsj ParamChangesJSON) ToParamChanges() []types.ParamChange {
	res := make([]types.ParamChange, len(pcsj))
	for i, pc := range pcsj {
		res[i] = pc.ToParamChange()
	}
	return res
}

// ParseTokenParamsChangeProposalJSON reads and parses a tokenParamsChangeProposalJSON from a file.
func ParseTokenParamsChangeProposalJSON(cdc *codec.LegacyAmino, proposalFile string) (TokenParamsChangeProposalJSON, error) {
	proposal := TokenParamsChangeProposalJSON{}
	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

// ParseDisableTokenProposalJSON reads and parses a disableTokenProposalJSON from a file.
func ParseDisableTokenProposalJSON(cdc *codec.LegacyAmino, proposalFile string) (DisableTokenProposalJSON, error) {
	proposal := DisableTokenProposalJSON{}
	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
