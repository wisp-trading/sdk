package options

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// OptionContract uniquely identifies an option
type OptionContract struct {
	Pair       portfolio.Pair // e.g., BTC/USDT
	Strike     float64        // e.g., 50000.0
	Expiration time.Time      // e.g., 2025-05-31
	OptionType string         // "CALL" or "PUT"
}

// Greeks represent the sensitivities of an option's price
type Greeks struct {
	Delta float64 // ∂P/∂S (price change per $1 underlying move)
	Gamma float64 // ∂²P/∂S² (delta change per $1 underlying move)
	Theta float64 // ∂P/∂t (daily time decay)
	Vega  float64 // ∂P/∂IV (price change per 1% IV move)
	Rho   float64 // ∂P/∂r (price change per 1% rate move)
}

// OptionData is the market data for a single option
// Everything needed by SDK is provided in one call
type OptionData struct {
	MarkPrice       float64   // Exchange's calculated mark price for this strike
	UnderlyingPrice float64   // The underlying spot price (BTC/USDT) used for Greeks
	IV              float64   // Implied Volatility
	Greeks          Greeks    // Pre-calculated by exchange
	BidAskSpread    float64   // Bid-ask spread
	Volume24h       float64   // 24-hour trading volume
	OpenInterest    float64   // Open interest
	Timestamp       time.Time // When this data was fetched/updated
}

// Connector represents an options exchange connection
// NOTE: Does NOT implement MarketDataReader (no FetchPrice/FetchKlines)
// Options have domain-specific data methods instead
type Connector interface {
	connector.Connector
	connector.OrderExecutor
	connector.AccountReader

	// Discovery: Get available expirations for an underlying pair
	GetExpirations(pair portfolio.Pair) ([]time.Time, error)

	// Discovery: Get available strikes for an expiration
	GetStrikes(pair portfolio.Pair, expiration time.Time) ([]float64, error)

	// Data fetching: Get mark price + Greeks for a specific option
	GetOptionData(contract OptionContract) (OptionData, error)

	// Data fetching: Get all option data for an expiration (all strikes, both call/put)
	// Returns map[strike][callOrPut]OptionData
	GetExpirationData(pair portfolio.Pair, expiration time.Time) (
		map[float64]map[string]OptionData,
		error,
	)
}
