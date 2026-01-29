package indicators

import (
	"fmt"
	"math"

	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// ATR calculates the Average True Range for the given price data.
func ATR(highs, lows, closes []float64, period int) (numerical.Decimal, error) {
	if len(highs) != len(lows) || len(highs) != len(closes) {
		return numerical.Zero(), fmt.Errorf("price arrays must have equal length")
	}
	if len(closes) < period+1 {
		return numerical.Zero(), fmt.Errorf("insufficient data: need %d prices, got %d", period+1, len(closes))
	}

	sum := 0.0
	for i := 1; i <= period; i++ {
		highLow := highs[i] - lows[i]
		highClose := math.Abs(highs[i] - closes[i-1])
		lowClose := math.Abs(lows[i] - closes[i-1])

		tr := highLow
		if highClose > tr {
			tr = highClose
		}
		if lowClose > tr {
			tr = lowClose
		}

		sum += tr
	}

	atr := sum / float64(period)

	for i := period + 1; i < len(closes); i++ {
		highLow := highs[i] - lows[i]
		highClose := math.Abs(highs[i] - closes[i-1])
		lowClose := math.Abs(lows[i] - closes[i-1])

		tr := highLow
		if highClose > tr {
			tr = highClose
		}
		if lowClose > tr {
			tr = lowClose
		}

		atr = (atr*float64(period-1) + tr) / float64(period)
	}

	return numerical.NewFromFloat(atr), nil
}
