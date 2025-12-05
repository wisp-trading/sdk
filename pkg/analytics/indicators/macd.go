package indicators

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

func macdFloat64(prices []float64, fastPeriod, slowPeriod, signalPeriod int) (float64, float64, float64, error) {
	if len(prices) < slowPeriod+signalPeriod {
		return 0, 0, 0, fmt.Errorf("insufficient data: need %d prices, got %d", slowPeriod+signalPeriod, len(prices))
	}

	fastEMA, err := emaFloat64(prices, fastPeriod)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to calculate fast EMA: %w", err)
	}

	slowEMA, err := emaFloat64(prices, slowPeriod)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to calculate slow EMA: %w", err)
	}

	macdLine := fastEMA - slowEMA

	multiplier := 2.0 / float64(signalPeriod+1)

	macdValues := make([]float64, len(prices)-slowPeriod+1)
	for i := slowPeriod - 1; i < len(prices); i++ {
		fastVal, _ := emaFloat64(prices[:i+1], fastPeriod)
		slowVal, _ := emaFloat64(prices[:i+1], slowPeriod)
		macdValues[i-slowPeriod+1] = fastVal - slowVal
	}

	if len(macdValues) < signalPeriod {
		return 0, 0, 0, fmt.Errorf("insufficient MACD values for signal calculation")
	}

	sum := 0.0
	for i := 0; i < signalPeriod; i++ {
		sum += macdValues[i]
	}
	signal := sum / float64(signalPeriod)

	for i := signalPeriod; i < len(macdValues); i++ {
		signal = (macdValues[i]-signal)*multiplier + signal
	}

	histogram := macdLine - signal

	return macdLine, signal, histogram, nil
}

func MACD(prices []numerical.Decimal, fastPeriod, slowPeriod, signalPeriod int) (analytics.MACDResult, error) {
	pricesFloat := make([]float64, len(prices))
	for i, p := range prices {
		pricesFloat[i], _ = p.Float64()
	}

	macdLine, signal, histogram, err := macdFloat64(pricesFloat, fastPeriod, slowPeriod, signalPeriod)
	if err != nil {
		return analytics.MACDResult{}, err
	}

	return analytics.MACDResult{
		MACD:      numerical.NewFromFloat(macdLine),
		Signal:    numerical.NewFromFloat(signal),
		Histogram: numerical.NewFromFloat(histogram),
	}, nil
}
