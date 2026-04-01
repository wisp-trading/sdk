package signal_test

import (
	"testing"

	temporal "github.com/wisp-trading/sdk/pkg/runtime/time"
	"github.com/wisp-trading/sdk/pkg/signal"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

var (
	benchFactory      strategy.SignalFactory
	benchStrategyName strategy.StrategyName
	benchPair         portfolio.Pair
	benchExchange     connector.ExchangeName
	benchQuantity     numerical.Decimal
	benchPrice        numerical.Decimal
)

func init() {
	timeProvider := temporal.NewTimeProvider()
	benchFactory = signal.NewFactory(timeProvider)
	benchStrategyName = strategy.StrategyName("benchmark-strategy")
	benchPair = portfolio.NewPair(
		portfolio.NewAsset("BTC"),
		portfolio.NewAsset("USDT"),
	)
	benchExchange = connector.ExchangeName("binance")
	benchQuantity = numerical.NewFromFloat(10.5)
	benchPrice = numerical.NewFromFloat(50000.0)
}

func BenchmarkSignalBuilder_Buy(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.NewSpot(benchStrategyName)
		_, _ = builder.Buy(benchPair, benchExchange, benchQuantity).Build()
	}
}

func BenchmarkSignalBuilder_BuyLimit(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.NewSpot(benchStrategyName)
		_, _ = builder.BuyLimit(benchPair, benchExchange, benchQuantity, benchPrice).Build()
	}
}

func BenchmarkSignalBuilder_Sell(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.NewSpot(benchStrategyName)
		_, _ = builder.Sell(benchPair, benchExchange, benchQuantity).Build()
	}
}

func BenchmarkSignalBuilder_SellLimit(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.NewSpot(benchStrategyName)
		_, _ = builder.SellLimit(benchPair, benchExchange, benchQuantity, benchPrice).Build()
	}
}

func BenchmarkSignalBuilder_SellShort(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.NewSpot(benchStrategyName)
		_, _ = builder.SellShort(benchPair, benchExchange, benchQuantity).Build()
	}
}

func BenchmarkSignalBuilder_SellShortLimit(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.NewSpot(benchStrategyName)
		_, _ = builder.SellShortLimit(benchPair, benchExchange, benchQuantity, benchPrice).Build()
	}
}

func BenchmarkSignalBuilder_ChainedActions(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.NewSpot(benchStrategyName)
		_, _ = builder.
			Buy(benchPair, benchExchange, benchQuantity).
			Sell(benchPair, benchExchange, benchQuantity).
			BuyLimit(benchPair, benchExchange, benchQuantity, benchPrice).
			Build()
	}
}

func BenchmarkSignalBuilder_MultipleActions(b *testing.B) {
	b.ReportAllocs()
	pair2 := portfolio.NewPair(
		portfolio.NewAsset("ETH"),
		portfolio.NewAsset("USDT"),
	)

	pair3 := portfolio.NewPair(
		portfolio.NewAsset("SOL"),
		portfolio.NewAsset("USDT"),
	)

	for i := 0; i < b.N; i++ {
		builder := benchFactory.NewSpot(benchStrategyName)
		_, _ = builder.
			Buy(benchPair, benchExchange, benchQuantity).
			Buy(pair2, benchExchange, benchQuantity).
			Buy(pair3, benchExchange, benchQuantity).
			SellLimit(benchPair, benchExchange, benchQuantity, benchPrice).
			SellLimit(pair2, benchExchange, benchQuantity, benchPrice).
			Build()
	}
}

func BenchmarkSignalBuilder_Build(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.NewSpot(benchStrategyName)
		_, _ = builder.Build()
	}
}

func BenchmarkSignalBuilder_LargeSignal(b *testing.B) {
	b.ReportAllocs()
	pairs := []portfolio.Pair{
		portfolio.NewPair(
			portfolio.NewAsset("BTC"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("ETH"),
			portfolio.NewAsset("USDT"),
		),

		portfolio.NewPair(
			portfolio.NewAsset("SOL"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("AVAX"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("MATIC"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("DOT"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("ATOM"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("LINK"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("UNI"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("AAVE"),
			portfolio.NewAsset("USDT"),
		),
	}

	for i := 0; i < b.N; i++ {
		builder := benchFactory.NewSpot(benchStrategyName)
		for _, asset := range pairs {
			builder = builder.Buy(asset, benchExchange, benchQuantity)
		}
		_, _ = builder.Build()
	}
}

func BenchmarkSignalFactory_NewSpot(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = benchFactory.NewSpot(benchStrategyName)
	}
}

// Parallel benchmarks
func BenchmarkSignalBuilder_Buy_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			builder := benchFactory.NewSpot(benchStrategyName)
			_, _ = builder.Buy(benchPair, benchExchange, benchQuantity).Build()
		}
	})
}

func BenchmarkSignalBuilder_ChainedActions_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			builder := benchFactory.NewSpot(benchStrategyName)
			_, _ = builder.
				Buy(benchPair, benchExchange, benchQuantity).
				Sell(benchPair, benchExchange, benchQuantity).
				BuyLimit(benchPair, benchExchange, benchQuantity, benchPrice).
				Build()
		}
	})
}

func BenchmarkSignalBuilder_LargeSignal_Parallel(b *testing.B) {
	b.ReportAllocs()
	pairs := []portfolio.Pair{
		portfolio.NewPair(
			portfolio.NewAsset("BTC"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("ETH"),
			portfolio.NewAsset("USDT"),
		),

		portfolio.NewPair(
			portfolio.NewAsset("SOL"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("AVAX"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("MATIC"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("DOT"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("ATOM"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("LINK"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("UNI"),
			portfolio.NewAsset("USDT"),
		),
		portfolio.NewPair(
			portfolio.NewAsset("AAVE"),
			portfolio.NewAsset("USDT"),
		),
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			builder := benchFactory.NewSpot(benchStrategyName)
			for _, asset := range pairs {
				builder = builder.Buy(asset, benchExchange, benchQuantity)
			}
			_, _ = builder.Build()
		}
	})
}
