// Package onchain defines the exchange-facing contract for AMM / on-chain venues.
// First implementation: Uniswap V3 via an EVM RPC (any chain id, e.g. Robinhood).
package onchain

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Connector is an on-chain venue: swaps (as market orders), balances, token registry.
//
// PlaceMarketOrder convention (exact-in AMM):
//   - BUY  — quantity is quote asset spent (e.g. WETH in) → receive base
//   - SELL — quantity is base asset sold (token in) → receive quote
//
// Limit orders are not supported by the pilot UniV3 path (return error).
type Connector interface {
	connector.Connector
	connector.OrderExecutor
	connector.AccountReader

	// RegisterToken maps a symbol (ticker or 0x address) to a chain token.
	// Call before trading pairs that use that symbol.
	RegisterToken(symbol string, address string, decimals uint8) error

	// ResolveToken returns the registered address and decimals for a symbol.
	ResolveToken(symbol string) (address string, decimals uint8, ok bool)

	// ChainID is the EVM chain id this connector is bound to.
	ChainID() uint64

	// NativeWrapped is the wrapped native symbol used as default quote (e.g. "WETH").
	NativeWrapped() string

	// QuoteMarket returns an exact-in quote without sending a transaction.
	// Side/quantity follow the same convention as PlaceMarketOrder.
	QuoteMarket(pair portfolio.Pair, side connector.OrderSide, quantity numerical.Decimal) (*Quote, error)
}

// Quote is a read-only swap estimate.
type Quote struct {
	AmountIn      numerical.Decimal `json:"amount_in"`
	AmountOut     numerical.Decimal `json:"amount_out"`
	FeeTier       uint32            `json:"fee_tier"`
	PriceImpactBps int              `json:"price_impact_bps,omitempty"`
	Pool          string            `json:"pool,omitempty"`
}
