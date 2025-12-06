package indicators_test

import (
	"testing"

	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

var (
	benchStochHighs50    []numerical.Decimal
	benchStochLows50     []numerical.Decimal
	benchStochCloses50   []numerical.Decimal
	benchStochHighs100   []numerical.Decimal
	benchStochLows100    []numerical.Decimal
	benchStochCloses100  []numerical.Decimal
	benchStochHighs500   []numerical.Decimal
	benchStochLows500    []numerical.Decimal
	benchStochCloses500  []numerical.Decimal
	benchStochHighs1000  []numerical.Decimal
	benchStochLows1000   []numerical.Decimal
	benchStochCloses1000 []numerical.Decimal
)

func init() {
	benchStochHighs50, benchStochLows50, benchStochCloses50 = generateStochPriceData(50)
	benchStochHighs100, benchStochLows100, benchStochCloses100 = generateStochPriceData(100)
	benchStochHighs500, benchStochLows500, benchStochCloses500 = generateStochPriceData(500)
	benchStochHighs1000, benchStochLows1000, benchStochCloses1000 = generateStochPriceData(1000)
}

func generateStochPriceData(count int) ([]numerical.Decimal, []numerical.Decimal, []numerical.Decimal) {
	highs := make([]numerical.Decimal, count)
	lows := make([]numerical.Decimal, count)
	closes := make([]numerical.Decimal, count)

	basePrice := 50000.0
	for i := 0; i < count; i++ {
		price := basePrice + float64(i)*10 + float64(i%10)*5
		highs[i] = numerical.NewFromFloat(price + 50)
		lows[i] = numerical.NewFromFloat(price - 50)
		closes[i] = numerical.NewFromFloat(price + 20)
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
		_, _ = indicators.Stochastic(highs, lows, closes, 5, 3)
	}
}

func BenchmarkStochastic_Allocations_Large(b *testing.B) {
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
