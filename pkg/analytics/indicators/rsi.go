package indicators

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// RSI calculates the Relative Strength Index
func RSI(prices []decimal.Decimal, period int) ([]decimal.Decimal, error) {
	if len(prices) < period+1 {
		return nil, fmt.Errorf("insufficient data: need %d prices, got %d", period+1, len(prices))
	}

	result := make([]decimal.Decimal, 0, len(prices)-period)

	// Calculate initial average gain and loss
	var gains, losses []decimal.Decimal
	for i := 1; i < len(prices); i++ {
		change := prices[i].Sub(prices[i-1])
		if change.GreaterThan(decimal.Zero) {
			gains = append(gains, change)
			losses = append(losses, decimal.Zero)
		} else {
			gains = append(gains, decimal.Zero)
			losses = append(losses, change.Abs())
		}
	}

	if len(gains) < period {
		return nil, fmt.Errorf("insufficient data after calculating changes")
	}

	// Calculate average gain and loss for the first period
	avgGain := decimal.Zero
	avgLoss := decimal.Zero
	for i := 0; i < period; i++ {
		avgGain = avgGain.Add(gains[i])
		avgLoss = avgLoss.Add(losses[i])
	}
	avgGain = avgGain.Div(decimal.NewFromInt(int64(period)))
	avgLoss = avgLoss.Div(decimal.NewFromInt(int64(period)))

	// Calculate first RSI
	if avgLoss.IsZero() {
		result = append(result, decimal.NewFromInt(100))
	} else {
		rs := avgGain.Div(avgLoss)
		rsi := decimal.NewFromInt(100).Sub(decimal.NewFromInt(100).Div(decimal.NewFromInt(1).Add(rs)))
		result = append(result, rsi)
	}

	// Calculate subsequent RSIs using smoothed averages
	for i := period; i < len(gains); i++ {
		avgGain = avgGain.Mul(decimal.NewFromInt(int64(period - 1))).Add(gains[i]).Div(decimal.NewFromInt(int64(period)))
		avgLoss = avgLoss.Mul(decimal.NewFromInt(int64(period - 1))).Add(losses[i]).Div(decimal.NewFromInt(int64(period)))

		if avgLoss.IsZero() {
			result = append(result, decimal.NewFromInt(100))
		} else {
			rs := avgGain.Div(avgLoss)
			rsi := decimal.NewFromInt(100).Sub(decimal.NewFromInt(100).Div(decimal.NewFromInt(1).Add(rs)))
			result = append(result, rsi)
		}
	}

	return result, nil
}
