package indicators

import (
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// emaFloat64 is the internal float64 implementation for use by other indicators
func emaFloat64(prices []float64, period int) (float64, error) {
	if len(prices) < period {
		return 0, fmt.Errorf("insufficient data: need %d prices, got %d", period, len(prices))
	}

	multiplier := 2.0 / float64(period+1)

	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema := sum / float64(period)

	for i := period; i < len(prices); i++ {
		ema = (prices[i]-ema)*multiplier + ema
	}

	return ema, nil
}

// EMA calculates the Exponential Moving Average for the given prices and period.
// Input prices are float64 for performance, output is numerical.Decimal for precision.
func EMA(prices []float64, period int) (numerical.Decimal, error) {
	ema, err := emaFloat64(prices, period)
	if err != nil {
		return numerical.Zero(), err
	}
	return numerical.NewFromFloat(ema), nil
}
