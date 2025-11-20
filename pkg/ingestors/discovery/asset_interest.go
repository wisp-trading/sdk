package discovery

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

type assetInterestService struct {
	strategies []strategy.Strategy
}

// NewAssetInterestService creates a new asset interest service from strategies
func NewAssetInterestService(strategies []strategy.Strategy) ingestors.AssetInterest {
	return &assetInterestService{
		strategies: strategies,
	}
}

func (ais *assetInterestService) GetRequiredAssets() []portfolio.Asset {
	assetSet := make(map[string]portfolio.Asset)

	// Collect assets from all enabled strategies
	for _, strat := range ais.strategies {
		if !strat.IsEnabled() {
			continue
		}

		// Get signals from strategy which contain the assets
		signals, err := strat.GetSignals()
		if err != nil {
			continue
		}

		for _, sig := range signals {
			if sig == nil {
				continue
			}

			// Extract assets from signal actions
			for _, action := range sig.Actions {
				assetSet[action.Asset.Symbol()] = action.Asset
			}
		}
	}

	// Convert to slice
	assets := make([]portfolio.Asset, 0, len(assetSet))
	for _, asset := range assetSet {
		assets = append(assets, asset)
	}

	return assets
}

func (ais *assetInterestService) GetAssetRequirements() []ingestors.AssetRequirement {
	assetMap := make(map[string]map[connector.Instrument]bool)

	// Collect asset requirements from all enabled strategies
	for _, strat := range ais.strategies {
		if !strat.IsEnabled() {
			continue
		}

		signals, err := strat.GetSignals()
		if err != nil {
			continue
		}

		for _, sig := range signals {
			if sig == nil {
				continue
			}

			// Extract assets and default to perpetual instrument type
			for _, action := range sig.Actions {
				symbol := action.Asset.Symbol()
				if assetMap[symbol] == nil {
					assetMap[symbol] = make(map[connector.Instrument]bool)
				}

				// Default to both spot and perpetual for now
				// This can be enhanced based on strategy requirements
				assetMap[symbol][connector.TypeSpot] = true
				assetMap[symbol][connector.TypePerpetual] = true
			}
		}
	}

	// Convert to slice of requirements
	requirements := make([]ingestors.AssetRequirement, 0, len(assetMap))
	for symbol, instruments := range assetMap {
		instrumentList := make([]connector.Instrument, 0, len(instruments))
		for inst := range instruments {
			instrumentList = append(instrumentList, inst)
		}

		requirements = append(requirements, ingestors.AssetRequirement{
			Asset:       portfolio.NewAsset(symbol),
			Instruments: instrumentList,
		})
	}

	return requirements
}

func (ais *assetInterestService) IsAssetRequired(symbol string) bool {
	for _, strat := range ais.strategies {
		if !strat.IsEnabled() {
			continue
		}

		signals, err := strat.GetSignals()
		if err != nil {
			continue
		}

		for _, sig := range signals {
			if sig == nil {
				continue
			}

			for _, action := range sig.Actions {
				if action.Asset.Symbol() == symbol {
					return true
				}
			}
		}
	}

	return false
}

// GetInstrumentTypes returns the instrument types needed for a specific asset
func (ais *assetInterestService) GetInstrumentTypes(asset portfolio.Asset) []connector.Instrument {
	requirements := ais.GetAssetRequirements()

	for _, req := range requirements {
		if req.Asset.Symbol() == asset.Symbol() {
			return req.Instruments
		}
	}

	// Default to both spot and perpetual if asset is found in signals
	if ais.IsAssetRequired(asset.Symbol()) {
		return []connector.Instrument{connector.TypeSpot, connector.TypePerpetual}
	}

	return []connector.Instrument{}
}
