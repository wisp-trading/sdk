package extensions

import (
	"sync"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	predictionTypes "github.com/wisp-trading/sdk/pkg/types/connector/prediction"
	"github.com/wisp-trading/sdk/pkg/types/data/stores/market/prediction"
)

type predictionOrderBookExtension struct {
	mu         sync.RWMutex
	orderBooks prediction.OrderBookMap
}

func NewPredictionOrderBookExtension() prediction.OrderBookStoreExtension {
	return &predictionOrderBookExtension{
		orderBooks: make(prediction.OrderBookMap),
	}
}

func (e *predictionOrderBookExtension) UpdateOrderBook(
	exchange connector.ExchangeName,
	marketID predictionTypes.MarketID,
	outcomeID predictionTypes.OutcomeID,
	orderBook connector.OrderBook,
) {
	e.mu.Lock()

	if e.orderBooks[exchange] == nil {
		e.orderBooks[exchange] = make(map[predictionTypes.MarketID]prediction.OutcomeOrderBookMap)
	}
	if e.orderBooks[exchange][marketID] == nil {
		e.orderBooks[exchange][marketID] = make(prediction.OutcomeOrderBookMap)
	}

	e.orderBooks[exchange][marketID][outcomeID] = &orderBook

	e.mu.Unlock()
}

func (e *predictionOrderBookExtension) GetOrderBook(
	exchange connector.ExchangeName,
	marketID predictionTypes.MarketID,
	outcomeID predictionTypes.OutcomeID,
) *connector.OrderBook {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if exMap, ok := e.orderBooks[exchange]; ok {
		if mMap, ok := exMap[marketID]; ok {
			return mMap[outcomeID]
		}
	}
	return nil
}

func (e *predictionOrderBookExtension) GetMarketOrderBooks(
	exchange connector.ExchangeName,
	marketID predictionTypes.MarketID,
) prediction.OutcomeOrderBookMap {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if exMap, ok := e.orderBooks[exchange]; ok {
		if mMap, ok := exMap[marketID]; ok {
			// shallow copy to avoid external mutation
			result := make(prediction.OutcomeOrderBookMap, len(mMap))
			for k, v := range mMap {
				result[k] = v
			}
			return result
		}
	}
	return make(prediction.OutcomeOrderBookMap)
}
