package perp

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// WebSocketConnector extends Connector with real-time capabilities for perpetual markets
type WebSocketConnector interface {
	Connector
	connector.WebSocketCapable

	// Subscription management
	SubscribeOrderBook(pair portfolio.Pair) error
	SubscribeTrades(pair portfolio.Pair) error
	SubscribePositions(pair portfolio.Pair) error
	SubscribeFundingRates(pair portfolio.Pair) error
	SubscribeKlines(pair portfolio.Pair, interval string) error
	SubscribeAccountBalance() error

	UnsubscribeOrderBook(pair portfolio.Pair) error
	UnsubscribeTrades(pair portfolio.Pair) error
	UnsubscribePositions(pair portfolio.Pair) error
	UnsubscribeFundingRates(pair portfolio.Pair) error
	UnsubscribeKlines(pair portfolio.Pair, interval string) error
	UnsubscribeAccountBalance() error

	// Data access channels
	GetOrderBookChannels() map[string]<-chan connector.OrderBook
	GetKlineChannels() map[string]<-chan connector.Kline
	TradeUpdates() <-chan connector.Trade
	PositionUpdates() <-chan Position
	FundingRateUpdates() <-chan FundingRate
	//AssetBalanceUpdates() <-chan AssetBalance
}
