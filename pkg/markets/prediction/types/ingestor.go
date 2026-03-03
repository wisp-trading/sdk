package types

import (
	"context"

	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// PredictionExtension allows market-specific WebSocket subscriptions for prediction markets (order book updates, etc.)
type PredictionExtension interface {
	Subscribe(wsConn interface{}, exchangeName connector.ExchangeName, market predictionconnector.Market) error
	Unsubscribe(wsConn interface{}, exchangeName connector.ExchangeName, market predictionconnector.Market) error
	ProcessChannels(wsConn interface{}, exchangeName connector.ExchangeName, ctx context.Context)
}

// PredictionSubscriber provides subscription methods
type PredictionSubscriber interface {
	SubscribeOrderBook(market predictionconnector.Market) error
	GetOrderBookChannels() map[string]<-chan connector.OrderBook
}
