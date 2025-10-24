package cash_carry

import (
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio/store"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CashCarryStrategy struct {
	*strategy.BaseStrategy
	assetStore store.Store
	logger     logging.ApplicationLogger

	// Simple config
	minFundingRate decimal.Decimal
}

func NewCashCarryStrategy(
	assetStore store.Store,
	logger logging.ApplicationLogger,
) *CashCarryStrategy {
	return &CashCarryStrategy{
		assetStore:     assetStore,
		logger:         logger,
		minFundingRate: decimal.RequireFromString("0.001"), // 0.1% minimum
	}
}

func (ccs *CashCarryStrategy) GetSignals() ([]*strategy.Signal, error) {
	assets := ccs.assetStore.GetAllAssetsWithFundingRates()
	if len(assets) == 0 {
		return nil, nil // No error, just no opportunities
	}

	var signals []*strategy.Signal

	for _, asset := range assets {
		fundingRates := ccs.assetStore.GetFundingRatesForAsset(asset)

		for exchange, fundingRate := range fundingRates {

			if fundingRate.CurrentRate.GreaterThan(ccs.minFundingRate) {
				signal := ccs.createSignal(asset, exchange, fundingRate.CurrentRate)
				signals = append(signals, signal)
			}
		}
	}

	return signals, nil
}

func (ccs *CashCarryStrategy) createSignal(asset portfolio.Asset, exchange connector.ExchangeName, fundingRate decimal.Decimal) *strategy.Signal {
	return &strategy.Signal{
		ID:       uuid.New(),
		Strategy: strategy.CashCarry,
		Actions: []strategy.TradeAction{
			{
				Action:   strategy.ActionBuy,
				Asset:    asset,
				Exchange: exchange,
				Quantity: decimal.NewFromInt(100),
				Price:    decimal.Zero,
			},
			{
				Action:   strategy.ActionSellShort,
				Asset:    asset,
				Exchange: exchange,
				Quantity: decimal.NewFromInt(100),
				Price:    decimal.Zero,
			},
		},
		Timestamp: time.Now(),
	}
}

// Minimal interface compliance
func (ccs *CashCarryStrategy) GetName() strategy.StrategyName {
	return strategy.CashCarry
}

func (ccs *CashCarryStrategy) GetDescription() string {
	return "Cash carry arbitrage"
}

func (ccs *CashCarryStrategy) GetRiskLevel() strategy.RiskLevel {
	return strategy.RiskLevelLow
}

func (ccs *CashCarryStrategy) GetStrategyType() strategy.StrategyType {
	return strategy.StrategyTypeCashCarry
}

var _ strategy.Strategy = (*CashCarryStrategy)(nil)
