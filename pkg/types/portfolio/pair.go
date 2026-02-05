package portfolio

type Pair struct {
	base      Asset
	quote     Asset
	separator string
}

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

func (p Pair) Base() Asset {
	return p.base
}

func (p Pair) Quote() Asset {
	return p.quote
}

func (p Pair) Symbol() string {
	return p.base.Symbol() + p.separator + p.quote.Symbol()
}

func (p Pair) Equals(other Pair) bool {
	return p.base.Equals(other.base) && p.quote.Equals(other.quote)
}
