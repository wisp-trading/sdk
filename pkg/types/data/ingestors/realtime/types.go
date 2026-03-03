package realtime

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// RealtimeIngestorFactory creates realtime ingestors dynamically based on registered connectors
type RealtimeIngestorFactory interface {
	CreateIngestors() []RealtimeIngestor
}

// RealtimeIngestor handles WebSocket data ingestion for a specific market type
type RealtimeIngestor interface {
	Start(ctx context.Context) error
	Stop() error
	IsActive() bool
	GetMarketType() connector.MarketType
	GetActiveConnections() map[connector.ExchangeName]interface{} // Returns type-specific connectors
}

// WebSocketExtension allows market-specific WebSocket subscriptions (funding rate updates, etc.)
type WebSocketExtension interface {
	Subscribe(wsConn interface{}, exchangeName connector.ExchangeName, assets []portfolio.Pair) error
	Unsubscribe(wsConn interface{}, exchangeName connector.ExchangeName) error
	ProcessChannels(wsConn interface{}, exchangeName connector.ExchangeName, ctx context.Context)
}

// WebSocketSubscriber provides subscription methods
type WebSocketSubscriber interface {
	SubscribeOrderBook(asset portfolio.Pair) error
	SubscribeKlines(asset portfolio.Pair, interval string) error
	GetOrderBookChannels() map[string]<-chan connector.OrderBook
	GetKlineChannels() map[string]<-chan connector.Kline
}
