package connector

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// WebSocketConnector extends Connector with real-time capabilities for prediction markets
type WebSocketConnector interface {
	Connector
	connector.WebSocketCapable

	SubscribeOrderBook(market Market) error
	SubscribePriceChanges(market Market) error
	SubscribeTrades(market Market) error
	SubscribeOrders(market Market) error

	UnsubscribeMarket(market Market) error
	UnsubscribeUserMarket(market Market) error

	CancelOrder(orderID string, outcome ...Outcome) (*connector.CancelResponse, error)

	FetchOrderBooks(market Market, outcome Outcome) (*OrderBook, error)

	GetOrderBookUpdates() <-chan OrderBook
	GetPriceChangeChannels() map[string]<-chan PriceChange
	GetTradesChannel() <-chan connector.Trade
	GetOrdersChannel() <-chan connector.Order
}
