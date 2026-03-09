package portfolio

import "encoding/json"

// Pair represents a trading pair consisting of a base asset and a quote asset.
// The base asset is the asset being bought or sold, while the quote asset is
// the asset used to price the base asset.
type Pair struct {
	base      Asset
	quote     Asset
	separator string
}

// NewPair creates a new Pair with the given base and quote assets.
// An optional separator can be provided to customize how the pair symbol is formatted.
// If no separator is provided, a hyphen ("-") is used by default.
func NewPair(base, quote Asset, separator ...string) Pair {
	pair := Pair{
		base:  base,
		quote: quote,
	}

	if len(separator) > 0 {
		pair.separator = separator[0]
	} else {
		pair.separator = "-"
	}

	return pair
}

// Base returns the base asset of the pair.
func (p Pair) Base() Asset {
	return p.base
}

// Quote returns the quote asset of the pair.
func (p Pair) Quote() Asset {
	return p.quote
}

// Symbol returns the string representation of the pair by combining
// the base and quote asset symbols with the separator.
func (p Pair) Symbol() string {
	return p.base.Symbol() + p.separator + p.quote.Symbol()
}

// Equals checks if this pair is equal to another pair by comparing
// both the base and quote assets.
func (p Pair) Equals(other Pair) bool {
	return p.base.Equals(other.base) && p.quote.Equals(other.quote)
}

// pairJSON is the wire format for Pair — uses exported fields so encoding/json can handle it.
type pairJSON struct {
	Base      Asset  `json:"base"`
	Quote     Asset  `json:"quote"`
	Separator string `json:"separator"`
}

// MarshalJSON encodes Pair as {"base":"BTC","quote":"USDC","separator":"-"}.
func (p Pair) MarshalJSON() ([]byte, error) {
	return json.Marshal(pairJSON{
		Base:      p.base,
		Quote:     p.quote,
		Separator: p.separator,
	})
}

// UnmarshalJSON reconstructs a Pair from the wire format, restoring all unexported fields.
func (p *Pair) UnmarshalJSON(data []byte) error {
	var w pairJSON
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	p.base = w.Base
	p.quote = w.Quote
	p.separator = w.Separator
	if p.separator == "" {
		p.separator = "-"
	}
	return nil
}
