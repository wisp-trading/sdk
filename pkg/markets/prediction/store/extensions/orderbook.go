package extensions

import (
	"sync"

	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

type predictionOrderBookExtension struct {
	mu         sync.RWMutex
	orderBooks types.OrderBookMap
}

func NewPredictionOrderBookExtension() types.OrderBookStoreExtension {
	return &predictionOrderBookExtension{
		orderBooks: make(types.OrderBookMap),
	}
}

func (e *predictionOrderBookExtension) UpdateOrderBook(
	exchange connector.ExchangeName,
	marketID predictionconnector.MarketID,
	outcomeID predictionconnector.OutcomeID,
	orderBook connector.OrderBook,
) {
	e.mu.Lock()

	if e.orderBooks[exchange] == nil {
		e.orderBooks[exchange] = make(map[predictionconnector.MarketID]types.OutcomeOrderBookMap)
	}
	if e.orderBooks[exchange][marketID] == nil {
		e.orderBooks[exchange][marketID] = make(types.OutcomeOrderBookMap)
	}

	e.orderBooks[exchange][marketID][outcomeID] = &orderBook

	e.mu.Unlock()
}

func (e *predictionOrderBookExtension) GetOrderBook(
	exchange connector.ExchangeName,
	marketID predictionconnector.MarketID,
	outcomeID predictionconnector.OutcomeID,
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

func (e *predictionOrderBookExtension) RemoveOrderBook(
	exchange connector.ExchangeName,
	marketID predictionconnector.MarketID,
) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if exMap, ok := e.orderBooks[exchange]; ok {
		delete(exMap, marketID)
	}
}

func (e *predictionOrderBookExtension) GetMarketOrderBooks(
	exchange connector.ExchangeName,
	marketID predictionconnector.MarketID,
) types.OutcomeOrderBookMap {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if exMap, ok := e.orderBooks[exchange]; ok {
		if mMap, ok := exMap[marketID]; ok {
			result := make(types.OutcomeOrderBookMap, len(mMap))
			for k, v := range mMap {
				result[k] = v
			}
			return result
		}
	}
	return make(types.OutcomeOrderBookMap)
}
