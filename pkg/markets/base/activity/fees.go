// Package activity holds shared trading-activity helpers (fees, etc.).
package activity

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// SumTradeFees totals Fee across trades.
func SumTradeFees(trades []connector.Trade) numerical.Decimal {
	total := numerical.Zero()
	for _, t := range trades {
		total = total.Add(t.Fee)
	}
	return total
}
