package analytics

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/wisp-trading/sdk/pkg/monitoring/profiling"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	analyticsTypes "github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Volatility calculates the standard deviation of returns for an asset.
// Returns annualized volatility as a percentage.
func (s *analytics) Volatility(ctx context.Context, asset portfolio.Pair, period int, opts ...analyticsTypes.AnalyticsOptions) (numerical.Decimal, error) {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("Volatility", time.Since(start))
		}
	}()

	options := s.parseOptions(opts...)

	prices, err := s.fetchClosePrices(ctx, asset, period+1, options)
	if err != nil {
		return numerical.Zero(), err
	}

	if len(prices) < 2 {
		return numerical.Zero(), fmt.Errorf("insufficient data for volatility calculation")
	}

	// Calculate returns
	returns := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		returns[i-1] = (prices[i] - prices[i-1]) / prices[i-1]
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
