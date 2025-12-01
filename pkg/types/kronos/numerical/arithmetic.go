package numerical

// Add returns the sum of this Decimal and another
func (d Decimal) Add(other Decimal) Decimal {
	return Decimal{d.Decimal.Add(other.Decimal)}
}

// Sub subtracts another Decimal from this one
func (d Decimal) Sub(other Decimal) Decimal {
	return Decimal{d.Decimal.Sub(other.Decimal)}
}

// Mul multiplies this Decimal by another
func (d Decimal) Mul(other Decimal) Decimal {
	return Decimal{d.Decimal.Mul(other.Decimal)}
}

// Div divides this Decimal by another
func (d Decimal) Div(other Decimal) Decimal {
	return Decimal{d.Decimal.Div(other.Decimal)}
}

// Mod returns the remainder of dividing this Decimal by another
func (d Decimal) Mod(other Decimal) Decimal {
	return Decimal{d.Decimal.Mod(other.Decimal)}
}

// Pow raises this Decimal to the power of the exponent
func (d Decimal) Pow(exponent Decimal) Decimal {
	return Decimal{d.Decimal.Pow(exponent.Decimal)}
}

// Abs returns the absolute value
func (d Decimal) Abs() Decimal {
	return Decimal{d.Decimal.Abs()}
}

// Neg returns the negation of this Decimal
func (d Decimal) Neg() Decimal {
	return Decimal{d.Decimal.Neg()}
}

// Round rounds this Decimal to the given number of decimal places
func (d Decimal) Round(places int32) Decimal {
	return Decimal{d.Decimal.Round(places)}
}

// Truncate truncates this Decimal to the given number of decimal places
func (d Decimal) Truncate(places int32) Decimal {
	return Decimal{d.Decimal.Truncate(places)}
}

// RoundUp rounds this Decimal up to the given number of decimal places
func (d Decimal) RoundUp(places int32) Decimal {
	return Decimal{d.Decimal.RoundUp(places)}
}

// RoundDown rounds this Decimal down to the given number of decimal places
func (d Decimal) RoundDown(places int32) Decimal {
	return Decimal{d.Decimal.RoundDown(places)}
}
