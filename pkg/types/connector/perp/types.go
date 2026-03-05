package perp

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// ContractType represents the type of a perpetual contract.
type ContractType string

const (
	ContractTypePerpetual ContractType = "PERPETUAL"
	ContractTypeFuture    ContractType = "FUTURE"
)

// RiskFundBalance represents the insurance/risk fund balance for a perpetual contract
// This is perp-specific as it's used to cover liquidations
type RiskFundBalance struct {
	Symbol    string            `json:"symbol"`
	Balance   numerical.Decimal `json:"balance"`
	Currency  string            `json:"currency"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type ContractInfo struct {
	Symbol       string            `json:"symbol"`
	BaseAsset    string            `json:"base_asset"`
	QuoteAsset   string            `json:"quote_asset"`
	ContractType ContractType      `json:"contract_type"`
	TickSize     numerical.Decimal `json:"tick_size"`
	StepSize     numerical.Decimal `json:"step_size"`
	MinOrderSize numerical.Decimal `json:"min_order_size"`
	MaxOrderSize numerical.Decimal `json:"max_order_size"`
	Status       string            `json:"status"` // "TRADING", "SUSPENDED", "DELISTED"
	UpdatedAt    time.Time         `json:"updated_at"`
}
