package aggregator

import (
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
)

// Market provides market data access across market types.
type Market interface {
	// Spot returns spot-specific market service
	Spot() spotTypes.SpotMarket

	// Perp returns perpetual-specific market service
	Perp() perpTypes.PerpMarket
}
