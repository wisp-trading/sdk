package position_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wisp-trading/sdk/pkg/data/stores/activity/position"
	activityTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/activity"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

var _ = Describe("Position Store - Trades", func() {
	var store activityTypes.Positions

	BeforeEach(func() {
		store = position.NewStore()
	})

	It("adds and retrieves trades", func() {
		store.AddTrade(connector.Trade{ID: "t1", Price: numerical.NewFromInt(100)})
		store.AddTrade(connector.Trade{ID: "t2", Price: numerical.NewFromInt(200)})
		trades := store.GetTrades()
		Expect(trades).To(HaveLen(2))
		Expect(trades[0].ID).To(Equal("t1"))
		Expect(trades[1].ID).To(Equal("t2"))
	})

	It("clears all trades", func() {
		store.AddTrade(connector.Trade{ID: "t1"})
		store.Clear()
		Expect(store.GetTrades()).To(BeEmpty())
	})
})
