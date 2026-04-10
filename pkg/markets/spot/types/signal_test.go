package types_test

import (
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

var _ = Describe("SpotSignal immutability", func() {
	var (
		btcUsdt portfolio.Pair
		binance connector.ExchangeName
		ts      time.Time
		id      uuid.UUID
		name    strategy.StrategyName
	)

	BeforeEach(func() {
		btcUsdt = portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
		binance = connector.ExchangeName("binance")
		ts = time.Now()
		id = uuid.New()
		name = strategy.StrategyName("test-strategy")
	})

	buildAction := func(t strategy.ActionType) spotTypes.SpotAction {
		return spotTypes.SpotAction{
			BaseAction: strategy.BaseAction{ActionType: t, Exchange: "binance"},
			Pair:       portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT")),
			Quantity:   numerical.NewFromInt(1),
			Price:      numerical.NewFromFloat(50000),
		}
	}

	Describe("GetActions", func() {
		It("returns a copy — mutating the result does not affect the signal", func() {
			sig := spotTypes.NewSpotSignal(id, name, ts, []spotTypes.SpotAction{buildAction(strategy.ActionBuy)})

			first := sig.GetActions()
			first[0].Price = numerical.NewFromFloat(99999)

			second := sig.GetActions()
			Expect(second[0].Price).NotTo(Equal(numerical.NewFromFloat(99999)),
				"mutating GetActions() result should not affect the internal state")
		})

		It("returns independent copies on repeated calls", func() {
			sig := spotTypes.NewSpotSignal(id, name, ts, []spotTypes.SpotAction{
				buildAction(strategy.ActionBuy),
				buildAction(strategy.ActionSell),
			})

			a := sig.GetActions()
			b := sig.GetActions()
			Expect(a).To(Equal(b))
			a[0].Price = numerical.NewFromFloat(1)
			Expect(b[0].Price).NotTo(Equal(numerical.NewFromFloat(1)))
		})

		It("returns an empty (non-nil) slice when no actions were provided", func() {
			sig := spotTypes.NewSpotSignal(id, name, ts, []spotTypes.SpotAction{})
			Expect(sig.GetActions()).NotTo(BeNil())
			Expect(sig.GetActions()).To(BeEmpty())
		})

		It("preserves all action fields", func() {
			action := buildAction(strategy.ActionBuy)
			sig := spotTypes.NewSpotSignal(id, name, ts, []spotTypes.SpotAction{action})

			got := sig.GetActions()
			Expect(got[0].ActionType).To(Equal(strategy.ActionBuy))
			Expect(got[0].Exchange).To(Equal(binance))
			Expect(got[0].Pair).To(Equal(btcUsdt))
			Expect(got[0].Quantity).To(Equal(numerical.NewFromInt(1)))
			Expect(got[0].Price).To(Equal(numerical.NewFromFloat(50000)))
		})
	})

	Describe("NewSpotSignal construction", func() {
		It("copies the source slice — mutating it after construction does not affect the signal", func() {
			actions := []spotTypes.SpotAction{buildAction(strategy.ActionBuy)}
			sig := spotTypes.NewSpotSignal(id, name, ts, actions)

			actions[0].Price = numerical.NewFromFloat(99999)

			got := sig.GetActions()
			Expect(got[0].Price).NotTo(Equal(numerical.NewFromFloat(99999)),
				"post-construction mutation of source slice should not affect the signal")
		})
	})

	Describe("metadata", func() {
		It("exposes the correct ID, strategy name, and timestamp", func() {
			sig := spotTypes.NewSpotSignal(id, name, ts, []spotTypes.SpotAction{buildAction(strategy.ActionBuy)})

			Expect(sig.GetID()).To(Equal(id))
			Expect(sig.GetStrategy()).To(Equal(name))
			Expect(sig.GetTimestamp()).To(BeTemporally("==", ts))
		})
	})
})
