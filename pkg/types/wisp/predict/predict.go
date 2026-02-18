package predict

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/prediction"
	"github.com/wisp-trading/sdk/pkg/types/logging"
)

type Predict interface {
	GetMarketBySlug(slug string, exchange connector.ExchangeName) (prediction.Market, error)
	GetRecurringMarketBySlug(slug string, recurrenceInterval prediction.RecurrenceInterval, exchange connector.ExchangeName) (prediction.Market, error)
	WatchMarket(market prediction.Market, exchange *connector.ExchangeName) error

	Markets() []prediction.Market
	Orderbooks(market prediction.Market) (map[prediction.Outcome]prediction.OrderBook, error)
	Orderbook(market prediction.Market, outcome prediction.Outcome) (*prediction.OrderBook, error)

	// Log returns the trading logger for strategy-specific logging.
	// Use for recording trading decisions and strategy events.
	Log() logging.TradingLogger
}
