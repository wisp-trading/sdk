package extensions

import (
	"sync"
	"sync/atomic"

	"github.com/wisp-trading/sdk/pkg/types/connector/perp"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// Type aliases for funding rate storage
type assetFundingRates map[portfolio.Pair]map[connector.ExchangeName]perp.FundingRate
type assetHistoricalFunding map[portfolio.Pair]map[connector.ExchangeName][]perp.HistoricalFundingRate

// FundingRateExtension stores perp-specific funding rate data
type FundingRateExtension struct {
	fundingRates           *atomic.Value // assetFundingRates
	historicalFundingRates *atomic.Value // assetHistoricalFunding
	mu                     sync.RWMutex
}

func NewFundingRateExtension() *FundingRateExtension {
	ext := &FundingRateExtension{
		fundingRates:           &atomic.Value{},
		historicalFundingRates: &atomic.Value{},
	}

	ext.fundingRates.Store(make(assetFundingRates))
	ext.historicalFundingRates.Store(make(assetHistoricalFunding))

	return ext
}

// Helper methods to get typed data
func (f *FundingRateExtension) getFundingRates() assetFundingRates {
	if v := f.fundingRates.Load(); v != nil {
		return v.(assetFundingRates)
	}
	return make(assetFundingRates)
}

func (f *FundingRateExtension) getHistoricalFundingRates() assetHistoricalFunding {
	if v := f.historicalFundingRates.Load(); v != nil {
		return v.(assetHistoricalFunding)
	}
	return make(assetHistoricalFunding)
}

// UpdateFundingRate updates a funding rate for a specific asset and exchange
func (f *FundingRateExtension) UpdateFundingRate(
	asset portfolio.Pair,
	exchange connector.ExchangeName,
	rate perp.FundingRate,
) {
	f.mu.Lock()
	defer f.mu.Unlock()

	current := f.getFundingRates()
	updated := make(assetFundingRates, len(current)+1)

	// Copy existing data
	for k, v := range current {
		updated[k] = v
	}

	// Update specific entry
	if updated[asset] == nil {
		updated[asset] = make(map[connector.ExchangeName]perp.FundingRate)
	}
	updated[asset][exchange] = rate

	f.fundingRates.Store(updated)
}

// UpdateFundingRates updates multiple funding rates for a specific exchange
func (f *FundingRateExtension) UpdateFundingRates(
	exchange connector.ExchangeName,
	rates map[portfolio.Pair]perp.FundingRate,
) {
	f.mu.Lock()
	defer f.mu.Unlock()

	current := f.getFundingRates()
	updated := make(assetFundingRates, len(current)+len(rates))

	// Copy existing data
	for k, v := range current {
		updated[k] = v
	}

	// Update entries for this exchange
	for asset, rate := range rates {
		if updated[asset] == nil {
			updated[asset] = make(map[connector.ExchangeName]perp.FundingRate)
		}
		updated[asset][exchange] = rate
	}

	f.fundingRates.Store(updated)
}

// GetFundingRate retrieves a funding rate for a specific asset and exchange
func (f *FundingRateExtension) GetFundingRate(
	asset portfolio.Pair,
	exchange connector.ExchangeName,
) *perp.FundingRate {
	f.mu.RLock()
	defer f.mu.RUnlock()

	rates := f.getFundingRates()
	if exchanges, ok := rates[asset]; ok {
		if rate, ok := exchanges[exchange]; ok {
			return &rate
		}
	}
	return nil
}

// GetFundingRatesForAsset retrieves all funding rates for a specific asset
func (f *FundingRateExtension) GetFundingRatesForAsset(
	asset portfolio.Pair,
) map[connector.ExchangeName]perp.FundingRate {
	f.mu.RLock()
	defer f.mu.RUnlock()

	rates := f.getFundingRates()
	if exchanges, ok := rates[asset]; ok {
		// Return a copy
		result := make(map[connector.ExchangeName]perp.FundingRate, len(exchanges))
		for k, v := range exchanges {
			result[k] = v
		}
		return result
	}
	return make(map[connector.ExchangeName]perp.FundingRate)
}

// GetAllAssetsWithFundingRates returns all assets that have funding rates
func (f *FundingRateExtension) GetAllAssetsWithFundingRates() []portfolio.Pair {
	f.mu.RLock()
	defer f.mu.RUnlock()

	rates := f.getFundingRates()
	assets := make([]portfolio.Pair, 0, len(rates))
	for asset := range rates {
		assets = append(assets, asset)
	}
	return assets
}

// UpdateHistoricalFundingRates updates historical funding rates
func (f *FundingRateExtension) UpdateHistoricalFundingRates(
	asset portfolio.Pair,
	exchange connector.ExchangeName,
	rates []perp.HistoricalFundingRate,
) {
	f.mu.Lock()
	defer f.mu.Unlock()

	current := f.getHistoricalFundingRates()
	updated := make(assetHistoricalFunding, len(current)+1)

	// Copy existing data
	for k, v := range current {
		updated[k] = v
	}

	// Update specific entry
	if updated[asset] == nil {
		updated[asset] = make(map[connector.ExchangeName][]perp.HistoricalFundingRate)
	}
	updated[asset][exchange] = rates

	f.historicalFundingRates.Store(updated)
}

// GetHistoricalFundingRates retrieves historical funding rates
func (f *FundingRateExtension) GetHistoricalFundingRates(
	asset portfolio.Pair,
	exchange connector.ExchangeName,
) []perp.HistoricalFundingRate {
	f.mu.RLock()
	defer f.mu.RUnlock()

	rates := f.getHistoricalFundingRates()
	if exchanges, ok := rates[asset]; ok {
		if rateList, ok := exchanges[exchange]; ok {
			return rateList
		}
	}
	return []perp.HistoricalFundingRate{}
}

// GetHistoricalFundingRatesForAsset retrieves all historical funding rates for an asset
func (f *FundingRateExtension) GetHistoricalFundingRatesForAsset(
	asset portfolio.Pair,
) map[connector.ExchangeName][]perp.HistoricalFundingRate {
	f.mu.RLock()
	defer f.mu.RUnlock()

	rates := f.getHistoricalFundingRates()
	if exchanges, ok := rates[asset]; ok {
		// Return a copy
		result := make(map[connector.ExchangeName][]perp.HistoricalFundingRate, len(exchanges))
		for k, v := range exchanges {
			result[k] = v
		}
		return result
	}
	return make(map[connector.ExchangeName][]perp.HistoricalFundingRate)
}

var _ marketTypes.StoreExtension = (*FundingRateExtension)(nil)
