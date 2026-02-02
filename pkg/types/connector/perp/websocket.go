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
	SubscribeOrderBook(asset portfolio.Asset) error
	SubscribeTrades(asset portfolio.Asset) error
	SubscribePositions(asset portfolio.Asset) error
	SubscribeFundingRates(asset portfolio.Asset) error
	SubscribeKlines(asset portfolio.Asset, interval string) error
	SubscribeAccountBalance() error

	UnsubscribeOrderBook(asset portfolio.Asset) error
	UnsubscribeTrades(asset portfolio.Asset) error
	UnsubscribePositions(asset portfolio.Asset) error
	UnsubscribeFundingRates(asset portfolio.Asset) error
	UnsubscribeKlines(asset portfolio.Asset, interval string) error
	UnsubscribeAccountBalance() error

	// Data access channels
	GetOrderBookChannels() map[string]<-chan connector.OrderBook
	GetKlineChannels() map[string]<-chan connector.Kline
	TradeUpdates() <-chan connector.Trade
	PositionUpdates() <-chan Position
	FundingRateUpdates() <-chan FundingRate
	AssetBalanceUpdates() <-chan connector.AssetBalance
}
