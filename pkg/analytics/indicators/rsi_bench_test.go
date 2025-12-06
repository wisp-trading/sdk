package indicators_test

import (
	"testing"

	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

var (
	benchRSIPrices50   []numerical.Decimal
	benchRSIPrices100  []numerical.Decimal
	benchRSIPrices500  []numerical.Decimal
	benchRSIPrices1000 []numerical.Decimal
)

func init() {
	benchRSIPrices50 = generateRSIPriceData(50)
	benchRSIPrices100 = generateRSIPriceData(100)
	benchRSIPrices500 = generateRSIPriceData(500)
	benchRSIPrices1000 = generateRSIPriceData(1000)
}

func generateRSIPriceData(count int) []numerical.Decimal {
	prices := make([]numerical.Decimal, count)

	basePrice := 50000.0
	for i := 0; i < count; i++ {
		price := basePrice + float64(i)*10 + float64(i%10)*5
		prices[i] = numerical.NewFromFloat(price)
	}

	return prices
}

func BenchmarkRSI_Period5_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(benchRSIPrices50, 5)
	}
}

func BenchmarkRSI_Period14_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(benchRSIPrices50, 14)
	}
}

func BenchmarkRSI_Period20_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(benchRSIPrices50, 20)
	}
}

func BenchmarkRSI_Period5_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(benchRSIPrices100, 5)
	}
}

func BenchmarkRSI_Period14_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(benchRSIPrices100, 14)
	}
}

func BenchmarkRSI_Period20_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(benchRSIPrices100, 20)
	}
}

func BenchmarkRSI_Period50_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(benchRSIPrices100, 50)
	}
}

func BenchmarkRSI_Period14_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(benchRSIPrices500, 14)
	}
}

func BenchmarkRSI_Period50_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(benchRSIPrices500, 50)
	}
}

func BenchmarkRSI_Period100_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(benchRSIPrices500, 100)
	}
}

func BenchmarkRSI_Period14_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(benchRSIPrices1000, 14)
	}
}

func BenchmarkRSI_Period50_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(benchRSIPrices1000, 50)
	}
}

func BenchmarkRSI_Period200_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(benchRSIPrices1000, 200)
	}
}

func BenchmarkRSI_Period14_Data100_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.RSI(benchRSIPrices100, 14)
		}
	})
}

func BenchmarkRSI_Period14_Data500_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.RSI(benchRSIPrices500, 14)
		}
	})
}

func BenchmarkRSI_Period14_Data1000_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.RSI(benchRSIPrices1000, 14)
		}
	})
}

func BenchmarkRSI_Allocations_Small(b *testing.B) {
	b.ReportAllocs()
	prices := make([]numerical.Decimal, 20)

	for i := 0; i < 20; i++ {
		prices[i] = numerical.NewFromInt(int64(100 + i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(prices, 5)
	}
}

func BenchmarkRSI_Allocations_Large(b *testing.B) {
	b.ReportAllocs()
	count := 5000
	prices := make([]numerical.Decimal, count)

	for i := 0; i < count; i++ {
		prices[i] = numerical.NewFromInt(int64(100 + i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(prices, 50)
	}
}

func BenchmarkRSI_Standard_14(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.RSI(benchRSIPrices100, 14)
	}
}

func BenchmarkRSI_Profile(b *testing.B) {
	b.Run("Period5", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.RSI(benchRSIPrices100, 5)
		}
	})

	b.Run("Period14", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.RSI(benchRSIPrices100, 14)
		}
	})

	b.Run("Period20", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.RSI(benchRSIPrices100, 20)
		}
	})

	b.Run("Period50", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.RSI(benchRSIPrices100, 50)
		}
	})
}
