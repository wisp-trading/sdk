package market

import (
	perpAnalytics "github.com/wisp-trading/sdk/pkg/markets/perp/analytics"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	spotAnalytics "github.com/wisp-trading/sdk/pkg/markets/spot/analytics"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics/aggregator"
)

type marketService struct {
	spot spotTypes.SpotMarket
	perp perpTypes.PerpMarket
}

// NewMarketService creates a new market aggregator from domain stores directly.
func NewMarketService(
	spotStore spotTypes.MarketStore,
	perpStore perpTypes.MarketStore,
) aggregator.Market {
	return &marketService{
		spot: spotAnalytics.New(spotStore),
		perp: perpAnalytics.New(perpStore),
	}
}

func (s *marketService) Spot() spotTypes.SpotMarket { return s.spot }
func (s *marketService) Perp() perpTypes.PerpMarket { return s.perp }
