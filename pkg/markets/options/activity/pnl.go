package activity

import (
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/logging"
)

type pnlCalculator struct {
	store  optionsTypes.OptionsStore
	logger logging.ApplicationLogger
}

// NewPNLCalculator creates a new PnL calculator for options
func NewPNLCalculator(
	store optionsTypes.OptionsStore,
	logger logging.ApplicationLogger,
) optionsTypes.OptionsPNL {
	return &pnlCalculator{
		store:  store,
		logger: logger,
	}
}

// CalculateUnrealizedPnL calculates the unrealized P&L for a single option contract
// Simple approximation: P&L = Position Quantity * Greeks * Price Move
// More sophisticated models would use Black-Scholes Greeks interpolation
func (p *pnlCalculator) CalculateUnrealizedPnL(contract optionsTypes.OptionContract) float64 {
	position := p.store.GetPosition(contract)
	if position == nil || position.Quantity == 0 {
		return 0
	}

	markPrice := p.store.GetMarkPrice(contract)
	if markPrice == 0 {
		return 0
	}

	// P&L = position quantity * (current price - entry price)
	pnl := position.Quantity * (markPrice - position.EntryPrice)

	return pnl
}

// CalculateDeltaExposure calculates the total delta exposure across all positions
// Delta represents the sensitivity of the option price to changes in the underlying asset price
// Portfolio Delta = Sum of (Position Quantity * Option Delta)
func (p *pnlCalculator) CalculateDeltaExposure() float64 {
	positions := p.store.GetAllPositions()
	if len(positions) == 0 {
		return 0
	}

	var totalDelta float64
	for _, pos := range positions {
		greeks := p.store.GetGreeks(pos.Contract)
		totalDelta += pos.Quantity * greeks.Delta
	}

	return totalDelta
}

// CalculateGammaExposure calculates the total gamma exposure across all positions
// Gamma represents the rate of change of delta with respect to underlying price changes
// Portfolio Gamma = Sum of (Position Quantity * Option Gamma)
func (p *pnlCalculator) CalculateGammaExposure() float64 {
	positions := p.store.GetAllPositions()
	if len(positions) == 0 {
		return 0
	}

	var totalGamma float64
	for _, pos := range positions {
		greeks := p.store.GetGreeks(pos.Contract)
		totalGamma += pos.Quantity * greeks.Gamma
	}

	return totalGamma
}

// CalculateThetaDecay calculates the daily theta decay across all positions
// Theta represents time decay - the daily change in option value due to passage of time
// Portfolio Theta = Sum of (Position Quantity * Option Theta)
func (p *pnlCalculator) CalculateThetaDecay() float64 {
	positions := p.store.GetAllPositions()
	if len(positions) == 0 {
		return 0
	}

	var totalTheta float64
	for _, pos := range positions {
		greeks := p.store.GetGreeks(pos.Contract)
		totalTheta += pos.Quantity * greeks.Theta
	}

	return totalTheta
}

// CalculateVegaExposure calculates the total vega exposure across all positions
// Vega represents the sensitivity of option price to changes in implied volatility
// Portfolio Vega = Sum of (Position Quantity * Option Vega)
// Per 1% change in IV
func (p *pnlCalculator) CalculateVegaExposure() float64 {
	positions := p.store.GetAllPositions()
	if len(positions) == 0 {
		return 0
	}

	var totalVega float64
	for _, pos := range positions {
		greeks := p.store.GetGreeks(pos.Contract)
		totalVega += pos.Quantity * greeks.Vega
	}

	return totalVega
}

// GetPortfolioGreeks returns the aggregated Greeks across all positions
func (p *pnlCalculator) GetPortfolioGreeks() optionsTypes.Greeks {
	return p.store.GetPortfolioGreeks()
}

var _ optionsTypes.OptionsPNL = (*pnlCalculator)(nil)
