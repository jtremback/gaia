package sale

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	params "github.com/cosmos/cosmos-sdk/x/params"
)

// NewHandler creates a new handler for sale module
func NewHandler(distrKeeper Keeper, bankKeeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {

		case MsgBuy:
			return handleMsgSend(ctx, distrKeeper, bankKeeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized sale message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgBuy
func handleMsgBuy(ctx sdk.Context, k Keeper, distrKeeper distribution.Keeper, bankKeeper bank.Keeper, msg MsgBuy) sdk.Result {
	// if !k.GetSendEnabled(ctx) {
	// 	return types.ErrSendDisabled(k.Codespace()).Result()
	// }

	// TODO: Check that coin is correct denom

	// Retreive params
	var buyingDenom string
	k.paramStore.Get(ctx, []byte("BuyingDenom"), &buyingDenom)

	var sellingDenom string
	k.paramStore.Get(ctx, []byte("SellingDenom"), &sellingDenom)

	var price sdk.Dec
	k.paramStore.Get(ctx, []byte("Price"), &price)

	// Throw if trying to buy with the wrong coins
	if msg.Coin.Denom != buyingDenom {
		// TODO: Error
	}

	// Take buying coins from buyer's account
	_, err := bankKeeper.SubtractCoins(ctx, msg.FromAddress, sdk.NewCoins(msg.Coin))
	if err != nil {
		return err.Result()
	}

	// Get community pool from storage
	feePool := distrKeeper.GetFeePool(ctx)

	// Add buying coins to community pool
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoins(sdk.NewCoins(msg.Coin)))
	
	// Find number of selling coins to sell
	sellingCoinsAmount := price.Mul(sdk.NewDecFromInt(msg.Coin.Amount))
	var sellingCoins sdk.DecCoins

	sellingCoins[0] = sdk.NewDecCoinFromDec(sellingDenom, sellingCoinsAmount)

	// Take selling coins from community pool
	newPool, negative := feePool.CommunityPool.SafeSub(sellingCoins)
	if negative {
		// TODO: Error
	}
	
	// Save new state of community pool
	feePool.CommunityPool = newPool
	distrKeeper.SetFeePool(ctx, feePool)

	// Transfer selling coins to buyer
	_, err := bankKeeper.AddCoins(ctx, msg.FromAddress, sellingCoins)
	if err != nil {
		return err.Result()
	}
	
	resTags := sdk.NewTags(
		sdk.TagCategory, "sale",
		"buyer", msg.FromAddress.String(),
	)

	return sdk.Result{
		Tags: resTags,
	}
}
