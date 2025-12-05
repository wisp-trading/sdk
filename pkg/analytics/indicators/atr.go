package indicators

import (
	"fmt"
	"math"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

func atrFloat64(highs, lows, closes []float64, period int) ([]float64, error) {
	if len(highs) != len(lows) || len(highs) != len(closes) {
		return nil, fmt.Errorf("price arrays must have equal length")
	}
	if len(closes) < period+1 {
		return nil, fmt.Errorf("insufficient data: need %d prices, got %d", period+1, len(closes))
	}

	trueRanges := make([]float64, len(closes)-1)

	for i := 1; i < len(closes); i++ {
		highLow := highs[i] - lows[i]
		highClose := math.Abs(highs[i] - closes[i-1])
		lowClose := math.Abs(lows[i] - closes[i-1])

		tr := highLow
		if highClose > tr {
			tr = highClose
		}
		if lowClose > tr {
			tr = lowClose
		}

		trueRanges[i-1] = tr
	}

	result := make([]float64, 0, len(trueRanges)-period+1)

	sum := 0.0
	for i := 0; i < period; i++ {
		sum += trueRanges[i]
	}
	atr := sum / float64(period)
	result = append(result, atr)

	for i := period; i < len(trueRanges); i++ {
		atr = (atr*float64(period-1) + trueRanges[i]) / float64(period)
		result = append(result, atr)
	}

	return result, nil
}

// ATR converts inputs to float64, calls atrFloat64, and converts results back to numerical.Decimal
func ATR(highs, lows, closes []numerical.Decimal, period int) ([]numerical.Decimal, error) {
	if len(highs) != len(lows) || len(highs) != len(closes) {
		return nil, fmt.Errorf("price arrays must have equal length")
	}

	highsFloat := make([]float64, len(highs))
	lowsFloat := make([]float64, len(lows))
	closesFloat := make([]float64, len(closes))

	for i := range highs {
		highsFloat[i], _ = highs[i].Float64()
		lowsFloat[i], _ = lows[i].Float64()
		closesFloat[i], _ = closes[i].Float64()
	}

	atrFloat, err := atrFloat64(highsFloat, lowsFloat, closesFloat, period)
	if err != nil {
		return nil, err
	}

	result := make([]numerical.Decimal, len(atrFloat))
	for i, v := range atrFloat {
		result[i] = numerical.NewFromFloat(v)
	}

	return result, nil
}
