package options

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// WebSocketConnector extends Connector with real-time capabilities for options markets
type WebSocketConnector interface {
	Connector
	connector.WebSocketCapable

	// Subscription management
	SubscribeExpirationUpdates(pair portfolio.Pair, expiration time.Time) error
	UnsubscribeExpirationUpdates(pair portfolio.Pair, expiration time.Time) error

	// Data channels
	GetOptionUpdateChannels() map[string]<-chan OptionUpdate
	GetTradeChannels() map[string]<-chan connector.Trade
	GetOrderBookChannels() map[string]<-chan connector.OrderBook
}

// OptionUpdate represents a real-time update for an option
type OptionUpdate struct {
	Contract        OptionContract
	MarkPrice       float64
	UnderlyingPrice float64
	Greeks          Greeks
	IV              float64
	Timestamp       time.Time
}
