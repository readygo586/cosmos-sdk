package types

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Validate performs a basic validation of the coin metadata fields. It checks:
//  - Base and Display denominations are valid coin denominations
//  - Base and Display denominations are present in the DenomUnit slice
//  - Base denomination has exponent 0
//  - Denomination units are sorted in ascending order
//  - Denomination units not duplicated
func (m Metadata) Validate() error {
	if err := sdk.ValidateDenom(m.Base); err != nil {
		return fmt.Errorf("invalid metadata base denom: %w", err)
	}

	if err := sdk.ValidateDenom(m.Display); err != nil {
		return fmt.Errorf("invalid metadata display denom: %w", err)
	}

	var (
		hasDisplay      bool
		currentExponent uint32 // check that the exponents are increasing
	)

	seenUnits := make(map[string]bool)

	for i, denomUnit := range m.DenomUnits {
		// The first denomination unit MUST be the base
		if i == 0 {
			// validate denomination and exponent
			if denomUnit.Denom != m.Base {
				return fmt.Errorf("metadata's first denomination unit must be the one with base denom '%s'", m.Base)
			}
			if denomUnit.Exponent != 0 {
				return fmt.Errorf("the exponent for base denomination unit %s must be 0", m.Base)
			}
		} else if currentExponent >= denomUnit.Exponent {
			return errors.New("denom units should be sorted asc by exponent")
		}

		currentExponent = denomUnit.Exponent

		if seenUnits[denomUnit.Denom] {
			return fmt.Errorf("duplicate denomination unit %s", denomUnit.Denom)
		}

		if denomUnit.Denom == m.Display {
			hasDisplay = true
		}

		if err := denomUnit.Validate(); err != nil {
			return err
		}

		seenUnits[denomUnit.Denom] = true
	}

	// in case of no DenomUnits, check Base and Display directly
	if !hasDisplay && m.Base != m.Display{
		return fmt.Errorf("metadata must contain a denomination unit with display denom '%s'", m.Display)
	}

	if m.Decimals > sdk.Precision {
		return ErrDecimalsOverFlow
	}
	return validateIsBool(m.SendEnabled)
}

//DefaultMetadata define default value
func DefaultMetadatas() []Metadata {
	defaultMetaData := Metadata{
		Description: sdk.DefaultDenom,
		Base:        sdk.DefaultDenom,
		Display:     sdk.DefaultDenom,
		Issuer:      "",
		Decimals:    18,
		SendEnabled: true,
	}
	return []Metadata{
		defaultMetaData,
	}
}

// Validate performs a basic validation of the denomination unit fields
func (du DenomUnit) Validate() error {
	if err := sdk.ValidateDenom(du.Denom); err != nil {
		return fmt.Errorf("invalid denom unit: %w", err)
	}

	seenAliases := make(map[string]bool)
	for _, alias := range du.Aliases {
		if seenAliases[alias] {
			return fmt.Errorf("duplicate denomination unit alias %s", alias)
		}

		if strings.TrimSpace(alias) == "" {
			return fmt.Errorf("alias for denom unit %s cannot be blank", du.Denom)
		}

		seenAliases[alias] = true
	}

	return nil
}

type Metadatas []Metadata

//ValidateBasic define
func (metadatas Metadatas) Validate() error {
	for _, m := range metadatas {
		if err := m.Validate(); err != nil {
			return err
		}
	}
	return nil
}

//-----------------------------------------------------------------------------
// Sort interface

// Len implements sort.Interface for Coins
func (metadatas Metadatas) Len() int { return len(metadatas) }

// Less implements sort.Interface for Coins
func (metadatas Metadatas) Less(i, j int) bool { return metadatas[i].Base < metadatas[j].Base }

// Swap implements sort.Interface for Coins
func (metadatas Metadatas) Swap(i, j int) { metadatas[i], metadatas[j] = metadatas[j], metadatas[i] }

var _ sort.Interface = Metadatas{}

// Sort is a helper function to sort the set of coins in-place
func (metadatas Metadatas) Sort() Metadatas {
	sort.Sort(metadatas)
	return metadatas
}

