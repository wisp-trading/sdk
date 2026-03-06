package analytics

import (
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	analyticsTypes "github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Trend analyzes the price trend for an asset using linear regression.
// Returns trend direction and strength.
func (s *analytics) Trend(klines []connector.Kline) (*analyticsTypes.TrendResult, error) {
	if len(klines) < 2 {
		return nil, fmt.Errorf("insufficient data for trend calculation")
	}

	prices := extractClose(klines)
	n := float64(len(prices))
	var sumX, sumY, sumXY, sumX2 float64

	for i, price := range prices {
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
	for i, price := range prices {
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
	const slopeThreshold = 0.01
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

// extractClose is a helper function to extract close prices from klines.
func extractClose(klines []connector.Kline) []float64 {
	prices := make([]float64, len(klines))
	for i, kline := range klines {
		prices[i] = kline.Close
	}
	return prices
}
