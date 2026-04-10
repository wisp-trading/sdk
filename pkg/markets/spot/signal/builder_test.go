package signal_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wisp-trading/sdk/pkg/markets/spot/signal"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	temporalType "github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
	temporal "github.com/wisp-trading/sdk/pkg/runtime/time"
)

var _ = Describe("SpotSignalBuilder", func() {
	var (
		timeProvider temporalType.TimeProvider
		strategyName strategy.StrategyName
		testAsset    portfolio.Pair
		testExchange connector.ExchangeName
	)

	BeforeEach(func() {
		timeProvider = temporal.NewTimeProvider()
		strategyName = strategy.StrategyName("test-strategy")
		testAsset = portfolio.NewPair(
			portfolio.NewAsset("BTC"),
			portfolio.NewAsset("USDT"),
		)
		testExchange = connector.ExchangeName("binance")
	})

	Describe("Buy", func() {
		It("should add a buy action with market order (price 0)", func() {
			quantity := numerical.NewFromInt(10)

			sig, err := signal.NewSpotBuilder(strategyName, timeProvider).Buy(testAsset, testExchange, quantity).Build()

			Expect(err).NotTo(HaveOccurred())
			Expect(sig).NotTo(BeNil())
			Expect(sig.GetStrategy()).To(Equal(strategyName))
			Expect(sig.GetActions()).To(HaveLen(1))
			Expect(sig.GetActions()[0].ActionType).To(Equal(strategy.ActionBuy))
			Expect(sig.GetActions()[0].Pair).To(Equal(testAsset))
			Expect(sig.GetActions()[0].Exchange).To(Equal(testExchange))
			Expect(sig.GetActions()[0].Quantity).To(Equal(quantity))
			Expect(sig.GetActions()[0].Price).To(Equal(numerical.NewFromInt(0)))
			Expect(sig.GetTimestamp()).NotTo(BeZero())
			Expect(sig.GetID()).NotTo(BeZero())
		})

		It("should allow chaining multiple buy actions", func() {
			qty1 := numerical.NewFromInt(10)
			qty2 := numerical.NewFromInt(20)
			asset2 := portfolio.NewPair(
				portfolio.NewAsset("ETH"),
				portfolio.NewAsset("USDT"),
			)

			sig, err := signal.NewSpotBuilder(strategyName, timeProvider).
				Buy(testAsset, testExchange, qty1).
				Buy(asset2, testExchange, qty2).
				Build()

			Expect(err).NotTo(HaveOccurred())
			Expect(sig.GetActions()).To(HaveLen(2))
			Expect(sig.GetActions()[0].Pair).To(Equal(testAsset))
			Expect(sig.GetActions()[0].Quantity).To(Equal(qty1))
			Expect(sig.GetActions()[1].Pair).To(Equal(asset2))
			Expect(sig.GetActions()[1].Quantity).To(Equal(qty2))
		})
	})

	Describe("BuyLimit", func() {
		It("should add a buy action with specified limit price", func() {
			quantity := numerical.NewFromInt(10)
			price := numerical.NewFromFloat(50000.50)

			sig, err := signal.NewSpotBuilder(strategyName, timeProvider).BuyLimit(testAsset, testExchange, quantity, price).Build()

			Expect(err).NotTo(HaveOccurred())
			Expect(sig.GetActions()).To(HaveLen(1))
			Expect(sig.GetActions()[0].ActionType).To(Equal(strategy.ActionBuy))
			Expect(sig.GetActions()[0].Price).To(Equal(price))
		})
	})

	Describe("Sell", func() {
		It("should add a sell action with market order (price 0)", func() {
			quantity := numerical.NewFromInt(5)

			sig, err := signal.NewSpotBuilder(strategyName, timeProvider).Sell(testAsset, testExchange, quantity).Build()

			Expect(err).NotTo(HaveOccurred())
			Expect(sig.GetActions()).To(HaveLen(1))
			Expect(sig.GetActions()[0].ActionType).To(Equal(strategy.ActionSell))
			Expect(sig.GetActions()[0].Price).To(Equal(numerical.NewFromInt(0)))
		})
	})

	Describe("SellLimit", func() {
		It("should add a sell action with specified limit price", func() {
			quantity := numerical.NewFromInt(5)
			price := numerical.NewFromFloat(51000.25)

			sig, err := signal.NewSpotBuilder(strategyName, timeProvider).SellLimit(testAsset, testExchange, quantity, price).Build()

			Expect(err).NotTo(HaveOccurred())
			Expect(sig.GetActions()).To(HaveLen(1))
			Expect(sig.GetActions()[0].ActionType).To(Equal(strategy.ActionSell))
			Expect(sig.GetActions()[0].Price).To(Equal(price))
		})
	})

	Describe("SellShort", func() {
		It("should add a short sell action with market order (price 0)", func() {
			quantity := numerical.NewFromInt(8)

			sig, err := signal.NewSpotBuilder(strategyName, timeProvider).SellShort(testAsset, testExchange, quantity).Build()

			Expect(err).NotTo(HaveOccurred())
			Expect(sig.GetActions()).To(HaveLen(1))
			Expect(sig.GetActions()[0].ActionType).To(Equal(strategy.ActionSellShort))
			Expect(sig.GetActions()[0].Price).To(Equal(numerical.NewFromInt(0)))
		})
	})

	Describe("SellShortLimit", func() {
		It("should add a short sell action with specified limit price", func() {
			quantity := numerical.NewFromInt(8)
			price := numerical.NewFromFloat(49500.75)

			sig, err := signal.NewSpotBuilder(strategyName, timeProvider).SellShortLimit(testAsset, testExchange, quantity, price).Build()

			Expect(err).NotTo(HaveOccurred())
			Expect(sig.GetActions()).To(HaveLen(1))
			Expect(sig.GetActions()[0].ActionType).To(Equal(strategy.ActionSellShort))
			Expect(sig.GetActions()[0].Price).To(Equal(price))
		})
	})

	Describe("Complex signal building", func() {
		It("should support mixed order types in a single signal", func() {
			qty1 := numerical.NewFromInt(10)
			price1 := numerical.NewFromFloat(50000)
			qty2 := numerical.NewFromInt(5)
			price2 := numerical.NewFromFloat(51000)

			sig, err := signal.NewSpotBuilder(strategyName, timeProvider).
				Buy(testAsset, testExchange, qty1).
				SellLimit(testAsset, testExchange, qty2, price2).
				BuyLimit(
					portfolio.NewPair(
						portfolio.NewAsset("ETH"),
						portfolio.NewAsset("USDT"),
					),
					testExchange,
					qty1,
					price1,
				).
				Build()

			Expect(err).NotTo(HaveOccurred())
			Expect(sig.GetActions()).To(HaveLen(3))
			Expect(sig.GetActions()[0].ActionType).To(Equal(strategy.ActionBuy))
			Expect(sig.GetActions()[0].Price).To(Equal(numerical.NewFromInt(0)))
			Expect(sig.GetActions()[1].ActionType).To(Equal(strategy.ActionSell))
			Expect(sig.GetActions()[1].Price).To(Equal(price2))
			Expect(sig.GetActions()[2].ActionType).To(Equal(strategy.ActionBuy))
			Expect(sig.GetActions()[2].Pair.Symbol()).To(Equal("ETH-USDT"))
			Expect(sig.GetActions()[2].Price).To(Equal(price1))
		})

		It("should support multiple exchanges", func() {
			exchange1 := connector.ExchangeName("binance")
			exchange2 := connector.ExchangeName("kraken")
			qty := numerical.NewFromInt(10)

			sig, err := signal.NewSpotBuilder(strategyName, timeProvider).
				Buy(testAsset, exchange1, qty).
				Sell(testAsset, exchange2, qty).
				Build()

			Expect(err).NotTo(HaveOccurred())
			Expect(sig.GetActions()).To(HaveLen(2))
			Expect(sig.GetActions()[0].Exchange).To(Equal(exchange1))
			Expect(sig.GetActions()[1].Exchange).To(Equal(exchange2))
		})
	})

	Describe("Build", func() {
		It("should create signal with unique ID", func() {
			sig1, err1 := signal.NewSpotBuilder(strategyName, timeProvider).Buy(testAsset, testExchange, numerical.NewFromInt(1)).Build()
			sig2, err2 := signal.NewSpotBuilder(strategyName, timeProvider).Buy(testAsset, testExchange, numerical.NewFromInt(1)).Build()

			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(sig1.GetID()).NotTo(Equal(sig2.GetID()))
		})

		It("should set timestamp when building signal", func() {
			beforeTime := time.Now()
			sig, err := signal.NewSpotBuilder(strategyName, timeProvider).Buy(testAsset, testExchange, numerical.NewFromInt(1)).Build()
			afterTime := time.Now()

			Expect(err).NotTo(HaveOccurred())
			Expect(sig.GetTimestamp()).To(BeTemporally(">=", beforeTime))
			Expect(sig.GetTimestamp()).To(BeTemporally("<=", afterTime))
		})

		It("should return an error when no actions have been added", func() {
			sig, err := signal.NewSpotBuilder(strategyName, timeProvider).Build()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("at least one action"))
			Expect(sig).To(BeNil())
		})

		It("should preserve strategy name", func() {
			customStrategy := strategy.StrategyName("momentum-strategy")
			sig, err := signal.NewSpotBuilder(customStrategy, timeProvider).Buy(testAsset, testExchange, numerical.NewFromInt(1)).Build()

			Expect(err).NotTo(HaveOccurred())
			Expect(sig.GetStrategy()).To(Equal(customStrategy))
		})

		It("should return an error when quantity is zero", func() {
			sig, err := signal.NewSpotBuilder(strategyName, timeProvider).Buy(testAsset, testExchange, numerical.NewFromInt(0)).Build()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("quantity must be positive"))
			Expect(sig).To(BeNil())
		})

		It("should return an error when exchange is empty", func() {
			sig, err := signal.NewSpotBuilder(strategyName, timeProvider).Buy(testAsset, connector.ExchangeName(""), numerical.NewFromInt(1)).Build()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("exchange is required"))
			Expect(sig).To(BeNil())
		})
	})

	Describe("Decimal handling", func() {
		It("should handle fractional quantities", func() {
			quantity := numerical.NewFromFloat(0.5)

			sig, err := signal.NewSpotBuilder(strategyName, timeProvider).Buy(testAsset, testExchange, quantity).Build()

			Expect(err).NotTo(HaveOccurred())
			Expect(sig.GetActions()[0].Quantity).To(Equal(quantity))
		})

		It("should handle large decimal values", func() {
			quantity := numerical.NewFromFloat(1000000.123456)
			price := numerical.NewFromFloat(50000.987654)

			sig, err := signal.NewSpotBuilder(strategyName, timeProvider).BuyLimit(testAsset, testExchange, quantity, price).Build()

			Expect(err).NotTo(HaveOccurred())
			Expect(sig.GetActions()[0].Quantity).To(Equal(quantity))
			Expect(sig.GetActions()[0].Price).To(Equal(price))
		})
	})
})
