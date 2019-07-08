package sale

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper is the model object for the package sale module
type Keeper struct {
	// storeKey   sdk.StoreKey
	codec      *codec.Codec
	paramStore params.Subspace
}

// NewKeeper creates a new keeper of the sale Keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec) Keeper {
	return Keeper{
		// storeKey,
		codec,
		paramStore.WithKeyTable(params.NewKeyTable()),
	}
}

func (k Keeper) SetParams(ctx sdk.Context, params GenesisState) {
	k.paramStore.SetParamSet(ctx, &params)
}
