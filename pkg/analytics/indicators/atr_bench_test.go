package indicators_test

import (
	"testing"

	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
)

// Benchmark setup
var (
	benchHighs50    []float64
	benchLows50     []float64
	benchCloses50   []float64
	benchHighs100   []float64
	benchLows100    []float64
	benchCloses100  []float64
	benchHighs500   []float64
	benchLows500    []float64
	benchCloses500  []float64
	benchHighs1000  []float64
	benchLows1000   []float64
	benchCloses1000 []float64
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

func generatePriceData(count int) ([]float64, []float64, []float64) {
	highs := make([]float64, count)
	lows := make([]float64, count)
	closes := make([]float64, count)

	basePrice := 50000.0
	for i := 0; i < count; i++ {
		// Simulate price movement
		price := basePrice + float64(i)*10 + float64(i%10)*5
		highs[i] = price + 50
		lows[i] = price - 50
		closes[i] = price + 20
	}

	return highs, lows, closes
}

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
	highs := make([]float64, 20)
	lows := make([]float64, 20)
	closes := make([]float64, 20)

	for i := 0; i < 20; i++ {
		highs[i] = float64(100 + i)
		lows[i] = float64(95 + i)
		closes[i] = float64(98 + i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.ATR(highs, lows, closes, 5)
	}
}

func BenchmarkATR_Allocations_Large(b *testing.B) {
	b.ReportAllocs()
	count := 5000
	highs := make([]float64, count)
	lows := make([]float64, count)
	closes := make([]float64, count)

	for i := 0; i < count; i++ {
		highs[i] = float64(100 + i)
		lows[i] = float64(95 + i)
		closes[i] = float64(98 + i)
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
