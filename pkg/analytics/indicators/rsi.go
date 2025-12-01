package indicators

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

// RSI calculates the Relative Strength Index
func RSI(prices []numerical.Decimal, period int) ([]numerical.Decimal, error) {
	if len(prices) < period+1 {
		return nil, fmt.Errorf("insufficient data: need %d prices, got %d", period+1, len(prices))
	}

	result := make([]numerical.Decimal, 0, len(prices)-period)

	// Calculate initial average gain and loss
	var gains, losses []numerical.Decimal
	for i := 1; i < len(prices); i++ {
		change := prices[i].Sub(prices[i-1])
		if change.GreaterThan(numerical.Zero()) {
			gains = append(gains, change)
			losses = append(losses, numerical.Zero())
		} else {
			gains = append(gains, numerical.Zero())
			losses = append(losses, change.Abs())
		}
	}

	if len(gains) < period {
		return nil, fmt.Errorf("insufficient data after calculating changes")
	}

	// Calculate average gain and loss for the first period
	avgGain := numerical.Zero()
	avgLoss := numerical.Zero()
	for i := 0; i < period; i++ {
		avgGain = avgGain.Add(gains[i])
		avgLoss = avgLoss.Add(losses[i])
	}
	avgGain = avgGain.Div(numerical.NewFromInt(int64(period)))
	avgLoss = avgLoss.Div(numerical.NewFromInt(int64(period)))

	// Calculate first RSI
	if avgLoss.IsZero() {
		result = append(result, numerical.NewFromInt(100))
	} else {
		rs := avgGain.Div(avgLoss)
		rsi := numerical.NewFromInt(100).Sub(numerical.NewFromInt(100).Div(numerical.NewFromInt(1).Add(rs)))
		result = append(result, rsi)
	}

	// Calculate subsequent RSIs using smoothed averages
	for i := period; i < len(gains); i++ {
		avgGain = avgGain.Mul(numerical.NewFromInt(int64(period - 1))).Add(gains[i]).Div(numerical.NewFromInt(int64(period)))
		avgLoss = avgLoss.Mul(numerical.NewFromInt(int64(period - 1))).Add(losses[i]).Div(numerical.NewFromInt(int64(period)))

		if avgLoss.IsZero() {
			result = append(result, numerical.NewFromInt(100))
		} else {
			rs := avgGain.Div(avgLoss)
			rsi := numerical.NewFromInt(100).Sub(numerical.NewFromInt(100).Div(numerical.NewFromInt(1).Add(rs)))
			result = append(result, rsi)
		}
	}

	return result, nil
}
