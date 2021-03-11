package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/token/types"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "default denom sendenable",
		DefaultDenomSendEnableInvariant(k))
	ir.RegisterRoute(types.ModuleName, "decimals overflow",
		DecimalsOverFlow(k))
}

func DecimalsOverFlow(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var (
			msg   string
			count uint64
		)
		k.bankKeeper.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
			if metadata.Decimals > sdk.Precision {
				count++
				msg += fmt.Sprintf("\t%s decimals overflow: %v\n", metadata.Base, metadata.Decimals)
			}
			return false
		})

		broken := count != 0
		return sdk.FormatInvariant(types.ModuleName, "decimals overflow",
			fmt.Sprintf("found %d denom with overflow decimals \n%s", count, msg)), broken
	}
}

func DefaultDenomSendEnableInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		sendenable := k.SendEnabled(ctx, sdk.DefaultDenom)
		broken := (sendenable != true)
		return sdk.FormatInvariant(types.ModuleName, "default denom send enable",
			"default denom send enable is false"), broken
	}
}
