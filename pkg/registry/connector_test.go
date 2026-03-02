package registry_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	perpMock "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/connector/perp"
	spotMock "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/connector/spot"
	sdkTesting "github.com/wisp-trading/sdk/pkg/testing"
	registryTypes "github.com/wisp-trading/sdk/pkg/types/registry"
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
				connectorReg.RegisterSpot("binance", spotConn1)

				conn, exists := connectorReg.Spot("binance")
				Expect(exists).To(BeTrue())
				Expect(conn).To(Equal(spotConn1))
			})

			It("should not be marked as ready initially", func() {
				connectorReg.RegisterSpot("binance", spotConn1)

				Expect(connectorReg.IsReady("binance")).To(BeFalse())
			})

			It("should be retrievable via generic GetConnector", func() {
				connectorReg.RegisterSpot("binance", spotConn1)

				conn, exists := connectorReg.Connector("binance")
				Expect(exists).To(BeTrue())
				Expect(conn).NotTo(BeNil())
			})
		})

		Context("when registering duplicate spot connector - BUG TEST", func() {
			It("should overwrite the previous connector", func() {
				connectorReg.RegisterSpot("binance", spotConn1)
				connectorReg.RegisterSpot("binance", spotConn2)

				conn, exists := connectorReg.Spot("binance")
				Expect(exists).To(BeTrue())
				Expect(conn).To(Equal(spotConn2), "should have the second connector, not the first")
			})

			It("should reset ready state on re-registration - THE BUG", func() {
				connectorReg.RegisterSpot("binance", spotConn1)
				err := connectorReg.MarkReady("binance")
				Expect(err).NotTo(HaveOccurred())
				Expect(connectorReg.IsReady("binance")).To(BeTrue())

				// Re-register same exchange (this is the bug scenario)
				connectorReg.RegisterSpot("binance", spotConn2)

				// BUG: Ready state is reset to false without warning
				Expect(connectorReg.IsReady("binance")).To(BeFalse(),
					"BUG: ready state is reset when re-registering - should either error or preserve state")
			})
		})

		Context("when marking connector as ready", func() {
			It("should update ready state", func() {
				connectorReg.RegisterSpot("binance", spotConn1)

				err := connectorReg.MarkReady("binance")
				Expect(err).NotTo(HaveOccurred())
				Expect(connectorReg.IsReady("binance")).To(BeTrue())
			})

			It("should include connector in ready list", func() {
				connectorReg.RegisterSpot("binance", spotConn1)
				connectorReg.RegisterSpot("coinbase", spotConn2)
				_ = connectorReg.MarkReady("binance")

				readySpot := connectorReg.FilterSpot(registryTypes.NewFilter().ReadyOnly().Build())
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
				connectorReg.RegisterPerp("hyperliquid", perpConn1)

				conn, exists := connectorReg.Perp("hyperliquid")
				Expect(exists).To(BeTrue())
				Expect(conn).To(Equal(perpConn1))
			})

			It("should not be marked as ready initially", func() {
				connectorReg.RegisterPerp("hyperliquid", perpConn1)

				Expect(connectorReg.IsReady("hyperliquid")).To(BeFalse())
			})

			It("should be retrievable via generic GetConnector", func() {
				connectorReg.RegisterPerp("hyperliquid", perpConn1)

				conn, exists := connectorReg.Connector("hyperliquid")
				Expect(exists).To(BeTrue())
				Expect(conn).NotTo(BeNil())
			})
		})

		Context("when registering duplicate perp connector - THE BUG", func() {
			It("should overwrite the previous connector", func() {
				connectorReg.RegisterPerp("hyperliquid", perpConn1)
				connectorReg.RegisterPerp("hyperliquid", perpConn2)

				conn, exists := connectorReg.Perp("hyperliquid")
				Expect(exists).To(BeTrue())
				Expect(conn).To(Equal(perpConn2), "BUG: should have the second connector")
			})

			It("should reset ready state on re-registration - THE BUG", func() {
				connectorReg.RegisterPerp("hyperliquid", perpConn1)
				err := connectorReg.MarkReady("hyperliquid")
				Expect(err).NotTo(HaveOccurred())
				Expect(connectorReg.IsReady("hyperliquid")).To(BeTrue())

				// Re-register same exchange (this is the production bug with hyperliquid)
				connectorReg.RegisterPerp("hyperliquid", perpConn2)

				// BUG: Ready state is reset to false, causing the connector to appear unready
				Expect(connectorReg.IsReady("hyperliquid")).To(BeFalse(),
					"BUG: ready state is reset when re-registering hyperliquid - this is the production bug!")
			})
		})

		Context("when marking connector as ready", func() {
			It("should update ready state", func() {
				connectorReg.RegisterPerp("hyperliquid", perpConn1)

				err := connectorReg.MarkReady("hyperliquid")
				Expect(err).NotTo(HaveOccurred())
				Expect(connectorReg.IsReady("hyperliquid")).To(BeTrue())
			})

			It("should include connector in ready list", func() {
				connectorReg.RegisterPerp("hyperliquid", perpConn1)
				connectorReg.RegisterPerp("dydx", perpConn2)
				_ = connectorReg.MarkReady("hyperliquid")

				readyPerp := connectorReg.FilterPerp(registryTypes.NewFilter().ReadyOnly().Build())
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
				connectorReg.RegisterSpot("binance", spotConn1)
				connectorReg.RegisterSpot("coinbase", spotConn2)
				connectorReg.RegisterPerp("hyperliquid", perpConn1)
				connectorReg.RegisterPerp("dydx", perpConn2)

				spotConnectors := connectorReg.FilterSpot(registryTypes.NewFilter().Build())
				perpConnectors := connectorReg.FilterPerp(registryTypes.NewFilter().Build())

				Expect(spotConnectors).To(HaveLen(2))
				Expect(perpConnectors).To(HaveLen(2))
			})

			It("should only return ready connectors when filtered", func() {
				connectorReg.RegisterSpot("binance", spotConn1)
				connectorReg.RegisterSpot("coinbase", spotConn2)
				connectorReg.RegisterPerp("hyperliquid", perpConn1)
				connectorReg.RegisterPerp("dydx", perpConn2)

				_ = connectorReg.MarkReady("binance")
				_ = connectorReg.MarkReady("hyperliquid")

				allReady := connectorReg.Filter(registryTypes.NewFilter().ReadyOnly().Build())
				Expect(allReady).To(HaveLen(2))
			})
		})
	})

	Describe("Non-Existent Connector Handling", func() {
		Context("when checking non-existent connector", func() {
			It("should return false for GetConnector", func() {
				_, exists := connectorReg.Connector("nonexistent")
				Expect(exists).To(BeFalse())
			})

			It("should return false for Spot", func() {
				_, exists := connectorReg.Spot("nonexistent")
				Expect(exists).To(BeFalse())
			})

			It("should return false for Perp", func() {
				_, exists := connectorReg.Perp("nonexistent")
				Expect(exists).To(BeFalse())
			})

			It("should return false for IsReady", func() {
				ready := connectorReg.IsReady("nonexistent")
				Expect(ready).To(BeFalse())
			})

			It("should return error when marking as ready", func() {
				err := connectorReg.MarkReady("nonexistent")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not found"))
			})
		})
	})

	Describe("Empty Registry Behavior", func() {
		Context("when registry is empty", func() {
			It("should return empty slices for all connector queries", func() {
				Expect(connectorReg.FilterSpot(registryTypes.NewFilter().Build())).To(BeEmpty())
				Expect(connectorReg.FilterPerp(registryTypes.NewFilter().Build())).To(BeEmpty())
				Expect(connectorReg.Filter(registryTypes.NewFilter().Build())).To(BeEmpty())
			})
		})
	})

})
