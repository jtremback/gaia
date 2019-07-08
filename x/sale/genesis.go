package sale

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	BuyingDenom  string  `json:"buying_denom"`
	SellingDenom string  `json:"selling_denom"`
	Price        sdk.Dec `json:"price"`
}

func (p *GenesisState) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{[]byte("BuyingDenom"), &p.BuyingDenom},
		{[]byte("SellingDenom"), &p.SellingDenom},
		{[]byte("Price"), &p.Price},
	}
}

// NewGenesisState creates a new genesis state.
// func NewGenesisState() GenesisState {
// 	return GenesisState{
// 		BuyingDenom:  string,
// 		SellingDenom: string,
// 		Price:        sdk.NewDec(1),
// 	}
// }

// DefaultGenesisState returns a default genesis state
// func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis initializes story state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		//		Sales: keeper.Sales(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {

	return nil
}
