package store_test

import (
	"testing"
	"time"

	"github.com/wisp-trading/sdk/pkg/markets/options/store"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	timeProvider "github.com/wisp-trading/sdk/pkg/runtime/time"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Options Store Suite")
}

var _ = Describe("Options Store", func() {
	var (
		optStore   optionsTypes.OptionsStore
		btcPair    portfolio.Pair
		ethPair    portfolio.Pair
		expiration time.Time
		contract1  optionsTypes.OptionContract
		contract2  optionsTypes.OptionContract
	)

	BeforeEach(func() {
		timeProviderInst := timeProvider.NewTimeProvider()
		optStore = store.NewStore(timeProviderInst)
		btcPair = portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
		ethPair = portfolio.NewPair(portfolio.NewAsset("ETH"), portfolio.NewAsset("USDT"))
		expiration = time.Now().AddDate(0, 0, 30)

		contract1 = optionsTypes.OptionContract{
			Pair:       btcPair,
			Strike:     50000,
			Expiration: expiration,
			OptionType: "CALL",
		}
		contract2 = optionsTypes.OptionContract{
			Pair:       ethPair,
			Strike:     2000,
			Expiration: expiration,
			OptionType: "PUT",
		}
	})

	Describe("Position Management", func() {
		It("should store and retrieve positions", func() {
			position := optionsTypes.Position{
				Contract:   contract1,
				Quantity:   5.0,
				EntryPrice: 1000.0,
			}

			optStore.SetPosition(contract1, position)
			retrieved := optStore.GetPosition(contract1)

			Expect(retrieved).NotTo(BeNil())
			Expect(retrieved.Quantity).To(Equal(5.0))
			Expect(retrieved.EntryPrice).To(Equal(1000.0))
		})

		It("should return nil for non-existent position", func() {
			retrieved := optStore.GetPosition(contract1)
			Expect(retrieved).To(BeNil())
		})

		It("should update position when set again", func() {
			pos1 := optionsTypes.Position{Contract: contract1, Quantity: 5.0, EntryPrice: 1000.0}
			pos2 := optionsTypes.Position{Contract: contract1, Quantity: 10.0, EntryPrice: 2000.0}

			optStore.SetPosition(contract1, pos1)
			optStore.SetPosition(contract1, pos2)
			retrieved := optStore.GetPosition(contract1)

			Expect(retrieved.Quantity).To(Equal(10.0))
			Expect(retrieved.EntryPrice).To(Equal(2000.0))
		})

		It("should retrieve all positions", func() {
			pos1 := optionsTypes.Position{Contract: contract1, Quantity: 5.0, EntryPrice: 1000.0}
			pos2 := optionsTypes.Position{Contract: contract2, Quantity: 10.0, EntryPrice: 2000.0}

			optStore.SetPosition(contract1, pos1)
			optStore.SetPosition(contract2, pos2)
			allPositions := optStore.GetAllPositions()

			Expect(allPositions).To(HaveLen(2))
		})
	})

	Describe("Mark Price Management", func() {
		It("should store and retrieve mark prices", func() {
			optStore.SetMarkPrice(contract1, 1500.0)
			price := optStore.GetMarkPrice(contract1)
			Expect(price).To(Equal(1500.0))
		})

		It("should return 0 for non-existent mark price", func() {
			price := optStore.GetMarkPrice(contract1)
			Expect(price).To(Equal(0.0))
		})

		It("should update mark price", func() {
			optStore.SetMarkPrice(contract1, 1500.0)
			optStore.SetMarkPrice(contract1, 1600.0)
			price := optStore.GetMarkPrice(contract1)
			Expect(price).To(Equal(1600.0))
		})
	})

	Describe("Greeks Management", func() {
		It("should store and retrieve Greeks", func() {
			greeks := optionsTypes.Greeks{
				Delta: 0.5,
				Gamma: 0.01,
				Theta: -0.05,
				Vega:  10.0,
				Rho:   0.2,
			}
			optStore.SetGreeks(contract1, greeks)
			retrieved := optStore.GetGreeks(contract1)

			Expect(retrieved.Delta).To(Equal(0.5))
			Expect(retrieved.Gamma).To(Equal(0.01))
			Expect(retrieved.Theta).To(Equal(-0.05))
			Expect(retrieved.Vega).To(Equal(10.0))
			Expect(retrieved.Rho).To(Equal(0.2))
		})

		It("should return zero Greeks for non-existent contract", func() {
			retrieved := optStore.GetGreeks(contract1)
			Expect(retrieved.Delta).To(Equal(0.0))
		})
	})

	Describe("Underlying Price Management", func() {
		It("should store and retrieve underlying price", func() {
			optStore.SetUnderlyingPrice(contract1, 50000.0)
			price := optStore.GetUnderlyingPrice(contract1)
			Expect(price).To(Equal(50000.0))
		})

		It("should return 0 for non-existent underlying price", func() {
			price := optStore.GetUnderlyingPrice(contract1)
			Expect(price).To(Equal(0.0))
		})
	})

	Describe("IV Management", func() {
		It("should store and retrieve IV", func() {
			optStore.SetIV(contract1, 0.25)
			iv := optStore.GetIV(contract1)
			Expect(iv).To(Equal(0.25))
		})

		It("should return 0 for non-existent IV", func() {
			iv := optStore.GetIV(contract1)
			Expect(iv).To(Equal(0.0))
		})
	})

	Describe("Portfolio Greeks Aggregation", func() {
		It("should aggregate Greeks across positions", func() {
			pos1 := optionsTypes.Position{Contract: contract1, Quantity: 2.0, EntryPrice: 1000.0}
			pos2 := optionsTypes.Position{Contract: contract2, Quantity: 3.0, EntryPrice: 2000.0}

			optStore.SetPosition(contract1, pos1)
			optStore.SetPosition(contract2, pos2)

			greeks1 := optionsTypes.Greeks{Delta: 0.5, Gamma: 0.01}
			greeks2 := optionsTypes.Greeks{Delta: 0.6, Gamma: 0.02}

			optStore.SetGreeks(contract1, greeks1)
			optStore.SetGreeks(contract2, greeks2)

			portfolioGreeks := optStore.GetPortfolioGreeks()
			// Contract1: 0.5 * 2.0 = 1.0, Contract2: 0.6 * 3.0 = 1.8
			Expect(portfolioGreeks.Delta).To(BeNumerically("~", 2.8, 0.0001))
			// Contract1: 0.01 * 2.0 = 0.02, Contract2: 0.02 * 3.0 = 0.06
			Expect(portfolioGreeks.Gamma).To(BeNumerically("~", 0.08, 0.0001))
		})

		It("should return zero Greeks when no positions exist", func() {
			portfolioGreeks := optStore.GetPortfolioGreeks()
			Expect(portfolioGreeks.Delta).To(Equal(0.0))
			Expect(portfolioGreeks.Gamma).To(Equal(0.0))
		})

		It("should only aggregate Greeks for positions that have them", func() {
			pos1 := optionsTypes.Position{Contract: contract1, Quantity: 2.0, EntryPrice: 1000.0}
			optStore.SetPosition(contract1, pos1)
			optStore.SetGreeks(contract1, optionsTypes.Greeks{Delta: 0.5})

			// Set position for contract2 but no Greeks
			pos2 := optionsTypes.Position{Contract: contract2, Quantity: 3.0, EntryPrice: 2000.0}
			optStore.SetPosition(contract2, pos2)

			portfolioGreeks := optStore.GetPortfolioGreeks()
			// Only contract1 contributes: 0.5 * 2.0 = 1.0
			Expect(portfolioGreeks.Delta).To(BeNumerically("~", 1.0, 0.0001))
		})
	})

	Describe("Market Type", func() {
		It("should return correct market type", func() {
			marketType := optStore.MarketType()
			Expect(marketType).To(Equal(connector.MarketTypeOptions))
		})
	})

	Describe("Concurrent Access", func() {
		It("should handle concurrent position updates", func() {
			done := make(chan bool)

			go func() {
				for i := 0; i < 100; i++ {
					pos := optionsTypes.Position{
						Contract:   contract1,
						Quantity:   float64(i),
						EntryPrice: float64(i * 100),
					}
					optStore.SetPosition(contract1, pos)
				}
				done <- true
			}()

			go func() {
				for i := 0; i < 100; i++ {
					_ = optStore.GetPosition(contract1)
				}
				done <- true
			}()

			<-done
			<-done

			retrieved := optStore.GetPosition(contract1)
			Expect(retrieved).NotTo(BeNil())
			Expect(retrieved.Quantity).To(BeNumerically(">=", 0))
		})
	})
})

// Benchmarks
func BenchmarkSetPosition(b *testing.B) {
	timeProviderInst := timeProvider.NewTimeProvider()
	optStore := store.NewStore(timeProviderInst)
	btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
	expiration := time.Now().AddDate(0, 0, 30)

	contract := optionsTypes.OptionContract{
		Pair:       btcPair,
		Strike:     50000,
		Expiration: expiration,
		OptionType: "CALL",
	}

	position := optionsTypes.Position{
		Contract:   contract,
		Quantity:   5.0,
		EntryPrice: 1000.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optStore.SetPosition(contract, position)
	}
}

func BenchmarkGetPosition(b *testing.B) {
	timeProviderInst := timeProvider.NewTimeProvider()
	optStore := store.NewStore(timeProviderInst)
	btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
	expiration := time.Now().AddDate(0, 0, 30)

	contract := optionsTypes.OptionContract{
		Pair:       btcPair,
		Strike:     50000,
		Expiration: expiration,
		OptionType: "CALL",
	}

	position := optionsTypes.Position{
		Contract:   contract,
		Quantity:   5.0,
		EntryPrice: 1000.0,
	}

	optStore.SetPosition(contract, position)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = optStore.GetPosition(contract)
	}
}

func BenchmarkGetPortfolioGreeks(b *testing.B) {
	timeProviderInst := timeProvider.NewTimeProvider()
	optStore := store.NewStore(timeProviderInst)
	btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
	expiration := time.Now().AddDate(0, 0, 30)

	// Set up 100 positions with Greeks
	for i := 0; i < 100; i++ {
		contract := optionsTypes.OptionContract{
			Pair:       btcPair,
			Strike:     float64(45000 + i*100),
			Expiration: expiration,
			OptionType: "CALL",
		}

		position := optionsTypes.Position{
			Contract:   contract,
			Quantity:   1.0,
			EntryPrice: float64(1000 + i*10),
		}

		greeks := optionsTypes.Greeks{
			Delta: 0.5,
			Gamma: 0.01,
			Theta: -0.05,
			Vega:  5.0,
			Rho:   0.15,
		}

		optStore.SetPosition(contract, position)
		optStore.SetGreeks(contract, greeks)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = optStore.GetPortfolioGreeks()
	}
}

func BenchmarkSetMarkPrice(b *testing.B) {
	timeProviderInst := timeProvider.NewTimeProvider()
	optStore := store.NewStore(timeProviderInst)
	btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
	expiration := time.Now().AddDate(0, 0, 30)

	contract := optionsTypes.OptionContract{
		Pair:       btcPair,
		Strike:     50000,
		Expiration: expiration,
		OptionType: "CALL",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optStore.SetMarkPrice(contract, 1500.0)
	}
}

func BenchmarkSetGreeks(b *testing.B) {
	timeProviderInst := timeProvider.NewTimeProvider()
	optStore := store.NewStore(timeProviderInst)
	btcPair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
	expiration := time.Now().AddDate(0, 0, 30)

	contract := optionsTypes.OptionContract{
		Pair:       btcPair,
		Strike:     50000,
		Expiration: expiration,
		OptionType: "CALL",
	}

	greeks := optionsTypes.Greeks{
		Delta: 0.5,
		Gamma: 0.01,
		Theta: -0.05,
		Vega:  5.0,
		Rho:   0.15,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optStore.SetGreeks(contract, greeks)
	}
}
