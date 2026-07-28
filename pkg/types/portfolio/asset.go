package portfolio

import (
	"database/sql/driver"
	"encoding/json"
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
// Length allows short tickers (BTC) and EVM contract addresses (0x + 40 hex = 42).
func (a Asset) IsValid() bool {
	n := len(a.symbol)
	return n > 0 && n <= 66
}

// Equals checks if this asset is equal to another asset by comparing their symbols.
func (a Asset) Equals(other Asset) bool {
	return a.symbol == other.symbol
}

// MarshalJSON encodes Asset as a plain JSON string of its symbol.
func (a Asset) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.symbol)
}

// UnmarshalJSON decodes Asset from a plain JSON string back into the unexported symbol field.
func (a *Asset) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("asset: cannot unmarshal %s: %w", data, err)
	}
	a.symbol = s
	return nil
}
