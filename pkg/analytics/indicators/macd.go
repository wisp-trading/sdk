package indicators

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
	"github.com/shopspring/decimal"
)

// MACD calculates the Moving Average Convergence Divergence
func MACD(prices []decimal.Decimal, fastPeriod, slowPeriod, signalPeriod int) ([]analytics.MACDResult, error) {
	if len(prices) < slowPeriod {
		return nil, fmt.Errorf("insufficient data: need %d prices, got %d", slowPeriod, len(prices))
	}

	// Calculate fast and slow EMAs
	fastEMA, err := EMA(prices, fastPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate fast EMA: %w", err)
	}

	slowEMA, err := EMA(prices, slowPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate slow EMA: %w", err)
	}

	// Calculate MACD line (difference between fast and slow EMA)
	// Align the arrays (slowEMA is shorter)
	startOffset := len(fastEMA) - len(slowEMA)
	macdLine := make([]decimal.Decimal, len(slowEMA))
	for i := 0; i < len(slowEMA); i++ {
		macdLine[i] = fastEMA[i+startOffset].Sub(slowEMA[i])
	}

	// Calculate signal line (EMA of MACD line)
	signalLine, err := EMA(macdLine, signalPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate signal line: %w", err)
	}

	// Calculate histogram (MACD - Signal)
	result := make([]analytics.MACDResult, len(signalLine))
	startOffset = len(macdLine) - len(signalLine)
	for i := 0; i < len(signalLine); i++ {
		macdValue := macdLine[i+startOffset]
		signalValue := signalLine[i]
		result[i] = analytics.MACDResult{
			MACD:      macdValue,
			Signal:    signalValue,
			Histogram: macdValue.Sub(signalValue),
		}
	}

	return result, nil
}
