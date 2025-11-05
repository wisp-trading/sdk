package indicators

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// EMA calculates the Exponential Moving Average
func EMA(prices []decimal.Decimal, period int) ([]decimal.Decimal, error) {
	if len(prices) < period {
		return nil, fmt.Errorf("insufficient data: need %d prices, got %d", period, len(prices))
	}

	result := make([]decimal.Decimal, len(prices))
	multiplier := decimal.NewFromInt(2).Div(decimal.NewFromInt(int64(period + 1)))

	// First EMA is SMA
	sum := decimal.Zero
	for i := 0; i < period; i++ {
		sum = sum.Add(prices[i])
	}
	result[period-1] = sum.Div(decimal.NewFromInt(int64(period)))

	// Calculate subsequent EMAs
	for i := period; i < len(prices); i++ {
		ema := prices[i].Sub(result[i-1]).Mul(multiplier).Add(result[i-1])
		result[i] = ema
	}

	return result[period-1:], nil
}
