package connector

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// AssetBalance represents the balance of a single asset in an account.
type AssetBalance struct {
	Asset     string            `json:"asset"`  // e.g., "BTC", "ETH", "USDT"
	Free      numerical.Decimal `json:"free"`   // Available balance
	Locked    numerical.Decimal `json:"locked"` // Balance locked in orders
	Total     numerical.Decimal `json:"total"`  // Total = Free + Locked
	UpdatedAt time.Time         `json:"updated_at"`
}
