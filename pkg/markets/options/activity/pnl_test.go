package activity_test

import (
	"testing"
	"time"

	"github.com/wisp-trading/sdk/pkg/markets/options/activity"
	"github.com/wisp-trading/sdk/pkg/markets/options/store"
	"github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	timeProvider "github.com/wisp-trading/sdk/pkg/runtime/time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPNL(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Options PnL Suite")
}

var _ = Describe("Options PnL Calculator", func() {
	var (
		pnl        types.OptionsPNL
		optStore   types.OptionsStore
		logger     logging.ApplicationLogger
		quote      = portfolio.NewAsset("USDT")
		btcPair    = portfolio.NewPair(portfolio.NewAsset("BTC"), quote)
		expiration = time.Now().AddDate(0, 0, 30)
	)

	BeforeEach(func() {
		logger = logging.NewNoOpLogger()
		optStore = store.NewStore(timeProvider.NewTimeProvider())
		pnl = activity.NewPNLCalculator(optStore, logger)
	})

	Describe("CalculateUnrealizedPnL", func() {
		It("should return zero for unknown contract", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			pnlValue := pnl.CalculateUnrealizedPnL(contract)
			Expect(pnlValue).To(Equal(0.0))
		})

		It("should calculate P&L from entry price to mark price", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			// Add position: 10 contracts at entry price 1000
			optStore.SetPosition(contract, types.Position{
				Contract:   contract,
				Quantity:   10.0,
				EntryPrice: 1000.0,
			})

			// Mark price is now 1200
			optStore.SetMarkPrice(contract, 1200.0)

			// Expected P&L = 10 * (1200 - 1000) = 2000
			pnlValue := pnl.CalculateUnrealizedPnL(contract)
			Expect(pnlValue).To(Equal(2000.0))
		})
	})

	Describe("CalculateDeltaExposure", func() {
		It("should return zero for no positions", func() {
			deltaExposure := pnl.CalculateDeltaExposure()
			Expect(deltaExposure).To(Equal(0.0))
		})

		It("should aggregate delta across positions", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			optStore.SetPosition(contract, types.Position{
				Contract:   contract,
				Quantity:   10.0,
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

			// Expected delta = 10 * 0.6 = 6.0
			deltaExposure := pnl.CalculateDeltaExposure()
			Expect(deltaExposure).To(Equal(6.0))
		})
	})

	Describe("CalculateGammaExposure", func() {
		It("should aggregate gamma across positions", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			optStore.SetPosition(contract, types.Position{
				Contract:   contract,
				Quantity:   10.0,
				EntryPrice: 1000.0,
			})

			greeks := types.Greeks{
				Delta: 0.6,
				Gamma: 0.01,
			}
			optStore.SetGreeks(contract, greeks)

			// Expected gamma = 10 * 0.01 = 0.1
			gammaExposure := pnl.CalculateGammaExposure()
			Expect(gammaExposure).To(Equal(0.1))
		})
	})

	Describe("CalculateThetaDecay", func() {
		It("should aggregate theta across positions", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			optStore.SetPosition(contract, types.Position{
				Contract:   contract,
				Quantity:   10.0,
				EntryPrice: 1000.0,
			})

			greeks := types.Greeks{
				Theta: -0.05,
			}
			optStore.SetGreeks(contract, greeks)

			// Expected theta = 10 * -0.05 = -0.5
			thetaDecay := pnl.CalculateThetaDecay()
			Expect(thetaDecay).To(Equal(-0.5))
		})
	})

	Describe("CalculateVegaExposure", func() {
		It("should aggregate vega across positions", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			optStore.SetPosition(contract, types.Position{
				Contract:   contract,
				Quantity:   10.0,
				EntryPrice: 1000.0,
			})

			greeks := types.Greeks{
				Vega: 5.0,
			}
			optStore.SetGreeks(contract, greeks)

			// Expected vega = 10 * 5.0 = 50.0
			vegaExposure := pnl.CalculateVegaExposure()
			Expect(vegaExposure).To(Equal(50.0))
		})
	})

	Describe("GetPortfolioGreeks", func() {
		It("should aggregate Greeks across positions", func() {
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

			portfolioGreeks := pnl.GetPortfolioGreeks()
			Expect(portfolioGreeks.Delta).To(BeNumerically("~", 1.2, 0.0001))
			Expect(portfolioGreeks.Gamma).To(BeNumerically("~", 0.02, 0.0001))
			Expect(portfolioGreeks.Theta).To(BeNumerically("~", -0.1, 0.0001))
		})
	})
})
