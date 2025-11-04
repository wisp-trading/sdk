package indicators

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// ATR calculates the Average True Range
func ATR(highs, lows, closes []decimal.Decimal, period int) ([]decimal.Decimal, error) {
	if len(highs) != len(lows) || len(highs) != len(closes) {
		return nil, fmt.Errorf("price arrays must have equal length")
	}
	if len(closes) < period+1 {
		return nil, fmt.Errorf("insufficient data: need %d prices, got %d", period+1, len(closes))
	}

	// Calculate True Range for each period
	trueRanges := make([]decimal.Decimal, len(closes)-1)

	for i := 1; i < len(closes); i++ {
		highLow := highs[i].Sub(lows[i])
		highClose := highs[i].Sub(closes[i-1]).Abs()
		lowClose := lows[i].Sub(closes[i-1]).Abs()

		tr := highLow
		if highClose.GreaterThan(tr) {
			tr = highClose
		}
		if lowClose.GreaterThan(tr) {
			tr = lowClose
		}

		trueRanges[i-1] = tr
	}

	// Calculate ATR using EMA-like smoothing
	result := make([]decimal.Decimal, 0, len(trueRanges)-period+1)

	// First ATR is simple average
	sum := decimal.Zero
	for i := 0; i < period; i++ {
		sum = sum.Add(trueRanges[i])
	}
	atr := sum.Div(decimal.NewFromInt(int64(period)))
	result = append(result, atr)

	// Subsequent ATRs use smoothing
	periodDecimal := decimal.NewFromInt(int64(period))
	for i := period; i < len(trueRanges); i++ {
		atr = atr.Mul(periodDecimal.Sub(decimal.NewFromInt(1))).Add(trueRanges[i]).Div(periodDecimal)
		result = append(result, atr)
	}

	return result, nil
}
