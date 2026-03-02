package market

// MarketType identifies the type of market
type MarketType string

const (
	MarketTypeUnknown    MarketType = "unknown"
	MarketTypeSpot       MarketType = "spot"
	MarketTypePerp       MarketType = "perp"
	MarketTypePrediction MarketType = "prediction"
)

// MarketRegistry provides access to all registered market stores
type MarketRegistry interface {
	// Get returns the store for a specific market type, or nil if not registered
	Get(marketType MarketType) MarketStore

	// GetAll returns all registered market stores
	GetAll() map[MarketType]MarketStore

	// Register adds a market store to the registry
	Register(store MarketStore)

	// Types returns all registered market types
	Types() []MarketType
}
