package indicators_test

import (
	"testing"

	"github.com/wisp-trading/sdk/pkg/analytics/indicators"
)

var (
	benchStochHighs50    []float64
	benchStochLows50     []float64
	benchStochCloses50   []float64
	benchStochHighs100   []float64
	benchStochLows100    []float64
	benchStochCloses100  []float64
	benchStochHighs500   []float64
	benchStochLows500    []float64
	benchStochCloses500  []float64
	benchStochHighs1000  []float64
	benchStochLows1000   []float64
	benchStochCloses1000 []float64
)

func init() {
	benchStochHighs50, benchStochLows50, benchStochCloses50 = generateStochPriceData(50)
	benchStochHighs100, benchStochLows100, benchStochCloses100 = generateStochPriceData(100)
	benchStochHighs500, benchStochLows500, benchStochCloses500 = generateStochPriceData(500)
	benchStochHighs1000, benchStochLows1000, benchStochCloses1000 = generateStochPriceData(1000)
}

func generateStochPriceData(count int) ([]float64, []float64, []float64) {
	highs := make([]float64, count)
	lows := make([]float64, count)
	closes := make([]float64, count)

	basePrice := 50000.0
	for i := 0; i < count; i++ {
		price := basePrice + float64(i)*10 + float64(i%10)*5
		highs[i] = price + 50
		lows[i] = price - 50
		closes[i] = price + 20
	}

	return highs, lows, closes
}

func BenchmarkStochastic_K5_D3_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.Stochastic(benchStochHighs50, benchStochLows50, benchStochCloses50, 5, 3)
	}
}

func BenchmarkStochastic_K14_D3_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.Stochastic(benchStochHighs50, benchStochLows50, benchStochCloses50, 14, 3)
	}
}

func BenchmarkStochastic_K20_D3_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.Stochastic(benchStochHighs50, benchStochLows50, benchStochCloses50, 20, 3)
	}
}

func BenchmarkStochastic_K5_D3_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.Stochastic(benchStochHighs100, benchStochLows100, benchStochCloses100, 5, 3)
	}
}

func BenchmarkStochastic_K14_D3_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.Stochastic(benchStochHighs100, benchStochLows100, benchStochCloses100, 14, 3)
	}
}

func BenchmarkStochastic_K20_D3_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.Stochastic(benchStochHighs100, benchStochLows100, benchStochCloses100, 20, 3)
	}
}

func BenchmarkStochastic_K50_D3_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.Stochastic(benchStochHighs100, benchStochLows100, benchStochCloses100, 50, 3)
	}
}

func BenchmarkStochastic_K14_D3_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.Stochastic(benchStochHighs500, benchStochLows500, benchStochCloses500, 14, 3)
	}
}

func BenchmarkStochastic_K50_D3_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.Stochastic(benchStochHighs500, benchStochLows500, benchStochCloses500, 50, 3)
	}
}

func BenchmarkStochastic_K100_D3_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.Stochastic(benchStochHighs500, benchStochLows500, benchStochCloses500, 100, 3)
	}
}

func BenchmarkStochastic_K14_D3_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.Stochastic(benchStochHighs1000, benchStochLows1000, benchStochCloses1000, 14, 3)
	}
}

func BenchmarkStochastic_K50_D3_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.Stochastic(benchStochHighs1000, benchStochLows1000, benchStochCloses1000, 50, 3)
	}
}

func BenchmarkStochastic_K200_D3_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.Stochastic(benchStochHighs1000, benchStochLows1000, benchStochCloses1000, 200, 3)
	}
}

func BenchmarkStochastic_K14_D3_Data100_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.Stochastic(benchStochHighs100, benchStochLows100, benchStochCloses100, 14, 3)
		}
	})
}

func BenchmarkStochastic_K14_D3_Data500_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.Stochastic(benchStochHighs500, benchStochLows500, benchStochCloses500, 14, 3)
		}
	})
}

func BenchmarkStochastic_K14_D3_Data1000_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.Stochastic(benchStochHighs1000, benchStochLows1000, benchStochCloses1000, 14, 3)
		}
	})
}

func BenchmarkStochastic_Allocations_Small(b *testing.B) {
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
		_, _ = indicators.Stochastic(highs, lows, closes, 5, 3)
	}
}

func BenchmarkStochastic_Allocations_Large(b *testing.B) {
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
		_, _ = indicators.Stochastic(highs, lows, closes, 50, 3)
	}
}

func BenchmarkStochastic_Standard_14_3(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.Stochastic(benchStochHighs100, benchStochLows100, benchStochCloses100, 14, 3)
	}
}

func BenchmarkStochastic_Profile(b *testing.B) {
	b.Run("K5_D3", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.Stochastic(benchStochHighs100, benchStochLows100, benchStochCloses100, 5, 3)
		}
	})

	b.Run("K14_D3", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.Stochastic(benchStochHighs100, benchStochLows100, benchStochCloses100, 14, 3)
		}
	})

	b.Run("K20_D3", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.Stochastic(benchStochHighs100, benchStochLows100, benchStochCloses100, 20, 3)
		}
	})

	b.Run("K50_D3", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.Stochastic(benchStochHighs100, benchStochLows100, benchStochCloses100, 50, 3)
		}
	})
}
