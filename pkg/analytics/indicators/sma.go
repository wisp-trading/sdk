package indicators

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// SMA calculates the Simple Moving Average
func SMA(prices []decimal.Decimal, period int) ([]decimal.Decimal, error) {
	if len(prices) < period {
		return nil, fmt.Errorf("insufficient data: need %d prices, got %d", period, len(prices))
	}

	result := make([]decimal.Decimal, 0, len(prices)-period+1)

	for i := period - 1; i < len(prices); i++ {
		sum := decimal.Zero
		for j := 0; j < period; j++ {
			sum = sum.Add(prices[i-j])
		}
		avg := sum.Div(decimal.NewFromInt(int64(period)))
		result = append(result, avg)
	}

	return result, nil
}
