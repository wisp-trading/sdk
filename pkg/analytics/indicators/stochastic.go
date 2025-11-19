package indicators

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// StochasticResult represents the output of the Stochastic Oscillator calculation.
// The Stochastic Oscillator is a momentum indicator that compares a security's closing price
// to its price range over a given time period.
type StochasticResult struct {
	// K represents the %K line (fast stochastic), which measures the current close
	// relative to the high-low range over the specified period. Values range from 0 to 100.
	K decimal.Decimal

	// D represents the %D line (slow stochastic), which is the simple moving average
	// of the %K line over the specified period. Values range from 0 to 100.
	// The %D line provides a smoothed signal line.
	D decimal.Decimal
}

// Stochastic calculates the Stochastic Oscillator for the given price data.
//
// The Stochastic Oscillator is a momentum indicator that shows the location of the close
// relative to the high-low range over a set number of periods. The indicator ranges between
// 0 and 100. Readings above 80 are considered overbought, while readings below 20 are
// considered oversold.
//
// Parameters:
//   - highs: slice of high prices for each period
//   - lows: slice of low prices for each period
//   - closes: slice of closing prices for each period
//   - kPeriod: lookback period for calculating %K (typically 14)
//   - dPeriod: smoothing period for calculating %D from %K (typically 3)
//
// Returns:
//   - []StochasticResult: slice of stochastic values, each containing %K and %D
//   - error: if input validation fails or insufficient data is provided
//
// The function calculates:
//   - %K = (Current Close - Lowest Low) / (Highest High - Lowest Low) * 100
//   - %D = SMA(%K, dPeriod)
//
// Example:
//
//	highs := []decimal.Decimal{...}
//	lows := []decimal.Decimal{...}
//	closes := []decimal.Decimal{...}
//	results, err := Stochastic(highs, lows, closes, 14, 3)
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
