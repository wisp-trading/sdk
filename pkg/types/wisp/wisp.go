package wisp

import (
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	predTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/activity"
	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
)

// Wisp is the main SDK context provided to strategies for accessing market data,
// indicators, analytics, and trading functionality. It is injected into strategy
// implementations and provides read-only access to all framework services.
type Wisp interface {

	// Indicators returns the indicators service for technical analysis.
	Indicators() analytics.Indicators

	// Analytics returns the analytics service for market analysis.
	Analytics() analytics.Analytics

	// Log returns the trading logger for strategy-specific logging.
	Log() logging.TradingLogger

	// Activity returns read-only access to positions, trades, and PNL data.
	Activity() activity.Activity

	// Asset creates a new portfolio.Asset from a symbol string.
	Asset(symbol string) portfolio.Asset

	// Pair creates a new portfolio.Pair from two assets.
	Pair(base, quote portfolio.Asset) portfolio.Pair

	// Emit routes a signal directly to the executor. Non-blocking.
	Emit(signal strategy.Signal)

	// Spot returns the spot market domain context.
	// Owns watchlist management, orderbooks, balances, positions, and signal creation.
	// Example: wisp.Spot().WatchPair(exchange, btc)
	// Example: wisp.Spot().Signal(strategyName).BuyLimit(pair, exchange, qty, price).Build()
	Spot() spotTypes.Spot

	// Perp returns the perpetual futures domain context.
	// Owns watchlist management, funding rates, positions, orderbooks, and signal creation.
	// Example: wisp.Perp().WatchPair(exchange, btc)
	// Example: wisp.Perp().Signal(strategyName).BuyLimit(pair, exchange, qty, price).Build()
	Perp() perpTypes.Perp

	// Predict returns the prediction market domain context.
	// Owns market discovery, orderbooks, balances, positions, and signal creation.
	// Example: wisp.Predict().WatchMarket(exchange, market)
	// Example: wisp.Predict().Signal(strategyName).Buy(market, outcome, exchange, shares, maxPrice, expiry).Build()
	Predict() predTypes.Predict

	// Options returns the options market domain context.
	// Owns watchlist management, Greeks, IV, positions, and signal creation.
	// Example: wisp.Options().WatchContract(exchange, contract)
	// Example: wisp.Options().MarkPrice(exchange, contract)
	Options() optionsTypes.Options

	// PriceFeeds returns the price feeds service for accessing external price data.
	// Strategies use this to query price feeds from sources like Pyth, Chainlink, etc.
	// Example: wisp.PriceFeeds().GetLatestPrice(feedID)
	PriceFeeds() types.PriceFeeds
}
