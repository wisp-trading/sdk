package executor_test

import (
	"errors"
	"testing"
	stdtime "time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	mockConnector "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/connector/options"
	mockOptionsStore "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/markets/options/types"
	optionsExecutor "github.com/wisp-trading/sdk/pkg/markets/options/executor"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/registry"
	rtime "github.com/wisp-trading/sdk/pkg/runtime/time"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	registryTypes "github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExecutor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Options Executor Suite")
}

var _ = Describe("Options SignalExecutor", func() {
	var (
		exec              optionsTypes.SignalExecutor
		connectorRegistry registryTypes.ConnectorRegistry
		store             *mockOptionsStore.OptionsStore
		btcPair           portfolio.Pair
		expiration        stdtime.Time
		strategyName      strategy.StrategyName
	)

	BeforeEach(func() {
		connectorRegistry = registry.NewConnectorRegistry()
		store = mockOptionsStore.NewOptionsStore(GinkgoT())
		exec = optionsExecutor.NewExecutor(
			connectorRegistry,
			store,
			logging.NewNoOpLogger(),
			rtime.NewTimeProvider(),
		)
		btcPair = portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
		expiration = stdtime.Now().AddDate(0, 0, 30)
		strategyName = strategy.StrategyName("test-strategy")
	})

	buildSignal := func(actions ...optionsTypes.OptionsAction) optionsTypes.OptionsSignal {
		return optionsTypes.NewOptionsSignal(uuid.New(), strategyName, stdtime.Now(), actions)
	}

	buildAction := func(actionType strategy.ActionType, price float64) optionsTypes.OptionsAction {
		return optionsTypes.OptionsAction{
			BaseAction: strategy.BaseAction{ActionType: actionType, Exchange: "test-exchange"},
			Contract: optionsTypes.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			},
			Quantity: numerical.NewFromFloat(1.0),
			Price:    numerical.NewFromFloat(price),
		}
	}

	Describe("ExecuteOptionsSignal", func() {
		var (
			mockConn *mockConnector.Connector
			ctx      *execution.ExecutionContext
			result   *execution.ExecutionResult
		)

		BeforeEach(func() {
			exchangeName := connector.ExchangeName("test-exchange")
			mockConn = mockConnector.NewConnector(GinkgoT())
			mockConn.EXPECT().GetConnectorInfo().Return(&connector.Info{Name: exchangeName}).Maybe()
			connectorRegistry.RegisterOptions(exchangeName, mockConn)
			connectorRegistry.MarkReady(exchangeName)

			ctx = &execution.ExecutionContext{Timestamp: stdtime.Now(), Metadata: make(map[string]interface{})}
			result = &execution.ExecutionResult{OrderIDs: make([]string, 0), Success: true}
		})

		It("places a limit order and records the order ID", func() {
			mockConn.EXPECT().PlaceLimitOrder(
				btcPair, connector.OrderSideBuy,
				mock.MatchedBy(func(q numerical.Decimal) bool { f, _ := q.Float64(); return f == 1.0 }),
				mock.MatchedBy(func(p numerical.Decimal) bool { f, _ := p.Float64(); return f == 50000.0 }),
			).Return(&connector.OrderResponse{OrderID: "order-1"}, nil)
			store.EXPECT().AddOrder(mock.Anything)

			sig := buildSignal(buildAction(strategy.ActionBuy, 50000.0))
			err := exec.ExecuteOptionsSignal(sig, ctx, result)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.OrderIDs).To(ConsistOf("order-1"))
		})

		It("places a market order when price is zero", func() {
			mockConn.EXPECT().PlaceMarketOrder(
				btcPair, connector.OrderSideSell,
				mock.MatchedBy(func(q numerical.Decimal) bool { f, _ := q.Float64(); return f == 1.0 }),
			).Return(&connector.OrderResponse{OrderID: "order-2"}, nil)
			store.EXPECT().AddOrder(mock.Anything)

			action := buildAction(strategy.ActionSell, 0)
			sig := buildSignal(action)
			err := exec.ExecuteOptionsSignal(sig, ctx, result)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.OrderIDs).To(ConsistOf("order-2"))
		})

		It("returns error when action validation fails", func() {
			action := optionsTypes.OptionsAction{
				BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: "test-exchange"},
				Contract: optionsTypes.OptionContract{
					Pair: btcPair, Strike: 0, Expiration: expiration, OptionType: "CALL", // invalid strike
				},
				Quantity: numerical.NewFromFloat(1.0),
				Price:    numerical.NewFromFloat(50000),
			}

			sig := buildSignal(action)
			err := exec.ExecuteOptionsSignal(sig, ctx, result)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid"))
		})

		It("returns error when exchange is unavailable", func() {
			action := buildAction(strategy.ActionBuy, 50000)
			action.Exchange = "unknown-exchange"
			sig := buildSignal(action)

			err := exec.ExecuteOptionsSignal(sig, ctx, result)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not available"))
		})

		It("returns error when the exchange rejects the order", func() {
			mockConn.EXPECT().PlaceLimitOrder(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(nil, errors.New("insufficient margin"))

			sig := buildSignal(buildAction(strategy.ActionBuy, 50000))
			err := exec.ExecuteOptionsSignal(sig, ctx, result)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to place"))
		})

		It("skips hold and close actions without placing orders", func() {
			holdSig := buildSignal(buildAction(strategy.ActionHold, 0))
			err := exec.ExecuteOptionsSignal(holdSig, ctx, result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.OrderIDs).To(BeEmpty())

			closeSig := buildSignal(buildAction(strategy.ActionClose, 0))
			err = exec.ExecuteOptionsSignal(closeSig, ctx, result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.OrderIDs).To(BeEmpty())
		})
	})
})
