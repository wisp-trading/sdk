package analytics_test

import (
	"testing"
	"time"

	"github.com/wisp-trading/sdk/pkg/markets/options/activity"
	"github.com/wisp-trading/sdk/pkg/markets/options/analytics"
	"github.com/wisp-trading/sdk/pkg/markets/options/store"
	"github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	timeProvider "github.com/wisp-trading/sdk/pkg/runtime/time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAnalytics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Options Analytics Suite")
}

var _ = Describe("Options Analytics Service", func() {
	var (
		analyticsService types.OptionsAnalytics
		pnl              types.OptionsPNL
		optStore         types.OptionsStore
		logger           logging.ApplicationLogger
		quote            = portfolio.NewAsset("USDT")
		btcPair          = portfolio.NewPair(portfolio.NewAsset("BTC"), quote)
		expiration       = time.Now().AddDate(0, 0, 30)
	)

	BeforeEach(func() {
		logger = logging.NewNoOpLogger()
		optStore = store.NewStore(timeProvider.NewTimeProvider())
		pnl = activity.NewPNLCalculator(optStore, logger)
		analyticsService = analytics.NewAnalyticsService(pnl, optStore, logger)
	})

	Describe("GetDeltaExposure", func() {
		It("should return aggregated delta exposure", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			optStore.SetPosition(contract, types.Position{
				Contract:   contract,
				Quantity:   5.0,
				EntryPrice: 1000.0,
			})

			greeks := types.Greeks{
				Delta: 0.6,
			}
			optStore.SetGreeks(contract, greeks)

			// Expected delta = 5 * 0.6 = 3.0
			deltaExposure := analyticsService.GetDeltaExposure()
			Expect(deltaExposure).To(Equal(3.0))
		})
	})

	Describe("GetGammaExposure", func() {
		It("should return aggregated gamma exposure", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			optStore.SetPosition(contract, types.Position{
				Contract:   contract,
				Quantity:   5.0,
				EntryPrice: 1000.0,
			})

			greeks := types.Greeks{
				Gamma: 0.02,
			}
			optStore.SetGreeks(contract, greeks)

			// Expected gamma = 5 * 0.02 = 0.1
			gammaExposure := analyticsService.GetGammaExposure()
			Expect(gammaExposure).To(Equal(0.1))
		})
	})

	Describe("GetThetaExposure", func() {
		It("should return aggregated theta exposure", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			optStore.SetPosition(contract, types.Position{
				Contract:   contract,
				Quantity:   5.0,
				EntryPrice: 1000.0,
			})

			greeks := types.Greeks{
				Theta: -0.1,
			}
			optStore.SetGreeks(contract, greeks)

			// Expected theta = 5 * -0.1 = -0.5
			thetaExposure := analyticsService.GetThetaExposure()
			Expect(thetaExposure).To(Equal(-0.5))
		})
	})

	Describe("GetVegaExposure", func() {
		It("should return aggregated vega exposure", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			optStore.SetPosition(contract, types.Position{
				Contract:   contract,
				Quantity:   5.0,
				EntryPrice: 1000.0,
			})

			greeks := types.Greeks{
				Vega: 10.0,
			}
			optStore.SetGreeks(contract, greeks)

			// Expected vega = 5 * 10.0 = 50.0
			vegaExposure := analyticsService.GetVegaExposure()
			Expect(vegaExposure).To(Equal(50.0))
		})
	})

	Describe("GetPortfolioGreeks", func() {
		It("should return aggregated Greeks", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			optStore.SetPosition(contract, types.Position{
				Contract:   contract,
				Quantity:   2.0,
				EntryPrice: 1000.0,
			})

			greeks := types.Greeks{
				Delta: 0.6,
				Gamma: 0.01,
				Theta: -0.05,
				Vega:  5.0,
				Rho:   0.15,
			}
			optStore.SetGreeks(contract, greeks)

			portfolioGreeks := analyticsService.GetPortfolioGreeks()
			Expect(portfolioGreeks.Delta).To(BeNumerically("~", 1.2, 0.0001))
			Expect(portfolioGreeks.Gamma).To(BeNumerically("~", 0.02, 0.0001))
			Expect(portfolioGreeks.Theta).To(BeNumerically("~", -0.1, 0.0001))
			Expect(portfolioGreeks.Vega).To(BeNumerically("~", 10.0, 0.0001))
			Expect(portfolioGreeks.Rho).To(BeNumerically("~", 0.3, 0.0001))
		})
	})

	Describe("GetDailyThetaDecay", func() {
		It("should return daily theta decay", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			optStore.SetPosition(contract, types.Position{
				Contract:   contract,
				Quantity:   5.0,
				EntryPrice: 1000.0,
			})

			greeks := types.Greeks{
				Theta: -0.1,
			}
			optStore.SetGreeks(contract, greeks)

			// Expected theta = 5 * -0.1 = -0.5
			thetaDecay := analyticsService.GetDailyThetaDecay()
			Expect(thetaDecay).To(Equal(-0.5))
		})
	})

	Describe("GetIVSensitivity", func() {
		It("should return IV sensitivity (vega exposure)", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			optStore.SetPosition(contract, types.Position{
				Contract:   contract,
				Quantity:   5.0,
				EntryPrice: 1000.0,
			})

			greeks := types.Greeks{
				Vega: 10.0,
			}
			optStore.SetGreeks(contract, greeks)

			// Expected vega = 5 * 10.0 = 50.0
			ivSensitivity := analyticsService.GetIVSensitivity()
			Expect(ivSensitivity).To(Equal(50.0))
		})
	})
})
