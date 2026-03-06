package market

import "github.com/wisp-trading/sdk/pkg/types/connector"

// MarketRegistry provides access to all registered market stores
type MarketRegistry interface {
	// Get returns the store for a specific market type, or nil if not registered
	Get(marketType connector.MarketType) MarketStore

	// GetAll returns all registered market stores
	GetAll() map[connector.MarketType]MarketStore

	// Register adds a market store to the registry
	Register(store MarketStore)

	// Types returns all registered market types
	Types() []connector.MarketType
}
