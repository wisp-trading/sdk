package market

import (
	spotTypes "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market/spot"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
)

// spotMarketService implements analytics.SpotMarket
// It embeds baseMarketService for all common operations.
// Spot markets have no additional methods beyond the common ones.
type spotMarketService struct {
	baseMarketService // Embedded - inherits Price, Prices, OrderBook, GetKlines, GetTradableQuantity
}

// newSpotMarketService creates a spot market service
func newSpotMarketService(store spotTypes.MarketStore) *spotMarketService {
	return &spotMarketService{
		baseMarketService: newBaseMarketService(store),
	}
}

// Spot returns the spot market service for spot-specific operations
func (s *marketService) Spot() analytics.SpotMarket {
	return s.spot
}
