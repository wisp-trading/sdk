package indicators

import (
	"fmt"
	"math"

	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// BollingerBands calculates Bollinger Bands for the given prices.
func BollingerBands(prices []float64, period int, stdDev float64) (analytics.BollingerBandsResult, error) {
	if len(prices) < period {
		return analytics.BollingerBandsResult{}, fmt.Errorf("insufficient data: need %d prices, got %d", period, len(prices))
	}

	// Calculate SMA
	sum := 0.0
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i]
	}
	sma := sum / float64(period)

	// Calculate standard deviation
	variance := 0.0
	for i := len(prices) - period; i < len(prices); i++ {
		diff := prices[i] - sma
		variance += diff * diff
	}
	variance /= float64(period)
	sd := math.Sqrt(variance) * stdDev

	return analytics.BollingerBandsResult{
		Upper:  numerical.NewFromFloat(sma + sd),
		Middle: numerical.NewFromFloat(sma),
		Lower:  numerical.NewFromFloat(sma - sd),
	}, nil
}
