package analytics

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/monitoring/profiling"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
	analyticsTypes "github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// analytics provides user-friendly methods for market analytics.
type analytics struct {
	store market.MarketData
}

// NewAnalyticsService creates a new analytics
func NewAnalyticsService(store market.MarketData) analyticsTypes.Analytics {
	return &analytics{
		store: store,
	}
}

// Volatility calculates the standard deviation of returns for an asset.
// Returns annualized volatility as a percentage.
func (s *analytics) Volatility(ctx context.Context, asset portfolio.Asset, period int, opts ...analyticsTypes.AnalyticsOptions) (numerical.Decimal, error) {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("Volatility", time.Since(start))
		}
	}()

	options := s.parseOptions(opts...)

	prices, err := s.fetchClosePrices(asset, period+1, opts...)
	if err != nil {
		return numerical.Zero(), err
	}

	if len(prices) < 2 {
		return numerical.Zero(), fmt.Errorf("insufficient data for volatility calculation")
	}

	// Calculate returns
	returns := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		ret := prices[i].Sub(prices[i-1]).Div(prices[i-1])
		retFloat, _ := ret.Float64()
		returns[i-1] = retFloat
	}

	// Calculate mean return
	var sum float64
	for _, r := range returns {
		sum += r
	}
	mean := sum / float64(len(returns))

	// Calculate variance
	var variance float64
	for _, r := range returns {
		diff := r - mean
		variance += diff * diff
	}
	variance /= float64(len(returns))

	// Standard deviation
	stdDev := math.Sqrt(variance)

	// Calculate annualization factor based on interval
	annualizationFactor := s.getAnnualizationFactor(options.Interval)
	annualizedVol := stdDev * annualizationFactor * 100 // Convert to percentage

	return numerical.NewFromFloat(annualizedVol), nil
}

// Trend analyzes the price trend for an asset using linear regression.
// Returns trend direction and strength.
func (s *analytics) Trend(ctx context.Context, asset portfolio.Asset, period int, opts ...analyticsTypes.AnalyticsOptions) (*analyticsTypes.TrendResult, error) {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("Trend", time.Since(start))
		}
	}()

	prices, err := s.fetchClosePrices(asset, period, opts...)
	if err != nil {
		return nil, err
	}

	if len(prices) < 2 {
		return nil, fmt.Errorf("insufficient data for trend calculation")
	}

	// Convert prices to float64 for calculation
	priceFloats := make([]float64, len(prices))
	for i, p := range prices {
		priceFloats[i], _ = p.Float64()
	}

	// Calculate linear regression
	n := float64(len(priceFloats))
	var sumX, sumY, sumXY, sumX2 float64

	for i, price := range priceFloats {
		x := float64(i)
		sumX += x
		sumY += price
		sumXY += x * price
		sumX2 += x * x
	}

	// Calculate slope (m) and intercept (b) for y = mx + b
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)

	// Calculate R-squared to measure trend strength
	meanY := sumY / n
	var ssTotal, ssResidual float64
	for i, price := range priceFloats {
		x := float64(i)
		predicted := slope*x + (sumY-slope*sumX)/n
		ssTotal += (price - meanY) * (price - meanY)
		ssResidual += (price - predicted) * (price - predicted)
	}

	rSquared := 1 - (ssResidual / ssTotal)
	if rSquared < 0 {
		rSquared = 0
	}

	strength := numerical.NewFromFloat(rSquared * 100) // Convert to percentage
	slopeDecimal := numerical.NewFromFloat(slope)

	// Determine direction based on slope
	// Use a threshold to avoid calling tiny slopes a trend
	slopeThreshold := 0.01
	var direction analyticsTypes.TrendDirection

	if slope > slopeThreshold && rSquared > 0.3 {
		direction = analyticsTypes.TrendBullish
	} else if slope < -slopeThreshold && rSquared > 0.3 {
		direction = analyticsTypes.TrendBearish
	} else {
		direction = analyticsTypes.TrendNeutral
	}

	return &analyticsTypes.TrendResult{
		Direction: direction,
		Strength:  strength,
		Slope:     slopeDecimal,
	}, nil
}

// VolumeAnalysis detects volume patterns and spikes.
func (s *analytics) VolumeAnalysis(ctx context.Context, asset portfolio.Asset, period int, opts ...analyticsTypes.AnalyticsOptions) (*analyticsTypes.VolumeAnalysis, error) {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("VolumeAnalysis", time.Since(start))
		}
	}()

	options := s.parseOptions(opts...)
	exchange := options.Exchange
	interval := options.Interval

	if exchange == "" {
		exchange = s.getDefaultExchange(asset)
		if exchange == "" {
			return nil, fmt.Errorf("no exchange data available for asset %s", asset.Symbol())
		}
	}

	klines := s.store.GetKlines(asset, exchange, interval, period+1)
	if len(klines) < 2 {
		return nil, fmt.Errorf("insufficient kline data for volume analysis")
	}

	// Get current volume (latest kline)
	currentVolume := klines[len(klines)-1].Volume

	// Calculate average volume over period
	var totalVolume numerical.Decimal
	for i := 0; i < len(klines)-1; i++ {
		totalVolume = totalVolume.Add(klines[i].Volume)
	}
	avgVolume := totalVolume.Div(numerical.NewFromInt(int64(len(klines) - 1)))

	// Calculate volume ratio
	var volumeRatio numerical.Decimal
	if avgVolume.IsZero() {
		volumeRatio = numerical.Zero()
	} else {
		volumeRatio = currentVolume.Div(avgVolume)
	}

	// Detect spike (current > 2x average)
	spikeThreshold := numerical.NewFromInt(2)
	isSpike := volumeRatio.GreaterThan(spikeThreshold)

	// Determine volume trend
	// Compare first half average to second half average
	midPoint := len(klines) / 2
	var firstHalfVolume, secondHalfVolume numerical.Decimal

	for i := 0; i < midPoint; i++ {
		firstHalfVolume = firstHalfVolume.Add(klines[i].Volume)
	}
	for i := midPoint; i < len(klines); i++ {
		secondHalfVolume = secondHalfVolume.Add(klines[i].Volume)
	}

	firstHalfAvg := firstHalfVolume.Div(numerical.NewFromInt(int64(midPoint)))
	secondHalfAvg := secondHalfVolume.Div(numerical.NewFromInt(int64(len(klines) - midPoint)))

	var volumeTrend analyticsTypes.TrendDirection
	trendThreshold := numerical.NewFromFloat(1.2) // 20% increase/decrease

	if !firstHalfAvg.IsZero() && secondHalfAvg.Div(firstHalfAvg).GreaterThan(trendThreshold) {
		volumeTrend = analyticsTypes.TrendBullish // Increasing volume
	} else if !secondHalfAvg.IsZero() && firstHalfAvg.Div(secondHalfAvg).GreaterThan(trendThreshold) {
		volumeTrend = analyticsTypes.TrendBearish // Decreasing volume
	} else {
		volumeTrend = analyticsTypes.TrendNeutral
	}

	return &analyticsTypes.VolumeAnalysis{
		CurrentVolume: currentVolume,
		AverageVolume: avgVolume,
		VolumeRatio:   volumeRatio,
		IsVolumeSpike: isSpike,
		VolumeTrend:   volumeTrend,
	}, nil
}

// GetPriceChange calculates price statistics over a period.
func (s *analytics) GetPriceChange(ctx context.Context, asset portfolio.Asset, period int, opts ...analyticsTypes.AnalyticsOptions) (*analyticsTypes.PriceChange, error) {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("GetPriceChange", time.Since(start))
		}
	}()

	options := s.parseOptions(opts...)
	exchange := options.Exchange
	interval := options.Interval

	if exchange == "" {
		exchange = s.getDefaultExchange(asset)
		if exchange == "" {
			return nil, fmt.Errorf("no exchange data available for asset %s", asset.Symbol())
		}
	}

	klines := s.store.GetKlines(asset, exchange, interval, period)
	if len(klines) < 2 {
		return nil, fmt.Errorf("insufficient kline data for price change calculation")
	}

	startPrice := klines[0].Open
	endPrice := klines[len(klines)-1].Close

	// Find high and low
	highPrice := klines[0].High
	lowPrice := klines[0].Low

	for _, kline := range klines {
		if kline.High.GreaterThan(highPrice) {
			highPrice = kline.High
		}
		if kline.Low.LessThan(lowPrice) {
			lowPrice = kline.Low
		}
	}

	change := endPrice.Sub(startPrice)
	changePercent := change.Div(startPrice).Mul(numerical.NewFromInt(100))
	priceRange := highPrice.Sub(lowPrice)
	priceRangePercent := priceRange.Div(startPrice).Mul(numerical.NewFromInt(100))

	return &analyticsTypes.PriceChange{
		StartPrice:        startPrice,
		EndPrice:          endPrice,
		Change:            change,
		ChangePercent:     changePercent,
		HighPrice:         highPrice,
		LowPrice:          lowPrice,
		PriceRange:        priceRange,
		PriceRangePercent: priceRangePercent,
	}, nil
}

// fetchClosePrices is a helper that fetches klines and extracts close prices
func (s *analytics) fetchClosePrices(asset portfolio.Asset, limit int, opts ...analyticsTypes.AnalyticsOptions) ([]numerical.Decimal, error) {
	options := s.parseOptions(opts...)
	exchange := options.Exchange
	interval := options.Interval

	if exchange == "" {
		exchange = s.getDefaultExchange(asset)
		if exchange == "" {
			return nil, fmt.Errorf("no exchange data available for asset %s", asset.Symbol())
		}
	}

	klines := s.store.GetKlines(asset, exchange, interval, limit)
	if len(klines) == 0 {
		return nil, fmt.Errorf("no kline data available for asset %s on exchange %s", asset.Symbol(), exchange)
	}

	prices := make([]numerical.Decimal, len(klines))
	for i, kline := range klines {
		prices[i] = kline.Close
	}

	return prices, nil
}

// getDefaultExchange returns the first available exchange for an asset
func (s *analytics) getDefaultExchange(asset portfolio.Asset) connector.ExchangeName {
	priceMap := s.store.GetAssetPrices(asset)
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
