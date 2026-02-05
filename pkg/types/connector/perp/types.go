package perp

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// RiskFundBalance represents the insurance/risk fund balance for a perpetual contract
// This is perp-specific as it's used to cover liquidations
type RiskFundBalance struct {
	Symbol    string            `json:"symbol"`
	Balance   numerical.Decimal `json:"balance"`
	Currency  string            `json:"currency"`
	UpdatedAt time.Time         `json:"updated_at"`
}
