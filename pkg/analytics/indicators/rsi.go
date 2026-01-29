package indicators

import (
	"fmt"
	"math"

	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// RSI calculates the Relative Strength Index for the given prices and period.
func RSI(prices []float64, period int) (numerical.Decimal, error) {
	if len(prices) < period+1 {
		return numerical.Zero(), fmt.Errorf("insufficient data: need %d prices, got %d", period+1, len(prices))
	}

	avgGain := 0.0
	avgLoss := 0.0

	for i := 1; i <= period; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			avgGain += change
		} else {
			avgLoss += math.Abs(change)
		}
	}

	avgGain /= float64(period)
	avgLoss /= float64(period)

	for i := period + 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]

		if change > 0 {
			avgGain = (avgGain*float64(period-1) + change) / float64(period)
			avgLoss = (avgLoss * float64(period-1)) / float64(period)
		} else {
			avgGain = (avgGain * float64(period-1)) / float64(period)
			avgLoss = (avgLoss*float64(period-1) + math.Abs(change)) / float64(period)
		}
	}

	if math.Abs(avgLoss) < 1e-10 {
		return numerical.NewFromFloat(100.0), nil
	}

	rs := avgGain / avgLoss
	rsi := 100.0 - (100.0 / (1.0 + rs))

	return numerical.NewFromFloat(rsi), nil
}
