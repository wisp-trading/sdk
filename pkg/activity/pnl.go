package activity

import (
	kronosActivity "github.com/backtesting-org/kronos-sdk/pkg/types/kronos/activity"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

// pnl provides PNL calculation functionality
// TODO: Implement full PNL calculation logic with market data integration
type pnl struct {
	positions kronosActivity.Positions
	trades    kronosActivity.Trades
}

// NewPNL creates a new PNL calculator
func NewPNL(positions kronosActivity.Positions, trades kronosActivity.Trades) kronosActivity.PNL {
	return &pnl{
		positions: positions,
		trades:    trades,
	}
}

// GetRealizedPNL returns the realized PNL for a strategy
// TODO: Implement realized PNL calculation from trades
func (p *pnl) GetRealizedPNL(strategyName strategy.StrategyName) numerical.Decimal {
	return numerical.Zero()
}

// GetRealizedPNLByAsset returns the realized PNL for an asset
// TODO: Implement realized PNL calculation by asset
func (p *pnl) GetRealizedPNLByAsset(asset portfolio.Asset) numerical.Decimal {
	return numerical.Zero()
}

// GetTotalRealizedPNL returns the total realized PNL across all strategies
// TODO: Implement total realized PNL calculation
func (p *pnl) GetTotalRealizedPNL() numerical.Decimal {
	return numerical.Zero()
}

// GetUnrealizedPNL returns the unrealized PNL for a strategy
// TODO: Implement unrealized PNL calculation (requires market data)
func (p *pnl) GetUnrealizedPNL(strategyName strategy.StrategyName) (numerical.Decimal, error) {
	return numerical.Zero(), nil
}

// GetTotalUnrealizedPNL returns the total unrealized PNL across all strategies
// TODO: Implement total unrealized PNL calculation (requires market data)
func (p *pnl) GetTotalUnrealizedPNL() (numerical.Decimal, error) {
	return numerical.Zero(), nil
}

// GetTotalPNL returns the total PNL (realized + unrealized)
// TODO: Implement total PNL calculation
func (p *pnl) GetTotalPNL() (numerical.Decimal, error) {
	return numerical.Zero(), nil
}

// GetTotalFees returns the total fees paid across all trades
// TODO: Implement total fees calculation
func (p *pnl) GetTotalFees() numerical.Decimal {
	return numerical.Zero()
}

// GetFeesByStrategy returns the total fees paid for a strategy
// TODO: Implement fees calculation per strategy
func (p *pnl) GetFeesByStrategy(strategyName strategy.StrategyName) numerical.Decimal {
	return numerical.Zero()
}

var _ kronosActivity.PNL = (*pnl)(nil)
