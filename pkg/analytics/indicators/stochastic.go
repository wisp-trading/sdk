package indicators

import (
	"fmt"
	"math"

	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Stochastic calculates the Stochastic Oscillator for the given price data.
// Input prices are float64 for performance, output is numerical.Decimal for precision.
func Stochastic(highs, lows, closes []float64, kPeriod, dPeriod int) (analytics.StochasticResult, error) {
	if len(highs) != len(lows) || len(highs) != len(closes) {
		return analytics.StochasticResult{}, fmt.Errorf("price arrays must have equal length")
	}
	if len(closes) < kPeriod+dPeriod-1 {
		return analytics.StochasticResult{}, fmt.Errorf("insufficient data: need %d prices, got %d", kPeriod+dPeriod-1, len(closes))
	}

	kValues := make([]float64, 0, len(closes)-kPeriod+1)

	for i := kPeriod - 1; i < len(closes); i++ {
		highestHigh := highs[i-kPeriod+1]
		lowestLow := lows[i-kPeriod+1]

		for j := i - kPeriod + 2; j <= i; j++ {
			if highs[j] > highestHigh {
				highestHigh = highs[j]
			}
			if lows[j] < lowestLow {
				lowestLow = lows[j]
			}
		}

		denominator := highestHigh - lowestLow
		var k float64
		if math.Abs(denominator) < 1e-10 {
			k = 50.0
		} else {
			k = ((closes[i] - lowestLow) / denominator) * 100.0
		}
		kValues = append(kValues, k)
	}

	if len(kValues) < dPeriod {
		return analytics.StochasticResult{}, fmt.Errorf("insufficient K values for D calculation: need %d, got %d", dPeriod, len(kValues))
	}

	sum := 0.0
	for i := len(kValues) - dPeriod; i < len(kValues); i++ {
		sum += kValues[i]
	}
	d := sum / float64(dPeriod)

	return analytics.StochasticResult{
		K: numerical.NewFromFloat(kValues[len(kValues)-1]),
		D: numerical.NewFromFloat(d),
	}, nil
}
