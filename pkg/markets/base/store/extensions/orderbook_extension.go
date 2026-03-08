package extensions

import (
	"sync"

	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type orderBookExtension struct {
	orderBooks       map[portfolio.Pair]market.OrderBookMap
	mu               sync.RWMutex
	onUpdatePrice    func(portfolio.Pair, connector.ExchangeName, connector.Price)
	onUpdateMetadata func(market.UpdateKey)
}

func NewOrderBookExtension(
	priceUpdater func(portfolio.Pair, connector.ExchangeName, connector.Price),
	metadataUpdater func(market.UpdateKey),
) market.OrderBookStoreExtension {
	return &orderBookExtension{
		orderBooks:       make(map[portfolio.Pair]market.OrderBookMap),
		onUpdatePrice:    priceUpdater,
		onUpdateMetadata: metadataUpdater,
	}
}

func (e *orderBookExtension) UpdateOrderBook(asset portfolio.Pair, exchangeName connector.ExchangeName, orderBook connector.OrderBook) {
	e.mu.Lock()

	if e.orderBooks[asset] == nil {
		e.orderBooks[asset] = make(market.OrderBookMap)
	}
	e.orderBooks[asset][exchangeName] = &orderBook

	e.mu.Unlock()

	if e.onUpdatePrice != nil && len(orderBook.Bids) > 0 && len(orderBook.Asks) > 0 {
		bestBid := orderBook.Bids[0].Price
		bestAsk := orderBook.Asks[0].Price
		midPrice := bestBid.Add(bestAsk).Div(numerical.NewFromInt(2))
		e.onUpdatePrice(asset, exchangeName, connector.Price{
			Pair:      asset,
			Price:     midPrice,
			BidPrice:  bestBid,
			AskPrice:  bestAsk,
			Source:    exchangeName,
			Timestamp: orderBook.Timestamp,
		})
	}

	if e.onUpdateMetadata != nil {
		e.onUpdateMetadata(market.UpdateKey{
			DataType: market.DataKeyOrderBooks,
			Pair:     asset,
			Exchange: exchangeName,
		})
	}
}

func (e *orderBookExtension) GetOrderBooks(asset portfolio.Pair) market.OrderBookMap {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if books, ok := e.orderBooks[asset]; ok {
		result := make(market.OrderBookMap, len(books))
		for k, v := range books {
			result[k] = v
		}
		return result
	}
	return make(market.OrderBookMap)
}

func (e *orderBookExtension) GetOrderBook(asset portfolio.Pair, exchangeName connector.ExchangeName) *connector.OrderBook {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if books, ok := e.orderBooks[asset]; ok {
		return books[exchangeName]
	}
	return nil
}

func (e *orderBookExtension) GetAllPairsWithOrderBooks() []portfolio.Pair {
	e.mu.RLock()
	defer e.mu.RUnlock()

	assets := make([]portfolio.Pair, 0, len(e.orderBooks))
	for asset := range e.orderBooks {
		assets = append(assets, asset)
	}
	return assets
}
