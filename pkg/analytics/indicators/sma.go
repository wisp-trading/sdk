package indicators

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

// SMA calculates the Simple Moving Average
func SMA(prices []numerical.Decimal, period int) ([]numerical.Decimal, error) {
	if len(prices) < period {
		return nil, fmt.Errorf("insufficient data: need %d prices, got %d", period, len(prices))
	}

	result := make([]numerical.Decimal, 0, len(prices)-period+1)

	for i := period - 1; i < len(prices); i++ {
		sum := numerical.Zero()
		for j := 0; j < period; j++ {
			sum = sum.Add(prices[i-j])
		}
		avg := sum.Div(numerical.NewFromInt(int64(period)))
		result = append(result, avg)
	}

	return result, nil
}
