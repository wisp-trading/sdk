package indicators_test

import (
	"testing"

	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

// Benchmark setup
var (
	benchHighs50    []numerical.Decimal
	benchLows50     []numerical.Decimal
	benchCloses50   []numerical.Decimal
	benchHighs100   []numerical.Decimal
	benchLows100    []numerical.Decimal
	benchCloses100  []numerical.Decimal
	benchHighs500   []numerical.Decimal
	benchLows500    []numerical.Decimal
	benchCloses500  []numerical.Decimal
	benchHighs1000  []numerical.Decimal
	benchLows1000   []numerical.Decimal
	benchCloses1000 []numerical.Decimal
)

func init() {
	// Initialize 50 data points
	benchHighs50, benchLows50, benchCloses50 = generatePriceData(50)
	// Initialize 100 data points
	benchHighs100, benchLows100, benchCloses100 = generatePriceData(100)
	// Initialize 500 data points
	benchHighs500, benchLows500, benchCloses500 = generatePriceData(500)
	// Initialize 1000 data points
	benchHighs1000, benchLows1000, benchCloses1000 = generatePriceData(1000)
}

func generatePriceData(count int) ([]numerical.Decimal, []numerical.Decimal, []numerical.Decimal) {
	highs := make([]numerical.Decimal, count)
	lows := make([]numerical.Decimal, count)
	closes := make([]numerical.Decimal, count)

	basePrice := 50000.0
	for i := 0; i < count; i++ {
		// Simulate price movement
		price := basePrice + float64(i)*10 + float64(i%10)*5
		highs[i] = numerical.NewFromFloat(price + 50)
		lows[i] = numerical.NewFromFloat(price - 50)
		closes[i] = numerical.NewFromFloat(price + 20)
	}

	return highs, lows, closes
}

// Benchmarks with different periods
func BenchmarkATR_Period5_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(benchHighs50, benchLows50, benchCloses50, 5)
	}
}

func BenchmarkATR_Period14_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(benchHighs50, benchLows50, benchCloses50, 14)
	}
}

func BenchmarkATR_Period20_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(benchHighs50, benchLows50, benchCloses50, 20)
	}
}

func BenchmarkATR_Period5_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(benchHighs100, benchLows100, benchCloses100, 5)
	}
}

func BenchmarkATR_Period14_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(benchHighs100, benchLows100, benchCloses100, 14)
	}
}

func BenchmarkATR_Period20_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(benchHighs100, benchLows100, benchCloses100, 20)
	}
}

func BenchmarkATR_Period50_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(benchHighs100, benchLows100, benchCloses100, 50)
	}
}

func BenchmarkATR_Period14_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(benchHighs500, benchLows500, benchCloses500, 14)
	}
}

func BenchmarkATR_Period50_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(benchHighs500, benchLows500, benchCloses500, 50)
	}
}

func BenchmarkATR_Period100_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(benchHighs500, benchLows500, benchCloses500, 100)
	}
}

func BenchmarkATR_Period14_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(benchHighs1000, benchLows1000, benchCloses1000, 14)
	}
}

func BenchmarkATR_Period50_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(benchHighs1000, benchLows1000, benchCloses1000, 50)
	}
}

func BenchmarkATR_Period200_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(benchHighs1000, benchLows1000, benchCloses1000, 200)
	}
}

// Parallel benchmarks
func BenchmarkATR_Period14_Data100_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.ATR(benchHighs100, benchLows100, benchCloses100, 14)
		}
	})
}

func BenchmarkATR_Period14_Data500_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.ATR(benchHighs500, benchLows500, benchCloses500, 14)
		}
	})
}

func BenchmarkATR_Period14_Data1000_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.ATR(benchHighs1000, benchLows1000, benchCloses1000, 14)
		}
	})
}

// Memory allocation benchmarks
func BenchmarkATR_Allocations_Small(b *testing.B) {
	b.ReportAllocs()
	highs := make([]numerical.Decimal, 20)
	lows := make([]numerical.Decimal, 20)
	closes := make([]numerical.Decimal, 20)

	for i := 0; i < 20; i++ {
		highs[i] = numerical.NewFromInt(int64(100 + i))
		lows[i] = numerical.NewFromInt(int64(95 + i))
		closes[i] = numerical.NewFromInt(int64(98 + i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(highs, lows, closes, 5)
	}
}

func BenchmarkATR_Allocations_Large(b *testing.B) {
	b.ReportAllocs()
	count := 5000
	highs := make([]numerical.Decimal, count)
	lows := make([]numerical.Decimal, count)
	closes := make([]numerical.Decimal, count)

	for i := 0; i < count; i++ {
		highs[i] = numerical.NewFromInt(int64(100 + i))
		lows[i] = numerical.NewFromInt(int64(95 + i))
		closes[i] = numerical.NewFromInt(int64(98 + i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(highs, lows, closes, 14)
	}
}

// Benchmark true range calculation specifically
func BenchmarkATR_TrueRangeCalculation(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(benchHighs100, benchLows100, benchCloses100, 14)
	}
}

// Benchmark smoothing calculation
func BenchmarkATR_SmoothingCalculation_LongPeriod(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(benchHighs1000, benchLows1000, benchCloses1000, 100)
	}
}
