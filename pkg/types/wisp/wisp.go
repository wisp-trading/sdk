package wisp

import (
	onchainTypes "github.com/wisp-trading/sdk/pkg/markets/onchain/types"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	predTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
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

	// Spot returns the spot market domain context.
	// Example: sig, _ := wisp.Spot().Signal(name).BuyLimit(...).Build(); wisp.Spot().Emit(sig)
	Spot() spotTypes.Spot

	// Perp returns the perpetual futures domain context.
	// Example: sig, _ := wisp.Perp().Signal(name).BuyLimit(...).Build(); wisp.Perp().Emit(sig)
	Perp() perpTypes.Perp

	// Predict returns the prediction market domain context.
	// Example: sig, _ := wisp.Predict().PredictionSignal(name).Buy(...).Build(); wisp.Predict().Emit(sig)
	Predict() predTypes.Predict

	// Options returns the options market domain context.
	// Example: sig, _ := wisp.Options().Signal(name)...Build(); wisp.Options().Emit(sig)
	Options() optionsTypes.Options

	// Onchain returns the on-chain AMM domain context (e.g. UniV3).
	// Buy quantity = quote spent; Sell quantity = base sold (exact-in).
	// Example: sig, _ := wisp.Onchain().Signal(name).Buy(...).Build(); wisp.Onchain().Emit(sig)
	Onchain() onchainTypes.Onchain

	// PriceFeeds returns the price feeds service for accessing external price data.
	// Strategies use this to query price feeds from sources like Pyth, Chainlink, etc.
	// Example: wisp.PriceFeeds().GetLatestPrice(feedID)
	PriceFeeds() types.PriceFeeds
}
