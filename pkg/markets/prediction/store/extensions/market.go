package extensions

import (
	"sync"

	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

type predictionMarketExtension struct {
	mu      sync.RWMutex
	markets types.MarketMap
}

func NewPredictionMarketExtension() types.MarketStoreExtension {
	return &predictionMarketExtension{
		markets: make(types.MarketMap),
	}
}

func (e *predictionMarketExtension) UpdateMarkets(
	exchange connector.ExchangeName,
	markets []predictionconnector.Market,
) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.markets[exchange] == nil {
		e.markets[exchange] = make(map[predictionconnector.MarketID]predictionconnector.Market)
	}

	for _, market := range markets {
		e.markets[exchange][market.MarketID] = market
	}
}

func (e *predictionMarketExtension) GetMarket(
	exchange connector.ExchangeName,
	marketID predictionconnector.MarketID,
) *predictionconnector.Market {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if exMap, ok := e.markets[exchange]; ok {
		if market, ok := exMap[marketID]; ok {
			return &market
		}
	}
	return nil
}

func (e *predictionMarketExtension) GetMarkets(exchange connector.ExchangeName) []predictionconnector.Market {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if exMap, ok := e.markets[exchange]; ok {
		result := make([]predictionconnector.Market, 0, len(exMap))
		for _, market := range exMap {
			result = append(result, market)
		}
		return result
	}
	return []predictionconnector.Market{}
}

func (e *predictionMarketExtension) GetAllMarkets() types.MarketMap {
	e.mu.RLock()
	defer e.mu.RUnlock()

	result := make(types.MarketMap)
	for ex, exMap := range e.markets {
		result[ex] = make(map[predictionconnector.MarketID]predictionconnector.Market)
		for id, market := range exMap {
			result[ex][id] = market
		}
	}
	return result
}

func (e *predictionMarketExtension) ClearMarkets(exchange connector.ExchangeName) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.markets, exchange)
}
