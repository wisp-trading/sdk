package numerical

// LessThan checks if this Decimal is less than another
func (d Decimal) LessThan(other Decimal) bool {
	return d.Decimal.LessThan(other.Decimal)
}

// LessThanOrEqual checks if this Decimal is less than or equal to another
func (d Decimal) LessThanOrEqual(other Decimal) bool {
	return d.Decimal.LessThanOrEqual(other.Decimal)
}

// GreaterThan checks if this Decimal is greater than another
func (d Decimal) GreaterThan(other Decimal) bool {
	return d.Decimal.GreaterThan(other.Decimal)
}

// GreaterThanOrEqual checks if this Decimal is greater than or equal to another
func (d Decimal) GreaterThanOrEqual(other Decimal) bool {
	return d.Decimal.GreaterThanOrEqual(other.Decimal)
}

// Equal checks if this Decimal is equal to another
func (d Decimal) Equal(other Decimal) bool {
	return d.Decimal.Equal(other.Decimal)
}

// Cmp compares this Decimal to another, returning -1, 0, or 1
func (d Decimal) Cmp(other Decimal) int {
	return d.Decimal.Cmp(other.Decimal)
}

// IsZero checks if the Decimal is zero
func (d Decimal) IsZero() bool {
	return d.Decimal.IsZero()
}

// IsPositive checks if the Decimal is positive
func (d Decimal) IsPositive() bool {
	return d.Decimal.IsPositive()
}

// IsNegative checks if the Decimal is negative
func (d Decimal) IsNegative() bool {
	return d.Decimal.IsNegative()
}
