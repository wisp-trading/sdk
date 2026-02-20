package data

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/prediction"
)

type PredictionWatchlist interface {
	RequireMarket(exchange connector.ExchangeName, market prediction.Market)
	ReleaseMarket(exchange connector.ExchangeName, marketID string)

	GetRequiredMarkets(exchange connector.ExchangeName) []prediction.Market

	Events() <-chan PredictionWatchEvent
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
