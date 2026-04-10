package options_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	optionsService "github.com/wisp-trading/sdk/pkg/markets/options/options"
	"github.com/wisp-trading/sdk/pkg/markets/options/activity"
	optionsStore "github.com/wisp-trading/sdk/pkg/markets/options/store"
	optionsWatchlistPkg "github.com/wisp-trading/sdk/pkg/markets/options"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	timeProvider "github.com/wisp-trading/sdk/pkg/runtime/time"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	mockLogging "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

var _ = Describe("Options Service", func() {
	var (
		svc        optionsTypes.Options
		store      optionsTypes.OptionsStore
		watchlist  optionsTypes.OptionsWatchlist
		btcPair    portfolio.Pair
		expiration time.Time
		contract   optionsTypes.OptionContract
		deribit    connector.ExchangeName
	)

	BeforeEach(func() {
		tp := timeProvider.NewTimeProvider()
		store = optionsStore.NewStore(tp)
		watchlist = optionsWatchlistPkg.NewOptionsWatchlist()
		logger := &mockLogging.TradingLogger{}
		pnl := activity.NewPNLCalculator(store, nil)
		svc = optionsService.NewOptions(logger, watchlist, store, tp, pnl)

		deribit = "deribit_options"
		btcPair = portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
		expiration = time.Now().UTC().Add(30 * 24 * time.Hour)
		contract = optionsTypes.OptionContract{
			Pair:       btcPair,
			Strike:     50000,
			Expiration: expiration,
			OptionType: "CALL",
		}
	})

	Describe("MarkPrice", func() {
		It("returns false when no price is stored", func() {
			_, found := svc.MarkPrice(deribit, contract)
			Expect(found).To(BeFalse())
		})

		It("returns the stored price", func() {
			store.SetMarkPrice(contract, 1250.50)
			price, found := svc.MarkPrice(deribit, contract)
			Expect(found).To(BeTrue())
			Expect(price.String()).To(Equal("1250.5"))
		})
	})

	Describe("UnderlyingPrice", func() {
		It("returns false when no price is stored", func() {
			_, found := svc.UnderlyingPrice(deribit, contract)
			Expect(found).To(BeFalse())
		})

		It("returns the stored underlying price", func() {
			store.SetUnderlyingPrice(contract, 83000.0)
			price, found := svc.UnderlyingPrice(deribit, contract)
			Expect(found).To(BeTrue())
			Expect(price.String()).To(Equal("83000"))
		})
	})

	Describe("Greeks", func() {
		It("returns false when no Greeks are stored", func() {
			_, found := svc.Greeks(deribit, contract)
			Expect(found).To(BeFalse())
		})

		It("returns stored Greeks", func() {
			greeks := optionsTypes.Greeks{Delta: 0.6, Gamma: 0.001, Theta: -50.0, Vega: 100.0, Rho: 5.0}
			store.SetGreeks(contract, greeks)
			result, found := svc.Greeks(deribit, contract)
			Expect(found).To(BeTrue())
			Expect(result.Delta).To(Equal(0.6))
			Expect(result.Gamma).To(Equal(0.001))
		})
	})

	Describe("ImpliedVolatility", func() {
		It("returns false when no IV is stored", func() {
			_, found := svc.ImpliedVolatility(deribit, contract)
			Expect(found).To(BeFalse())
		})

		It("returns the stored IV", func() {
			store.SetIV(contract, 0.75)
			iv, found := svc.ImpliedVolatility(deribit, contract)
			Expect(found).To(BeTrue())
			Expect(iv).To(Equal(0.75))
		})
	})

	Describe("WatchContract / UnwatchContract / Expirations", func() {
		It("Expirations returns false before any contract is watched", func() {
			_, found := svc.Expirations(deribit, btcPair)
			Expect(found).To(BeFalse())
		})

		It("Expirations returns the expiration after WatchContract", func() {
			svc.WatchContract(deribit, contract)
			exps, found := svc.Expirations(deribit, btcPair)
			Expect(found).To(BeTrue())
			Expect(exps).To(HaveLen(1))
			Expect(exps[0].Equal(expiration)).To(BeTrue())
		})

		It("Expirations returns false after UnwatchContract", func() {
			svc.WatchContract(deribit, contract)
			svc.UnwatchContract(deribit, contract)
			_, found := svc.Expirations(deribit, btcPair)
			Expect(found).To(BeFalse())
		})

		It("does not mix expirations across pairs", func() {
			ethPair := portfolio.NewPair(portfolio.NewAsset("ETH"), portfolio.NewAsset("USDT"))
			svc.WatchContract(deribit, contract)
			_, found := svc.Expirations(deribit, ethPair)
			Expect(found).To(BeFalse())
		})
	})

	Describe("Strikes", func() {
		It("returns false when no strikes are registered for the expiration", func() {
			svc.WatchContract(deribit, contract)
			_, found := svc.Strikes(deribit, btcPair, expiration)
			Expect(found).To(BeFalse())
		})

		It("returns false for an unwatched expiration", func() {
			_, found := svc.Strikes(deribit, btcPair, expiration)
			Expect(found).To(BeFalse())
		})
	})

	Describe("Trades", func() {
		It("returns empty slice when no trades exist", func() {
			trades := svc.Trades()
			Expect(trades).To(BeEmpty())
		})

		It("returns all trades when no query is given", func() {
			store.AddTrade(connector.Trade{ID: "t1", Exchange: deribit, Pair: btcPair})
			store.AddTrade(connector.Trade{ID: "t2", Exchange: deribit, Pair: btcPair})
			trades := svc.Trades()
			Expect(trades).To(HaveLen(2))
		})

		It("filters trades by exchange", func() {
			other := connector.ExchangeName("other")
			store.AddTrade(connector.Trade{ID: "t1", Exchange: deribit, Pair: btcPair})
			store.AddTrade(connector.Trade{ID: "t2", Exchange: other, Pair: btcPair})
			trades := svc.Trades(market.ActivityQuery{Exchange: &deribit})
			Expect(trades).To(HaveLen(1))
			Expect(trades[0].ID).To(Equal("t1"))
		})

		It("filters trades by pair", func() {
			ethPair := portfolio.NewPair(portfolio.NewAsset("ETH"), portfolio.NewAsset("USDT"))
			store.AddTrade(connector.Trade{ID: "t1", Exchange: deribit, Pair: btcPair})
			store.AddTrade(connector.Trade{ID: "t2", Exchange: deribit, Pair: ethPair})
			trades := svc.Trades(market.ActivityQuery{Pair: &btcPair})
			Expect(trades).To(HaveLen(1))
			Expect(trades[0].ID).To(Equal("t1"))
		})
	})

	Describe("Positions", func() {
		It("returns empty slice when no positions exist", func() {
			positions := svc.Positions()
			Expect(positions).To(BeEmpty())
		})

		It("returns all positions when no query is given", func() {
			store.SetPosition(contract, optionsTypes.Position{Exchange: deribit, Contract: contract, Quantity: 5})
			positions := svc.Positions()
			Expect(positions).To(HaveLen(1))
		})

		It("filters positions by exchange", func() {
			other := connector.ExchangeName("other")
			c2 := optionsTypes.OptionContract{Pair: btcPair, Strike: 55000, Expiration: expiration, OptionType: "PUT"}
			store.SetPosition(contract, optionsTypes.Position{Exchange: deribit, Contract: contract, Quantity: 5})
			store.SetPosition(c2, optionsTypes.Position{Exchange: other, Contract: c2, Quantity: 3})
			positions := svc.Positions(market.ActivityQuery{Exchange: &deribit})
			Expect(positions).To(HaveLen(1))
			Expect(positions[0].Exchange).To(Equal(deribit))
		})

		It("filters positions by pair", func() {
			ethPair := portfolio.NewPair(portfolio.NewAsset("ETH"), portfolio.NewAsset("USDT"))
			ethContract := optionsTypes.OptionContract{Pair: ethPair, Strike: 2000, Expiration: expiration, OptionType: "CALL"}
			store.SetPosition(contract, optionsTypes.Position{Exchange: deribit, Contract: contract, Quantity: 5})
			store.SetPosition(ethContract, optionsTypes.Position{Exchange: deribit, Contract: ethContract, Quantity: 2})
			positions := svc.Positions(market.ActivityQuery{Pair: &btcPair})
			Expect(positions).To(HaveLen(1))
			Expect(positions[0].Contract.Pair.Symbol()).To(Equal(btcPair.Symbol()))
		})
	})

	Describe("PNL", func() {
		It("returns a non-nil PNL calculator", func() {
			Expect(svc.PNL()).ToNot(BeNil())
		})
	})
})
