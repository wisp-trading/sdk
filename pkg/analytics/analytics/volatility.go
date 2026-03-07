package analytics

import (
	"fmt"
	"math"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	analyticsTypes "github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Volatility calculates annualised volatility (std dev of returns) from the provided klines.
// interval is used to derive the annualisation factor (e.g. "1h", "4h", "1d").
func (s *analytics) Volatility(klines []connector.Kline, interval string) (numerical.Decimal, error) {
	if len(klines) < 2 {
		return numerical.Zero(), fmt.Errorf("insufficient data for volatility calculation")
	}

	prices := extractClose(klines)

	returns := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		returns[i-1] = (prices[i] - prices[i-1]) / prices[i-1]
	}

	var sum float64
	for _, r := range returns {
		sum += r
	}
	mean := sum / float64(len(returns))

	var variance float64
	for _, r := range returns {
		diff := r - mean
		variance += diff * diff
	}
	variance /= float64(len(returns))

	stdDev := math.Sqrt(variance)
	annualizedVol := stdDev * annualizationFactor(interval) * 100

	return numerical.NewFromFloat(annualizedVol), nil
}

// annualizationFactor returns sqrt(periods per year) for the given interval.
func annualizationFactor(interval string) float64 {
	periods, ok := analyticsTypes.PeriodsPerYear[interval]
	if !ok {
		periods = analyticsTypes.PeriodsPerYear[analyticsTypes.DefaultInterval]
	}
	return math.Sqrt(periods)
}
