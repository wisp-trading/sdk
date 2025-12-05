package indicators

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

func emaFloat64(prices []float64, period int) (float64, error) {
	if len(prices) < period {
		return 0, fmt.Errorf("insufficient data: need %d prices, got %d", period, len(prices))
	}

	multiplier := 2.0 / float64(period+1)

	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema := sum / float64(period)

	for i := period; i < len(prices); i++ {
		ema = (prices[i]-ema)*multiplier + ema
	}

	return ema, nil
}

func EMA(prices []numerical.Decimal, period int) (numerical.Decimal, error) {
	pricesFloat := make([]float64, len(prices))
	for i, p := range prices {
		pricesFloat[i], _ = p.Float64()
	}

	emaFloat, err := emaFloat64(pricesFloat, period)
	if err != nil {
		return numerical.Zero(), err
	}

	return numerical.NewFromFloat(emaFloat), nil
}
