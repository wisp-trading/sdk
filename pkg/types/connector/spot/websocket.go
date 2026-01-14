package spot

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector/common"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// WebSocketConnector extends Connector with real-time capabilities for spot markets
type WebSocketConnector interface {
	Connector
	common.WebSocketCapable

	// Subscription management
	SubscribeOrderBook(asset portfolio.Asset) error
	SubscribeTrades(asset portfolio.Asset) error
	SubscribeKlines(asset portfolio.Asset, interval string) error
	SubscribeAccountBalance() error

	UnsubscribeOrderBook(asset portfolio.Asset) error
	UnsubscribeTrades(asset portfolio.Asset) error
	UnsubscribeKlines(asset portfolio.Asset, interval string) error
	UnsubscribeAccountBalance() error

	// Data access channels
	GetOrderBookChannels() map[string]<-chan connector.OrderBook
	GetKlineChannels() map[string]<-chan connector.Kline
	TradeUpdates() <-chan connector.Trade
	AccountBalanceUpdates() <-chan connector.AccountBalance
}
