package store

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

func (ds *dataStore) UpdateOrderBook(asset portfolio.Pair, exchangeName connector.ExchangeName, orderBook connector.OrderBook) {
	ds.mutex.Lock()

	current := ds.getOrderBooks()
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

	ds.orderBooks.Store(updated)

	ds.mutex.Unlock()

	if len(orderBook.Bids) > 0 && len(orderBook.Asks) > 0 {
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

		ds.UpdateAssetPrice(asset, exchangeName, price)
	}

	ds.UpdateLastUpdated(marketTypes.UpdateKey{
		DataType: marketTypes.DataKeyOrderBooks,
		Asset:    asset,
		Exchange: exchangeName,
	})
}

func (ds *dataStore) GetOrderBooks(asset portfolio.Pair) marketTypes.OrderBookMap {
	current := ds.getOrderBooks()
	if books, ok := current[asset]; ok {
		return books
	}
	return make(marketTypes.OrderBookMap)
}

func (ds *dataStore) GetOrderBook(asset portfolio.Pair, exchangeName connector.ExchangeName) *connector.OrderBook {
	current := ds.getOrderBooks()
	if books, ok := current[asset]; ok {
		if book, ok := books[exchangeName]; ok {
			return book
		}
	}
	return nil
}

func (ds *dataStore) GetAllAssetsWithOrderBooks() []portfolio.Pair {
	current := ds.getOrderBooks()
	assets := make([]portfolio.Pair, 0, len(current))
	for asset := range current {
		assets = append(assets, asset)
	}
	return assets
}
