package market

import (
	marketTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	perpAnalytics "github.com/wisp-trading/sdk/pkg/markets/perp/analytics"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	spotAnalytics "github.com/wisp-trading/sdk/pkg/markets/spot/analytics"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics/aggregator"
)

// marketService is the concrete implementation of aggregator.Market.
// It holds typed spot and perp services — all data access goes through them directly.
type marketService struct {
	spot spotTypes.SpotMarket
	perp perpTypes.PerpMarket
}

// NewMarketService creates a new market aggregator from the market registry.
func NewMarketService(store marketTypes.MarketRegistry) aggregator.Market {
	var spot spotTypes.SpotMarket
	var perp perpTypes.PerpMarket

	if spotStore := store.Get(connector.MarketTypeSpot); spotStore != nil {
		if typed, ok := spotStore.(spotTypes.MarketStore); ok {
			spot = spotAnalytics.New(typed)
		}
	}

	if perpStore := store.Get(connector.MarketTypePerp); perpStore != nil {
		if typed, ok := perpStore.(perpTypes.MarketStore); ok {
			perp = perpAnalytics.New(typed)
		}
	}

	return &marketService{spot: spot, perp: perp}
}

func (s *marketService) Spot() spotTypes.SpotMarket { return s.spot }
func (s *marketService) Perp() perpTypes.PerpMarket { return s.perp }
