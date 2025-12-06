package indicators

import (
	"fmt"
	"math"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

func stochasticFloat64(highs, lows, closes []float64, kPeriod, dPeriod int) (float64, float64, error) {
	if len(highs) != len(lows) || len(highs) != len(closes) {
		return 0, 0, fmt.Errorf("price arrays must have equal length")
	}
	if len(closes) < kPeriod+dPeriod-1 {
		return 0, 0, fmt.Errorf("insufficient data: need %d prices, got %d", kPeriod+dPeriod-1, len(closes))
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
		return 0, 0, fmt.Errorf("insufficient K values for D calculation: need %d, got %d", dPeriod, len(kValues))
	}

	sum := 0.0
	for i := len(kValues) - dPeriod; i < len(kValues); i++ {
		sum += kValues[i]
	}
	d := sum / float64(dPeriod)

	return kValues[len(kValues)-1], d, nil
}

func Stochastic(highs, lows, closes []numerical.Decimal, kPeriod, dPeriod int) (analytics.StochasticResult, error) {
	if len(highs) != len(lows) || len(highs) != len(closes) {
		return analytics.StochasticResult{}, fmt.Errorf("price arrays must have equal length")
	}

	highsFloat := make([]float64, len(highs))
	lowsFloat := make([]float64, len(lows))
	closesFloat := make([]float64, len(closes))

	for i := range highs {
		highsFloat[i], _ = highs[i].Float64()
		lowsFloat[i], _ = lows[i].Float64()
		closesFloat[i], _ = closes[i].Float64()
	}

	k, d, err := stochasticFloat64(highsFloat, lowsFloat, closesFloat, kPeriod, dPeriod)
	if err != nil {
		return analytics.StochasticResult{}, err
	}

	return analytics.StochasticResult{
		K: numerical.NewFromFloat(k),
		D: numerical.NewFromFloat(d),
	}, nil
}
