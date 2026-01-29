package numerical

import (
	"github.com/shopspring/decimal"
)

// Decimal is the wrapped decimal type for all financial calculations in Wisp
// It embeds shopspring/decimal and ensures version consistency
type Decimal struct {
	decimal.Decimal
}

// NewFromInt creates a Decimal from an int64
func NewFromInt(value int64) Decimal {
	return Decimal{decimal.NewFromInt(value)}
}

// Zero returns a zero Decimal
func Zero() Decimal {
	return Decimal{decimal.Zero}
}

// NewFromString creates a Decimal from a string
func NewFromString(value string) (Decimal, error) {
	dec, err := decimal.NewFromString(value)
	return Decimal{dec}, err
}

// NewFromFloat creates a Decimal from a float64
func NewFromFloat(value float64) Decimal {
	return Decimal{decimal.NewFromFloat(value)}
}
