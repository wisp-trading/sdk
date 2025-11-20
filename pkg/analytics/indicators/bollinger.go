package indicators

import (
	"fmt"
	"math"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
	"github.com/shopspring/decimal"
)

// BollingerBands calculates Bollinger Bands
func BollingerBands(prices []decimal.Decimal, period int, stdDev float64) ([]analytics.BollingerBandsResult, error) {
	if len(prices) < period {
		return nil, fmt.Errorf("insufficient data: need %d prices, got %d", period, len(prices))
	}

	result := make([]analytics.BollingerBandsResult, 0, len(prices)-period+1)

	for i := period - 1; i < len(prices); i++ {
		// Calculate SMA (middle band)
		sum := decimal.Zero
		for j := 0; j < period; j++ {
			sum = sum.Add(prices[i-j])
		}
		sma := sum.Div(decimal.NewFromInt(int64(period)))

		// Calculate standard deviation
		variance := 0.0
		for j := 0; j < period; j++ {
			diff, _ := prices[i-j].Sub(sma).Float64()
			variance += diff * diff
		}
		variance /= float64(period)
		sd := math.Sqrt(variance)

		stdDevDecimal := decimal.NewFromFloat(sd * stdDev)

		result = append(result, analytics.BollingerBandsResult{
			Upper:  sma.Add(stdDevDecimal),
			Middle: sma,
			Lower:  sma.Sub(stdDevDecimal),
		})
	}

	return result, nil
}
