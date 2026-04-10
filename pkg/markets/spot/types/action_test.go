package types_test

import (
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SpotAction", func() {
	var (
		btcUsdt portfolio.Pair
		binance connector.ExchangeName
	)

	BeforeEach(func() {
		btcUsdt = portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
		binance = connector.ExchangeName("binance")
	})

	Describe("GetMarketType", func() {
		It("returns spot market type", func() {
			action := &spotTypes.SpotAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: binance},
				Pair:       btcUsdt,
				Quantity:   numerical.NewFromFloat(1),
				Price:      numerical.NewFromFloat(50000),
			}
			Expect(action.GetMarketType()).To(Equal(connector.MarketTypeSpot))
		})
	})

	Describe("Validate", func() {
		It("passes for a valid buy action", func() {
			action := &spotTypes.SpotAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: binance},
				Pair:       btcUsdt,
				Quantity:   numerical.NewFromFloat(0.5),
				Price:      numerical.NewFromFloat(50000),
			}
			Expect(action.Validate()).To(Succeed())
		})

		It("passes for a market order (zero price)", func() {
			action := &spotTypes.SpotAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: binance},
				Pair:       btcUsdt,
				Quantity:   numerical.NewFromFloat(1),
				Price:      numerical.NewFromInt(0),
			}
			Expect(action.Validate()).To(Succeed())
		})

		It("fails when exchange is empty", func() {
			action := &spotTypes.SpotAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: ""},
				Pair:       btcUsdt,
				Quantity:   numerical.NewFromFloat(1),
				Price:      numerical.NewFromFloat(50000),
			}
			Expect(action.Validate()).To(MatchError(ContainSubstring("exchange is required")))
		})

		It("fails when quantity is zero", func() {
			action := &spotTypes.SpotAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: binance},
				Pair:       btcUsdt,
				Quantity:   numerical.NewFromInt(0),
				Price:      numerical.NewFromFloat(50000),
			}
			Expect(action.Validate()).To(MatchError(ContainSubstring("quantity must be positive")))
		})

		It("fails when quantity is negative", func() {
			action := &spotTypes.SpotAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: binance},
				Pair:       btcUsdt,
				Quantity:   numerical.NewFromInt(-1),
				Price:      numerical.NewFromFloat(50000),
			}
			Expect(action.Validate()).To(MatchError(ContainSubstring("quantity must be positive")))
		})

		It("fails when price is negative", func() {
			action := &spotTypes.SpotAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: binance},
				Pair:       btcUsdt,
				Quantity:   numerical.NewFromFloat(1),
				Price:      numerical.NewFromInt(-1),
			}
			Expect(action.Validate()).To(MatchError(ContainSubstring("price must not be negative")))
		})
	})

	Describe("BaseAction embedding", func() {
		It("inherits GetType and GetExchange from BaseAction", func() {
			action := &spotTypes.SpotAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionSell, Exchange: binance},
				Pair:       btcUsdt,
				Quantity:   numerical.NewFromFloat(1),
				Price:      numerical.NewFromFloat(50000),
			}
			Expect(action.GetType()).To(Equal(strategy.ActionSell))
			Expect(action.GetExchange()).To(Equal(binance))
		})
	})
})
