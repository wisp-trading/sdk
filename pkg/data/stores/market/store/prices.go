package store

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

func (ds *dataStore) UpdatePairPrice(asset portfolio.Pair, exchangeName connector.ExchangeName, price connector.Price) {
	ds.mu.Lock()

	if ds.prices[asset] == nil {
		ds.prices[asset] = make(marketTypes.PriceMap)
	}
	ds.prices[asset][exchangeName] = price

	ds.mu.Unlock()

	ds.UpdateLastUpdated(marketTypes.UpdateKey{
		DataType: marketTypes.DataKeyPairPrice,
		Pair:     asset,
		Exchange: exchangeName,
	})
}

func (ds *dataStore) UpdatePairPrices(asset portfolio.Pair, prices marketTypes.PriceMap) {
	ds.mu.Lock()

	if ds.prices[asset] == nil {
		ds.prices[asset] = make(marketTypes.PriceMap)
	}

	// Collect exchanges to update after releasing the lock
	exchangesToUpdate := make([]connector.ExchangeName, 0, len(prices))

	for exchangeName, price := range prices {
		ds.prices[asset][exchangeName] = price
		exchangesToUpdate = append(exchangesToUpdate, exchangeName)
	}

	ds.mu.Unlock()

	// Update timestamps after releasing the lock to avoid deadlock
	for _, exchangeName := range exchangesToUpdate {
		ds.UpdateLastUpdated(marketTypes.UpdateKey{
			DataType: marketTypes.DataKeyPairPrice,
			Pair:     asset,
			Exchange: exchangeName,
		})
	}
}

func (ds *dataStore) GetPairPrice(asset portfolio.Pair, exchangeName connector.ExchangeName) *connector.Price {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	if priceMap, ok := ds.prices[asset]; ok {
		if price, ok := priceMap[exchangeName]; ok {
			priceCopy := price
			return &priceCopy
		}
	}
	return nil
}

func (ds *dataStore) GetPairPrices(asset portfolio.Pair) marketTypes.PriceMap {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	if prices, ok := ds.prices[asset]; ok {
		// Return shallow copy to prevent external mutation
		result := make(marketTypes.PriceMap, len(prices))
		for k, v := range prices {
			result[k] = v
		}
		return result
	}
	return make(marketTypes.PriceMap)
}
