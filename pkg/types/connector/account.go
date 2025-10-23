package connector

import (
	"kronos/sdk/pkg/types/portfolio"
	"time"

	"github.com/shopspring/decimal"
)

// Position represents a trading position.
type Position struct {
	Symbol           portfolio.Asset `json:"symbol"`
	Exchange         ExchangeName    `json:"exchange"`
	Side             OrderSide       `json:"side"`
	Size             decimal.Decimal `json:"size"`
	EntryPrice       decimal.Decimal `json:"entry_price"`
	MarkPrice        decimal.Decimal `json:"mark_price"`
	UnrealizedPnL    decimal.Decimal `json:"unrealized_pnl"`
	RealizedPnL      decimal.Decimal `json:"realized_pnl"`
	Leverage         decimal.Decimal `json:"leverage,omitempty"`
	MarginType       string          `json:"margin_type,omitempty"` // "ISOLATED" or "CROSS"
	LiquidationPrice decimal.Decimal `json:"liquidation_price,omitempty"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// AccountBalance represents the balance and positions of an account.
type AccountBalance struct {
	TotalBalance     decimal.Decimal `json:"total_balance"`
	AvailableBalance decimal.Decimal `json:"available_balance"`
	UsedMargin       decimal.Decimal `json:"used_margin,omitempty"`
	UnrealizedPnL    decimal.Decimal `json:"unrealized_pnl,omitempty"`
	Currency         string          `json:"currency"`
	UpdatedAt        time.Time       `json:"updated_at"`
}
