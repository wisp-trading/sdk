package types

import (
	predictiontypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"

	"github.com/wisp-trading/sdk/pkg/types/connector"
)

type PredictionWatchlist interface {
	RequireMarket(exchange connector.ExchangeName, market predictiontypes.Market)
	ReleaseMarket(exchange connector.ExchangeName, marketID predictiontypes.MarketID)

	GetRequiredMarkets(exchange connector.ExchangeName) []predictiontypes.Market

	// GetAllMarkets returns all watched markets grouped by exchange.
	GetAllMarkets() map[connector.ExchangeName][]predictiontypes.Market

	Subscribe(exchange connector.ExchangeName) chan PredictionWatchEvent
	Unsubscribe(exchange connector.ExchangeName)
}

type PredictionWatchEventType int

const (
	PredictionMarketAdded PredictionWatchEventType = iota
	PredictionMarketRemoved
)

type PredictionWatchEvent struct {
	Exchange connector.ExchangeName
	Market   predictiontypes.Market
	Type     PredictionWatchEventType
}
