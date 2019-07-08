package sale

// "time"
import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines sale module constants
const (
	DefaultParamspace = ModuleName
	StoreKey          = ModuleName
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
)

// Sale stores data about a sale
type Sale struct {
	//  fill me in
}

type MsgBuy struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	Coin      sdk.Coin       `json:"amount"`
}

// GetSignBytes Implements Msg.
func (msg MsgBuy) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgBuy) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}
