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
	benchAsset        portfolio.Asset
	benchExchange     connector.ExchangeName
	benchQuantity     numerical.Decimal
	benchPrice        numerical.Decimal
)

func init() {
	timeProvider := temporal.NewTimeProvider()
	benchFactory = signal.NewFactory(timeProvider)
	benchStrategyName = strategy.StrategyName("benchmark-strategy")
	benchAsset = portfolio.NewAsset("BTC")
	benchExchange = connector.ExchangeName("binance")
	benchQuantity = numerical.NewFromFloat(10.5)
	benchPrice = numerical.NewFromFloat(50000.0)
}

func BenchmarkSignalBuilder_Buy(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.New(benchStrategyName)
		_ = builder.Buy(benchAsset, benchExchange, benchQuantity).Build()
	}
}

func BenchmarkSignalBuilder_BuyLimit(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.New(benchStrategyName)
		_ = builder.BuyLimit(benchAsset, benchExchange, benchQuantity, benchPrice).Build()
	}
}

func BenchmarkSignalBuilder_Sell(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.New(benchStrategyName)
		_ = builder.Sell(benchAsset, benchExchange, benchQuantity).Build()
	}
}

func BenchmarkSignalBuilder_SellLimit(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.New(benchStrategyName)
		_ = builder.SellLimit(benchAsset, benchExchange, benchQuantity, benchPrice).Build()
	}
}

func BenchmarkSignalBuilder_SellShort(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.New(benchStrategyName)
		_ = builder.SellShort(benchAsset, benchExchange, benchQuantity).Build()
	}
}

func BenchmarkSignalBuilder_SellShortLimit(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.New(benchStrategyName)
		_ = builder.SellShortLimit(benchAsset, benchExchange, benchQuantity, benchPrice).Build()
	}
}

func BenchmarkSignalBuilder_ChainedActions(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.New(benchStrategyName)
		_ = builder.
			Buy(benchAsset, benchExchange, benchQuantity).
			Sell(benchAsset, benchExchange, benchQuantity).
			BuyLimit(benchAsset, benchExchange, benchQuantity, benchPrice).
			Build()
	}
}

func BenchmarkSignalBuilder_MultipleActions(b *testing.B) {
	b.ReportAllocs()
	asset2 := portfolio.NewAsset("ETH")
	asset3 := portfolio.NewAsset("SOL")

	for i := 0; i < b.N; i++ {
		builder := benchFactory.New(benchStrategyName)
		_ = builder.
			Buy(benchAsset, benchExchange, benchQuantity).
			Buy(asset2, benchExchange, benchQuantity).
			Buy(asset3, benchExchange, benchQuantity).
			SellLimit(benchAsset, benchExchange, benchQuantity, benchPrice).
			SellLimit(asset2, benchExchange, benchQuantity, benchPrice).
			Build()
	}
}

func BenchmarkSignalBuilder_Build(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := benchFactory.New(benchStrategyName)
		_ = builder.Build()
	}
}

func BenchmarkSignalBuilder_LargeSignal(b *testing.B) {
	b.ReportAllocs()
	assets := []portfolio.Asset{
		portfolio.NewAsset("BTC"),
		portfolio.NewAsset("ETH"),
		portfolio.NewAsset("SOL"),
		portfolio.NewAsset("AVAX"),
		portfolio.NewAsset("MATIC"),
		portfolio.NewAsset("DOT"),
		portfolio.NewAsset("ATOM"),
		portfolio.NewAsset("LINK"),
		portfolio.NewAsset("UNI"),
		portfolio.NewAsset("AAVE"),
	}

	for i := 0; i < b.N; i++ {
		builder := benchFactory.New(benchStrategyName)
		for _, asset := range assets {
			builder = builder.Buy(asset, benchExchange, benchQuantity)
		}
		_ = builder.Build()
	}
}

func BenchmarkSignalFactory_New(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = benchFactory.New(benchStrategyName)
	}
}

// Parallel benchmarks
func BenchmarkSignalBuilder_Buy_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			builder := benchFactory.New(benchStrategyName)
			_ = builder.Buy(benchAsset, benchExchange, benchQuantity).Build()
		}
	})
}

func BenchmarkSignalBuilder_ChainedActions_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			builder := benchFactory.New(benchStrategyName)
			_ = builder.
				Buy(benchAsset, benchExchange, benchQuantity).
				Sell(benchAsset, benchExchange, benchQuantity).
				BuyLimit(benchAsset, benchExchange, benchQuantity, benchPrice).
				Build()
		}
	})
}

func BenchmarkSignalBuilder_LargeSignal_Parallel(b *testing.B) {
	b.ReportAllocs()
	assets := []portfolio.Asset{
		portfolio.NewAsset("BTC"),
		portfolio.NewAsset("ETH"),
		portfolio.NewAsset("SOL"),
		portfolio.NewAsset("AVAX"),
		portfolio.NewAsset("MATIC"),
		portfolio.NewAsset("DOT"),
		portfolio.NewAsset("ATOM"),
		portfolio.NewAsset("LINK"),
		portfolio.NewAsset("UNI"),
		portfolio.NewAsset("AAVE"),
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			builder := benchFactory.New(benchStrategyName)
			for _, asset := range assets {
				builder = builder.Buy(asset, benchExchange, benchQuantity)
			}
			_ = builder.Build()
		}
	})
}
