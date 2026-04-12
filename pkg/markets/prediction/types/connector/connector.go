package connector

import (
	"math/big"

	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// MarketsFilter specifies criteria for filtering markets returned by Markets().
type MarketsFilter struct {
	// Pagination
	Limit  *int
	Offset *int

	// Volume filters
	MinVolume string
	MaxVolume string

	// Liquidity filters
	MinLiquidity string
	MaxLiquidity string

	// Date range filters (ISO 8601 strings)
	MinEndDate string
	MaxEndDate string

	// Status filter
	Active *bool
	Closed *bool
}

// Connector represents a prediction market exchange connection
// Follows the same pattern as spot.Connector and perp.Connector
// Uses standard connector.OrderExecutor and connector.AccountReader for trading
type Connector interface {
	connector.Connector

	OrderExecutor
	AccountReader

	Redeem(market Market) (string, error)
	GetMarket(slug string) (Market, error)
	GetRecurringMarket(slug string, interval RecurrenceInterval) (Market, error)
	GetOutcome(marketID, outcomeID string) Outcome

	// Markets returns markets available on this exchange, optionally filtered.
	// Pass nil filters to get all active markets.
	Markets(filter *MarketsFilter) ([]Market, error)

	// FetchOrderBooksForMarket fetches all orderbooks for a market in a single batch request.
	FetchOrderBooksForMarket(market Market) (map[string]*OrderBook, error)

	// SplitPosition deposits amountUSDC (6 decimal units) into the CTF contract
	// and mints 1 YES + 1 NO token per unit. $1.00 = big.NewInt(1_000_000).
	// Returns the tx hash immediately and a ready channel that closes once the tx
	// is mined and the CLOB balance cache is refreshed. Callers MUST drain ready
	// before placing SELL orders — the CLOB checks on-chain balance at submission.
	SplitPosition(market Market, amountUSDC *big.Int) (txHash string, ready <-chan error, err error)

	// MergePositions burns amountUSDC worth of YES+NO tokens and returns USDC.
	MergePositions(market Market, amountUSDC *big.Int) (txHash string, err error)

	// GetLockedPositions returns all CTF ERC-1155 conditional token positions
	// currently held on-chain by the signing EOA, grouped by condition ID.
	// Each LockedPosition.Market can be passed directly to MergePositions.
	// Requires a Polygon RPC URL (Alchemy endpoint) to be configured; returns
	// an empty slice (not an error) when no on-chain backend is available.
	GetLockedPositions() ([]LockedPosition, error)
}
