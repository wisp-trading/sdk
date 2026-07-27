package analytics

import (
	"context"
	"fmt"

	baseAnalytics "github.com/wisp-trading/sdk/pkg/markets/base/analytics"
	"github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// service implements types.PerpMarket: shared pair analytics + funding.
type service struct {
	*baseAnalytics.Service
	store types.MarketStore
}

// New creates a new perp analytics service.
func New(store types.MarketStore) types.PerpMarket {
	return &service{
		Service: baseAnalytics.NewService(store),
		store:   store,
	}
}

// FundingRate returns the latest funding rate for an asset on the specified exchange.
func (s *service) FundingRate(ctx context.Context, asset portfolio.Pair, exchange connector.ExchangeName) (*perpConn.FundingRate, error) {
	rate := s.store.GetFundingRate(asset, exchange)
	if rate == nil {
		return nil, fmt.Errorf("no funding rate found for %s on %s", asset.Symbol(), exchange)
	}
	return rate, nil
}

// FundingRates returns funding rates across all perp exchanges for an asset.
func (s *service) FundingRates(ctx context.Context, asset portfolio.Pair) map[connector.ExchangeName]perpConn.FundingRate {
	return s.store.GetFundingRatesForAsset(asset)
}

// GetAllAssetsWithFundingRates returns all assets that have funding rate data.
func (s *service) GetAllAssetsWithFundingRates(ctx context.Context) []portfolio.Pair {
	return s.store.GetAllAssetsWithFundingRates()
}

var _ types.PerpMarket = (*service)(nil)
