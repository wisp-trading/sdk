package extensions

import (
	"sync"
	"sync/atomic"

	marketTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"

	domainTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
)

// Type aliases for funding rate storage
type assetFundingRates map[portfolio.Pair]map[connector.ExchangeName]perpConn.FundingRate
type assetHistoricalFunding map[portfolio.Pair]map[connector.ExchangeName][]perpConn.HistoricalFundingRate

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
	rate perpConn.FundingRate,
) {
	f.mu.Lock()
	defer f.mu.Unlock()

	current := f.getFundingRates()
	updated := make(assetFundingRates, len(current)+1)

	for k, v := range current {
		updated[k] = v
	}

	if updated[asset] == nil {
		updated[asset] = make(map[connector.ExchangeName]perpConn.FundingRate)
	}
	updated[asset][exchange] = rate

	f.fundingRates.Store(updated)
}

// UpdateFundingRates updates multiple funding rates for a specific exchange
func (f *FundingRateExtension) UpdateFundingRates(
	exchange connector.ExchangeName,
	rates map[portfolio.Pair]perpConn.FundingRate,
) {
	f.mu.Lock()
	defer f.mu.Unlock()

	current := f.getFundingRates()
	updated := make(assetFundingRates, len(current)+len(rates))

	for k, v := range current {
		updated[k] = v
	}

	for asset, rate := range rates {
		if updated[asset] == nil {
			updated[asset] = make(map[connector.ExchangeName]perpConn.FundingRate)
		}
		updated[asset][exchange] = rate
	}

	f.fundingRates.Store(updated)
}

// GetFundingRate retrieves a funding rate for a specific asset and exchange
func (f *FundingRateExtension) GetFundingRate(
	asset portfolio.Pair,
	exchange connector.ExchangeName,
) *perpConn.FundingRate {
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
) map[connector.ExchangeName]perpConn.FundingRate {
	f.mu.RLock()
	defer f.mu.RUnlock()

	rates := f.getFundingRates()
	if exchanges, ok := rates[asset]; ok {
		result := make(map[connector.ExchangeName]perpConn.FundingRate, len(exchanges))
		for k, v := range exchanges {
			result[k] = v
		}
		return result
	}
	return make(map[connector.ExchangeName]perpConn.FundingRate)
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
	rates []perpConn.HistoricalFundingRate,
) {
	f.mu.Lock()
	defer f.mu.Unlock()

	current := f.getHistoricalFundingRates()
	updated := make(assetHistoricalFunding, len(current)+1)

	for k, v := range current {
		updated[k] = v
	}

	if updated[asset] == nil {
		updated[asset] = make(map[connector.ExchangeName][]perpConn.HistoricalFundingRate)
	}
	updated[asset][exchange] = rates

	f.historicalFundingRates.Store(updated)
}

// GetHistoricalFundingRates retrieves historical funding rates
func (f *FundingRateExtension) GetHistoricalFundingRates(
	asset portfolio.Pair,
	exchange connector.ExchangeName,
) []perpConn.HistoricalFundingRate {
	f.mu.RLock()
	defer f.mu.RUnlock()

	rates := f.getHistoricalFundingRates()
	if exchanges, ok := rates[asset]; ok {
		if rateList, ok := exchanges[exchange]; ok {
			return rateList
		}
	}
	return []perpConn.HistoricalFundingRate{}
}

// GetHistoricalFundingRatesForAsset retrieves all historical funding rates for an asset
func (f *FundingRateExtension) GetHistoricalFundingRatesForAsset(
	asset portfolio.Pair,
) map[connector.ExchangeName][]perpConn.HistoricalFundingRate {
	f.mu.RLock()
	defer f.mu.RUnlock()

	rates := f.getHistoricalFundingRates()
	if exchanges, ok := rates[asset]; ok {
		result := make(map[connector.ExchangeName][]perpConn.HistoricalFundingRate, len(exchanges))
		for k, v := range exchanges {
			result[k] = v
		}
		return result
	}
	return make(map[connector.ExchangeName][]perpConn.HistoricalFundingRate)
}

var _ domainTypes.FundingRateStoreExtension = (*FundingRateExtension)(nil)
var _ marketTypes.StoreExtension = (*FundingRateExtension)(nil)
