package types

import (
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

	Markets(exchange connector.ExchangeName, filter *predictionconnector.MarketsFilter) ([]predictionconnector.Market, error)

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

	// PNL returns profit and loss calculations for this prediction market instance.
	PNL() PredictionPNL
}
