package types

import (
	"math/big"

	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type Predict interface {
	GetMarketBySlug(slug string, exchange connector.ExchangeName) (predictionconnector.Market, error)
	GetRecurringMarketBySlug(slug string, recurrenceInterval predictionconnector.RecurrenceInterval, exchange connector.ExchangeName) (predictionconnector.Market, error)
	WatchMarket(exchange connector.ExchangeName, market predictionconnector.Market)

	// Markets returns markets from store if available, fetches from API only if store is empty.
	// Use RefreshMarkets() for explicit refresh from API.
	Markets(exchange connector.ExchangeName, filter *predictionconnector.MarketsFilter) ([]predictionconnector.Market, error)

	// RefreshMarkets fetches markets fresh from API and updates the store.
	RefreshMarkets(exchange connector.ExchangeName, filter *predictionconnector.MarketsFilter) ([]predictionconnector.Market, error)

	// LoadMarkets starts background pagination of all markets matching the filter.
	// Returns immediately; actual loading happens asynchronously in background.
	// Strategy calls Markets() later to get cached data.
	LoadMarkets(exchange connector.ExchangeName, filter *predictionconnector.MarketsFilter) error

	// IsLoadingMarkets returns true if background market loading is in progress for the exchange
	IsLoadingMarkets(exchange connector.ExchangeName) bool

	// GetLoadProgress returns count of markets loaded so far (for monitoring)
	GetLoadProgress(exchange connector.ExchangeName) int

	Orderbook(
		exchange connector.ExchangeName,
		market predictionconnector.Market,
		outcome predictionconnector.Outcome,
	) (*connector.OrderBook, error)

	// Balance returns the current balance for an asset on an exchange.
	Balance(exchange connector.ExchangeName, asset portfolio.Asset) (numerical.Decimal, bool)

	// Positions returns orders recorded for this instance.
	// Optionally filter by exchange and/or market slug:
	//   wisp.Predict().Positions()
	//   wisp.Predict().Positions(PredictionActivityQuery{Exchange: &exchange})
	//   wisp.Predict().Positions(PredictionActivityQuery{MarketID: &slug})
	Positions(q ...PredictionActivityQuery) []PredictionOrder

	// Log returns the trading logger for strategy-specific logging.
	Log() logging.TradingLogger

	// PredictionSignal creates a new signal builder for prediction market trading signals.
	PredictionSignal(strategyName strategy.StrategyName) PredictionSignalBuilder

	GetTokensToRedeem(market predictionconnector.Market) ([]predictionconnector.Balance, error)

	// Redeem attempts to redeem winnings for a market. Returns an error if redemption fails.
	Redeem(market predictionconnector.Market) error

	// MergePositions burns one bundle of YES tokens across all outcomes in a NegRisk group
	// and returns USDC. amountUSDC is in 6-decimal units ($1.00 = 1_000_000).
	// Use this to close a BuyMerge arbitrage position instantly without waiting for resolution.
	MergePositions(market predictionconnector.Market, amountUSDC *big.Int) (string, error)

	// SplitPosition deposits amountUSDC and mints one YES token per outcome in the market.
	// amountUSDC is in 6-decimal units ($1.00 = 1_000_000).
	// Use this to open a MintSell arbitrage position.
	SplitPosition(market predictionconnector.Market, amountUSDC *big.Int) (string, error)

	// PNL returns profit and loss calculations for this prediction market instance.
	PNL() PredictionPNL
}
