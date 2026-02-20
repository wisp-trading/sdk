package extensions

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// Type alias for kline storage
type assetKlines map[portfolio.Pair]marketTypes.KlineMap

// klineExtension stores kline data
type klineExtension struct {
	klines *atomic.Value // assetKlines
	mu     sync.RWMutex

	// Dependency injected at construction
	onUpdateMetadata func(marketTypes.UpdateKey)
}

// NewKlineExtension creates a new kline extension
// Optional metadata updater can be provided for tracking updates
func NewKlineExtension(metadataUpdater func(marketTypes.UpdateKey)) marketTypes.KlineStoreExtension {
	ext := &klineExtension{
		klines:           &atomic.Value{},
		onUpdateMetadata: metadataUpdater,
	}
	ext.klines.Store(make(assetKlines))
	return ext
}

// Helper methods to get typed data
func (e *klineExtension) getKlines() assetKlines {
	if v := e.klines.Load(); v != nil {
		return v.(assetKlines)
	}
	return make(assetKlines)
}

func (e *klineExtension) UpdateKline(asset portfolio.Pair, exchangeName connector.ExchangeName, kline connector.Kline) {
	e.mu.Lock()

	current := e.getKlines()
	updated := make(assetKlines, len(current))
	for k, v := range current {
		updated[k] = v
	}

	if updated[asset] == nil {
		updated[asset] = make(marketTypes.KlineMap)
	}

	assetKlineMap := make(marketTypes.KlineMap, len(updated[asset]))
	for k, v := range updated[asset] {
		assetKlineMap[k] = v
	}

	if assetKlineMap[exchangeName] == nil {
		assetKlineMap[exchangeName] = make(map[string][]connector.Kline)
	}

	exchangeKlines := make(map[string][]connector.Kline, len(assetKlineMap[exchangeName]))
	for k, v := range assetKlineMap[exchangeName] {
		exchangeKlines[k] = v
	}

	// Add or update kline for this interval
	interval := kline.Interval
	existingKlines := exchangeKlines[interval]

	// Check if kline already exists (by open time)
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

	exchangeKlines[interval] = existingKlines
	assetKlineMap[exchangeName] = exchangeKlines
	updated[asset] = assetKlineMap

	e.klines.Store(updated)
	e.mu.Unlock()

	// Trigger metadata update callback
	if e.onUpdateMetadata != nil {
		e.onUpdateMetadata(marketTypes.UpdateKey{
			DataType: marketTypes.DataKeyKlines,
			Pair:     asset,
			Exchange: exchangeName,
		})
	}
}

func (e *klineExtension) GetKlines(asset portfolio.Pair, exchangeName connector.ExchangeName, interval string, limit int) []connector.Kline {
	current := e.getKlines()
	if klineMap, ok := current[asset]; ok {
		if exchangeKlines, ok := klineMap[exchangeName]; ok {
			if klines, ok := exchangeKlines[interval]; ok {
				if limit > 0 && len(klines) > limit {
					return klines[len(klines)-limit:]
				}
				return klines
			}
		}
	}
	return []connector.Kline{}
}

func (e *klineExtension) GetKlinesSince(asset portfolio.Pair, exchangeName connector.ExchangeName, interval string, since time.Time) []connector.Kline {
	current := e.getKlines()
	if klineMap, ok := current[asset]; ok {
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
