package indicators_test

import (
	"testing"

	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

var (
	benchMACDPrices50   []numerical.Decimal
	benchMACDPrices100  []numerical.Decimal
	benchMACDPrices500  []numerical.Decimal
	benchMACDPrices1000 []numerical.Decimal
)

func init() {
	benchMACDPrices50 = generateMACDPriceData(50)
	benchMACDPrices100 = generateMACDPriceData(100)
	benchMACDPrices500 = generateMACDPriceData(500)
	benchMACDPrices1000 = generateMACDPriceData(1000)
}

func generateMACDPriceData(count int) []numerical.Decimal {
	prices := make([]numerical.Decimal, count)

	basePrice := 50000.0
	for i := 0; i < count; i++ {
		price := basePrice + float64(i)*10 + float64(i%10)*5
		prices[i] = numerical.NewFromFloat(price)
	}

	return prices
}

func BenchmarkMACD_12_26_9_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.MACD(benchMACDPrices50, 12, 26, 9)
	}
}

func BenchmarkMACD_5_13_5_Data50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.MACD(benchMACDPrices50, 5, 13, 5)
	}
}

func BenchmarkMACD_12_26_9_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.MACD(benchMACDPrices100, 12, 26, 9)
	}
}

func BenchmarkMACD_5_13_5_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.MACD(benchMACDPrices100, 5, 13, 5)
	}
}

func BenchmarkMACD_19_39_9_Data100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.MACD(benchMACDPrices100, 19, 39, 9)
	}
}

func BenchmarkMACD_12_26_9_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.MACD(benchMACDPrices500, 12, 26, 9)
	}
}

func BenchmarkMACD_19_39_9_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.MACD(benchMACDPrices500, 19, 39, 9)
	}
}

func BenchmarkMACD_26_52_9_Data500(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.MACD(benchMACDPrices500, 26, 52, 9)
	}
}

func BenchmarkMACD_12_26_9_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.MACD(benchMACDPrices1000, 12, 26, 9)
	}
}

func BenchmarkMACD_19_39_9_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.MACD(benchMACDPrices1000, 19, 39, 9)
	}
}

func BenchmarkMACD_26_52_9_Data1000(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.MACD(benchMACDPrices1000, 26, 52, 9)
	}
}

func BenchmarkMACD_12_26_9_Data100_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.MACD(benchMACDPrices100, 12, 26, 9)
		}
	})
}

func BenchmarkMACD_12_26_9_Data500_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.MACD(benchMACDPrices500, 12, 26, 9)
		}
	})
}

func BenchmarkMACD_12_26_9_Data1000_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = indicators.MACD(benchMACDPrices1000, 12, 26, 9)
		}
	})
}

func BenchmarkMACD_Allocations_Small(b *testing.B) {
	b.ReportAllocs()
	prices := make([]numerical.Decimal, 50)

	for i := 0; i < 50; i++ {
		prices[i] = numerical.NewFromInt(int64(100 + i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.MACD(prices, 12, 26, 9)
	}
}

func BenchmarkMACD_Allocations_Large(b *testing.B) {
	b.ReportAllocs()
	count := 5000
	prices := make([]numerical.Decimal, count)

	for i := 0; i < count; i++ {
		prices[i] = numerical.NewFromInt(int64(100 + i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.MACD(prices, 12, 26, 9)
	}
}

func BenchmarkMACD_Standard_12_26_9(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = indicators.MACD(benchMACDPrices100, 12, 26, 9)
	}
}

func BenchmarkMACD_Profile(b *testing.B) {
	b.Run("Fast_5_13_5", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.MACD(benchMACDPrices100, 5, 13, 5)
		}
	})

	b.Run("Standard_12_26_9", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.MACD(benchMACDPrices100, 12, 26, 9)
		}
	})

	b.Run("Slow_19_39_9", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.MACD(benchMACDPrices100, 19, 39, 9)
		}
	})

	b.Run("Slower_26_52_9", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = indicators.MACD(benchMACDPrices100, 26, 52, 9)
		}
	})
}
