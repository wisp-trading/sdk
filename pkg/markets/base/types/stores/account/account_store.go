package account

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// AccountStore handles account balance and margin data per exchange
// Used by both spot and perp - balances are exchange-level, not instrument-level
type AccountStore interface {
	// Balance per exchange
	UpdateBalance(exchange connector.ExchangeName, balance connector.AssetBalance)
	GetBalance(exchange connector.ExchangeName) *connector.AssetBalance
	GetAllBalances() map[connector.ExchangeName]connector.AssetBalance

	UpdateMarginInfo(exchange connector.ExchangeName, margin MarginInfo)
	GetMarginInfo(exchange connector.ExchangeName) *MarginInfo
	GetAllMarginInfo() map[connector.ExchangeName]MarginInfo

	// Clear all data
	Clear()
}

// MarginInfo contains margin-specific account data
type MarginInfo struct {
	TotalMargin     numerical.Decimal `json:"total_margin"`
	UsedMargin      numerical.Decimal `json:"used_margin"`
	AvailableMargin numerical.Decimal `json:"available_margin"`
	MarginRatio     numerical.Decimal `json:"margin_ratio"`
	UpdatedAt       time.Time         `json:"updated_at"`
}
