package spot

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// WebSocketConnector extends Connector with real-time capabilities for spot markets
type WebSocketConnector interface {
	Connector
	connector.WebSocketCapable

	// Subscription management
	SubscribeOrderBook(asset portfolio.Pair) error
	SubscribeTrades(asset portfolio.Pair) error
	SubscribeKlines(asset portfolio.Pair, interval string) error
	SubscribeAccountBalance() error

	UnsubscribeOrderBook(asset portfolio.Pair) error
	UnsubscribeTrades(asset portfolio.Pair) error
	UnsubscribeKlines(asset portfolio.Pair, interval string) error
	UnsubscribeAccountBalance() error

	// Data access channels
	GetOrderBookChannels() map[string]<-chan connector.OrderBook
	GetKlineChannels() map[string]<-chan connector.Kline
	TradeUpdates() <-chan connector.Trade
	AssetBalanceUpdates() <-chan connector.AssetBalance
}
