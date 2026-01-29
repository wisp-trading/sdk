package indicators

import (
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// SMA calculates the Simple Moving Average for the given prices and period.
func SMA(prices []float64, period int) (numerical.Decimal, error) {
	if len(prices) < period {
		return numerical.Zero(), fmt.Errorf("insufficient data: need %d prices, got %d", period, len(prices))
	}

	sum := 0.0
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i]
	}

	return numerical.NewFromFloat(sum / float64(period)), nil
}
