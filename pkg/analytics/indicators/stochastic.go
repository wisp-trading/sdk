package indicators

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type StochasticResult struct {
	K decimal.Decimal // %K line (fast)
	D decimal.Decimal // %D line (slow - SMA of %K)
}

// Stochastic calculates the Stochastic Oscillator
// Requires high, low, and close prices
func Stochastic(highs, lows, closes []decimal.Decimal, kPeriod, dPeriod int) ([]StochasticResult, error) {
	if len(highs) != len(lows) || len(highs) != len(closes) {
		return nil, fmt.Errorf("price arrays must have equal length")
	}
	if len(closes) < kPeriod {
		return nil, fmt.Errorf("insufficient data: need %d prices, got %d", kPeriod, len(closes))
	}

	// Calculate %K values
	kValues := make([]decimal.Decimal, 0, len(closes)-kPeriod+1)

	for i := kPeriod - 1; i < len(closes); i++ {
		// Find highest high and lowest low in the period
		highestHigh := highs[i-kPeriod+1]
		lowestLow := lows[i-kPeriod+1]

		for j := i - kPeriod + 2; j <= i; j++ {
			if highs[j].GreaterThan(highestHigh) {
				highestHigh = highs[j]
			}
			if lows[j].LessThan(lowestLow) {
				lowestLow = lows[j]
			}
		}

		// Calculate %K
		denominator := highestHigh.Sub(lowestLow)
		var k decimal.Decimal
		if denominator.IsZero() {
			k = decimal.NewFromInt(50) // Default to middle if no range
		} else {
			k = closes[i].Sub(lowestLow).Div(denominator).Mul(decimal.NewFromInt(100))
		}
		kValues = append(kValues, k)
	}

	// Calculate %D (SMA of %K)
	if len(kValues) < dPeriod {
		return nil, fmt.Errorf("insufficient K values for D calculation: need %d, got %d", dPeriod, len(kValues))
	}

	result := make([]StochasticResult, 0, len(kValues)-dPeriod+1)

	for i := dPeriod - 1; i < len(kValues); i++ {
		sum := decimal.Zero
		for j := 0; j < dPeriod; j++ {
			sum = sum.Add(kValues[i-j])
		}
		d := sum.Div(decimal.NewFromInt(int64(dPeriod)))

		result = append(result, StochasticResult{
			K: kValues[i],
			D: d,
		})
	}

	return result, nil
}
