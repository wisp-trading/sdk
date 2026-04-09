package types

import (
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// MarketLoader handles background pagination and loading of markets into the store
type MarketLoader interface {
	// LoadMarkets starts background pagination of all markets matching the filter.
	// Returns immediately; actual loading happens asynchronously.
	LoadMarkets(exchange connector.ExchangeName, filter *predictionconnector.MarketsFilter) error

	// IsLoading returns true if background market loading is currently in progress for the exchange
	IsLoading(exchange connector.ExchangeName) bool

	// GetLoadProgress returns count of markets loaded so far for an exchange
	GetLoadProgress(exchange connector.ExchangeName) int
}
