package indicators

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

func macdFloat64(prices []float64, fastPeriod, slowPeriod, signalPeriod int) (float64, float64, float64, error) {
	if len(prices) < slowPeriod {
		return 0, 0, 0, fmt.Errorf("insufficient data: need %d prices, got %d", slowPeriod, len(prices))
	}

	// Step 1: Calculate fast and slow EMAs ONCE - O(n)
	fastEMA, err := emaFloat64(prices, fastPeriod)
	if err != nil {
		return 0, 0, 0, err
	}

	slowEMA, err := emaFloat64(prices, slowPeriod)
	if err != nil {
		return 0, 0, 0, err
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
		return 0, 0, 0, err
	}

	// Step 5: Histogram is MACD - Signal
	histogram := macdValue - signalValue

	return macdValue, signalValue, histogram, nil
}

func MACD(prices []numerical.Decimal, fastPeriod, slowPeriod, signalPeriod int) (analytics.MACDResult, error) {
	pricesFloat := make([]float64, len(prices))
	for i, p := range prices {
		pricesFloat[i], _ = p.Float64()
	}

	macd, signal, histogram, err := macdFloat64(pricesFloat, fastPeriod, slowPeriod, signalPeriod)
	if err != nil {
		return analytics.MACDResult{}, err
	}

	return analytics.MACDResult{
		MACD:      numerical.NewFromFloat(macd),
		Signal:    numerical.NewFromFloat(signal),
		Histogram: numerical.NewFromFloat(histogram),
	}, nil
}
