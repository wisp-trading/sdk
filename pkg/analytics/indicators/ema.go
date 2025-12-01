package indicators

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

// EMA calculates the Exponential Moving Average
func EMA(prices []numerical.Decimal, period int) ([]numerical.Decimal, error) {
	if len(prices) < period {
		return nil, fmt.Errorf("insufficient data: need %d prices, got %d", period, len(prices))
	}

	result := make([]numerical.Decimal, len(prices))
	multiplier := numerical.NewFromInt(2).Div(numerical.NewFromInt(int64(period + 1)))

	// First EMA is SMA
	sum := numerical.Zero()
	for i := 0; i < period; i++ {
		sum = sum.Add(prices[i])
	}
	result[period-1] = sum.Div(numerical.NewFromInt(int64(period)))

	// Calculate subsequent EMAs
	for i := period; i < len(prices); i++ {
		ema := prices[i].Sub(result[i-1]).Mul(multiplier).Add(result[i-1])
		result[i] = ema
	}

	return result[period-1:], nil
}
