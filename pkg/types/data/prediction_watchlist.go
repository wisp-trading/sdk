package data

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	prediction "github.com/wisp-trading/sdk/pkg/types/connector/prediction"
)

type PredictionWatchlist interface {
	RequireMarket(exchange connector.ExchangeName, market prediction.Market)
	ReleaseMarket(exchange connector.ExchangeName, marketID prediction.MarketID)

	GetRequiredMarkets(exchange connector.ExchangeName) []prediction.Market

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
	Market   prediction.Market
	Type     PredictionWatchEventType
}
