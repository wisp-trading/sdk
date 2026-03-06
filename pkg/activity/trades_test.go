package activity_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	sdkTesting "github.com/wisp-trading/sdk/pkg/testing"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	wispActivity "github.com/wisp-trading/sdk/pkg/types/wisp/activity"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("Trades", func() {
	var (
		app        *fxtest.App
		trades     wispActivity.Trades
		spotTrades spotTypes.SpotTrades
	)

	BeforeEach(func() {
		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(
				&trades,
				&spotTrades,
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

			spotTrades.AddTrade(connector.Trade{ID: "trade-1", Pair: btcPair, Price: numerical.NewFromFloat(50000)})
			spotTrades.AddTrade(connector.Trade{ID: "trade-2", Pair: ethPair, Price: numerical.NewFromFloat(3000)})

			result := trades.GetAllTrades(ctx)

			Expect(result).To(HaveLen(2))
		})

		It("should return empty slice when no trades exist", func() {
			Expect(trades.GetAllTrades(context.Background())).To(BeEmpty())
		})
	})

	Describe("GetTradesByExchange", func() {
		It("should filter trades by exchange", func() {
			ctx := context.Background()
			exchange := connector.ExchangeName("binance")
			btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))

			spotTrades.AddTrade(connector.Trade{ID: "trade-1", Exchange: exchange, Pair: btcPair})
			spotTrades.AddTrade(connector.Trade{ID: "trade-2", Exchange: exchange, Pair: btcPair})
			spotTrades.AddTrade(connector.Trade{ID: "trade-3", Exchange: "coinbase", Pair: btcPair})

			result := trades.GetTradesByExchange(ctx, exchange)

			Expect(result).To(HaveLen(2))
		})
	})

	Describe("GetTradesByPair", func() {
		It("should filter trades by pair", func() {
			ctx := context.Background()
			btc := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
			eth := portfolio.NewPair(portfolio.NewAsset("ETH"), portfolio.NewAsset("USDT"))

			spotTrades.AddTrade(connector.Trade{ID: "trade-1", Pair: btc})
			spotTrades.AddTrade(connector.Trade{ID: "trade-2", Pair: btc})
			spotTrades.AddTrade(connector.Trade{ID: "trade-3", Pair: eth})

			result := trades.GetTradesByPair(ctx, btc)

			Expect(result).To(HaveLen(2))
		})
	})

	Describe("GetTradesSince", func() {
		It("should filter trades by time", func() {
			ctx := context.Background()
			since := time.Now().Add(-1 * time.Hour)
			now := time.Now()
			btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))

			spotTrades.AddTrade(connector.Trade{ID: "trade-1", Timestamp: now, Pair: btcPair})
			spotTrades.AddTrade(connector.Trade{ID: "trade-2", Timestamp: now.Add(-30 * time.Minute), Pair: btcPair})
			spotTrades.AddTrade(connector.Trade{ID: "trade-3", Timestamp: now.Add(-2 * time.Hour), Pair: btcPair})

			result := trades.GetTradesSince(ctx, since)

			Expect(result).To(HaveLen(2))
		})
	})

	Describe("GetTradeByID", func() {
		It("should return trade for existing ID", func() {
			ctx := context.Background()
			btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
			spotTrades.AddTrade(connector.Trade{ID: "trade-123", Pair: btcPair})

			result := trades.GetTradeByID(ctx, "trade-123")

			Expect(result).NotTo(BeNil())
			Expect(result.ID).To(Equal("trade-123"))
		})

		It("should return nil for unknown ID", func() {
			Expect(trades.GetTradeByID(context.Background(), "unknown")).To(BeNil())
		})
	})

	Describe("GetTradeCount", func() {
		It("should return count from underlying store", func() {
			ctx := context.Background()
			btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))

			spotTrades.AddTrade(connector.Trade{ID: "trade-1", Pair: btcPair})
			spotTrades.AddTrade(connector.Trade{ID: "trade-2", Pair: btcPair})
			spotTrades.AddTrade(connector.Trade{ID: "trade-3", Pair: btcPair})

			Expect(trades.GetTradeCount(ctx)).To(Equal(3))
		})

		It("should return zero when no trades exist", func() {
			Expect(trades.GetTradeCount(context.Background())).To(Equal(0))
		})
	})

	Describe("GetTotalVolume", func() {
		It("should return volume from underlying store", func() {
			ctx := context.Background()
			btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))

			spotTrades.AddTrade(connector.Trade{Pair: btcPair, Quantity: numerical.NewFromFloat(5.5), ID: "trade-1"})
			spotTrades.AddTrade(connector.Trade{Pair: btcPair, Quantity: numerical.NewFromFloat(3.0), ID: "trade-2"})
			spotTrades.AddTrade(connector.Trade{Pair: btcPair, Quantity: numerical.NewFromFloat(2.0), ID: "trade-3"})

			result := trades.GetTotalVolume(ctx, btcPair)

			Expect(result.Equal(numerical.NewFromFloat(10.5))).To(BeTrue())
		})

		It("should return zero volume when no trades for pair", func() {
			unknownPair := portfolio.NewPair(portfolio.NewAsset("UNKNOWN"), portfolio.NewAsset("USDT"))
			Expect(trades.GetTotalVolume(context.Background(), unknownPair).IsZero()).To(BeTrue())
		})
	})
})
