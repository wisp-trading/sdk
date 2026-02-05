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
	SubscribeOrderBook(asset portfolio.Pair) error
	SubscribeTrades(asset portfolio.Pair) error
	SubscribePositions(asset portfolio.Pair) error
	SubscribeFundingRates(asset portfolio.Pair) error
	SubscribeKlines(asset portfolio.Pair, interval string) error
	SubscribeAccountBalance() error

	UnsubscribeOrderBook(asset portfolio.Pair) error
	UnsubscribeTrades(asset portfolio.Pair) error
	UnsubscribePositions(asset portfolio.Pair) error
	UnsubscribeFundingRates(asset portfolio.Pair) error
	UnsubscribeKlines(asset portfolio.Pair, interval string) error
	UnsubscribeAccountBalance() error

	// Data access channels
	GetOrderBookChannels() map[string]<-chan connector.OrderBook
	GetKlineChannels() map[string]<-chan connector.Kline
	TradeUpdates() <-chan connector.Trade
	PositionUpdates() <-chan Position
	FundingRateUpdates() <-chan FundingRate
	AssetBalanceUpdates() <-chan connector.AssetBalance
}
