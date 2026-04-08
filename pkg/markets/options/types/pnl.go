package types

// OptionsPNL calculates profit/loss using Greeks-based sensitivity analysis
type OptionsPNL interface {
	// Position P&L
	CalculateUnrealizedPnL(contract OptionContract) float64

	// Greeks sensitivities (float64 for hot path)
	CalculateDeltaExposure() float64     // ∂P/∂S - price change per $1 underlying move
	CalculateGammaExposure() float64     // ∂²P/∂S² - delta change per $1 underlying move
	CalculateThetaDecay() float64        // Daily theta decay
	CalculateVegaExposure() float64      // Price change per 1% IV move

	// Aggregate
	GetPortfolioGreeks() Greeks
}
