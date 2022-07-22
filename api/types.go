package api

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Candle struct {
	Price  sdk.Dec
	Volume sdk.Dec
	Time   int64
}

type Ticker struct {
	Price  sdk.Dec
	Volume sdk.Dec
}
