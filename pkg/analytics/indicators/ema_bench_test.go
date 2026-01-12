package indicators_test

import (
	"testing"

	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
)

var (
	benchEMAPrices50   []float64
	benchEMAPrices100  []float64
	benchEMAPrices500  []float64
	benchEMAPrices1000 []float64
)

func init() {
	benchEMAPrices50 = generateEMAPriceData(50)
	benchEMAPrices100 = generateEMAPriceData(100)
	benchEMAPrices500 = generateEMAPriceData(500)
	benchEMAPrices1000 = generateEMAPriceData(1000)
}

func generateEMAPriceData(count int) []float64 {
	prices := make([]float64, count)

	basePrice := 50000.0
	for i := 0; i < count; i++ {
		prices[i] = basePrice + float64(i)*10 + float64(i%10)*5
	}

	return prices
}

func BenchmarkEMA_Period5_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices50, 5)
	}
}

func BenchmarkEMA_Period12_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices50, 12)
	}
}

func BenchmarkEMA_Period20_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices50, 20)
	}
}

func BenchmarkEMA_Period5_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices100, 5)
	}
}

func BenchmarkEMA_Period12_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices100, 12)
	}
}

func BenchmarkEMA_Period20_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices100, 20)
	}
}

func BenchmarkEMA_Period50_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices100, 50)
	}
}

func BenchmarkEMA_Period12_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices500, 12)
	}
}

func BenchmarkEMA_Period26_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices500, 26)
	}
}

func BenchmarkEMA_Period50_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices500, 50)
	}
}

func BenchmarkEMA_Period100_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices500, 100)
	}
}

func BenchmarkEMA_Period12_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices1000, 12)
	}
}

func BenchmarkEMA_Period50_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices1000, 50)
	}
}

func BenchmarkEMA_Period100_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices1000, 100)
	}
}

func BenchmarkEMA_Period200_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices1000, 200)
	}
}

func BenchmarkEMA_Period12_Data100_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.EMA(benchEMAPrices100, 12)
		}
	})
}

func BenchmarkEMA_Period12_Data500_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.EMA(benchEMAPrices500, 12)
		}
	})
}

func BenchmarkEMA_Period12_Data1000_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.EMA(benchEMAPrices1000, 12)
		}
	})
}

func BenchmarkEMA_Allocations_Small(b *testing.B) {
	b.ReportAllocs()
	prices := make([]float64, 20)

	for i := 0; i < 20; i++ {
		prices[i] = float64(100 + i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(prices, 5)
	}
}

func BenchmarkEMA_Allocations_Medium(b *testing.B) {
	b.ReportAllocs()
	count := 200
	prices := make([]float64, count)

	for i := 0; i < count; i++ {
		prices[i] = float64(100 + i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(prices, 20)
	}
}

func BenchmarkEMA_Allocations_Large(b *testing.B) {
	b.ReportAllocs()
	count := 5000
	prices := make([]float64, count)

	for i := 0; i < count; i++ {
		prices[i] = float64(100 + i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(prices, 50)
	}
}

func BenchmarkEMA_Standard_12(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices100, 12)
	}
}

func BenchmarkEMA_Standard_26(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices100, 26)
	}
}

func BenchmarkEMA_Standard_50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(benchEMAPrices100, 50)
	}
}

func BenchmarkEMA_Standard_200(b *testing.B) {
	b.ReportAllocs()
	prices := generateEMAPriceData(400)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.EMA(prices, 200)
	}
}

func BenchmarkEMA_Profile(b *testing.B) {
	b.Run("Period5", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.EMA(benchEMAPrices100, 5)
		}
	})

	b.Run("Period12", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.EMA(benchEMAPrices100, 12)
		}
	})

	b.Run("Period26", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.EMA(benchEMAPrices100, 26)
		}
	})

	b.Run("Period50", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.EMA(benchEMAPrices100, 50)
		}
	})
}
