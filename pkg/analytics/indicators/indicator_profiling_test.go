package indicators_test

import (
	"context"

	mockMarket "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
	mockProfiling "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/monitoring/profiling"
	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
	"github.com/backtesting-org/kronos-sdk/pkg/monitoring/profiling"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Indicator Profiling", func() {
	var (
		mockStore     *mockMarket.MarketData
		indicatorsSvc analytics.Indicators
		mockProfCtx   *mockProfiling.Context
		ctx           context.Context
	)

	BeforeEach(func() {
		mockStore = mockMarket.NewMarketData(GinkgoT())
		indicatorsSvc = indicators.NewIndicators(mockStore)
		mockProfCtx = mockProfiling.NewContext(GinkgoT())

		// Setup mock expectations for klines
		mockStore.On("GetKlines", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(makeMockKlines(100)).Maybe()

		// Setup mock expectations for GetAssetPrices
		mockStore.On("GetAssetPrices", mock.Anything).
			Return(market.PriceMap{
				"binance": connector.Price{Price: numerical.NewFromInt(100)},
			}).Maybe()

		// Attach profiling context
		ctx = profiling.WithContext(context.Background(), mockProfCtx)
	})

	Context("RSI indicator", func() {
		It("should record timing metrics when profiling context is present", func() {
			asset := portfolio.NewAsset("BTC")

			// Expect RecordIndicator to be called
			mockProfCtx.On("RecordIndicator", "RSI", mock.AnythingOfType("time.Duration")).Return().Once()

			// Call RSI with profiling context
			_, err := indicatorsSvc.RSI(ctx, asset, 14)
			Expect(err).ToNot(HaveOccurred())

			// Verify expectations
			mockProfCtx.AssertExpectations(GinkgoT())
		})

		It("should work normally without profiling context", func() {
			asset := portfolio.NewAsset("BTC")

			// Call RSI without profiling context - should not panic
			_, err := indicatorsSvc.RSI(context.Background(), asset, 14)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("Multiple indicators", func() {
		It("should record timing for all called indicators", func() {
			asset := portfolio.NewAsset("BTC")

			// Expect RecordIndicator to be called for each indicator
			mockProfCtx.On("RecordIndicator", "RSI", mock.AnythingOfType("time.Duration")).Return().Once()
			mockProfCtx.On("RecordIndicator", "SMA", mock.AnythingOfType("time.Duration")).Return().Once()
			mockProfCtx.On("RecordIndicator", "EMA", mock.AnythingOfType("time.Duration")).Return().Once()

			// Call multiple indicators
			_, _ = indicatorsSvc.RSI(ctx, asset, 14)
			_, _ = indicatorsSvc.SMA(ctx, asset, 20)
			_, _ = indicatorsSvc.EMA(ctx, asset, 50)

			// Verify all were called
			mockProfCtx.AssertExpectations(GinkgoT())
		})

		It("should accumulate timing when same indicator is called multiple times", func() {
			asset := portfolio.NewAsset("BTC")

			// Expect RecordIndicator to be called twice
			mockProfCtx.On("RecordIndicator", "RSI", mock.AnythingOfType("time.Duration")).Return().Twice()

			// Call RSI multiple times
			_, _ = indicatorsSvc.RSI(ctx, asset, 14)
			_, _ = indicatorsSvc.RSI(ctx, asset, 21)

			// Verify it was called twice
			mockProfCtx.AssertExpectations(GinkgoT())
		})
	})
})

// Helper function to create mock klines
func makeMockKlines(count int) []connector.Kline {
	klines := make([]connector.Kline, count)
	for i := 0; i < count; i++ {
		klines[i] = connector.Kline{
			Open:  numerical.NewFromInt(int64(100 + i)),
			High:  numerical.NewFromInt(int64(105 + i)),
			Low:   numerical.NewFromInt(int64(95 + i)),
			Close: numerical.NewFromInt(int64(100 + i)),
		}
	}
	return klines
}
