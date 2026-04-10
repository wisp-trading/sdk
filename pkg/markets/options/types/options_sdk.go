package types

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Options is the domain-scoped context object for options market strategies.
// Injected via wisp.Options() — owns watchlist management, market data,
// and signal creation for the options domain.
type Options interface {
	// WatchContract registers an options contract on the watchlist so data ingestors
	// begin collecting market data (Greeks, IV, prices) for it.
	WatchContract(exchange connector.ExchangeName, contract OptionContract)

	// UnwatchContract removes an options contract from the watchlist.
	UnwatchContract(exchange connector.ExchangeName, contract OptionContract)

	// MarkPrice returns the current mark price for an option contract.
	MarkPrice(exchange connector.ExchangeName, contract OptionContract) (numerical.Decimal, bool)

	// UnderlyingPrice returns the current underlying asset price for an option contract.
	UnderlyingPrice(exchange connector.ExchangeName, contract OptionContract) (numerical.Decimal, bool)

	// Greeks returns the Greeks (delta, gamma, theta, vega, rho) for an option contract.
	Greeks(exchange connector.ExchangeName, contract OptionContract) (Greeks, bool)

	// ImpliedVolatility returns the implied volatility for an option contract.
	ImpliedVolatility(exchange connector.ExchangeName, contract OptionContract) (float64, bool)

	// Expirations returns all available expiration dates for a pair on an exchange.
	Expirations(exchange connector.ExchangeName, pair portfolio.Pair) ([]time.Time, bool)

	// Strikes returns all available strike prices for a pair and expiration on an exchange.
	Strikes(exchange connector.ExchangeName, pair portfolio.Pair, expiration time.Time) ([]float64, bool)

	// Signal creates a new options signal builder for the given strategy.
	Signal(strategyName strategy.StrategyName) OptionsSignalBuilder

	// Log returns the trading logger for strategy-specific logging.
	Log() logging.TradingLogger

	// Trades returns all trades executed in the options domain.
	// Optionally filter by exchange and/or contract.
	Trades(q ...market.ActivityQuery) []connector.Trade

	// Positions returns all open positions in the options domain.
	// Optionally filter by exchange and/or contract.
	Positions(q ...market.ActivityQuery) []Position

	// PNL returns profit and loss calculations for the options domain.
	PNL() OptionsPNL
}
