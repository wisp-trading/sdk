package market

import (
	spotTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market/spot"
)

// spotMarketService implements analytics.SpotMarket
// It embeds baseMarketService for all common operations.
// Spot markets have no additional methods beyond the common ones.
type spotMarketService struct {
	baseMarketService
}

// newSpotMarketService creates a spot market service
func newSpotMarketService(store spotTypes.MarketStore) *spotMarketService {
	return &spotMarketService{
		baseMarketService: newBaseMarketService(store),
	}
}
