package extensions

import (
	"sync"
	"time"

	marketTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

type klineExtension struct {
	klines map[portfolio.Pair]marketTypes.KlineMap
	mu     sync.RWMutex
}

func NewKlineExtension() marketTypes.KlineStoreExtension {
	return &klineExtension{
		klines: make(map[portfolio.Pair]marketTypes.KlineMap),
	}
}

func (e *klineExtension) UpdateKline(asset portfolio.Pair, exchangeName connector.ExchangeName, kline connector.Kline) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.klines[asset] == nil {
		e.klines[asset] = make(marketTypes.KlineMap)
	}
	if e.klines[asset][exchangeName] == nil {
		e.klines[asset][exchangeName] = make(map[string][]connector.Kline)
	}

	interval := kline.Interval
	existingKlines := e.klines[asset][exchangeName][interval]

	// Update or append
	found := false
	for i, existing := range existingKlines {
		if existing.OpenTime.Equal(kline.OpenTime) {
			existingKlines[i] = kline
			found = true
			break
		}
	}

	if !found {
		existingKlines = append(existingKlines, kline)
	}

	e.klines[asset][exchangeName][interval] = existingKlines
}

func (e *klineExtension) GetKlines(asset portfolio.Pair, exchangeName connector.ExchangeName, interval string, limit int) []connector.Kline {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if klineMap, ok := e.klines[asset]; ok {
		if exchangeKlines, ok := klineMap[exchangeName]; ok {
			if klines, ok := exchangeKlines[interval]; ok {
				if limit > 0 && len(klines) > limit {
					// Return a copy to prevent external modification
					return append([]connector.Kline(nil), klines[len(klines)-limit:]...)
				}
				return append([]connector.Kline(nil), klines...)
			}
		}
	}
	return []connector.Kline{}
}

func (e *klineExtension) GetKlinesSince(asset portfolio.Pair, exchangeName connector.ExchangeName, interval string, since time.Time) []connector.Kline {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if klineMap, ok := e.klines[asset]; ok {
		if exchangeKlines, ok := klineMap[exchangeName]; ok {
			if klines, ok := exchangeKlines[interval]; ok {
				result := make([]connector.Kline, 0)
				for _, kline := range klines {
					if kline.OpenTime.After(since) || kline.OpenTime.Equal(since) {
						result = append(result, kline)
					}
				}
				return result
			}
		}
	}
	return []connector.Kline{}
}
