package types

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/exported"
)

// Implements Supplys interface
var _ exported.SupplysI = (*Supplys)(nil)
// NewSupply creates a new Supply instance
func NewSupplys(total sdk.Coins) *Supplys {
	return &Supplys{total}
}
// DefaultSupply creates an empty Supply
func DefaultSupplys() *Supplys {
	return NewSupplys(sdk.NewCoins())
}
// SetTotal sets the total supply.
func (supplys *Supplys) SetTotal(total sdk.Coins) {
	supplys.Total = total
}
// GetTotal returns the supply total.
func (supplys Supplys) GetTotal() sdk.Coins {
	return supplys.Total
}
// Inflate adds coins to the total supply
func (supplys *Supplys) Inflate(amount sdk.Coins) {
	supplys.Total = supplys.Total.Add(amount...)
}
// Deflate subtracts coins from the total supply.
func (supplys *Supplys) Deflate(amount sdk.Coins) {
	supplys.Total = supplys.Total.Sub(amount)
}
// String returns a human readable string representation of a supplier.
func (supplys Supplys) String() string {
	bz, _ := yaml.Marshal(supplys)
	return string(bz)
}
// ValidateBasic validates the Supply coins and returns error if invalid
func (supplys Supplys) ValidateBasic() error {
	if !supplys.Total.IsValid() {
		return fmt.Errorf("invalid total supply: %s", supplys.Total.String())
	}
	return nil
}
// Implements Supply interface
var _ exported.SupplyI = (*Supply)(nil)

// NewSupply creates a new Supply instance
func NewSupply(total sdk.Coin) *Supply {
	return &Supply{total}
}

// DefaultSupply creates an empty Supply
func DefaultSupply() *Supply {
	return NewSupply(sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()))
}

// SetTotal sets the total supply.
func (supply *Supply) SetTotal(total sdk.Coin) {
	supply.Total = total
}

// GetTotal returns the supply total.
func (supply Supply) GetTotal() sdk.Coin {
	return supply.Total
}

// Inflate adds coins to the total supply
func (supply *Supply) Inflate(amount sdk.Coin) {
	supply.Total = supply.Total.Add(amount)
}

// Deflate subtracts coins from the total supply.
func (supply *Supply) Deflate(amount sdk.Coin) {
	supply.Total = supply.Total.Sub(amount)
}

// String returns a human readable string representation of a supplier.
func (supply Supply) String() string {
	bz, _ := yaml.Marshal(supply)
	return string(bz)
}

// ValidateBasic validates the Supply coins and returns error if invalid
func (supply Supply) ValidateBasic() error {
	if !supply.Total.IsValid() {
		return fmt.Errorf("invalid total supply: %s", supply.Total.String())
	}

	return nil
}
