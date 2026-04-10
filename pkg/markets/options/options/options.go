package options

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	optionsSignal "github.com/wisp-trading/sdk/pkg/markets/options/signal"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type options struct {
	tradingLogger logging.TradingLogger
	watchlist     optionsTypes.OptionsWatchlist
	store         optionsTypes.OptionsStore
	timeProvider  temporal.TimeProvider
	pnl           optionsTypes.OptionsPNL
}

// NewOptions creates a new options service
func NewOptions(
	tradingLogger logging.TradingLogger,
	watchlist optionsTypes.OptionsWatchlist,
	store optionsTypes.OptionsStore,
	timeProvider temporal.TimeProvider,
	pnl optionsTypes.OptionsPNL,
) optionsTypes.Options {
	return &options{
		tradingLogger: tradingLogger,
		watchlist:     watchlist,
		store:         store,
		timeProvider:  timeProvider,
		pnl:           pnl,
	}
}

// WatchContract registers an options contract's expiration on the watchlist so
// data ingestors begin collecting market data for all strikes on that expiration.
func (o *options) WatchContract(exchange connector.ExchangeName, contract optionsTypes.OptionContract) {
	_ = o.watchlist.RequireExpiration(exchange, contract.Pair, contract.Expiration)
}

// UnwatchContract removes an options contract's expiration from the watchlist.
func (o *options) UnwatchContract(exchange connector.ExchangeName, contract optionsTypes.OptionContract) {
	_ = o.watchlist.ReleaseExpiration(exchange, contract.Pair, contract.Expiration)
}

// MarkPrice returns the current mark price for an option contract from the store.
func (o *options) MarkPrice(exchange connector.ExchangeName, contract optionsTypes.OptionContract) (numerical.Decimal, bool) {
	price := o.store.GetMarkPrice(contract)
	if price == 0 {
		return numerical.Zero(), false
	}
	return numerical.NewFromFloat(price), true
}

// UnderlyingPrice returns the current underlying asset price for an option contract.
func (o *options) UnderlyingPrice(exchange connector.ExchangeName, contract optionsTypes.OptionContract) (numerical.Decimal, bool) {
	price := o.store.GetUnderlyingPrice(contract)
	if price == 0 {
		return numerical.Zero(), false
	}
	return numerical.NewFromFloat(price), true
}

// Greeks returns the Greeks for an option contract.
func (o *options) Greeks(exchange connector.ExchangeName, contract optionsTypes.OptionContract) (optionsTypes.Greeks, bool) {
	greeks := o.store.GetGreeks(contract)
	if greeks == (optionsTypes.Greeks{}) {
		return optionsTypes.Greeks{}, false
	}
	return greeks, true
}

// ImpliedVolatility returns the implied volatility for an option contract.
func (o *options) ImpliedVolatility(exchange connector.ExchangeName, contract optionsTypes.OptionContract) (float64, bool) {
	iv := o.store.GetIV(contract)
	if iv == 0 {
		return 0, false
	}
	return iv, true
}

// Expirations returns all watched expiration dates for a pair on an exchange.
// Only expirations registered via WatchContract are returned.
func (o *options) Expirations(exchange connector.ExchangeName, pair portfolio.Pair) ([]time.Time, bool) {
	watched := o.watchlist.GetWatchedExpirations(exchange)
	expirations, ok := watched[pair]
	if !ok || len(expirations) == 0 {
		return nil, false
	}
	return expirations, true
}

// Strikes returns all available strike prices for a watched expiration.
// Strikes are populated by the ingestor after WatchContract is called.
func (o *options) Strikes(exchange connector.ExchangeName, pair portfolio.Pair, expiration time.Time) ([]float64, bool) {
	strikes := o.watchlist.GetAvailableStrikes(exchange, pair, expiration)
	if len(strikes) == 0 {
		return nil, false
	}
	return strikes, true
}

// Log returns the trading logger for strategy-specific logging.
// Signal creates a new options signal builder for the given strategy.
func (o *options) Signal(strategyName strategy.StrategyName) optionsTypes.OptionsSignalBuilder {
	return optionsSignal.NewOptionsBuilder(strategyName, o.timeProvider)
}

func (o *options) Log() logging.TradingLogger {
	return o.tradingLogger
}

// Trades returns all trades executed in the options domain.
// Pass an ActivityQuery to filter by exchange and/or pair.
func (o *options) Trades(q ...market.ActivityQuery) []connector.Trade {
	if len(q) > 0 {
		return o.store.QueryTrades(q[0])
	}
	return o.store.GetAllTrades()
}

// Positions returns all open positions in the options domain.
// Pass an ActivityQuery to filter by exchange and/or pair.
func (o *options) Positions(q ...market.ActivityQuery) []optionsTypes.Position {
	if len(q) > 0 {
		return o.store.QueryPositions(q[0])
	}
	return o.store.GetAllPositions()
}

// PNL returns profit and loss calculations for the options domain.
func (o *options) PNL() optionsTypes.OptionsPNL {
	return o.pnl
}

var _ optionsTypes.Options = (*options)(nil)
