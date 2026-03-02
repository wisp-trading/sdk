package prediction

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/prediction"
	"github.com/wisp-trading/sdk/pkg/types/data/stores/market"
)

// One market has N outcomes, each with an order book
type OutcomeOrderBookMap map[prediction.OutcomeID]*connector.OrderBook

// Exchange -> Market -> Outcome -> OrderBook
type OrderBookMap map[connector.ExchangeName]map[prediction.MarketID]OutcomeOrderBookMap

// OrderBookWriter is a narrower interface for components that only write order books
type OrderBookWriter interface {
	UpdateOrderBook(
		exchange connector.ExchangeName,
		market prediction.Market,
		outcome prediction.Outcome,
		orderBook connector.OrderBook,
	)
}

// OrderBookStoreExtension provides order book data storage
type OrderBookStoreExtension interface {
	market.StoreExtension

	// UpdateOrderBook Update the order book for a specific market outcome on a specific exchange
	UpdateOrderBook(
		exchange connector.ExchangeName,
		marketID prediction.MarketID,
		outcomeID prediction.OutcomeID,
		orderBook connector.OrderBook,
	)

	// GetOrderBook Get a single outcome’s order book for a given market+exchange
	GetOrderBook(
		exchange connector.ExchangeName,
		marketID prediction.MarketID,
		outcomeID prediction.OutcomeID,
	) *connector.OrderBook

	// GetMarketOrderBooks Get all outcome orderbooks for a specific market on a specific exchange
	GetMarketOrderBooks(
		exchange connector.ExchangeName,
		marketID prediction.MarketID,
	) OutcomeOrderBookMap
}
