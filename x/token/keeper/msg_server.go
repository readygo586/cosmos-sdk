package keeper

import (
	"context"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/token/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

//NewToken, new tokens and send it to recipient
func (k msgServer) NewToken(goCtx context.Context, msg *types.MsgNewToken) (*types.MsgNewTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	denom := msg.Amount.GetDenom()
	from, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	to, _ := sdk.AccAddressFromBech32(msg.ToAddress)
	amount := sdk.NewCoins(msg.Amount)
	metadata := k.bankKeeper.GetDenomMetaData(ctx, denom)
	if metadata.Base == msg.Amount.GetDenom() {
		return nil, sdkerrors.Wrapf(types.ErrDenomAlreadyExist, "%s", denom)
	}

	issueFee := sdk.NewCoins(sdk.NewCoin(sdk.DefaultDenom, k.GetParams(ctx).NewTokenFee))
	if err := k.distrKeeper.AddCoinsFromAccountToFeePool(ctx, from, issueFee); err != nil {
		return nil, err
	}

	k.bankKeeper.MintCoins(ctx, types.ModuleName, amount)
	k.bankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
		Description: denom,
		Base:        denom,
		Display:     denom,
		Decimals:    msg.Decimals,
		Issuer:      msg.FromAddress,
		SendEnabled: true,
	})
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, amount)

	//ignore events in SendCoinsFromAccountToModule and SendCoinsFromModuleToAccount
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeNewToken,
			sdk.NewAttribute(types.AttributeKeyIssuer, msg.FromAddress),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.ToAddress),
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Amount.GetDenom()),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyIssueFee, issueFee.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgNewTokenResponse{}, nil
}

//InflateToken, inflate tokens and send it to recipient
func (k msgServer) InflateToken(goCtx context.Context, msg *types.MsgInflateToken) (*types.MsgInflateTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	denom := msg.Amount.GetDenom()
	recipientAddr, _ := sdk.AccAddressFromBech32(msg.ToAddress)
	amount := sdk.NewCoins(msg.Amount)
	metadata := k.bankKeeper.GetDenomMetaData(ctx, denom)
	if metadata.Base != msg.Amount.GetDenom() {
		return nil, sdkerrors.Wrapf(types.ErrDenomNotExist, "%s", denom)
	}

	if metadata.Issuer != msg.FromAddress {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not authorized to inflate %s ", msg.FromAddress, denom)
	}

	k.bankKeeper.MintCoins(ctx, types.ModuleName, amount)
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddr, amount)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeInflateToken,
			sdk.NewAttribute(types.AttributeKeyIssuer, msg.FromAddress),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.ToAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})
	return &types.MsgInflateTokenResponse{}, nil
}

//BurnToken burn owned token
func (k msgServer) BurnToken(goCtx context.Context, msg *types.MsgBurnToken) (*types.MsgBurnTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	amount := sdk.NewCoins(msg.Amount)
	senderAddr, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, amount); err != nil {
		return nil, err
	}
	k.bankKeeper.BurnCoins(ctx, types.ModuleName, amount)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnToken,
			sdk.NewAttribute(types.AttributeKeyIssuer, msg.FromAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})
	return &types.MsgBurnTokenResponse{}, nil
}
