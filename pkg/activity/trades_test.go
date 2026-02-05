package activity_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	sdkTesting "github.com/wisp-trading/sdk/pkg/testing"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	storeActivity "github.com/wisp-trading/sdk/pkg/types/data/stores/activity"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	wispActivity "github.com/wisp-trading/sdk/pkg/types/wisp/activity"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("Trades", func() {
	var (
		app    *fxtest.App
		trades wispActivity.Trades
		store  storeActivity.Trades
	)

	BeforeEach(func() {
		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(
				&trades,
				&store,
			),
			fx.NopLogger,
		)

		app.RequireStart()
	})

	AfterEach(func() {
		app.RequireStop()
	})

	Describe("GetAllTrades", func() {
		It("should return all trades from underlying store", func() {
			ctx := context.Background()

			btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
			ethPair := portfolio.NewPair(portfolio.NewAsset("ETH"), portfolio.NewAsset("USDT"))

			// Add trades to store
			store.AddTrade(connector.Trade{ID: "trade-1", Pair: btcPair, Price: numerical.NewFromFloat(50000)})
			store.AddTrade(connector.Trade{ID: "trade-2", Pair: ethPair, Price: numerical.NewFromFloat(3000)})

			result := trades.GetAllTrades(ctx)

			Expect(result).To(HaveLen(2))
			Expect(result[0].ID).To(Equal("trade-1"))
			Expect(result[1].ID).To(Equal("trade-2"))
		})

		It("should return empty slice when no trades exist", func() {
			ctx := context.Background()

			result := trades.GetAllTrades(ctx)

			Expect(result).To(BeEmpty())
		})
	})

	Describe("GetTradesByExchange", func() {
		It("should filter trades by exchange", func() {
			ctx := context.Background()

			exchange := connector.ExchangeName("binance")
			btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
			ethPair := portfolio.NewPair(portfolio.NewAsset("ETH"), portfolio.NewAsset("USDT"))

			// Add trades to store
			store.AddTrade(connector.Trade{ID: "trade-1", Exchange: exchange, Pair: btcPair})
			store.AddTrade(connector.Trade{ID: "trade-2", Exchange: exchange, Pair: ethPair})
			store.AddTrade(connector.Trade{ID: "trade-3", Exchange: "coinbase", Pair: btcPair})

			result := trades.GetTradesByExchange(ctx, exchange)

			Expect(result).To(HaveLen(2))
			Expect(result[0].Exchange).To(Equal(exchange))
			Expect(result[1].Exchange).To(Equal(exchange))
		})

		It("should return empty slice when no trades for exchange", func() {
			ctx := context.Background()

			exchange := connector.ExchangeName("unknown-exchange")

			result := trades.GetTradesByExchange(ctx, exchange)

			Expect(result).To(BeEmpty())
		})
	})

	Describe("GetTradesByPair", func() {
		It("should filter trades by pair", func() {
			ctx := context.Background()

			btc := portfolio.NewPair(
				portfolio.NewAsset("BTC"),
				portfolio.NewAsset("USDT"),
			)

			eth := portfolio.NewPair(
				portfolio.NewAsset("ETH"),
				portfolio.NewAsset("USDT"),
			)

			// Add trades to store
			store.AddTrade(connector.Trade{ID: "trade-1", Pair: btc, Price: numerical.NewFromFloat(50000)})
			store.AddTrade(connector.Trade{ID: "trade-2", Pair: btc, Price: numerical.NewFromFloat(51000)})
			store.AddTrade(connector.Trade{ID: "trade-3", Pair: eth, Price: numerical.NewFromFloat(3000)})

			result := trades.GetTradesByPair(ctx, btc)

			Expect(result).To(HaveLen(2))
			Expect(result[0].Pair.Symbol()).To(Equal("BTC-USDT"))
			Expect(result[1].Pair.Symbol()).To(Equal("BTC-USDT"))
		})

		It("should return empty slice when no trades for pair", func() {
			ctx := context.Background()

			pair := portfolio.NewPair(
				portfolio.NewAsset("UNKNOWN"),
				portfolio.NewAsset("USDT"),
			)

			result := trades.GetTradesByPair(ctx, pair)

			Expect(result).To(BeEmpty())
		})
	})

	Describe("GetTradesSince", func() {
		It("should filter trades by time", func() {
			ctx := context.Background()

			since := time.Now().Add(-1 * time.Hour)
			now := time.Now()
			btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))

			// Add trades to store with different timestamps
			store.AddTrade(connector.Trade{ID: "trade-1", Timestamp: now, Pair: btcPair})
			store.AddTrade(connector.Trade{ID: "trade-2", Timestamp: now.Add(-30 * time.Minute), Pair: btcPair})
			store.AddTrade(connector.Trade{ID: "trade-3", Timestamp: now.Add(-2 * time.Hour), Pair: btcPair}) // Before 'since'

			result := trades.GetTradesSince(ctx, since)

			Expect(result).To(HaveLen(2))
		})

		It("should return empty slice when no trades since time", func() {
			ctx := context.Background()

			since := time.Now()

			result := trades.GetTradesSince(ctx, since)

			Expect(result).To(BeEmpty())
		})
	})

	Describe("GetTradeByID", func() {
		It("should return trade for existing ID", func() {
			ctx := context.Background()

			tradeID := "trade-123"
			btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))

			// Add trade to store
			store.AddTrade(connector.Trade{
				ID:    tradeID,
				Pair:  btcPair,
				Price: numerical.NewFromFloat(50000),
			})

			result := trades.GetTradeByID(ctx, tradeID)

			Expect(result).NotTo(BeNil())
			Expect(result.ID).To(Equal(tradeID))
		})

		It("should return nil for unknown ID", func() {
			ctx := context.Background()

			result := trades.GetTradeByID(ctx, "unknown")

			Expect(result).To(BeNil())
		})
	})

	Describe("GetTradeCount", func() {
		It("should return count from underlying store", func() {
			ctx := context.Background()

			btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))

			// Add trades to store
			store.AddTrade(connector.Trade{ID: "trade-1", Pair: btcPair})
			store.AddTrade(connector.Trade{ID: "trade-2", Pair: btcPair})
			store.AddTrade(connector.Trade{ID: "trade-3", Pair: btcPair})

			result := trades.GetTradeCount(ctx)

			Expect(result).To(Equal(3))
		})

		It("should return zero when no trades exist", func() {
			ctx := context.Background()

			result := trades.GetTradeCount(ctx)

			Expect(result).To(Equal(0))
		})
	})

	Describe("GetTotalVolume", func() {
		It("should return volume from underlying store", func() {
			ctx := context.Background()

			btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))

			// Add trades with quantities
			store.AddTrade(connector.Trade{
				Pair:     btcPair,
				Quantity: numerical.NewFromFloat(5.5),
				ID:       "trade-1",
			})
			store.AddTrade(connector.Trade{
				Pair:     btcPair,
				Quantity: numerical.NewFromFloat(3.0),
				ID:       "trade-2",
			})
			store.AddTrade(connector.Trade{
				Pair:     btcPair,
				Quantity: numerical.NewFromFloat(2.0),
				ID:       "trade-3",
			})

			result := trades.GetTotalVolume(ctx, btcPair)

			Expect(result.Equal(numerical.NewFromFloat(10.5))).To(BeTrue())
		})

		It("should return zero volume when no trades for pair", func() {
			ctx := context.Background()

			unknownPair := portfolio.NewPair(portfolio.NewAsset("UNKNOWN"), portfolio.NewAsset("USDT"))

			result := trades.GetTotalVolume(ctx, unknownPair)

			Expect(result.IsZero()).To(BeTrue())
		})
	})
})
