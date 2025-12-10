package activity_test

import (
	"context"
	"time"

	mockStoreActivity "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/activity"
	"github.com/backtesting-org/kronos-sdk/pkg/activity"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	kronosActivity "github.com/backtesting-org/kronos-sdk/pkg/types/kronos/activity"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Trades", func() {
	var (
		mockStore *mockStoreActivity.Trades
		trades    kronosActivity.Trades
	)

	BeforeEach(func() {
		mockStore = mockStoreActivity.NewTrades(GinkgoT())
		trades = activity.NewTrades(mockStore)
	})

	Describe("GetAllTrades", func() {
		It("should return all trades from underlying store", func() {
			ctx := context.Background()

			expectedTrades := []connector.Trade{
				{ID: "trade-1", Symbol: "BTC", Price: numerical.NewFromFloat(50000)},
				{ID: "trade-2", Symbol: "ETH", Price: numerical.NewFromFloat(3000)},
			}

			mockStore.EXPECT().GetAllTrades().Return(expectedTrades)

			result := trades.GetAllTrades(ctx)

			Expect(result).To(HaveLen(2))
			Expect(result[0].ID).To(Equal("trade-1"))
			Expect(result[1].ID).To(Equal("trade-2"))
		})

		It("should return empty slice when no trades exist", func() {
			ctx := context.Background()

			mockStore.EXPECT().GetAllTrades().Return([]connector.Trade{})

			result := trades.GetAllTrades(ctx)

			Expect(result).To(BeEmpty())
		})
	})

	Describe("GetTradesByExchange", func() {
		It("should filter trades by exchange", func() {
			ctx := context.Background()

			exchange := connector.ExchangeName("binance")
			expectedTrades := []connector.Trade{
				{ID: "trade-1", Exchange: exchange, Symbol: "BTC"},
				{ID: "trade-2", Exchange: exchange, Symbol: "ETH"},
			}

			mockStore.EXPECT().GetTradesByExchange(exchange).Return(expectedTrades)

			result := trades.GetTradesByExchange(ctx, exchange)

			Expect(result).To(HaveLen(2))
			Expect(result[0].Exchange).To(Equal(exchange))
			Expect(result[1].Exchange).To(Equal(exchange))
		})

		It("should return empty slice when no trades for exchange", func() {
			ctx := context.Background()

			exchange := connector.ExchangeName("unknown-exchange")

			mockStore.EXPECT().GetTradesByExchange(exchange).Return([]connector.Trade{})

			result := trades.GetTradesByExchange(ctx, exchange)

			Expect(result).To(BeEmpty())
		})
	})

	Describe("GetTradesByAsset", func() {
		It("should filter trades by asset", func() {
			ctx := context.Background()

			asset := portfolio.NewAsset("BTC")
			expectedTrades := []connector.Trade{
				{ID: "trade-1", Symbol: "BTC", Price: numerical.NewFromFloat(50000)},
				{ID: "trade-2", Symbol: "BTC", Price: numerical.NewFromFloat(51000)},
			}

			mockStore.EXPECT().GetTradesByAsset(asset).Return(expectedTrades)

			result := trades.GetTradesByAsset(ctx, asset)

			Expect(result).To(HaveLen(2))
			Expect(result[0].Symbol).To(Equal("BTC"))
			Expect(result[1].Symbol).To(Equal("BTC"))
		})

		It("should return empty slice when no trades for asset", func() {
			ctx := context.Background()

			asset := portfolio.NewAsset("UNKNOWN")

			mockStore.EXPECT().GetTradesByAsset(asset).Return([]connector.Trade{})

			result := trades.GetTradesByAsset(ctx, asset)

			Expect(result).To(BeEmpty())
		})
	})

	Describe("GetTradesSince", func() {
		It("should filter trades by time", func() {
			ctx := context.Background()

			since := time.Now().Add(-1 * time.Hour)
			expectedTrades := []connector.Trade{
				{ID: "trade-1", Timestamp: time.Now()},
				{ID: "trade-2", Timestamp: time.Now().Add(-30 * time.Minute)},
			}

			mockStore.EXPECT().GetTradesSince(since).Return(expectedTrades)

			result := trades.GetTradesSince(ctx, since)

			Expect(result).To(HaveLen(2))
		})

		It("should return empty slice when no trades since time", func() {
			ctx := context.Background()

			since := time.Now()

			mockStore.EXPECT().GetTradesSince(since).Return([]connector.Trade{})

			result := trades.GetTradesSince(ctx, since)

			Expect(result).To(BeEmpty())
		})
	})

	Describe("GetTradeByID", func() {
		It("should return trade for existing ID", func() {
			ctx := context.Background()

			tradeID := "trade-123"
			expectedTrade := &connector.Trade{
				ID:    tradeID,
				Price: numerical.NewFromFloat(50000),
			}

			mockStore.EXPECT().GetTradeByID(tradeID).Return(expectedTrade)

			result := trades.GetTradeByID(ctx, tradeID)

			Expect(result).NotTo(BeNil())
			Expect(result.ID).To(Equal(tradeID))
		})

		It("should return nil for unknown ID", func() {
			ctx := context.Background()

			mockStore.EXPECT().GetTradeByID("unknown").Return(nil)

			result := trades.GetTradeByID(ctx, "unknown")

			Expect(result).To(BeNil())
		})
	})

	Describe("GetTradeCount", func() {
		It("should return count from underlying store", func() {
			ctx := context.Background()

			mockStore.EXPECT().GetTradeCount().Return(100)

			result := trades.GetTradeCount(ctx)

			Expect(result).To(Equal(100))
		})

		It("should return zero when no trades exist", func() {
			ctx := context.Background()

			mockStore.EXPECT().GetTradeCount().Return(0)

			result := trades.GetTradeCount(ctx)

			Expect(result).To(Equal(0))
		})
	})

	Describe("GetTotalVolume", func() {
		It("should return volume from underlying store", func() {
			ctx := context.Background()

			asset := portfolio.NewAsset("BTC")
			expectedVolume := numerical.NewFromFloat(10.5)

			mockStore.EXPECT().GetTotalVolume(asset).Return(expectedVolume)

			result := trades.GetTotalVolume(ctx, asset)

			Expect(result.Equal(expectedVolume)).To(BeTrue())
		})

		It("should return zero volume when no trades for asset", func() {
			ctx := context.Background()

			asset := portfolio.NewAsset("UNKNOWN")
			zeroVolume := numerical.Zero()

			mockStore.EXPECT().GetTotalVolume(asset).Return(zeroVolume)

			result := trades.GetTotalVolume(ctx, asset)

			Expect(result.IsZero()).To(BeTrue())
		})
	})
})
