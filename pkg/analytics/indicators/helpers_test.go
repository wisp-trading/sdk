package indicators_test

import "github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"

// Helper functions
func makeDecimals(values ...float64) []numerical.Decimal {
	result := make([]numerical.Decimal, len(values))
	for i, v := range values {
		result[i] = numerical.NewFromFloat(v)
	}
	return result
}

func makeDecimalsFloat(values ...float64) []numerical.Decimal {
	result := make([]numerical.Decimal, len(values))
	for i, v := range values {
		result[i] = numerical.NewFromFloat(v)
	}
	return result
}
