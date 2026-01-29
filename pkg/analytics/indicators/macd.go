package indicators

import (
	"fmt"

	"github.com/wisp-trading/wisp/pkg/types/wisp/analytics"
	"github.com/wisp-trading/wisp/pkg/types/wisp/numerical"
)

// MACD calculates the Moving Average Convergence Divergence for the given prices.
func MACD(prices []float64, fastPeriod, slowPeriod, signalPeriod int) (analytics.MACDResult, error) {
	if len(prices) < slowPeriod {
		return analytics.MACDResult{}, fmt.Errorf("insufficient data: need %d prices, got %d", slowPeriod, len(prices))
	}

	// Step 1: Calculate fast and slow EMAs
	fastEMA, err := emaFloat64(prices, fastPeriod)
	if err != nil {
		return analytics.MACDResult{}, err
	}

	slowEMA, err := emaFloat64(prices, slowPeriod)
	if err != nil {
		return analytics.MACDResult{}, err
	}

	// Step 2: MACD line is just the difference
	macdValue := fastEMA - slowEMA

	// Step 3: Calculate signal line (EMA of MACD values)
	// For signal, we need a series of MACD values
	// Calculate incrementally to avoid O(n²)
	macdSeries := make([]float64, 0, len(prices)-slowPeriod+1)

	// Initialize EMAs for incremental calculation
	fastEMAval := 0.0
	slowEMAval := 0.0

	// Calculate initial EMAs
	for i := 0; i < slowPeriod; i++ {
		if i < fastPeriod {
			fastSum := 0.0
			for j := 0; j <= i && j < fastPeriod; j++ {
				fastSum += prices[j]
			}
			if i+1 >= fastPeriod {
				fastEMAval = fastSum / float64(fastPeriod)
			}
		}

		slowSum := 0.0
		for j := 0; j <= i; j++ {
			slowSum += prices[j]
		}
		if i+1 == slowPeriod {
			slowEMAval = slowSum / float64(slowPeriod)
		}
	}

	// Now calculate MACD incrementally
	fastMultiplier := 2.0 / float64(fastPeriod+1)
	slowMultiplier := 2.0 / float64(slowPeriod+1)

	for i := slowPeriod; i < len(prices); i++ {
		fastEMAval = (prices[i]-fastEMAval)*fastMultiplier + fastEMAval
		slowEMAval = (prices[i]-slowEMAval)*slowMultiplier + slowEMAval
		macdSeries = append(macdSeries, fastEMAval-slowEMAval)
	}

	// Step 4: Signal line is EMA of MACD series
	signalValue, err := emaFloat64(macdSeries, signalPeriod)
	if err != nil {
		return analytics.MACDResult{}, err
	}

	// Step 5: Histogram is MACD - Signal
	histogram := macdValue - signalValue

	return analytics.MACDResult{
		MACD:      numerical.NewFromFloat(macdValue),
		Signal:    numerical.NewFromFloat(signalValue),
		Histogram: numerical.NewFromFloat(histogram),
	}, nil
}
