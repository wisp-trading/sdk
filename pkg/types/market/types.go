package market

// MarketType represents the type of market being traded
type MarketType string

const (
	// MarketTypeSpot represents spot trading markets (immediate settlement)
	MarketTypeSpot MarketType = "spot"

	// MarketTypePerp represents perpetual futures markets
	MarketTypePerp MarketType = "perpetual"

	// MarketTypePrediction represents prediction/betting markets (Polymarket, etc.)
	MarketTypePrediction MarketType = "prediction"

	// MarketTypeOptions represents options markets
	MarketTypeOptions MarketType = "options"

	// MarketTypeNFT represents NFT markets
	MarketTypeNFT MarketType = "nft"
)

// String returns the string representation of the market type
func (m MarketType) String() string {
	return string(m)
}

// IsValid checks if the market type is a valid value
func (m MarketType) IsValid() bool {
	switch m {
	case MarketTypeSpot, MarketTypePerp, MarketTypePrediction, MarketTypeOptions, MarketTypeNFT:
		return true
	default:
		return false
	}
}
