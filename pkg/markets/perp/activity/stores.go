package activity

import (
	positionStore "github.com/wisp-trading/sdk/pkg/markets/base/store/activity/position"
	tradeStore "github.com/wisp-trading/sdk/pkg/markets/base/store/activity/trade"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
)

func NewPerpPositions() perpTypes.PerpPositions {
	return positionStore.NewStore()
}

func NewPerpTrades() perpTypes.PerpTrades {
	return tradeStore.NewStore()
}
