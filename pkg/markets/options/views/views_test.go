package views_test

import (
	"testing"
	"time"

	"github.com/wisp-trading/sdk/pkg/markets/options/store"
	"github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/markets/options/views"
	timeProvider "github.com/wisp-trading/sdk/pkg/runtime/time"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestViews(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Options Views Suite")
}

var _ = Describe("Options View", func() {
	var (
		view       types.OptionsView
		optStore   types.OptionsStore
		logger     logging.ApplicationLogger
		quote      = portfolio.NewAsset("USDT")
		btcPair    = portfolio.NewPair(portfolio.NewAsset("BTC"), quote)
		expiration = time.Now().AddDate(0, 0, 30)
	)

	BeforeEach(func() {
		logger = logging.NewNoOpLogger()
		optStore = store.NewStore(timeProvider.NewTimeProvider())
		view = views.NewView(optStore, logger)
	})

	Describe("GetMarkPrice", func() {
		It("should return zero for unknown contract", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			price := view.GetMarkPrice(contract)
			Expect(price).To(Equal(0.0))
		})

		It("should return stored mark price", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}
			expectedPrice := 2000.0

			optStore.SetMarkPrice(contract, expectedPrice)
			price := view.GetMarkPrice(contract)
			Expect(price).To(Equal(expectedPrice))
		})
	})

	Describe("GetUnderlyingPrice", func() {
		It("should return stored underlying price", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}
			expectedPrice := 50000.0

			optStore.SetUnderlyingPrice(contract, expectedPrice)
			price := view.GetUnderlyingPrice(contract)
			Expect(price).To(Equal(expectedPrice))
		})
	})

	Describe("GetGreeks", func() {
		It("should return zero Greeks for unknown contract", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}

			greeks := view.GetGreeks(contract)
			Expect(greeks.Delta).To(Equal(0.0))
			Expect(greeks.Gamma).To(Equal(0.0))
			Expect(greeks.Theta).To(Equal(0.0))
		})

		It("should return stored Greeks", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}
			expectedGreeks := types.Greeks{
				Delta: 0.6,
				Gamma: 0.01,
				Theta: -0.05,
				Vega:  5.0,
				Rho:   0.15,
			}

			optStore.SetGreeks(contract, expectedGreeks)
			greeks := view.GetGreeks(contract)
			Expect(greeks).To(Equal(expectedGreeks))
		})
	})

	Describe("GetIV", func() {
		It("should return stored IV", func() {
			contract := types.OptionContract{
				Pair:       btcPair,
				Strike:     50000,
				Expiration: expiration,
				OptionType: "CALL",
			}
			expectedIV := 0.25

			optStore.SetIV(contract, expectedIV)
			iv := view.GetIV(contract)
			Expect(iv).To(Equal(expectedIV))
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
				Quantity:   3.0,
				EntryPrice: 1000.0,
			})

			greeks := types.Greeks{
				Delta: 0.6,
				Gamma: 0.01,
			}

			optStore.SetGreeks(contract, greeks)
			portfolioGreeks := view.GetPortfolioGreeks()
			Expect(portfolioGreeks.Delta).To(BeNumerically("~", 1.8, 0.0001))
			Expect(portfolioGreeks.Gamma).To(BeNumerically("~", 0.03, 0.0001))
		})
	})
})
