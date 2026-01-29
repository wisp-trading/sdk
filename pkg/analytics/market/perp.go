package market

import (
	"context"
	"fmt"

	"github.com/wisp-trading/wisp/pkg/types/connector"
	"github.com/wisp-trading/wisp/pkg/types/connector/perp"
	perpTypes "github.com/wisp-trading/wisp/pkg/types/data/stores/market/perp"
	"github.com/wisp-trading/wisp/pkg/types/portfolio"
)

// perpMarketService implements analytics.PerpMarket
// It embeds baseMarketService for common operations and adds perp-specific methods.
type perpMarketService struct {
	baseMarketService
	store perpTypes.MarketStore
}

// newPerpMarketService creates a perp market service
func newPerpMarketService(store perpTypes.MarketStore) *perpMarketService {
	return &perpMarketService{
		baseMarketService: newBaseMarketService(store),
		store:             store,
	}
}

// ========== Perp-specific methods (not inherited from base) ==========
func (s *perpMarketService) FundingRate(ctx context.Context, asset portfolio.Asset, exchange connector.ExchangeName) (*perp.FundingRate, error) {
	rate := s.store.GetFundingRate(asset, exchange)
	if rate == nil {
		return nil, fmt.Errorf("no funding rate found for %s on %s", asset.Symbol(), exchange)
	}
	return rate, nil
}

// FundingRates returns funding rates across all perp exchanges
func (s *perpMarketService) FundingRates(ctx context.Context, asset portfolio.Asset) map[connector.ExchangeName]perp.FundingRate {
	return s.store.GetFundingRatesForAsset(asset)
}

// GetAllAssetsWithFundingRates returns all assets that have funding rate data
func (s *perpMarketService) GetAllAssetsWithFundingRates(ctx context.Context) []portfolio.Asset {
	return s.store.GetAllAssetsWithFundingRates()
}
