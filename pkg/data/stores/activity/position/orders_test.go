package position_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wisp-trading/sdk/pkg/data/stores/activity/position"
	activityTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/activity"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

var _ = Describe("Position Store - Orders", func() {
	var store activityTypes.Positions

	BeforeEach(func() {
		store = position.NewStore()
	})

	It("adds and counts orders", func() {
		store.AddOrder(connector.Order{ID: "o1", Quantity: numerical.NewFromInt(1)})
		store.AddOrder(connector.Order{ID: "o2", Quantity: numerical.NewFromInt(2)})
		Expect(store.GetTotalOrderCount()).To(Equal(int64(2)))
	})

	It("updates order status", func() {
		store.AddOrder(connector.Order{ID: "o1", Status: connector.OrderStatusNew})
		err := store.UpdateOrderStatus("o1", connector.OrderStatusFilled)
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns error for unknown order", func() {
		err := store.UpdateOrderStatus("unknown", connector.OrderStatusFilled)
		Expect(err).To(HaveOccurred())
	})

	It("clears all orders", func() {
		store.AddOrder(connector.Order{ID: "o1"})
		store.Clear()
		Expect(store.GetTotalOrderCount()).To(Equal(int64(0)))
	})
})
