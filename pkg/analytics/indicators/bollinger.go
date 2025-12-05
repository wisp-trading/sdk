package indicators

import (
	"fmt"
	"math"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

// bollingerBandsFloat64 is the internal high-performance implementation using float64
func bollingerBandsFloat64(prices []float64, period int, stdDev float64) ([]analytics.BollingerBandsResult, error) {
	if len(prices) < period {
		return nil, fmt.Errorf("insufficient data: need %d prices, got %d", period, len(prices))
	}

	result := make([]analytics.BollingerBandsResult, 0, len(prices)-period+1)

	for i := period - 1; i < len(prices); i++ {
		// Calculate SMA (middle band)
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += prices[i-j]
		}
		sma := sum / float64(period)

		// Calculate standard deviation
		variance := 0.0
		for j := 0; j < period; j++ {
			diff := prices[i-j] - sma
			variance += diff * diff
		}
		variance /= float64(period)
		sd := math.Sqrt(variance) * stdDev

		result = append(result, analytics.BollingerBandsResult{
			Upper:  numerical.NewFromFloat(sma + sd),
			Middle: numerical.NewFromFloat(sma),
			Lower:  numerical.NewFromFloat(sma - sd),
		})
	}

	return result, nil
}

// BollingerBands calculates Bollinger Bands with automatic conversion from numerical.Decimal
func BollingerBands(prices []numerical.Decimal, period int, stdDev float64) ([]analytics.BollingerBandsResult, error) {
	// Convert to float64 for high-performance calculation
	pricesFloat := make([]float64, len(prices))
	for i, p := range prices {
		pricesFloat[i], _ = p.Float64()
	}

	return bollingerBandsFloat64(pricesFloat, period, stdDev)
}
