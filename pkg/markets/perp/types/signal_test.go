package types_test

import (
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

var _ = Describe("PerpSignal immutability", func() {
	var (
		ethUsdt portfolio.Pair
		bybit   connector.ExchangeName
		ts      time.Time
		id      uuid.UUID
		name    strategy.StrategyName
	)

	BeforeEach(func() {
		ethUsdt = portfolio.NewPair(portfolio.NewAsset("ETH"), portfolio.NewAsset("USDT"))
		bybit = connector.ExchangeName("bybit")
		ts = time.Now()
		id = uuid.New()
		name = strategy.StrategyName("test-strategy")
	})

	buildAction := func(t strategy.ActionType) perpTypes.PerpAction {
		return perpTypes.PerpAction{
			BaseAction: strategy.BaseAction{ActionType: t, Exchange: "bybit"},
			Pair:       portfolio.NewPair(portfolio.NewAsset("ETH"), portfolio.NewAsset("USDT")),
			Quantity:   numerical.NewFromInt(1),
			Price:      numerical.NewFromFloat(3000),
		}
	}

	Describe("GetActions", func() {
		It("returns a copy — mutating the result does not affect the signal", func() {
			sig := perpTypes.NewPerpSignal(id, name, ts, []perpTypes.PerpAction{buildAction(strategy.ActionBuy)})

			first := sig.GetActions()
			first[0].Price = numerical.NewFromFloat(99999)

			second := sig.GetActions()
			Expect(second[0].Price).NotTo(Equal(numerical.NewFromFloat(99999)),
				"mutating GetActions() result should not affect the internal state")
		})

		It("returns independent copies on repeated calls", func() {
			sig := perpTypes.NewPerpSignal(id, name, ts, []perpTypes.PerpAction{
				buildAction(strategy.ActionBuy),
				buildAction(strategy.ActionSellShort),
			})

			a := sig.GetActions()
			b := sig.GetActions()
			Expect(a).To(Equal(b))
			a[0].Leverage = numerical.NewFromInt(10)
			Expect(b[0].Leverage).NotTo(Equal(numerical.NewFromInt(10)))
		})

		It("returns an empty (non-nil) slice when no actions were provided", func() {
			sig := perpTypes.NewPerpSignal(id, name, ts, []perpTypes.PerpAction{})
			Expect(sig.GetActions()).NotTo(BeNil())
			Expect(sig.GetActions()).To(BeEmpty())
		})

		It("preserves all action fields including leverage", func() {
			action := perpTypes.PerpAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: bybit},
				Pair:       ethUsdt,
				Quantity:   numerical.NewFromFloat(0.5),
				Price:      numerical.NewFromFloat(3000),
				Leverage:   numerical.NewFromInt(5),
			}
			sig := perpTypes.NewPerpSignal(id, name, ts, []perpTypes.PerpAction{action})

			got := sig.GetActions()
			Expect(got[0].ActionType).To(Equal(strategy.ActionBuy))
			Expect(got[0].Exchange).To(Equal(bybit))
			Expect(got[0].Pair).To(Equal(ethUsdt))
			Expect(got[0].Quantity).To(Equal(numerical.NewFromFloat(0.5)))
			Expect(got[0].Price).To(Equal(numerical.NewFromFloat(3000)))
			Expect(got[0].Leverage).To(Equal(numerical.NewFromInt(5)))
		})
	})

	Describe("NewPerpSignal construction", func() {
		It("copies the source slice — mutating it after construction does not affect the signal", func() {
			actions := []perpTypes.PerpAction{buildAction(strategy.ActionBuy)}
			sig := perpTypes.NewPerpSignal(id, name, ts, actions)

			actions[0].Price = numerical.NewFromFloat(99999)

			got := sig.GetActions()
			Expect(got[0].Price).NotTo(Equal(numerical.NewFromFloat(99999)),
				"post-construction mutation of source slice should not affect the signal")
		})
	})

	Describe("metadata", func() {
		It("exposes the correct ID, strategy name, and timestamp", func() {
			sig := perpTypes.NewPerpSignal(id, name, ts, []perpTypes.PerpAction{buildAction(strategy.ActionBuy)})

			Expect(sig.GetID()).To(Equal(id))
			Expect(sig.GetStrategy()).To(Equal(name))
			Expect(sig.GetTimestamp()).To(BeTemporally("==", ts))
		})
	})
})
