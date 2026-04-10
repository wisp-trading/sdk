package types_test

import (
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PerpAction", func() {
	var (
		ethUsdt portfolio.Pair
		bybit   connector.ExchangeName
	)

	BeforeEach(func() {
		ethUsdt = portfolio.NewPair(portfolio.NewAsset("ETH"), portfolio.NewAsset("USDT"))
		bybit = connector.ExchangeName("bybit")
	})

	Describe("GetMarketType", func() {
		It("returns perp market type", func() {
			action := &perpTypes.PerpAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: bybit},
				Pair:       ethUsdt,
				Quantity:   numerical.NewFromFloat(1),
				Price:      numerical.NewFromFloat(3000),
			}
			Expect(action.GetMarketType()).To(Equal(connector.MarketTypePerp))
		})
	})

	Describe("Validate", func() {
		It("passes for a valid buy action with leverage", func() {
			action := &perpTypes.PerpAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: bybit},
				Pair:       ethUsdt,
				Quantity:   numerical.NewFromFloat(10),
				Price:      numerical.NewFromFloat(3000),
				Leverage:   numerical.NewFromInt(10),
			}
			Expect(action.Validate()).To(Succeed())
		})

		It("passes for a market order (zero price)", func() {
			action := &perpTypes.PerpAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: bybit},
				Pair:       ethUsdt,
				Quantity:   numerical.NewFromFloat(1),
				Price:      numerical.NewFromInt(0),
			}
			Expect(action.Validate()).To(Succeed())
		})

		It("fails when exchange is empty", func() {
			action := &perpTypes.PerpAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: ""},
				Pair:       ethUsdt,
				Quantity:   numerical.NewFromFloat(1),
				Price:      numerical.NewFromFloat(3000),
			}
			Expect(action.Validate()).To(MatchError(ContainSubstring("exchange is required")))
		})

		It("fails when quantity is zero", func() {
			action := &perpTypes.PerpAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: bybit},
				Pair:       ethUsdt,
				Quantity:   numerical.NewFromInt(0),
				Price:      numerical.NewFromFloat(3000),
			}
			Expect(action.Validate()).To(MatchError(ContainSubstring("quantity must be positive")))
		})

		It("fails when price is negative", func() {
			action := &perpTypes.PerpAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: bybit},
				Pair:       ethUsdt,
				Quantity:   numerical.NewFromFloat(1),
				Price:      numerical.NewFromInt(-1),
			}
			Expect(action.Validate()).To(MatchError(ContainSubstring("price must not be negative")))
		})
	})

	Describe("BaseAction embedding", func() {
		It("inherits GetType and GetExchange from BaseAction", func() {
			action := &perpTypes.PerpAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionSellShort, Exchange: bybit},
				Pair:       ethUsdt,
				Quantity:   numerical.NewFromFloat(1),
				Price:      numerical.NewFromFloat(3000),
			}
			Expect(action.GetType()).To(Equal(strategy.ActionSellShort))
			Expect(action.GetExchange()).To(Equal(bybit))
		})
	})
})
