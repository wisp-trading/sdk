package types

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
)

// OptionsStore manages options market data, positions, and Greeks
type OptionsStore interface {
	market.MarketStore

	// Positions
	GetPosition(contract OptionContract) *Position
	SetPosition(contract OptionContract, position Position)

	// Market data (float64 for hot path)
	GetMarkPrice(contract OptionContract) float64
	SetMarkPrice(contract OptionContract, price float64)

	GetUnderlyingPrice(contract OptionContract) float64
	SetUnderlyingPrice(contract OptionContract, price float64)

	// Greeks
	GetGreeks(contract OptionContract) Greeks
	SetGreeks(contract OptionContract, greeks Greeks)

	// IV
	GetIV(contract OptionContract) float64
	SetIV(contract OptionContract, iv float64)

	// Portfolio aggregates
	GetPortfolioGreeks() Greeks
	GetAllPositions() []Position
}

// Position represents a holder's position in an option contract
type Position struct {
	Contract OptionContract
	Quantity float64 // Number of contracts held
	EntryPrice float64 // Average entry price
	Unrealized float64 // Unrealized P&L
}
