package indicators_test

import (
	"testing"

	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
)

var (
	benchSMAPrices50   []float64
	benchSMAPrices100  []float64
	benchSMAPrices500  []float64
	benchSMAPrices1000 []float64
)

func init() {
	benchSMAPrices50 = generateSMAPriceData(50)
	benchSMAPrices100 = generateSMAPriceData(100)
	benchSMAPrices500 = generateSMAPriceData(500)
	benchSMAPrices1000 = generateSMAPriceData(1000)
}

func generateSMAPriceData(count int) []float64 {
	prices := make([]float64, count)

	basePrice := 50000.0
	for i := 0; i < count; i++ {
		prices[i] = basePrice + float64(i)*10 + float64(i%10)*5
	}

	return prices
}

func BenchmarkSMA_Period5_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices50, 5)
	}
}

func BenchmarkSMA_Period10_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices50, 10)
	}
}

func BenchmarkSMA_Period20_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices50, 20)
	}
}

func BenchmarkSMA_Period5_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices100, 5)
	}
}

func BenchmarkSMA_Period10_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices100, 10)
	}
}

func BenchmarkSMA_Period20_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices100, 20)
	}
}

func BenchmarkSMA_Period50_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices100, 50)
	}
}

func BenchmarkSMA_Period10_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices500, 10)
	}
}

func BenchmarkSMA_Period20_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices500, 20)
	}
}

func BenchmarkSMA_Period50_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices500, 50)
	}
}

func BenchmarkSMA_Period100_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices500, 100)
	}
}

func BenchmarkSMA_Period20_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices1000, 20)
	}
}

func BenchmarkSMA_Period50_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices1000, 50)
	}
}

func BenchmarkSMA_Period100_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices1000, 100)
	}
}

func BenchmarkSMA_Period200_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices1000, 200)
	}
}

func BenchmarkSMA_Period20_Data100_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.SMA(benchSMAPrices100, 20)
		}
	})
}

func BenchmarkSMA_Period20_Data500_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.SMA(benchSMAPrices500, 20)
		}
	})
}

func BenchmarkSMA_Period20_Data1000_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.SMA(benchSMAPrices1000, 20)
		}
	})
}

func BenchmarkSMA_Allocations_Small(b *testing.B) {
	b.ReportAllocs()
	prices := make([]float64, 20)

	for i := 0; i < 20; i++ {
		prices[i] = float64(100 + i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(prices, 5)
	}
}

func BenchmarkSMA_Allocations_Medium(b *testing.B) {
	b.ReportAllocs()
	count := 200
	prices := make([]float64, count)

	for i := 0; i < count; i++ {
		prices[i] = float64(100 + i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(prices, 20)
	}
}

func BenchmarkSMA_Allocations_Large(b *testing.B) {
	b.ReportAllocs()
	count := 5000
	prices := make([]float64, count)

	for i := 0; i < count; i++ {
		prices[i] = float64(100 + i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(prices, 50)
	}
}

func BenchmarkSMA_Standard_20(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices100, 20)
	}
}

func BenchmarkSMA_Standard_50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(benchSMAPrices100, 50)
	}
}

func BenchmarkSMA_Standard_200(b *testing.B) {
	b.ReportAllocs()
	prices := generateSMAPriceData(400)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.SMA(prices, 200)
	}
}

func BenchmarkSMA_Profile(b *testing.B) {
	b.Run("Period5", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.SMA(benchSMAPrices100, 5)
		}
	})

	b.Run("Period10", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.SMA(benchSMAPrices100, 10)
		}
	})

	b.Run("Period20", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.SMA(benchSMAPrices100, 20)
		}
	})

	b.Run("Period50", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.SMA(benchSMAPrices100, 50)
		}
	})
}
