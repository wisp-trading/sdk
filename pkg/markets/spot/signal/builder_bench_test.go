package signal_test

import (
	"testing"

	"github.com/wisp-trading/sdk/pkg/markets/spot/signal"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
	temporal "github.com/wisp-trading/sdk/pkg/runtime/time"
	temporalType "github.com/wisp-trading/sdk/pkg/types/temporal"
)

var (
	benchTimeProvider temporalType.TimeProvider
	benchStrategyName strategy.StrategyName
	benchPair         portfolio.Pair
	benchExchange     connector.ExchangeName
	benchQuantity     numerical.Decimal
	benchPrice        numerical.Decimal
)

func init() {
	benchTimeProvider = temporal.NewTimeProvider()
	benchStrategyName = strategy.StrategyName("benchmark-strategy")
	benchPair = portfolio.NewPair(
		portfolio.NewAsset("BTC"),
		portfolio.NewAsset("USDT"),
	)
	benchExchange = connector.ExchangeName("binance")
	benchQuantity = numerical.NewFromFloat(10.5)
	benchPrice = numerical.NewFromFloat(50000.0)
}

func BenchmarkSpotBuilder_Buy(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = signal.NewSpotBuilder(benchStrategyName, benchTimeProvider).
			Buy(benchPair, benchExchange, benchQuantity).Build()
	}
}

func BenchmarkSpotBuilder_BuyLimit(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = signal.NewSpotBuilder(benchStrategyName, benchTimeProvider).
			BuyLimit(benchPair, benchExchange, benchQuantity, benchPrice).Build()
	}
}

func BenchmarkSpotBuilder_ChainedActions(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = signal.NewSpotBuilder(benchStrategyName, benchTimeProvider).
			Buy(benchPair, benchExchange, benchQuantity).
			Sell(benchPair, benchExchange, benchQuantity).
			BuyLimit(benchPair, benchExchange, benchQuantity, benchPrice).
			Build()
	}
}

func BenchmarkSpotBuilder_LargeSignal(b *testing.B) {
	b.ReportAllocs()
	pairs := []portfolio.Pair{
		portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT")),
		portfolio.NewPair(portfolio.NewAsset("ETH"), portfolio.NewAsset("USDT")),
		portfolio.NewPair(portfolio.NewAsset("SOL"), portfolio.NewAsset("USDT")),
		portfolio.NewPair(portfolio.NewAsset("AVAX"), portfolio.NewAsset("USDT")),
		portfolio.NewPair(portfolio.NewAsset("MATIC"), portfolio.NewAsset("USDT")),
	}

	for i := 0; i < b.N; i++ {
		builder := signal.NewSpotBuilder(benchStrategyName, benchTimeProvider)
		for _, pair := range pairs {
			builder = builder.Buy(pair, benchExchange, benchQuantity)
		}
		_, _ = builder.Build()
	}
}

func BenchmarkSpotBuilder_Buy_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = signal.NewSpotBuilder(benchStrategyName, benchTimeProvider).
				Buy(benchPair, benchExchange, benchQuantity).Build()
		}
	})
}
