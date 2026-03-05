package connector

// AccountReader provides prediction-specific account and position information
// Uses connector.AccountReader for standard balance/trading history
type AccountReader interface {
	// GetPositions Get all prediction positions
	GetPositions() ([]Position, error)

	// GetPositionsByMarket Get positions for a specific market
	GetPositionsByMarket(marketID string) ([]Position, error)

	// GetTokensToRedeem returns the balances of outcome tokens for a market, which can be redeemed after resolution
	GetTokensToRedeem(market Market) ([]Balance, error)
}
