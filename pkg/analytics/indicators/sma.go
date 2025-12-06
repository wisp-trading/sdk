package indicators

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

func smaFloat64(prices []float64, period int) (float64, error) {
	if len(prices) < period {
		return 0, fmt.Errorf("insufficient data: need %d prices, got %d", period, len(prices))
	}

	sum := 0.0
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i]
	}

	return sum / float64(period), nil
}

func SMA(prices []numerical.Decimal, period int) (numerical.Decimal, error) {
	pricesFloat := make([]float64, len(prices))
	for i, p := range prices {
		pricesFloat[i], _ = p.Float64()
	}

	smaFloat, err := smaFloat64(pricesFloat, period)
	if err != nil {
		return numerical.Zero(), err
	}

	return numerical.NewFromFloat(smaFloat), nil
}
