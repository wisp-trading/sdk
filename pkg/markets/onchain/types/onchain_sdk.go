package types

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Onchain is the domain-scoped context for AMM / on-chain strategies.
// Injected via wisp.Onchain().
//
// Buy quantity = quote spent (exact-in). Sell quantity = base sold (exact-in).
type Onchain interface {
	// WatchPair registers a pair for local price/order bookkeeping (no CEX WS).
	WatchPair(exchange connector.ExchangeName, pair portfolio.Pair)
	UnwatchPair(exchange connector.ExchangeName, pair portfolio.Pair)

	// Price returns a last-known price if the strategy or connector recorded one.
	Price(exchange connector.ExchangeName, pair portfolio.Pair) (numerical.Decimal, bool)
	Prices(pair portfolio.Pair) map[connector.ExchangeName]numerical.Decimal

	// Signal creates an onchain signal builder.
	Signal(strategyName strategy.StrategyName) OnchainSignalBuilder

	// Emit routes an onchain signal to the domain executor (swap path).
	Emit(signal OnchainSignal) execution.ExecutionCallback

	Log() logging.TradingLogger

	Trades(q ...market.ActivityQuery) []connector.Trade
	// Positions returns open/placed swap records (same surface as spot orders).
	Positions(q ...market.ActivityQuery) []connector.Order
	PNL() OnchainPNL
}
