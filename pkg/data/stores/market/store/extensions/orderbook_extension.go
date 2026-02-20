package extensions

import (
	"sync"
	"sync/atomic"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Type alias for order book storage
type assetOrderBooks map[portfolio.Pair]marketTypes.OrderBookMap

// orderBookExtension stores order book data
type orderBookExtension struct {
	orderBooks *atomic.Value // assetOrderBooks
	mu         sync.RWMutex

	// Dependencies injected at construction
	onUpdatePrice    func(portfolio.Pair, connector.ExchangeName, connector.Price)
	onUpdateMetadata func(marketTypes.UpdateKey)
}

// NewOrderBookExtension creates a new order book extension
// Optional callbacks can be provided for side effects (price updates, metadata updates)
func NewOrderBookExtension(
	priceUpdater func(portfolio.Pair, connector.ExchangeName, connector.Price),
	metadataUpdater func(marketTypes.UpdateKey),
) marketTypes.OrderBookStoreExtension {
	ext := &orderBookExtension{
		orderBooks:       &atomic.Value{},
		onUpdatePrice:    priceUpdater,
		onUpdateMetadata: metadataUpdater,
	}
	ext.orderBooks.Store(make(assetOrderBooks))
	return ext
}

// Helper methods to get typed data
func (e *orderBookExtension) getOrderBooks() assetOrderBooks {
	if v := e.orderBooks.Load(); v != nil {
		return v.(assetOrderBooks)
	}
	return make(assetOrderBooks)
}

func (e *orderBookExtension) UpdateOrderBook(asset portfolio.Pair, exchangeName connector.ExchangeName, orderBook connector.OrderBook) {
	e.mu.Lock()

	current := e.getOrderBooks()
	updated := make(assetOrderBooks, len(current))
	for k, v := range current {
		updated[k] = v
	}

	if updated[asset] == nil {
		updated[asset] = make(marketTypes.OrderBookMap)
	}

	assetBooks := make(marketTypes.OrderBookMap, len(updated[asset]))
	for k, v := range updated[asset] {
		assetBooks[k] = v
	}
	assetBooks[exchangeName] = &orderBook
	updated[asset] = assetBooks

	e.orderBooks.Store(updated)

	e.mu.Unlock()

	// Calculate mid-price and trigger price update callback if provided
	if e.onUpdatePrice != nil && len(orderBook.Bids) > 0 && len(orderBook.Asks) > 0 {
		bestBid := orderBook.Bids[0].Price
		bestAsk := orderBook.Asks[0].Price
		midPrice := bestBid.Add(bestAsk).Div(numerical.NewFromInt(2))

		price := connector.Price{
			Symbol:    asset.Symbol(),
			Price:     midPrice,
			BidPrice:  bestBid,
			AskPrice:  bestAsk,
			Source:    exchangeName,
			Timestamp: orderBook.Timestamp,
		}

		e.onUpdatePrice(asset, exchangeName, price)
	}

	// Trigger metadata update callback if provided
	if e.onUpdateMetadata != nil {
		e.onUpdateMetadata(marketTypes.UpdateKey{
			DataType: marketTypes.DataKeyOrderBooks,
			Pair:     asset,
			Exchange: exchangeName,
		})
	}
}

func (e *orderBookExtension) GetOrderBooks(asset portfolio.Pair) marketTypes.OrderBookMap {
	current := e.getOrderBooks()
	if books, ok := current[asset]; ok {
		return books
	}
	return make(marketTypes.OrderBookMap)
}

func (e *orderBookExtension) GetOrderBook(asset portfolio.Pair, exchangeName connector.ExchangeName) *connector.OrderBook {
	current := e.getOrderBooks()
	if books, ok := current[asset]; ok {
		if book, ok := books[exchangeName]; ok {
			return book
		}
	}
	return nil
}

func (e *orderBookExtension) GetAllPairsWithOrderBooks() []portfolio.Pair {
	current := e.getOrderBooks()
	assets := make([]portfolio.Pair, 0, len(current))
	for asset := range current {
		assets = append(assets, asset)
	}
	return assets
}
