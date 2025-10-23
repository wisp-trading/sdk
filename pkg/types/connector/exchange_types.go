package connector

import (
	"time"

	"github.com/shopspring/decimal"
)

type RiskFundBalance struct {
	Symbol    string          `json:"symbol"`
	Balance   decimal.Decimal `json:"balance"`
	Currency  string          `json:"currency"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ContractInfo struct {
	Symbol       string          `json:"symbol"`
	BaseAsset    string          `json:"base_asset"`
	QuoteAsset   string          `json:"quote_asset"`
	ContractType string          `json:"contract_type"` // "PERPETUAL", "SPOT", "FUTURE"
	TickSize     decimal.Decimal `json:"tick_size"`
	StepSize     decimal.Decimal `json:"step_size"`
	MinOrderSize decimal.Decimal `json:"min_order_size"`
	MaxOrderSize decimal.Decimal `json:"max_order_size"`
	Status       string          `json:"status"` // "TRADING", "SUSPENDED", "DELISTED"
	UpdatedAt    time.Time       `json:"updated_at"`
}

type ExchangeStatus struct {
	Status    string    `json:"status"` // "NORMAL", "MAINTENANCE", "DEGRADED"
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
