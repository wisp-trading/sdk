package portfolio

import (
	"database/sql/driver"
	"fmt"
)

// Asset represents a financial asset identified by its symbol.
// It implements the sql.Scanner and driver.Valuer interfaces for database operations.
type Asset struct {
	symbol string
}

// NewAsset creates a new Asset with the given symbol.
func NewAsset(symbol string) Asset {
	return Asset{symbol: symbol}
}

// Symbol returns the symbol identifier of the asset.
func (a Asset) Symbol() string {
	return a.symbol
}

// Value implements the driver.Valuer interface, allowing the Asset
// to be stored in a database by returning its symbol as the value.
func (a Asset) Value() (driver.Value, error) {
	return a.symbol, nil
}

// Scan implements the sql.Scanner interface, allowing the Asset
// to be populated from a database value. It accepts string or []byte values.
// Returns an error if the value is of an unsupported type.
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
		return fmt.Errorf("cannot scan %T into Pair", value)
	}
}

// IsValid checks if the asset has a valid symbol.
// A valid symbol must be between 1 and 20 characters in length.
func (a Asset) IsValid() bool {
	return len(a.symbol) > 0 && len(a.symbol) <= 20
}

// Equals checks if this asset is equal to another asset by comparing their symbols.
func (a Asset) Equals(other Asset) bool {
	return a.symbol == other.symbol
}
