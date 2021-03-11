package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	yaml "gopkg.in/yaml.v2"
	"strings"
)

const (
	// ProposalTypeAddToken defines the type for a AddToken
	ProposalTypeTokenParamsChange = "TokenParamsChange"
	ProposalTypeDisableToken      = "DisableToken"
)

// Assert proposl implements govtypes.Content at compile-time
var _ govtypes.Content = &TokenParamsChangeProposal{}
var _ govtypes.Content = &DisableTokenProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeTokenParamsChange)
	govtypes.RegisterProposalTypeCodec(&TokenParamsChangeProposal{}, "cosmos-sdk/token/TokenParamsChangeProposal")
	govtypes.RegisterProposalType(ProposalTypeDisableToken)
	govtypes.RegisterProposalTypeCodec(&DisableTokenProposal{}, "cosmos-sdk/token/DisableTokenProposal")
}

//NewParamChange create a paramchange(k,v)
func NewParamChange(key, value string) ParamChange {
	return ParamChange{key, value}
}

func (pc ParamChange) String() string {
	out, _ := yaml.Marshal(pc)
	return string(out)
}

// NewTokenParamsChangeProposal creates a new add token proposal.
func NewTokenParamsChangeProposal(title, description, denom string, changes []ParamChange) *TokenParamsChangeProposal {
	return &TokenParamsChangeProposal{
		Title:       title,
		Description: description,
		Denom:       denom,
		Changes:     changes,
	}
}

// GetTitle returns the title of a token parameter change proposal.
func (ctpp *TokenParamsChangeProposal) GetTitle() string { return ctpp.Title }

// GetDescription returns the description of a token parameter change proposal.
func (ctpp *TokenParamsChangeProposal) GetDescription() string { return ctpp.Description }

// GetDescription returns the routing key of a token parameter change proposal.
func (ctpp *TokenParamsChangeProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a token parameter change proposal.
func (ctpp *TokenParamsChangeProposal) ProposalType() string { return ProposalTypeTokenParamsChange }

// ValidateBasic runs basic stateless validity checks
func (ctpp *TokenParamsChangeProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(ctpp); err != nil {
		return err
	}

	if err := sdk.ValidateDenom(ctpp.Denom); err != nil {
		return err
	}

	if ctpp.Denom == sdk.DefaultDenom {
		return fmt.Errorf("Not allowed to change native token's params")
	}

	//dectect duplicated keys if any
	keysMap := map[string]interface{}{}

	for _, pc := range ctpp.Changes {
		_, ok := keysMap[pc.Key]
		if !ok {
			keysMap[pc.Key] = nil
		} else {
			return fmt.Errorf("Duplicated key in TokenParamsChangeProposal")
		}

		if len(pc.Key) == 0 {
			return fmt.Errorf("Empty key found in TokenParamsChangeProposal")
		}
		if len(pc.Value) == 0 {
			return fmt.Errorf("Empty value found in TokenParamsChangeProposal")
		}
	}
	return nil
}

// String implements the Stringer interface.
func (ctpp *TokenParamsChangeProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Change Token Param Proposal:
 Title:       %s
 Description: %s
 Denom:      %s
 Changes:
`, ctpp.Title, ctpp.Description, ctpp.Denom))

	for _, pc := range ctpp.Changes {
		b.WriteString(fmt.Sprintf("%s: %s\t", pc.Key, pc.Value))
	}
	return b.String()
}

// NewDisableTokenProposal creates a new disable token proposal.
func NewDisableTokenProposal(title, description, denom string) *DisableTokenProposal {
	return &DisableTokenProposal{
		Title:       title,
		Description: description,
		Denom:       denom,
	}
}

// GetTitle returns the title of a disable token proposal..
func (dtp *DisableTokenProposal) GetTitle() string { return dtp.Title }

// GetDescription returns the description of a disable token proposal..
func (dtp *DisableTokenProposal) GetDescription() string { return dtp.Description }

// GetDescription returns the routing key of a disable token proposal..
func (dtp *DisableTokenProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a disable token proposal.
func (dtp *DisableTokenProposal) ProposalType() string { return ProposalTypeDisableToken }

// ValidateBasic runs basic stateless validity checks
func (dtp *DisableTokenProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(dtp)
	if err != nil {
		return err
	}

	if err := sdk.ValidateDenom(dtp.Denom); err != nil {
		return err
	}

	if dtp.Denom == sdk.DefaultDenom {
		return fmt.Errorf("Not allowed to change native token's params")
	}

	return err
}

// String implements the Stringer interface.
func (dtp *DisableTokenProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Disable Token Proposal:
 Title:       %s
 Description: %s
 Denom:      %s
`, dtp.Title, dtp.Description, dtp.Denom))
	return b.String()
}
