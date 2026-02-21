package prediction

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type MarketID string

func (m MarketID) String() string {
	return string(m)
}

func MarketIDFromString(s string) MarketID {
	return MarketID(s)
}

type OutcomeID string

func (o OutcomeID) String() string {
	return string(o)
}

func OutcomeIDFromString(s string) OutcomeID {
	return OutcomeID(s)
}

// Position represents a position in a prediction market outcome
type Position struct {
	Market        Market                 `json:"market"`
	Exchange      connector.ExchangeName `json:"exchange"`
	Shares        numerical.Decimal      `json:"shares"`
	AvgCost       numerical.Decimal      `json:"avg_cost"`
	CurrentPrice  numerical.Decimal      `json:"current_price"`
	UnrealizedPnL numerical.Decimal      `json:"unrealized_pnl"`
	MaxPayout     numerical.Decimal      `json:"max_payout"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// OrderBook represents an order book for a specific market outcome
type OrderBook struct {
	Outcome   Market                 `json:"outcome"`
	Bids      []connector.PriceLevel `json:"bids"`
	Asks      []connector.PriceLevel `json:"asks"`
	Timestamp time.Time              `json:"timestamp"`
}

type PriceChange struct {
	Outcome   Outcome           `json:"outcome"`
	Timestamp string            `json:"timestamp"`
	Price     numerical.Decimal `json:"price"`
	Size      string            `json:"size"`
	Side      string            `json:"side"`
	BestBid   string            `json:"best_bid"`
	BestAsk   string            `json:"best_ask"`
}

//
//// Trade represents a trade on a prediction market
//type Trade struct {
//	ID        string                 `json:"id"`
//	Pair      Market                 `json:"pair"`
//	Exchange  connector.ExchangeName `json:"exchange"`
//	Side      connector.OrderSide    `json:"side"`
//	Shares    numerical.Decimal      `json:"shares"`
//	Price     numerical.Decimal      `json:"price"`
//	Fee       numerical.Decimal      `json:"fee"`
//	Timestamp time.Time              `json:"timestamp"`
//}
//
//// Settlement represents the settlement of a resolved market
//type Settlement struct {
//	Pair             Market                 `json:"pair"`
//	Exchange         connector.ExchangeName `json:"exchange"`
//	WinningOutcome   string                 `json:"winning_outcome"`
//	Shares           numerical.Decimal      `json:"shares"`
//	Payout           numerical.Decimal      `json:"payout"`
//	SettlementTime   time.Time              `json:"settlement_time"`
//	ResolutionSource string                 `json:"resolution_source"`
//}
//
//// Resolution represents a market resolution event
//type Resolution struct {
//	MarketID         string    `json:"market_id"`
//	WinningOutcome   string    `json:"winning_outcome"`
//	ResolutionTime   time.Time `json:"resolution_time"`
//	ResolutionSource string    `json:"resolution_source"`
//}
//
//// PricePoint represents a historical price data point for an outcome
//type PricePoint struct {
//	Pair      Market            `json:"pair"`
//	Price     numerical.Decimal `json:"price"`
//	Volume    numerical.Decimal `json:"volume"`
//	Timestamp time.Time         `json:"timestamp"`
//}
