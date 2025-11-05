package main

import (
	"github.com/backtesting-org/kronos-sdk/pkg/kronos"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	"github.com/shopspring/decimal"
)

// CashCarryStrategy implements cash and carry arbitrage using Kronos SDK
type CashCarryStrategy struct {
	*strategy.BaseStrategy
	k      *kronos.Kronos
	config CashCarryConfig
}

// CashCarryConfig holds cash and carry strategy parameters
type CashCarryConfig struct {
	MinFundingRate decimal.Decimal
	OrderSizeUSD   decimal.Decimal // Order size in USD
}

// NewCashCarryStrategy creates a new cash carry strategy instance
func NewCashCarryStrategy(k *kronos.Kronos, config CashCarryConfig) *CashCarryStrategy {
	base := strategy.NewBaseStrategy(
		strategy.StrategyName("Cash and Carry"),
		"Cash and carry arbitrage strategy exploiting funding rates",
		strategy.RiskLevelLow,
		strategy.StrategyTypeCashCarry,
	)

	return &CashCarryStrategy{
		BaseStrategy: base,
		k:            k,
		config:       config,
	}
}

// GetSignals generates trading signals for cash and carry strategy
func (ccs *CashCarryStrategy) GetSignals() ([]*strategy.Signal, error) {
	if !ccs.IsEnabled() {
		return nil, nil
	}

	ccs.k.Log().Info("CashCarry", "", "Scanning for cash carry opportunities...")

	// Get all assets with funding rate data using Kronos Market service
	assets := ccs.k.Market.GetAllAssetsWithFundingRates()
	if len(assets) == 0 {
		ccs.k.Log().Info("CashCarry", "", "No assets with funding rate data available")
		return nil, nil
	}

	ccs.k.Log().Info("CashCarry", "", "Checking %d assets for funding rate opportunities", len(assets))

	var signals []*strategy.Signal
	for _, asset := range assets {
		// Get funding rates across all exchanges for this asset
		fundingRates := ccs.k.Market.FundingRates(asset)

		for exchange, fundingRate := range fundingRates {
			// Check if funding rate exceeds minimum threshold
			if fundingRate.CurrentRate.GreaterThan(ccs.config.MinFundingRate) {
				// Get current price to calculate quantity
				price, err := ccs.k.Market.Price(asset)
				if err != nil {
					ccs.k.Log().Debug("CashCarry", asset.Symbol(), "Failed to get price: %v", err)
					continue
				}

				// Calculate quantity based on USD order size
				quantity := ccs.config.OrderSizeUSD.Div(price)

				ccs.k.Log().Opportunity(
					"CashCarry",
					asset.Symbol(),
					"High funding rate on %s: %s%% (threshold: %s%%), quantity: %s",
					exchange,
					fundingRate.CurrentRate.Mul(decimal.NewFromInt(100)).String(),
					ccs.config.MinFundingRate.Mul(decimal.NewFromInt(100)).String(),
					quantity.String(),
				)

				signal := ccs.k.Signal(ccs.GetName()).
					Buy(asset, exchange, quantity).
					SellShort(asset, exchange, quantity).
					Build()

				signals = append(signals, signal)
			}
		}
	}

	if len(signals) > 0 {
		ccs.k.Log().Success("CashCarry", "", "Found %d opportunities", len(signals))
	}

	return signals, nil
}

// NewStrategy creates a new strategy instance for plugin loading
// This is called by the plugin manager to extract metadata
func NewStrategy() strategy.Strategy {
	// Create with nil Kronos and default config for metadata extraction
	return NewCashCarryStrategy(
		nil, // Kronos not needed for metadata
		CashCarryConfig{
			MinFundingRate: decimal.NewFromFloat(0.0001), // 0.01% funding rate threshold
			OrderSizeUSD:   decimal.NewFromFloat(100),    // $100 per order
		},
	)
}

// Plugin export - required for Go plugin system
var Strategy strategy.Strategy = NewStrategy()
