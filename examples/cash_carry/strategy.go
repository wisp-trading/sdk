package cash_carry

import (
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/kronos"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CashCarryStrategy struct {
	*strategy.BaseStrategy
	k *kronos.Kronos

	// Simple config
	minFundingRate decimal.Decimal
}

func NewCashCarryStrategy(k *kronos.Kronos) *CashCarryStrategy {
	return &CashCarryStrategy{
		k:              k,
		minFundingRate: decimal.RequireFromString("0.001"), // 0.1% minimum (10 bps)
	}
}

func (ccs *CashCarryStrategy) GetSignals() ([]*strategy.Signal, error) {
	ccs.k.Log().Info("🔍 Scanning for cash carry opportunities...")

	// Get all assets with funding rate data using Kronos Market service
	assets := ccs.k.Market.GetAllAssetsWithFundingRates()
	if len(assets) == 0 {
		ccs.k.Log().Info("No assets with funding rate data available")
		return nil, nil
	}

	ccs.k.Log().Info("Checking %d assets for funding rate opportunities", len(assets))

	var signals []*strategy.Signal
	for _, asset := range assets {
		// Get funding rates across all exchanges for this asset
		fundingRates := ccs.k.Market.FundingRates(asset)

		for exchange, fundingRate := range fundingRates {
			// Check if funding rate exceeds minimum threshold
			if fundingRate.CurrentRate.GreaterThan(ccs.minFundingRate) {
				ccs.k.Log().Opportunity(
					"CashCarry",
					asset.Symbol(),
					"High funding rate on %s: %s%% (threshold: %s%%)",
					exchange,
					fundingRate.CurrentRate.Mul(decimal.NewFromInt(100)).StringFixed(4),
					ccs.minFundingRate.Mul(decimal.NewFromInt(100)).StringFixed(4),
				)

				// Get current price for context
				price, err := ccs.k.Market.Price(asset)
				if err != nil {
					ccs.k.Log().Failed("CashCarry", asset.Symbol(), "Failed to get price: %v", err)
				} else {
					ccs.k.Log().Debug(
						"CashCarry",
						asset.Symbol(),
						"Creating cash carry signal - Price: %s, Funding: %s%%, Exchange: %s",
						price.String(),
						fundingRate.CurrentRate.Mul(decimal.NewFromInt(100)).StringFixed(4),
						exchange,
					)
				}

				signal := &strategy.Signal{
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
				signals = append(signals, signal)
			}
		}
	}

	if len(signals) > 0 {
		ccs.k.Log().Success("CashCarry", "", "Found %d opportunities", len(signals))
	} else {
		ccs.k.Log().Info("No opportunities found - all funding rates below threshold")
	}

	return signals, nil
}

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
