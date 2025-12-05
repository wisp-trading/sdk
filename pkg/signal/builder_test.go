package signal_test

import (
	"time"

	temporal "github.com/backtesting-org/kronos-sdk/pkg/runtime/time"
	"github.com/backtesting-org/kronos-sdk/pkg/signal"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	temporalType "github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SignalBuilder", func() {
	var (
		factory      strategy.SignalFactory
		timeProvider temporalType.TimeProvider
		strategyName strategy.StrategyName
		testAsset    portfolio.Asset
		testExchange connector.ExchangeName
	)

	BeforeEach(func() {
		timeProvider = temporal.NewTimeProvider()
		factory = signal.NewFactory(timeProvider)
		strategyName = strategy.StrategyName("test-strategy")
		testAsset = portfolio.NewAsset("BTC")
		testExchange = connector.ExchangeName("binance")
	})

	Describe("Buy", func() {
		It("should add a buy action with market order (price 0)", func() {
			quantity := numerical.NewFromInt(10)

			builder := factory.New(strategyName)
			signal := builder.Buy(testAsset, testExchange, quantity).Build()

			Expect(signal).NotTo(BeNil())
			Expect(signal.Strategy).To(Equal(strategyName))
			Expect(signal.Actions).To(HaveLen(1))
			Expect(signal.Actions[0].Action).To(Equal(strategy.ActionBuy))
			Expect(signal.Actions[0].Asset).To(Equal(testAsset))
			Expect(signal.Actions[0].Exchange).To(Equal(testExchange))
			Expect(signal.Actions[0].Quantity).To(Equal(quantity))
			Expect(signal.Actions[0].Price).To(Equal(numerical.NewFromInt(0)))
			Expect(signal.Timestamp).NotTo(BeZero())
			Expect(signal.ID).NotTo(Equal(""))
		})

		It("should allow chaining multiple buy actions", func() {
			qty1 := numerical.NewFromInt(10)
			qty2 := numerical.NewFromInt(20)
			asset2 := portfolio.NewAsset("ETH")

			builder := factory.New(strategyName)
			signal := builder.
				Buy(testAsset, testExchange, qty1).
				Buy(asset2, testExchange, qty2).
				Build()

			Expect(signal.Actions).To(HaveLen(2))
			Expect(signal.Actions[0].Asset).To(Equal(testAsset))
			Expect(signal.Actions[0].Quantity).To(Equal(qty1))
			Expect(signal.Actions[1].Asset).To(Equal(asset2))
			Expect(signal.Actions[1].Quantity).To(Equal(qty2))
		})
	})

	Describe("BuyLimit", func() {
		It("should add a buy action with specified limit price", func() {
			quantity := numerical.NewFromInt(10)
			price := numerical.NewFromFloat(50000.50)

			builder := factory.New(strategyName)
			signal := builder.BuyLimit(testAsset, testExchange, quantity, price).Build()

			Expect(signal.Actions).To(HaveLen(1))
			Expect(signal.Actions[0].Action).To(Equal(strategy.ActionBuy))
			Expect(signal.Actions[0].Asset).To(Equal(testAsset))
			Expect(signal.Actions[0].Exchange).To(Equal(testExchange))
			Expect(signal.Actions[0].Quantity).To(Equal(quantity))
			Expect(signal.Actions[0].Price).To(Equal(price))
		})
	})

	Describe("Sell", func() {
		It("should add a sell action with market order (price 0)", func() {
			quantity := numerical.NewFromInt(5)

			builder := factory.New(strategyName)
			signal := builder.Sell(testAsset, testExchange, quantity).Build()

			Expect(signal.Actions).To(HaveLen(1))
			Expect(signal.Actions[0].Action).To(Equal(strategy.ActionSell))
			Expect(signal.Actions[0].Asset).To(Equal(testAsset))
			Expect(signal.Actions[0].Exchange).To(Equal(testExchange))
			Expect(signal.Actions[0].Quantity).To(Equal(quantity))
			Expect(signal.Actions[0].Price).To(Equal(numerical.NewFromInt(0)))
		})
	})

	Describe("SellLimit", func() {
		It("should add a sell action with specified limit price", func() {
			quantity := numerical.NewFromInt(5)
			price := numerical.NewFromFloat(51000.25)

			builder := factory.New(strategyName)
			signal := builder.SellLimit(testAsset, testExchange, quantity, price).Build()

			Expect(signal.Actions).To(HaveLen(1))
			Expect(signal.Actions[0].Action).To(Equal(strategy.ActionSell))
			Expect(signal.Actions[0].Asset).To(Equal(testAsset))
			Expect(signal.Actions[0].Exchange).To(Equal(testExchange))
			Expect(signal.Actions[0].Quantity).To(Equal(quantity))
			Expect(signal.Actions[0].Price).To(Equal(price))
		})
	})

	Describe("SellShort", func() {
		It("should add a short sell action with market order (price 0)", func() {
			quantity := numerical.NewFromInt(8)

			builder := factory.New(strategyName)
			signal := builder.SellShort(testAsset, testExchange, quantity).Build()

			Expect(signal.Actions).To(HaveLen(1))
			Expect(signal.Actions[0].Action).To(Equal(strategy.ActionSellShort))
			Expect(signal.Actions[0].Asset).To(Equal(testAsset))
			Expect(signal.Actions[0].Exchange).To(Equal(testExchange))
			Expect(signal.Actions[0].Quantity).To(Equal(quantity))
			Expect(signal.Actions[0].Price).To(Equal(numerical.NewFromInt(0)))
		})
	})

	Describe("SellShortLimit", func() {
		It("should add a short sell action with specified limit price", func() {
			quantity := numerical.NewFromInt(8)
			price := numerical.NewFromFloat(49500.75)

			builder := factory.New(strategyName)
			signal := builder.SellShortLimit(testAsset, testExchange, quantity, price).Build()

			Expect(signal.Actions).To(HaveLen(1))
			Expect(signal.Actions[0].Action).To(Equal(strategy.ActionSellShort))
			Expect(signal.Actions[0].Asset).To(Equal(testAsset))
			Expect(signal.Actions[0].Exchange).To(Equal(testExchange))
			Expect(signal.Actions[0].Quantity).To(Equal(quantity))
			Expect(signal.Actions[0].Price).To(Equal(price))
		})
	})

	Describe("Complex signal building", func() {
		It("should support mixed order types in a single signal", func() {
			qty1 := numerical.NewFromInt(10)
			price1 := numerical.NewFromFloat(50000)
			qty2 := numerical.NewFromInt(5)
			price2 := numerical.NewFromFloat(51000)

			builder := factory.New(strategyName)
			signal := builder.
				Buy(testAsset, testExchange, qty1).
				SellLimit(testAsset, testExchange, qty2, price2).
				BuyLimit(portfolio.NewAsset("ETH"), testExchange, qty1, price1).
				Build()

			Expect(signal.Actions).To(HaveLen(3))
			Expect(signal.Actions[0].Action).To(Equal(strategy.ActionBuy))
			Expect(signal.Actions[0].Price).To(Equal(numerical.NewFromInt(0))) // market order
			Expect(signal.Actions[1].Action).To(Equal(strategy.ActionSell))
			Expect(signal.Actions[1].Price).To(Equal(price2)) // limit order
			Expect(signal.Actions[2].Action).To(Equal(strategy.ActionBuy))
			Expect(signal.Actions[2].Asset).To(Equal(portfolio.NewAsset("ETH")))
			Expect(signal.Actions[2].Price).To(Equal(price1)) // limit order
		})

		It("should support multiple exchanges", func() {
			exchange1 := connector.ExchangeName("binance")
			exchange2 := connector.ExchangeName("kraken")
			qty := numerical.NewFromInt(10)

			builder := factory.New(strategyName)
			signal := builder.
				Buy(testAsset, exchange1, qty).
				Sell(testAsset, exchange2, qty).
				Build()

			Expect(signal.Actions).To(HaveLen(2))
			Expect(signal.Actions[0].Exchange).To(Equal(exchange1))
			Expect(signal.Actions[1].Exchange).To(Equal(exchange2))
		})
	})

	Describe("Build", func() {
		It("should create signal with unique ID", func() {
			builder := factory.New(strategyName)
			signal1 := builder.Buy(testAsset, testExchange, numerical.NewFromInt(1)).Build()

			builder2 := factory.New(strategyName)
			signal2 := builder2.Buy(testAsset, testExchange, numerical.NewFromInt(1)).Build()

			Expect(signal1.ID).NotTo(Equal(signal2.ID))
		})

		It("should set timestamp when building signal", func() {
			builder := factory.New(strategyName)
			beforeTime := time.Now()
			signal := builder.Buy(testAsset, testExchange, numerical.NewFromInt(1)).Build()
			afterTime := time.Now()

			Expect(signal.Timestamp).To(BeTemporally(">=", beforeTime))
			Expect(signal.Timestamp).To(BeTemporally("<=", afterTime))
		})

		It("should handle empty actions list", func() {
			builder := factory.New(strategyName)
			signal := builder.Build()

			Expect(signal.Actions).To(HaveLen(0))
			Expect(signal.Strategy).To(Equal(strategyName))
			Expect(signal.Timestamp).NotTo(BeZero())
		})

		It("should preserve strategy name", func() {
			customStrategy := strategy.StrategyName("momentum-strategy")
			builder := factory.New(customStrategy)
			signal := builder.Buy(testAsset, testExchange, numerical.NewFromInt(1)).Build()

			Expect(signal.Strategy).To(Equal(customStrategy))
		})
	})

	Describe("Decimal handling", func() {
		It("should handle fractional quantities", func() {
			quantity := numerical.NewFromFloat(0.5)

			builder := factory.New(strategyName)
			signal := builder.Buy(testAsset, testExchange, quantity).Build()

			Expect(signal.Actions[0].Quantity).To(Equal(quantity))
		})

		It("should handle large decimal values", func() {
			quantity := numerical.NewFromFloat(1000000.123456)
			price := numerical.NewFromFloat(50000.987654)

			builder := factory.New(strategyName)
			signal := builder.BuyLimit(testAsset, testExchange, quantity, price).Build()

			Expect(signal.Actions[0].Quantity).To(Equal(quantity))
			Expect(signal.Actions[0].Price).To(Equal(price))
		})
	})
})
