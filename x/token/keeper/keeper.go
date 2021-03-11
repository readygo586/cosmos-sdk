package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/token/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey      sdk.StoreKey          // Unexposed key to access store from sdk.Context
	cdc           codec.BinaryMarshaler // The wire codec for binary encoding/decoding
	legacyAmino   *codec.LegacyAmino
	distrKeeper   types.DistrKeeper
	bankKeeper    types.BankKeeper
	paramSubSpace paramtypes.Subspace
}

//NewKeeper create token's Keeper
func NewKeeper(cdc codec.BinaryMarshaler, legacyAmino *codec.LegacyAmino, storeKey sdk.StoreKey, distrKeeper types.DistrKeeper, bankKeeper types.BankKeeper, paramSubSpace paramtypes.Subspace) Keeper {
	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		legacyAmino:   legacyAmino,
		distrKeeper:   distrKeeper,
		bankKeeper:    bankKeeper,
		paramSubSpace: paramSubSpace.WithKeyTable(types.ParamKeyTable()),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

//GetTokenInfo get a specified  tokeninfo, whose totalsupply is stored in supply module
func (k *Keeper) GetTokenInfo(ctx sdk.Context, denom string) banktypes.Metadata {
	return k.bankKeeper.GetDenomMetaData(ctx, denom)
}

func (k *Keeper) SetTokenInfo(ctx sdk.Context, meta banktypes.Metadata) {
	k.bankKeeper.SetDenomMetaData(ctx, meta)
}

//GetAllTokenInfo get all token info from bank module
func (k *Keeper) GetAllTokenInfo(ctx sdk.Context) []banktypes.Metadata {
	tis := make([]banktypes.Metadata, 0)
	k.bankKeeper.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
		tis = append(tis, metadata)
		return false
	})
	return tis
}

//GetIssuer ...
func (k *Keeper) GetIssuer(ctx sdk.Context, denom string) string {
	return k.GetTokenInfo(ctx, denom).Issuer
}

//IsTokenSupported ...
func (k *Keeper) IsSupported(ctx sdk.Context, denom string) bool {
	return k.GetTokenInfo(ctx, denom).Base == denom
}

//SendEnabled ...
func (k *Keeper) SendEnabled(ctx sdk.Context, denom string) bool {
	return k.GetTokenInfo(ctx, denom).SendEnabled
}

//GetDecimals ...
func (k *Keeper) GetDecimals(ctx sdk.Context, denom string) uint64 {
	return k.GetTokenInfo(ctx, denom).Decimals
}

//GetTotalSupply ...
func (k *Keeper) GetTotalSupply(ctx sdk.Context, denom string) sdk.Int {
	return k.bankKeeper.GetSupply(ctx, denom).GetTotal().Amount
}

//EnableSend ...
func (k *Keeper) EnableSend(ctx sdk.Context, denom string) {
	meta := k.bankKeeper.GetDenomMetaData(ctx, denom)
	if meta.Base == denom {
		meta.SendEnabled = true
		k.bankKeeper.SetDenomMetaData(ctx, meta)
	}
}

//DisableSend ...
func (k *Keeper) DisableSend(ctx sdk.Context, denom string) {
	meta := k.bankKeeper.GetDenomMetaData(ctx, denom)
	if meta.Base == denom {
		meta.SendEnabled = false
		k.bankKeeper.SetDenomMetaData(ctx, meta)
	}
}

//GetSymbols ...
func (k *Keeper) GetSymbols(ctx sdk.Context) []string {
	symbols := make([]string, 0)
	k.bankKeeper.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
		symbols = append(symbols, metadata.Base)
		return false
	})
	return symbols
}

// SetParams sets the token module's parameters.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubSpace.SetParamSet(ctx, &params)
}

// GetParams gets the token module's parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSubSpace.GetParamSet(ctx, &params)
	return
}
