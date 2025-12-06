package indicators

import (
	"fmt"
	"math"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

// bollingerBandsFloat64 is the internal high-performance implementation using float64
func bollingerBandsFloat64(prices []float64, period int, stdDev float64) (analytics.BollingerBandsResult, error) {
	if len(prices) < period {
		return analytics.BollingerBandsResult{}, fmt.Errorf("insufficient data: need %d prices, got %d", period, len(prices))
	}

	sma, err := smaFloat64(prices, period)
	if err != nil {
		return analytics.BollingerBandsResult{}, err
	}

	variance := 0.0
	for i := len(prices) - period; i < len(prices); i++ {
		diff := prices[i] - sma
		variance += diff * diff
	}
	variance /= float64(period)
	sd := math.Sqrt(variance) * stdDev

	return analytics.BollingerBandsResult{
		Upper:  numerical.NewFromFloat(sma + sd),
		Middle: numerical.NewFromFloat(sma),
		Lower:  numerical.NewFromFloat(sma - sd),
	}, nil
}

func BollingerBands(prices []numerical.Decimal, period int, stdDev float64) (analytics.BollingerBandsResult, error) {
	pricesFloat := make([]float64, len(prices))
	for i, p := range prices {
		pricesFloat[i], _ = p.Float64()
	}

	return bollingerBandsFloat64(pricesFloat, period, stdDev)
}
