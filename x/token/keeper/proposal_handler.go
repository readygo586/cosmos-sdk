package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/token/types"
)

// HandleTokenParamsChangeProposal is a handler for executing a token param change proposal
func (k Keeper) HandleTokenParamsChangeProposal(ctx sdk.Context, proposal *types.TokenParamsChangeProposal) error {
	ctx.Logger().Info("HandleTokenParamsChangeProposal", "proposal", proposal)

	if proposal.Denom == sdk.DefaultDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "Not allowed to change native token's params")
	}

	meta := k.bankKeeper.GetDenomMetaData(ctx, proposal.Denom)
	if meta.Base == "" {
		return sdkerrors.Wrapf(types.ErrDenomNotExist, "%s does not exist", proposal.Denom)
	}

	attr := []sdk.Attribute{}
	for _, pc := range proposal.Changes {
		err := processChangeParam(pc.Key, pc.Value, &meta, k)
		if err != nil {
			return err
		}
		attr = append(attr, sdk.NewAttribute(types.AttributeKeyTokenParam, pc.Key), sdk.NewAttribute(types.AttributeKeyTokenParamValue, pc.Value))
	}
	k.bankKeeper.SetDenomMetaData(ctx, meta)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeExecuteTokenParamsChangeProposal, attr...),
	)
	return nil
}

func processChangeParam(key, value string, meta *banktypes.Metadata, k Keeper) error {
	switch key {
	case "send_enabled":
		val := false
		err := k.legacyAmino.UnmarshalJSON([]byte(value), &val)
		if err != nil {
			return err
		}
		meta.SendEnabled = val

	default:
		return fmt.Errorf("Unkonwn parameter:%v", key)
	}

	return nil
}

func (k Keeper) HandleDisableTokenProposal(ctx sdk.Context, proposal *types.DisableTokenProposal) error {
	ctx.Logger().Info("HandleDisableTokenProposal", "proposal", proposal)

	if proposal.Denom == sdk.DefaultDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "Not allowed to change native token's params")
	}

	meta := k.bankKeeper.GetDenomMetaData(ctx, proposal.Denom)
	if meta.Base == "" {
		return sdkerrors.Wrapf(types.ErrDenomNotExist, "%s does not exist", proposal.Denom)
	}

	meta.SendEnabled = false
	k.bankKeeper.SetDenomMetaData(ctx, meta)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeExecuteDisableTokenProposal,
			sdk.NewAttribute(types.AttributeKeyToken, proposal.Denom),
		),
	)
	return nil
}
