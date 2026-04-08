package analytics

import (
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/logging"
)

type service struct {
	pnl    optionsTypes.OptionsPNL
	store  optionsTypes.OptionsStore
	logger logging.ApplicationLogger
}

// NewAnalyticsService creates a new analytics service for options
func NewAnalyticsService(
	pnl optionsTypes.OptionsPNL,
	store optionsTypes.OptionsStore,
	logger logging.ApplicationLogger,
) optionsTypes.OptionsAnalytics {
	return &service{
		pnl:    pnl,
		store:  store,
		logger: logger,
	}
}

// GetDeltaExposure returns the portfolio delta exposure
// Delta represents sensitivity to underlying asset price changes
func (s *service) GetDeltaExposure() float64 {
	return s.pnl.CalculateDeltaExposure()
}

// GetGammaExposure returns the portfolio gamma exposure
// Gamma represents the rate of change of delta
func (s *service) GetGammaExposure() float64 {
	return s.pnl.CalculateGammaExposure()
}

// GetThetaExposure returns the portfolio theta exposure
// Theta represents daily time decay
func (s *service) GetThetaExposure() float64 {
	return s.pnl.CalculateThetaDecay()
}

// GetVegaExposure returns the portfolio vega exposure
// Vega represents sensitivity to implied volatility changes (per 1% IV move)
func (s *service) GetVegaExposure() float64 {
	return s.pnl.CalculateVegaExposure()
}

// GetPortfolioGreeks returns the aggregated Greeks across all positions
func (s *service) GetPortfolioGreeks() optionsTypes.Greeks {
	return s.pnl.GetPortfolioGreeks()
}

// GetDailyThetaDecay returns the estimated daily theta decay
// Useful for understanding daily P&L impact from time decay
func (s *service) GetDailyThetaDecay() float64 {
	// Theta is already in daily terms (∂P/∂t where t is in days)
	return s.pnl.CalculateThetaDecay()
}

// GetIVSensitivity returns the portfolio sensitivity to implied volatility changes
// Per 1% change in IV
func (s *service) GetIVSensitivity() float64 {
	// Vega is the sensitivity per 1% IV change
	return s.pnl.CalculateVegaExposure()
}

var _ optionsTypes.OptionsAnalytics = (*service)(nil)
