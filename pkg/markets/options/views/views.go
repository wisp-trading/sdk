package views

import (
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/logging"
)

type view struct {
	store  optionsTypes.OptionsStore
	logger logging.ApplicationLogger
}

// NewView creates a new options view that provides read-only access to market data
func NewView(
	store optionsTypes.OptionsStore,
	logger logging.ApplicationLogger,
) optionsTypes.OptionsView {
	return &view{
		store:  store,
		logger: logger,
	}
}

// GetMarkPrice returns the mark price for an option contract
func (v *view) GetMarkPrice(contract optionsTypes.OptionContract) float64 {
	return v.store.GetMarkPrice(contract)
}

// GetUnderlyingPrice returns the underlying asset price for an option contract
func (v *view) GetUnderlyingPrice(contract optionsTypes.OptionContract) float64 {
	return v.store.GetUnderlyingPrice(contract)
}

// GetGreeks returns the Greeks (sensitivities) for an option contract
func (v *view) GetGreeks(contract optionsTypes.OptionContract) optionsTypes.Greeks {
	return v.store.GetGreeks(contract)
}

// GetIV returns the implied volatility for an option contract
func (v *view) GetIV(contract optionsTypes.OptionContract) float64 {
	return v.store.GetIV(contract)
}

// GetPosition returns the position for an option contract
func (v *view) GetPosition(contract optionsTypes.OptionContract) *optionsTypes.Position {
	return v.store.GetPosition(contract)
}

// GetAllPositions returns all open positions
func (v *view) GetAllPositions() []optionsTypes.Position {
	return v.store.GetAllPositions()
}

// GetPortfolioGreeks returns the aggregated Greeks across all positions
func (v *view) GetPortfolioGreeks() optionsTypes.Greeks {
	return v.store.GetPortfolioGreeks()
}

var _ optionsTypes.OptionsView = (*view)(nil)
