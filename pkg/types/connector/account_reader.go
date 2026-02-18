package connector

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// AccountReader provides account information
type AccountReader interface {
	GetBalances() ([]AssetBalance, error)
	GetBalance(asset portfolio.Asset) (*AssetBalance, error)
	GetTradingHistory(pair portfolio.Pair, limit int) ([]Trade, error)
}

// AssetBalance represents the balance of a single asset in an account.
type AssetBalance struct {
	Asset     portfolio.Asset   `json:"asset"`  // e.g., "BTC", "ETH", "USDT"
	Free      numerical.Decimal `json:"free"`   // Available balance
	Locked    numerical.Decimal `json:"locked"` // Balance locked in orders
	Total     numerical.Decimal `json:"total"`  // Total = Free + Locked
	UpdatedAt time.Time         `json:"updated_at"`
}
