package indicators_test

import (
	"testing"

	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

// Benchmark setup
var (
	benchPrices50   []numerical.Decimal
	benchPrices100  []numerical.Decimal
	benchPrices500  []numerical.Decimal
	benchPrices1000 []numerical.Decimal
)

func init() {
	// Initialize 50 data points
	benchPrices50 = generateBBPriceData(50)
	// Initialize 100 data points
	benchPrices100 = generateBBPriceData(100)
	// Initialize 500 data points
	benchPrices500 = generateBBPriceData(500)
	// Initialize 1000 data points
	benchPrices1000 = generateBBPriceData(1000)
}

func generateBBPriceData(count int) []numerical.Decimal {
	prices := make([]numerical.Decimal, count)

	basePrice := 50000.0
	for i := 0; i < count; i++ {
		// Simulate price movement with some volatility
		price := basePrice + float64(i)*10 + float64(i%10)*5
		prices[i] = numerical.NewFromFloat(price)
	}

	return prices
}

// Benchmarks with different periods
func BenchmarkBollingerBands_Period5_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices50, 5, 2.0)
	}
}

func BenchmarkBollingerBands_Period10_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices50, 10, 2.0)
	}
}

func BenchmarkBollingerBands_Period20_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices50, 20, 2.0)
	}
}

func BenchmarkBollingerBands_Period5_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices100, 5, 2.0)
	}
}

func BenchmarkBollingerBands_Period10_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices100, 10, 2.0)
	}
}

func BenchmarkBollingerBands_Period20_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices100, 20, 2.0)
	}
}

func BenchmarkBollingerBands_Period50_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices100, 50, 2.0)
	}
}

func BenchmarkBollingerBands_Period10_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices500, 10, 2.0)
	}
}

func BenchmarkBollingerBands_Period20_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices500, 20, 2.0)
	}
}

func BenchmarkBollingerBands_Period50_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices500, 50, 2.0)
	}
}

func BenchmarkBollingerBands_Period100_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices500, 100, 2.0)
	}
}

func BenchmarkBollingerBands_Period20_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices1000, 20, 2.0)
	}
}

func BenchmarkBollingerBands_Period50_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices1000, 50, 2.0)
	}
}

func BenchmarkBollingerBands_Period100_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices1000, 100, 2.0)
	}
}

func BenchmarkBollingerBands_Period200_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices1000, 200, 2.0)
	}
}

// Benchmarks with different stdDev multipliers
func BenchmarkBollingerBands_StdDev1_Period20_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices100, 20, 1.0)
	}
}

func BenchmarkBollingerBands_StdDev2_Period20_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices100, 20, 2.0)
	}
}

func BenchmarkBollingerBands_StdDev3_Period20_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices100, 20, 3.0)
	}
}

// Parallel benchmarks
func BenchmarkBollingerBands_Period20_Data100_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.BollingerBands(benchPrices100, 20, 2.0)
		}
	})
}

func BenchmarkBollingerBands_Period20_Data500_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.BollingerBands(benchPrices500, 20, 2.0)
		}
	})
}

func BenchmarkBollingerBands_Period20_Data1000_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.BollingerBands(benchPrices1000, 20, 2.0)
		}
	})
}

// Memory allocation benchmarks
func BenchmarkBollingerBands_Allocations_Small(b *testing.B) {
	b.ReportAllocs()
	prices := make([]numerical.Decimal, 20)

	for i := 0; i < 20; i++ {
		prices[i] = numerical.NewFromInt(int64(100 + i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(prices, 5, 2.0)
	}
}

func BenchmarkBollingerBands_Allocations_Medium(b *testing.B) {
	b.ReportAllocs()
	count := 200
	prices := make([]numerical.Decimal, count)

	for i := 0; i < count; i++ {
		prices[i] = numerical.NewFromInt(int64(100 + i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(prices, 20, 2.0)
	}
}

func BenchmarkBollingerBands_Allocations_Large(b *testing.B) {
	b.ReportAllocs()
	count := 5000
	prices := make([]numerical.Decimal, count)

	for i := 0; i < count; i++ {
		prices[i] = numerical.NewFromInt(int64(100 + i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(prices, 50, 2.0)
	}
}

// Standard configuration benchmark (most common use case)
func BenchmarkBollingerBands_Standard_20_2(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.BollingerBands(benchPrices100, 20, 2.0)
	}
}

// Sub-benchmarks for profiling
func BenchmarkBollingerBands_Profile(b *testing.B) {
	b.Run("Period5", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.BollingerBands(benchPrices100, 5, 2.0)
		}
	})

	b.Run("Period10", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.BollingerBands(benchPrices100, 10, 2.0)
		}
	})

	b.Run("Period20", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.BollingerBands(benchPrices100, 20, 2.0)
		}
	})

	b.Run("Period50", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.BollingerBands(benchPrices100, 50, 2.0)
		}
	})
}
