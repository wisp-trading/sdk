package analytics

import (
	"context"
	"fmt"
	"math"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	analyticsTypes "github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
)

// fetchKlines gets klines from the market service, which automatically searches all stores
func (s *analytics) fetchKlines(ctx context.Context, asset portfolio.Asset, limit int, options analyticsTypes.AnalyticsOptions) ([]connector.Kline, error) {
	exchange := options.Exchange
	interval := options.Interval

	if exchange == "" {
		exchange = s.getDefaultExchange(ctx, asset)
		if exchange == "" {
			return nil, fmt.Errorf("no exchange data available for asset %s", asset.Symbol())
		}
	}

	klines := s.market.Klines(asset, exchange, interval, limit)

	if len(klines) == 0 {
		return nil, fmt.Errorf("no kline data available for asset %s on exchange %s", asset.Symbol(), exchange)
	}

	return klines, nil
}

// fetchClosePrices is a helper that fetches klines and extracts close prices as float64
func (s *analytics) fetchClosePrices(ctx context.Context, asset portfolio.Asset, limit int, options analyticsTypes.AnalyticsOptions) ([]float64, error) {
	klines, err := s.fetchKlines(ctx, asset, limit, options)
	if err != nil {
		return nil, err
	}

	prices := make([]float64, len(klines))
	for i, kline := range klines {
		prices[i] = kline.Close
	}

	return prices, nil
}

// getDefaultExchange returns the first available exchange for an asset
func (s *analytics) getDefaultExchange(ctx context.Context, asset portfolio.Asset) connector.ExchangeName {
	// Try spot exchanges first
	priceMap := s.market.Spot().Prices(ctx, asset)
	for exchange := range priceMap {
		return exchange
	}

	// Try perp exchanges
	priceMap = s.market.Perp().Prices(ctx, asset)
	for exchange := range priceMap {
		return exchange
	}

	return ""
}

// parseOptions extracts options with defaults
func (s *analytics) parseOptions(opts ...analyticsTypes.AnalyticsOptions) analyticsTypes.AnalyticsOptions {
	if len(opts) > 0 {
		options := opts[0]
		if options.Interval == "" {
			options.Interval = analyticsTypes.DefaultInterval
		}
		return options
	}
	return analyticsTypes.AnalyticsOptions{
		Interval: analyticsTypes.DefaultInterval,
	}
}

// getAnnualizationFactor returns the factor to annualize volatility based on interval
// Formula: sqrt(periods_per_year)
func (s *analytics) getAnnualizationFactor(interval string) float64 {
	periods, ok := analyticsTypes.PeriodsPerYear[interval]
	if !ok {
		// Default to hourly if unknown interval
		periods = analyticsTypes.PeriodsPerYear[analyticsTypes.DefaultInterval]
	}

	return math.Sqrt(periods)
}
