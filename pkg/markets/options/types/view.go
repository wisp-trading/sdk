package types

// OptionsView provides read-only access to options market data
type OptionsView interface {
	// Get market data for a contract
	GetMarkPrice(contract OptionContract) float64
	GetUnderlyingPrice(contract OptionContract) float64
	GetGreeks(contract OptionContract) Greeks
	GetIV(contract OptionContract) float64

	// Get position data
	GetPosition(contract OptionContract) *Position
	GetAllPositions() []Position

	// Get portfolio Greeks
	GetPortfolioGreeks() Greeks
}
