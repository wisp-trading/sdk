package connector

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Position represents a trading position.
type Position struct {
	Symbol           portfolio.Asset   `json:"symbol"`
	Exchange         ExchangeName      `json:"exchange"`
	Side             OrderSide         `json:"side"`
	Size             numerical.Decimal `json:"size"`
	EntryPrice       numerical.Decimal `json:"entry_price"`
	MarkPrice        numerical.Decimal `json:"mark_price"`
	UnrealizedPnL    numerical.Decimal `json:"unrealized_pnl"`
	RealizedPnL      numerical.Decimal `json:"realized_pnl"`
	Leverage         numerical.Decimal `json:"leverage,omitempty"`
	MarginType       string            `json:"margin_type,omitempty"` // "ISOLATED" or "CROSS"
	LiquidationPrice numerical.Decimal `json:"liquidation_price,omitempty"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

// AccountBalance represents the balance and positions of an account.
type AccountBalance struct {
	TotalBalance     numerical.Decimal `json:"total_balance"`
	AvailableBalance numerical.Decimal `json:"available_balance"`
	UsedMargin       numerical.Decimal `json:"used_margin,omitempty"`
	UnrealizedPnL    numerical.Decimal `json:"unrealized_pnl,omitempty"`
	Currency         string            `json:"currency"`
	UpdatedAt        time.Time         `json:"updated_at"`
}
