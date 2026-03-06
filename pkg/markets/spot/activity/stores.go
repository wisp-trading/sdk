package activity

import (
	positionStore "github.com/wisp-trading/sdk/pkg/markets/base/store/activity/position"
	tradeStore "github.com/wisp-trading/sdk/pkg/markets/base/store/activity/trade"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
)

func NewSpotPositions() spotTypes.SpotPositions {
	return positionStore.NewStore()
}

func NewSpotTrades() spotTypes.SpotTrades {
	return tradeStore.NewStore()
}
