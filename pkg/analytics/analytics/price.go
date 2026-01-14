package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/monitoring/profiling"
	analyticsTypes "github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// GetPriceChange calculates price statistics over a period.
func (s *analytics) GetPriceChange(ctx context.Context, asset portfolio.Asset, period int, opts ...analyticsTypes.AnalyticsOptions) (*analyticsTypes.PriceChange, error) {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("GetPriceChange", time.Since(start))
		}
	}()

	options := s.parseOptions(opts...)

	klines, err := s.fetchKlines(ctx, asset, period, options)
	if err != nil {
		return nil, err
	}

	if len(klines) < 2 {
		return nil, fmt.Errorf("insufficient kline data for price change calculation")
	}

	startPrice := klines[0].Open
	endPrice := klines[len(klines)-1].Close

	// Find high and low
	highPrice := klines[0].High
	lowPrice := klines[0].Low

	for _, kline := range klines {
		if kline.High > highPrice {
			highPrice = kline.High
		}
		if kline.Low < lowPrice {
			lowPrice = kline.Low
		}
	}

	change := endPrice - startPrice
	changePercent := (change / startPrice) * 100
	priceRange := highPrice - lowPrice
	priceRangePercent := (priceRange / startPrice) * 100

	return &analyticsTypes.PriceChange{
		StartPrice:        numerical.NewFromFloat(startPrice),
		EndPrice:          numerical.NewFromFloat(endPrice),
		Change:            numerical.NewFromFloat(change),
		ChangePercent:     numerical.NewFromFloat(changePercent),
		HighPrice:         numerical.NewFromFloat(highPrice),
		LowPrice:          numerical.NewFromFloat(lowPrice),
		PriceRange:        numerical.NewFromFloat(priceRange),
		PriceRangePercent: numerical.NewFromFloat(priceRangePercent),
	}, nil
}
