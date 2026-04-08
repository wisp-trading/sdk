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

	// Slow watch — subscribe to ticker updates for all contracts in an expiration.
	// contracts is resolved by the realtime ingestor from the watchlist (no REST calls).
	SubscribeExpirationUpdates(pair portfolio.Pair, expiration time.Time, contracts []OptionContract) error
	UnsubscribeExpirationUpdates(pair portfolio.Pair, expiration time.Time, contracts []OptionContract) error

	// Fast watch — subscribe to real-time order book depth for a specific contract.
	// Called when the user explicitly watches an instrument (pre-trade or position management).
	SubscribeOrderBook(contract *OptionContract) error
	UnsubscribeOrderBook(contract *OptionContract) error

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
