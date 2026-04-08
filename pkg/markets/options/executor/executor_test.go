package executor_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	mockConnector "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/connector/options"
	execution "github.com/wisp-trading/sdk/pkg/markets/options/executor"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/registry"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	registryTypes "github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExecutor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Options Executor Suite")
}

func setupMockConnector(t GinkgoTInterface, name connector.ExchangeName) *mockConnector.MockConnector {
	m := mockConnector.NewMockConnector(t)
	m.EXPECT().GetConnectorInfo().Return(&connector.Info{
		Name: name,
	}).Maybe()
	return m
}

var _ = Describe("Options Executor", func() {
	var (
		executor          optionsTypes.OptionsExecutor
		connectorRegistry registryTypes.ConnectorRegistry
		logger            logging.ApplicationLogger
		btcPair           portfolio.Pair
		expiration        time.Time
	)

	BeforeEach(func() {
		logger = logging.NewNoOpLogger()
		connectorRegistry = registry.NewConnectorRegistry()
		executor = execution.NewExecutor(connectorRegistry, logger)
		btcPair = portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
		expiration = time.Now().AddDate(0, 0, 30)
	})

	Describe("PlaceOrder", func() {
		var mockConn *mockConnector.MockConnector

		BeforeEach(func() {
			exchangeName := connector.ExchangeName("test-exchange")
			mockConn = setupMockConnector(GinkgoT(), exchangeName)
			connectorRegistry.RegisterOptions(exchangeName, mockConn)
			connectorRegistry.MarkReady(exchangeName)
		})

		It("should place a limit order successfully", func() {
			mockConn.EXPECT().PlaceLimitOrder(
				btcPair,
				connector.OrderSideBuy,
				mock.MatchedBy(func(q numerical.Decimal) bool {
					f, _ := q.Float64()
					return f == 1.0
				}),
				mock.MatchedBy(func(p numerical.Decimal) bool {
					f, _ := p.Float64()
					return f == 50000.0
				}),
			).Return(&connector.OrderResponse{
				OrderID: "test-order-1",
			}, nil)

			contract := optionsTypes.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			order := optionsTypes.OptionOrder{
				Exchange: connector.ExchangeName("test-exchange"),
				Contract: contract,
				Side:     connector.OrderSideBuy,
				Quantity: 1.0,
				Price:    50000.0,
			}

			resp, err := executor.PlaceOrder(order)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.OrderID).To(Equal("test-order-1"))
		})

		It("should place a market order successfully", func() {
			mockConn.EXPECT().PlaceMarketOrder(
				btcPair,
				connector.OrderSideSell,
				mock.MatchedBy(func(q numerical.Decimal) bool {
					f, _ := q.Float64()
					return f == 2.0
				}),
			).Return(&connector.OrderResponse{
				OrderID: "test-order-2",
			}, nil)

			contract := optionsTypes.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "PUT",
			}

			order := optionsTypes.OptionOrder{
				Exchange: connector.ExchangeName("test-exchange"),
				Contract: contract,
				Side:     connector.OrderSideSell,
				Quantity: 2.0,
				Price:    0, // Market order
			}

			resp, err := executor.PlaceOrder(order)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.OrderID).To(Equal("test-order-2"))
		})

		It("should reject order with missing exchange", func() {
			contract := optionsTypes.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			order := optionsTypes.OptionOrder{
				Exchange: "", // Missing exchange
				Contract: contract,
				Side:     connector.OrderSideBuy,
				Quantity: 1.0,
				Price:    50000.0,
			}

			resp, err := executor.PlaceOrder(order)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("exchange"))
		})

		It("should reject order with invalid strike", func() {
			contract := optionsTypes.OptionContract{
				Pair:       btcPair,
				Strike:     0, // Invalid strike
				Expiration: expiration,
				OptionType: "CALL",
			}

			order := optionsTypes.OptionOrder{
				Exchange: connector.ExchangeName("test-exchange"),
				Contract: contract,
				Side:     connector.OrderSideBuy,
				Quantity: 1.0,
				Price:    50000.0,
			}

			resp, err := executor.PlaceOrder(order)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("strike"))
		})

		It("should reject order with invalid option type", func() {
			contract := optionsTypes.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "INVALID", // Invalid type
			}

			order := optionsTypes.OptionOrder{
				Exchange: connector.ExchangeName("test-exchange"),
				Contract: contract,
				Side:     connector.OrderSideBuy,
				Quantity: 1.0,
				Price:    50000.0,
			}

			_, err := executor.PlaceOrder(order)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("option type"))
		})

		It("should reject order with zero quantity", func() {
			contract := optionsTypes.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			order := optionsTypes.OptionOrder{
				Exchange: connector.ExchangeName("test-exchange"),
				Contract: contract,
				Side:     connector.OrderSideBuy,
				Quantity: 0, // Invalid quantity
				Price:    50000.0,
			}

			_, err := executor.PlaceOrder(order)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("quantity"))
		})

		It("should reject order with negative price", func() {
			contract := optionsTypes.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			order := optionsTypes.OptionOrder{
				Exchange: connector.ExchangeName("test-exchange"),
				Contract: contract,
				Side:     connector.OrderSideBuy,
				Quantity: 1.0,
				Price:    -100.0, // Invalid price
			}

			_, err := executor.PlaceOrder(order)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("price"))
		})

		It("should reject order for non-existent connector", func() {
			contract := optionsTypes.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			order := optionsTypes.OptionOrder{
				Exchange: connector.ExchangeName("non-existent"),
				Contract: contract,
				Side:     connector.OrderSideBuy,
				Quantity: 1.0,
				Price:    50000.0,
			}

			resp, err := executor.PlaceOrder(order)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})
	})

	Describe("CancelOrder", func() {
		var mockConn *mockConnector.MockConnector

		BeforeEach(func() {
			exchangeName := connector.ExchangeName("test-exchange")
			mockConn = setupMockConnector(GinkgoT(), exchangeName)
			connectorRegistry.RegisterOptions(exchangeName, mockConn)
			connectorRegistry.MarkReady(exchangeName)
		})

		It("should cancel an order successfully", func() {
			mockConn.EXPECT().CancelOrder("test-order-1").Return(&connector.CancelResponse{
				OrderID: "test-order-1",
				Status:  "cancelled",
			}, nil)

			resp, err := executor.CancelOrder("test-order-1", connector.ExchangeName("test-exchange"))
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.OrderID).To(Equal("test-order-1"))
		})

		It("should reject cancel with empty order ID", func() {
			resp, err := executor.CancelOrder("", connector.ExchangeName("test-exchange"))
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("order ID"))
		})

		It("should reject cancel for non-existent connector", func() {
			resp, err := executor.CancelOrder("test-order-1", connector.ExchangeName("non-existent"))
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})
	})
})
