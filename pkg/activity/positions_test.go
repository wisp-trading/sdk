package activity_test

import (
	"context"

	sdkTesting "github.com/backtesting-org/kronos-sdk/pkg/testing"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	storeActivity "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/activity"
	kronosActivity "github.com/backtesting-org/kronos-sdk/pkg/types/kronos/activity"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("Positions", func() {
	var (
		app       *fxtest.App
		positions kronosActivity.Positions
		store     storeActivity.Positions
	)

	BeforeEach(func() {
		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(
				&positions,
				&store,
			),
			fx.NopLogger,
		)

		app.RequireStart()
	})

	AfterEach(func() {
		app.RequireStop()
	})

	Describe("GetStrategyExecution", func() {
		It("should return execution for strategy in context", func() {
			strategyName := strategy.StrategyName("test-strategy")
			ctx := strategy.NewStrategyContext(context.Background(), strategyName)

			// Populate store with test data
			expectedExecution := &strategy.StrategyExecution{
				Orders: []connector.Order{{ID: "order-1"}},
				Trades: []connector.Trade{{ID: "trade-1"}},
			}
			store.StoreStrategyExecution(strategyName, expectedExecution)

			result := positions.GetStrategyExecution(ctx)

			Expect(result).To(Equal(expectedExecution))
		})

		It("should return nil when no strategy in context", func() {
			ctx := strategy.NewStrategyContext(context.Background(), "")

			result := positions.GetStrategyExecution(ctx)

			Expect(result).To(BeNil())
		})

		It("should return nil for unknown strategy", func() {
			strategyName := strategy.StrategyName("unknown")
			ctx := strategy.NewStrategyContext(context.Background(), strategyName)

			result := positions.GetStrategyExecution(ctx)

			Expect(result).To(BeNil())
		})
	})

	Describe("GetAllStrategyExecutions", func() {
		It("should return all executions from store", func() {
			ctx := strategy.NewStrategyContext(context.Background(), "")

			// Populate store with multiple strategies
			store.StoreStrategyExecution("strategy-1", &strategy.StrategyExecution{
				Orders: []connector.Order{{ID: "order-1"}},
			})
			store.StoreStrategyExecution("strategy-2", &strategy.StrategyExecution{
				Orders: []connector.Order{{ID: "order-2"}},
			})

			result := positions.GetAllStrategyExecutions(ctx)

			Expect(result).To(HaveLen(2))
			Expect(result["strategy-1"].Orders[0].ID).To(Equal("order-1"))
			Expect(result["strategy-2"].Orders[0].ID).To(Equal("order-2"))
		})

		It("should return empty map when no executions exist", func() {
			ctx := strategy.NewStrategyContext(context.Background(), "")

			result := positions.GetAllStrategyExecutions(ctx)

			Expect(result).To(BeEmpty())
		})
	})

	Describe("GetStrategyForOrder", func() {
		It("should find strategy for existing order", func() {
			ctx := strategy.NewStrategyContext(context.Background(), "")

			strategyName := strategy.StrategyName("test-strategy")
			order := connector.Order{ID: "order-123"}
			store.AddOrderToStrategy(strategyName, order)

			name, found := positions.GetStrategyForOrder(ctx, "order-123")

			Expect(found).To(BeTrue())
			Expect(name).To(Equal(strategyName))
		})

		It("should return false for unknown order", func() {
			ctx := strategy.NewStrategyContext(context.Background(), "")

			_, found := positions.GetStrategyForOrder(ctx, "unknown-order")

			Expect(found).To(BeFalse())
		})
	})

	Describe("GetTotalOrderCount", func() {
		It("should return count from store", func() {
			ctx := strategy.NewStrategyContext(context.Background(), "")

			// Add orders to multiple strategies
			store.AddOrderToStrategy("strategy-1", connector.Order{ID: "order-1"})
			store.AddOrderToStrategy("strategy-1", connector.Order{ID: "order-2"})
			store.AddOrderToStrategy("strategy-2", connector.Order{ID: "order-3"})

			result := positions.GetTotalOrderCount(ctx)

			Expect(result).To(Equal(int64(3)))
		})

		It("should return zero when no orders exist", func() {
			ctx := strategy.NewStrategyContext(context.Background(), "")

			result := positions.GetTotalOrderCount(ctx)

			Expect(result).To(Equal(int64(0)))
		})
	})

	Describe("GetTradesForStrategy", func() {
		It("should return trades for strategy in context", func() {
			strategyName := strategy.StrategyName("test-strategy")
			ctx := strategy.NewStrategyContext(context.Background(), strategyName)

			// Add trades to store
			store.AddTradeToStrategy(strategyName, connector.Trade{
				ID:    "trade-1",
				Price: numerical.NewFromFloat(100),
			})
			store.AddTradeToStrategy(strategyName, connector.Trade{
				ID:    "trade-2",
				Price: numerical.NewFromFloat(200),
			})

			result := positions.GetTradesForStrategy(ctx)

			Expect(result).To(HaveLen(2))
			Expect(result[0].ID).To(Equal("trade-1"))
			Expect(result[1].ID).To(Equal("trade-2"))
		})

		It("should return nil when no strategy in context", func() {
			ctx := strategy.NewStrategyContext(context.Background(), "")

			result := positions.GetTradesForStrategy(ctx)

			Expect(len(result)).To(BeZero())
		})

		It("should return empty slice when no trades exist for strategy", func() {
			strategyName := strategy.StrategyName("test-strategy")
			ctx := strategy.NewStrategyContext(context.Background(), strategyName)

			// Ensure strategy exists but with no trades
			store.StoreStrategyExecution(strategyName, &strategy.StrategyExecution{
				Orders: []connector.Order{},
				Trades: []connector.Trade{},
			})

			result := positions.GetTradesForStrategy(ctx)

			Expect(result).To(BeEmpty())
		})
	})
})
