package registry_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	perpMock "github.com/wisp-trading/wisp/mocks/github.com/wisp-trading/wisp/pkg/types/connector/perp"
	spotMock "github.com/wisp-trading/wisp/mocks/github.com/wisp-trading/wisp/pkg/types/connector/spot"
	sdkTesting "github.com/wisp-trading/wisp/pkg/testing"
	registryTypes "github.com/wisp-trading/wisp/pkg/types/registry"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("ConnectorRegistry", func() {
	var (
		app          *fxtest.App
		connectorReg registryTypes.ConnectorRegistry
	)

	BeforeEach(func() {
		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(&connectorReg),
			fx.NopLogger,
		)
		app.RequireStart()
	})

	AfterEach(func() {
		app.RequireStop()
	})

	Describe("Spot Connector Registration", func() {
		var (
			spotConn1 *spotMock.Connector
			spotConn2 *spotMock.Connector
		)

		BeforeEach(func() {
			spotConn1 = spotMock.NewConnector(GinkgoT())
			spotConn2 = spotMock.NewConnector(GinkgoT())
		})

		Context("when registering a new spot connector", func() {
			It("should store the connector", func() {
				connectorReg.RegisterSpotConnector("binance", spotConn1)

				conn, exists := connectorReg.GetSpotConnector("binance")
				Expect(exists).To(BeTrue())
				Expect(conn).To(Equal(spotConn1))
			})

			It("should not be marked as ready initially", func() {
				connectorReg.RegisterSpotConnector("binance", spotConn1)

				Expect(connectorReg.IsConnectorReady("binance")).To(BeFalse())
			})

			It("should be retrievable via generic GetConnector", func() {
				connectorReg.RegisterSpotConnector("binance", spotConn1)

				conn, exists := connectorReg.GetConnector("binance")
				Expect(exists).To(BeTrue())
				Expect(conn).NotTo(BeNil())
			})
		})

		Context("when registering duplicate spot connector - BUG TEST", func() {
			It("should overwrite the previous connector", func() {
				connectorReg.RegisterSpotConnector("binance", spotConn1)
				connectorReg.RegisterSpotConnector("binance", spotConn2)

				conn, exists := connectorReg.GetSpotConnector("binance")
				Expect(exists).To(BeTrue())
				Expect(conn).To(Equal(spotConn2), "should have the second connector, not the first")
			})

			It("should reset ready state on re-registration - THE BUG", func() {
				connectorReg.RegisterSpotConnector("binance", spotConn1)
				err := connectorReg.MarkConnectorReady("binance")
				Expect(err).NotTo(HaveOccurred())
				Expect(connectorReg.IsConnectorReady("binance")).To(BeTrue())

				// Re-register same exchange (this is the bug scenario)
				connectorReg.RegisterSpotConnector("binance", spotConn2)

				// BUG: Ready state is reset to false without warning
				Expect(connectorReg.IsConnectorReady("binance")).To(BeFalse(),
					"BUG: ready state is reset when re-registering - should either error or preserve state")
			})
		})

		Context("when marking connector as ready", func() {
			It("should update ready state", func() {
				connectorReg.RegisterSpotConnector("binance", spotConn1)

				err := connectorReg.MarkConnectorReady("binance")
				Expect(err).NotTo(HaveOccurred())
				Expect(connectorReg.IsConnectorReady("binance")).To(BeTrue())
			})

			It("should include connector in ready list", func() {
				connectorReg.RegisterSpotConnector("binance", spotConn1)
				connectorReg.RegisterSpotConnector("coinbase", spotConn2)
				_ = connectorReg.MarkConnectorReady("binance")

				readySpot := connectorReg.GetReadySpotConnectors()
				Expect(readySpot).To(HaveLen(1))
			})
		})
	})

	Describe("Perp Connector Registration", func() {
		var (
			perpConn1 *perpMock.Connector
			perpConn2 *perpMock.Connector
		)

		BeforeEach(func() {
			perpConn1 = perpMock.NewConnector(GinkgoT())
			perpConn2 = perpMock.NewConnector(GinkgoT())
		})

		Context("when registering a new perp connector", func() {
			It("should store the connector", func() {
				connectorReg.RegisterPerpConnector("hyperliquid", perpConn1)

				conn, exists := connectorReg.GetPerpConnector("hyperliquid")
				Expect(exists).To(BeTrue())
				Expect(conn).To(Equal(perpConn1))
			})

			It("should not be marked as ready initially", func() {
				connectorReg.RegisterPerpConnector("hyperliquid", perpConn1)

				Expect(connectorReg.IsConnectorReady("hyperliquid")).To(BeFalse())
			})

			It("should be retrievable via generic GetConnector", func() {
				connectorReg.RegisterPerpConnector("hyperliquid", perpConn1)

				conn, exists := connectorReg.GetConnector("hyperliquid")
				Expect(exists).To(BeTrue())
				Expect(conn).NotTo(BeNil())
			})
		})

		Context("when registering duplicate perp connector - THE BUG", func() {
			It("should overwrite the previous connector", func() {
				connectorReg.RegisterPerpConnector("hyperliquid", perpConn1)
				connectorReg.RegisterPerpConnector("hyperliquid", perpConn2)

				conn, exists := connectorReg.GetPerpConnector("hyperliquid")
				Expect(exists).To(BeTrue())
				Expect(conn).To(Equal(perpConn2), "BUG: should have the second connector")
			})

			It("should reset ready state on re-registration - THE BUG", func() {
				connectorReg.RegisterPerpConnector("hyperliquid", perpConn1)
				err := connectorReg.MarkConnectorReady("hyperliquid")
				Expect(err).NotTo(HaveOccurred())
				Expect(connectorReg.IsConnectorReady("hyperliquid")).To(BeTrue())

				// Re-register same exchange (this is the production bug with hyperliquid)
				connectorReg.RegisterPerpConnector("hyperliquid", perpConn2)

				// BUG: Ready state is reset to false, causing the connector to appear unready
				Expect(connectorReg.IsConnectorReady("hyperliquid")).To(BeFalse(),
					"BUG: ready state is reset when re-registering hyperliquid - this is the production bug!")
			})
		})

		Context("when marking connector as ready", func() {
			It("should update ready state", func() {
				connectorReg.RegisterPerpConnector("hyperliquid", perpConn1)

				err := connectorReg.MarkConnectorReady("hyperliquid")
				Expect(err).NotTo(HaveOccurred())
				Expect(connectorReg.IsConnectorReady("hyperliquid")).To(BeTrue())
			})

			It("should include connector in ready list", func() {
				connectorReg.RegisterPerpConnector("hyperliquid", perpConn1)
				connectorReg.RegisterPerpConnector("dydx", perpConn2)
				_ = connectorReg.MarkConnectorReady("hyperliquid")

				readyPerp := connectorReg.GetReadyPerpConnectors()
				Expect(readyPerp).To(HaveLen(1))
			})
		})
	})

	Describe("Multi-Connector Operations", func() {
		var (
			spotConn1 *spotMock.Connector
			spotConn2 *spotMock.Connector
			perpConn1 *perpMock.Connector
			perpConn2 *perpMock.Connector
		)

		BeforeEach(func() {
			spotConn1 = spotMock.NewConnector(GinkgoT())
			spotConn2 = spotMock.NewConnector(GinkgoT())
			perpConn1 = perpMock.NewConnector(GinkgoT())
			perpConn2 = perpMock.NewConnector(GinkgoT())
		})

		Context("when multiple connector types are registered", func() {
			It("should return all spot and perp connectors separately", func() {
				connectorReg.RegisterSpotConnector("binance", spotConn1)
				connectorReg.RegisterSpotConnector("coinbase", spotConn2)
				connectorReg.RegisterPerpConnector("hyperliquid", perpConn1)
				connectorReg.RegisterPerpConnector("dydx", perpConn2)

				spotConnectors := connectorReg.GetSpotConnectors()
				perpConnectors := connectorReg.GetPerpConnectors()

				Expect(spotConnectors).To(HaveLen(2))
				Expect(perpConnectors).To(HaveLen(2))
			})

			It("should return all connectors combined", func() {
				connectorReg.RegisterSpotConnector("binance", spotConn1)
				connectorReg.RegisterSpotConnector("coinbase", spotConn2)
				connectorReg.RegisterPerpConnector("hyperliquid", perpConn1)
				connectorReg.RegisterPerpConnector("dydx", perpConn2)

				allConnectors := connectorReg.GetAllBaseConnectors()
				Expect(allConnectors).To(HaveLen(4))
			})

			It("should only return ready connectors when filtered", func() {
				connectorReg.RegisterSpotConnector("binance", spotConn1)
				connectorReg.RegisterSpotConnector("coinbase", spotConn2)
				connectorReg.RegisterPerpConnector("hyperliquid", perpConn1)
				connectorReg.RegisterPerpConnector("dydx", perpConn2)

				_ = connectorReg.MarkConnectorReady("binance")
				_ = connectorReg.MarkConnectorReady("hyperliquid")

				allReady := connectorReg.GetAllReadyConnectors()
				Expect(allReady).To(HaveLen(2))
			})
		})
	})

	Describe("Non-Existent Connector Handling", func() {
		Context("when checking non-existent connector", func() {
			It("should return false for GetConnector", func() {
				_, exists := connectorReg.GetConnector("nonexistent")
				Expect(exists).To(BeFalse())
			})

			It("should return false for GetSpotConnector", func() {
				_, exists := connectorReg.GetSpotConnector("nonexistent")
				Expect(exists).To(BeFalse())
			})

			It("should return false for GetPerpConnector", func() {
				_, exists := connectorReg.GetPerpConnector("nonexistent")
				Expect(exists).To(BeFalse())
			})

			It("should return false for IsConnectorReady", func() {
				ready := connectorReg.IsConnectorReady("nonexistent")
				Expect(ready).To(BeFalse())
			})

			It("should return error when marking as ready", func() {
				err := connectorReg.MarkConnectorReady("nonexistent")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not found"))
			})
		})
	})

	Describe("Empty Registry Behavior", func() {
		Context("when registry is empty", func() {
			It("should return empty slices for all connector queries", func() {
				Expect(connectorReg.GetSpotConnectors()).To(BeEmpty())
				Expect(connectorReg.GetPerpConnectors()).To(BeEmpty())
				Expect(connectorReg.GetAllBaseConnectors()).To(BeEmpty())
				Expect(connectorReg.GetAllReadyConnectors()).To(BeEmpty())
				Expect(connectorReg.GetReadySpotConnectors()).To(BeEmpty())
				Expect(connectorReg.GetReadyPerpConnectors()).To(BeEmpty())
				Expect(connectorReg.GetSpotWebSocketConnectors()).To(BeEmpty())
				Expect(connectorReg.GetPerpWebSocketConnectors()).To(BeEmpty())
				Expect(connectorReg.GetReadySpotWebSocketConnectors()).To(BeEmpty())
				Expect(connectorReg.GetReadyPerpWebSocketConnectors()).To(BeEmpty())
			})
		})
	})

})
