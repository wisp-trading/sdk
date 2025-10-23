package portfolio

import (
	"database/sql/driver"
	"fmt"
)

type Asset struct {
	symbol string
}

func NewAsset(symbol string) Asset {
	return Asset{symbol: symbol}
}

func (a Asset) Symbol() string {
	return a.symbol
}

func (a Asset) Value() (driver.Value, error) {
	return a.symbol, nil
}

func (a *Asset) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		a.symbol = v
		return nil
	case []byte:
		a.symbol = string(v)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Asset", value)
	}
}

func (a Asset) IsValid() bool {
	return len(a.symbol) > 0 && len(a.symbol) <= 20
}

func (a Asset) Equals(other Asset) bool {
	return a.symbol == other.symbol
}
