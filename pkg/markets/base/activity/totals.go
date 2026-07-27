package activity

import (
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// SumDecimals totals a slice of decimals.
func SumDecimals(vals []numerical.Decimal) numerical.Decimal {
	total := numerical.Zero()
	for _, v := range vals {
		total = total.Add(v)
	}
	return total
}
