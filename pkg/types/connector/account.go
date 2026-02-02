package connector

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// AccountBalance represents the balance and positions of an account.
type AccountBalance struct {
	TotalBalance     numerical.Decimal `json:"total_balance"`
	AvailableBalance numerical.Decimal `json:"available_balance"`
	UsedMargin       numerical.Decimal `json:"used_margin,omitempty"`
	UnrealizedPnL    numerical.Decimal `json:"unrealized_pnl,omitempty"`
	Currency         string            `json:"currency"`
	UpdatedAt        time.Time         `json:"updated_at"`
}
